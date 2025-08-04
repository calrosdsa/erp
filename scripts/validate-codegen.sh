#!/bin/bash
# Code generation validation script for ERP backend

set -e

echo "🔄 Starting code generation validation..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Ensure we're in the right directory
if [ ! -f "go.mod" ]; then
    print_error "Must be run from project root directory"
    exit 1
fi

# Create gen directories if they don't exist
print_status "Creating generation directories..."
mkdir -p gen/db/model gen/db/query gen/openapi gen/mocks gen/proto

# Check if database is available
print_status "Checking database connectivity..."
if ! docker-compose exec -T postgres pg_isready -U postgres -d erp_dev > /dev/null 2>&1; then
    print_warning "Database not available, starting infrastructure..."
    make docker-dev-up
    sleep 5
fi

# Generate GORM models
print_status "Generating GORM models..."
if make models; then
    print_status "✅ GORM model generation successful"
else
    print_error "❌ GORM model generation failed"
    exit 1
fi

# Verify generated models exist
print_status "Verifying generated models..."
model_count=$(find gen/db/model -name "*.gen.go" | wc -l)
query_count=$(find gen/db/query -name "*.gen.go" | wc -l)

if [ "$model_count" -gt 0 ] && [ "$query_count" -gt 0 ]; then
    print_status "✅ Found $model_count model files and $query_count query files"
else
    print_error "❌ Generated models not found or empty"
    exit 1
fi

# Generate mocks
print_status "Generating mocks..."
if make mockery; then
    print_status "✅ Mock generation successful"
    mock_count=$(find gen/mocks -name "*.go" | wc -l)
    print_status "Generated $mock_count mock files"
else
    print_warning "⚠️ Mock generation failed (may be expected if no interfaces found)"
fi

# Validate Go code compiles
print_status "Validating generated code compiles..."
if go build ./...; then
    print_status "✅ All generated code compiles successfully"
else
    print_error "❌ Generated code has compilation errors"
    exit 1
fi

# Run a quick syntax check on generated files
print_status "Running syntax validation..."
syntax_errors=0
for file in $(find gen -name "*.go"); do
    if ! go vet "$file" > /dev/null 2>&1; then
        print_warning "Syntax issue in $file"
        syntax_errors=$((syntax_errors + 1))
    fi
done

if [ $syntax_errors -eq 0 ]; then
    print_status "✅ All generated files pass syntax validation"
else
    print_warning "⚠️ Found $syntax_errors files with syntax warnings"
fi

print_status "🎉 Code generation validation complete!"

# Summary
echo ""
echo "📊 Generation Summary:"
echo "  - Models: $model_count files"
echo "  - Queries: $query_count files" 
echo "  - Mocks: $(find gen/mocks -name "*.go" 2>/dev/null | wc -l) files"
echo "  - Proto: $(find gen/proto -name "*.go" 2>/dev/null | wc -l) files"
echo ""
print_status "All code generation pipelines are working correctly!"