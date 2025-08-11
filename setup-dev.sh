#!/bin/bash

# Development setup script for local-k8s-cluster project

set -e

echo "🚀 Setting up local-k8s-cluster development environment..."

# Check if we're in the right directory
if [ ! -f "decision_log.md" ]; then
    echo "❌ Please run this script from the project root directory"
    exit 1
fi

# Install pre-commit if not installed
if ! command -v pre-commit &> /dev/null; then
    echo "📦 Installing pre-commit..."
    if command -v pip3 &> /dev/null; then
        pip3 install pre-commit
    elif command -v brew &> /dev/null; then
        brew install pre-commit
    else
        echo "❌ Please install pre-commit manually: https://pre-commit.com/#installation"
        exit 1
    fi
fi

# Install golangci-lint if not installed
if ! command -v golangci-lint &> /dev/null; then
    echo "📦 Installing golangci-lint..."
    if command -v brew &> /dev/null; then
        brew install golangci-lint
    else
        echo "Installing golangci-lint via curl..."
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
    fi
fi

# Install goimports if not installed
if ! command -v goimports &> /dev/null; then
    echo "📦 Installing goimports..."
    go install golang.org/x/tools/cmd/goimports@latest
fi

# Setup Go project dependencies
echo "📦 Setting up Go dependencies..."
cd local-k8s-cluster-go
go mod download
go mod verify
cd ..

# Install pre-commit hooks
echo "🪝 Installing pre-commit hooks..."
pre-commit install
pre-commit install --hook-type commit-msg

# Setup git commit template
echo "📝 Setting up git commit template..."
git config commit.template .gitmessage

# Create secrets baseline for detect-secrets
if [ ! -f ".secrets.baseline" ]; then
    echo "🔒 Creating secrets baseline..."
    detect-secrets scan --baseline .secrets.baseline
fi

# Setup conventional commits helper
if command -v npm &> /dev/null; then
    echo "📦 Installing commitizen for conventional commits..."
    npm install -g commitizen cz-conventional-changelog
    echo '{ "path": "cz-conventional-changelog" }' > .czrc
fi

echo "✅ Development environment setup complete!"
echo ""
echo "📋 Next steps:"
echo "  1. Run 'git cz' instead of 'git commit' for conventional commits"
echo "  2. Pre-commit hooks will run automatically on commits"
echo "  3. Use 'make help' in local-k8s-cluster-go/ for available commands"
echo "  4. Run 'pre-commit run --all-files' to test all hooks"
echo ""
echo "🔧 Useful commands:"
echo "  make build          # Build the Go binary"
echo "  make test           # Run tests"  
echo "  make lint           # Run linting"
echo "  git cz              # Conventional commit helper"
echo "  pre-commit run -a   # Run all pre-commit hooks"