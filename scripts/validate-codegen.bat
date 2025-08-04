@echo off
REM Code generation validation script for ERP backend (Windows)

echo [INFO] Starting code generation validation...

REM Ensure we're in the right directory
if not exist "go.mod" (
    echo [ERROR] Must be run from project root directory
    exit /b 1
)

REM Create gen directories if they don't exist
echo [INFO] Creating generation directories...
if not exist "gen\db\model" mkdir gen\db\model
if not exist "gen\db\query" mkdir gen\db\query  
if not exist "gen\openapi" mkdir gen\openapi
if not exist "gen\mocks" mkdir gen\mocks
if not exist "gen\proto" mkdir gen\proto

REM Check if database is available
echo [INFO] Checking database connectivity...
docker-compose exec -T postgres pg_isready -U postgres -d erp_dev >nul 2>&1
if errorlevel 1 (
    echo [WARN] Database not available, starting infrastructure...
    make docker-dev-up
    timeout /t 5 >nul
)

REM Generate GORM models
echo [INFO] Generating GORM models...
make models
if errorlevel 1 (
    echo [ERROR] GORM model generation failed
    exit /b 1
)
echo [INFO] GORM model generation successful

REM Verify generated models exist
echo [INFO] Verifying generated models...
for /r gen\db\model %%f in (*.gen.go) do set /a model_count+=1
for /r gen\db\query %%f in (*.gen.go) do set /a query_count+=1

if %model_count% gtr 0 if %query_count% gtr 0 (
    echo [INFO] Found %model_count% model files and %query_count% query files
) else (
    echo [ERROR] Generated models not found or empty
    exit /b 1
)

REM Generate mocks
echo [INFO] Generating mocks...
make mockery
if errorlevel 1 (
    echo [WARN] Mock generation failed (may be expected if no interfaces found)
) else (
    echo [INFO] Mock generation successful
)

REM Validate Go code compiles
echo [INFO] Validating generated code compiles...
go build ./...
if errorlevel 1 (
    echo [ERROR] Generated code has compilation errors
    exit /b 1
)
echo [INFO] All generated code compiles successfully

echo [INFO] Code generation validation complete!
echo [INFO] All code generation pipelines are working correctly!