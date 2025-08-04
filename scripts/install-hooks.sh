#!/bin/bash
# Git hooks installation script for ERP backend

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

print_header() {
    echo -e "${BLUE}[HOOKS]${NC} $1"
}

# Ensure we're in a git repository
if [ ! -d ".git" ]; then
    print_error "This script must be run from the root of a git repository"
    exit 1
fi

print_header "Installing Git Hooks for ERP Backend Development"

# Check if pre-commit is installed
if ! command -v pre-commit &> /dev/null; then
    print_status "Installing pre-commit..."
    
    # Try different installation methods
    if command -v pip3 &> /dev/null; then
        pip3 install pre-commit
    elif command -v pip &> /dev/null; then
        pip install pre-commit
    elif command -v brew &> /dev/null; then
        brew install pre-commit
    elif command -v curl &> /dev/null; then
        # Install using curl (direct download)
        curl https://pre-commit.com/install-local.py | python3 -
    else
        print_error "Cannot install pre-commit. Please install it manually:"
        echo "  pip install pre-commit"
        echo "  # or"
        echo "  brew install pre-commit"
        exit 1
    fi
else
    print_status "pre-commit is already installed"
fi

# Install Go tools if not present
print_status "Installing required Go tools..."

# golangci-lint
if ! command -v golangci-lint &> /dev/null; then
    print_status "Installing golangci-lint..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
fi

# goimports
if ! command -v goimports &> /dev/null; then
    print_status "Installing goimports..."
    go install golang.org/x/tools/cmd/goimports@latest
fi

# gosec
if ! command -v gosec &> /dev/null; then
    print_status "Installing gosec..."
    go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
fi

# Install pre-commit hooks
print_status "Installing pre-commit hooks..."
pre-commit install

# Install commit-msg hook for conventional commits (optional)
print_status "Installing commit-msg hook..."
pre-commit install --hook-type commit-msg

# Run hooks on all files to test installation
print_status "Testing hooks installation..."
if pre-commit run --all-files; then
    print_status "✅ All hooks passed successfully!"
else
    print_warning "⚠️ Some hooks failed. This is normal on first run."
    print_status "The hooks will now run on each commit."
fi

# Create a simple manual hook script as backup
print_status "Creating backup manual hook script..."
cat > scripts/run-hooks-manual.sh << 'EOF'
#!/bin/bash
# Manual hook execution script
# Run this if you want to check your code before committing

echo "🔍 Running code quality checks..."

echo "📝 Formatting Go code..."
go fmt ./...

echo "🔧 Organizing imports..."
if command -v goimports &> /dev/null; then
    goimports -w .
fi

echo "🔍 Running go vet..."
go vet ./...

echo "📊 Running golangci-lint..."
if command -v golangci-lint &> /dev/null; then
    golangci-lint run
fi

echo "🧪 Running unit tests..."
go test -short ./...

echo "✅ All checks completed!"
EOF

chmod +x scripts/run-hooks-manual.sh

print_header "Git Hooks Installation Complete!"
echo ""
echo "📋 What was installed:"
echo "  ✅ pre-commit framework"
echo "  ✅ golangci-lint (Go linter)"
echo "  ✅ goimports (import organizer)"
echo "  ✅ gosec (security scanner)"
echo "  ✅ Pre-commit hooks configured"
echo ""
echo "🚀 Next steps:"
echo "  • Hooks will run automatically on each commit"
echo "  • Run 'pre-commit run --all-files' to check all files"
echo "  • Run 'scripts/run-hooks-manual.sh' for manual checking"
echo "  • Configure your editor to run gofmt on save"
echo ""
print_status "Happy coding! 🎉"