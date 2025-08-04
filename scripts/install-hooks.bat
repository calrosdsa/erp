@echo off
REM Git hooks installation script for ERP backend (Windows)

echo [HOOKS] Installing Git Hooks for ERP Backend Development

REM Ensure we're in a git repository
if not exist ".git" (
    echo [ERROR] This script must be run from the root of a git repository
    exit /b 1
)

REM Check if pre-commit is installed
where pre-commit >nul 2>&1
if errorlevel 1 (
    echo [INFO] Installing pre-commit...
    pip install pre-commit
    if errorlevel 1 (
        echo [ERROR] Failed to install pre-commit. Please install it manually:
        echo   pip install pre-commit
        exit /b 1
    )
) else (
    echo [INFO] pre-commit is already installed
)

REM Install Go tools if not present
echo [INFO] Installing required Go tools...

REM Check golangci-lint
where golangci-lint >nul 2>&1
if errorlevel 1 (
    echo [INFO] Installing golangci-lint...
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
)

REM Check goimports
where goimports >nul 2>&1
if errorlevel 1 (
    echo [INFO] Installing goimports...
    go install golang.org/x/tools/cmd/goimports@latest
)

REM Check gosec
where gosec >nul 2>&1
if errorlevel 1 (
    echo [INFO] Installing gosec...
    go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
)

REM Install pre-commit hooks
echo [INFO] Installing pre-commit hooks...
pre-commit install
if errorlevel 1 (
    echo [ERROR] Failed to install pre-commit hooks
    exit /b 1
)

REM Install commit-msg hook
echo [INFO] Installing commit-msg hook...
pre-commit install --hook-type commit-msg

REM Test hooks installation
echo [INFO] Testing hooks installation...
pre-commit run --all-files
if errorlevel 1 (
    echo [WARN] Some hooks failed. This is normal on first run.
)

REM Create manual hook script
echo [INFO] Creating backup manual hook script...
(
echo @echo off
echo REM Manual hook execution script
echo echo [INFO] Running code quality checks...
echo.
echo echo [INFO] Formatting Go code...
echo go fmt ./...
echo.
echo echo [INFO] Organizing imports...
echo where goimports ^>nul 2^>^&1
echo if not errorlevel 1 goimports -w .
echo.
echo echo [INFO] Running go vet...
echo go vet ./...
echo.
echo echo [INFO] Running golangci-lint...
echo where golangci-lint ^>nul 2^>^&1
echo if not errorlevel 1 golangci-lint run
echo.
echo echo [INFO] Running unit tests...
echo go test -short ./...
echo.
echo echo [INFO] All checks completed!
) > scripts\run-hooks-manual.bat

echo [HOOKS] Git Hooks Installation Complete!
echo.
echo What was installed:
echo   - pre-commit framework
echo   - golangci-lint (Go linter)
echo   - goimports (import organizer)  
echo   - gosec (security scanner)
echo   - Pre-commit hooks configured
echo.
echo Next steps:
echo   - Hooks will run automatically on each commit
echo   - Run 'pre-commit run --all-files' to check all files
echo   - Run 'scripts\run-hooks-manual.bat' for manual checking
echo.
echo [INFO] Happy coding!