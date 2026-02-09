#!/bin/bash

# Render Deployment Quick Start Script
# This script helps you prepare and validate your deployment to Render

set -e

echo "🚀 Render Deployment Quick Start"
echo "=================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_info() {
    echo -e "ℹ️  $1"
}

# Check if required files exist
echo "📋 Step 1: Checking required files..."
echo ""

files_ok=true

if [ -f "Dockerfile" ]; then
    print_success "Dockerfile found"
else
    print_error "Dockerfile not found"
    files_ok=false
fi

if [ -f "render.yaml" ]; then
    print_success "render.yaml found"
else
    print_error "render.yaml not found"
    files_ok=false
fi

if [ -f "go.mod" ]; then
    print_success "go.mod found"
else
    print_error "go.mod not found"
    files_ok=false
fi

if [ -f ".env.example" ]; then
    print_success ".env.example found"
else
    print_warning ".env.example not found (optional)"
fi

if [ -f ".gitignore" ]; then
    if grep -q "^.env$" .gitignore; then
        print_success ".env is in .gitignore"
    else
        print_warning ".env not in .gitignore - add it to prevent committing secrets"
    fi
else
    print_warning ".gitignore not found"
fi

echo ""

# Check Go installation
echo "🔧 Step 2: Checking Go installation..."
echo ""

if command -v go &> /dev/null; then
    go_version=$(go version | awk '{print $3}')
    print_success "Go is installed: $go_version"
else
    print_error "Go is not installed. Install from https://golang.org/dl/"
    exit 1
fi

echo ""

# Verify Go modules
echo "📦 Step 3: Verifying Go modules..."
echo ""

if go mod verify; then
    print_success "Go modules verified"
else
    print_error "Go modules verification failed"
    print_info "Run: go mod tidy"
    exit 1
fi

echo ""

# Check Git status
echo "📝 Step 4: Checking Git status..."
echo ""

if git rev-parse --git-dir > /dev/null 2>&1; then
    print_success "Git repository detected"
    
    # Check if there are uncommitted changes
    if [[ -n $(git status -s) ]]; then
        print_warning "You have uncommitted changes:"
        git status -s
        echo ""
        print_info "Commit your changes before deploying"
    else
        print_success "No uncommitted changes"
    fi
    
    # Check current branch
    current_branch=$(git branch --show-current)
    print_info "Current branch: $current_branch"
    
    if [ "$current_branch" != "main" ]; then
        print_warning "You're not on the 'main' branch"
        print_info "Render auto-deploys from 'main' branch by default"
    fi
else
    print_error "Not a Git repository"
    print_info "Initialize with: git init"
    exit 1
fi

echo ""

# Check Docker (optional)
echo "🐳 Step 5: Checking Docker (optional for local testing)..."
echo ""

if command -v docker &> /dev/null; then
    docker_version=$(docker --version)
    print_success "Docker is installed: $docker_version"
    
    if command -v docker-compose &> /dev/null; then
        print_success "Docker Compose is installed"
        print_info "You can test locally with: make docker-run"
    else
        print_warning "Docker Compose not found (optional)"
    fi
else
    print_warning "Docker not installed (optional for local testing)"
    print_info "Install from https://www.docker.com/get-started"
fi

echo ""

# Environment variables check
echo "🔐 Step 6: Checking environment variables..."
echo ""

if [ -f ".env" ]; then
    print_warning ".env file found - make sure it's in .gitignore"
    print_info "Never commit .env to Git!"
else
    print_info "No .env file found (will use environment variables)"
fi

required_vars=("APPWRITE_API_KEY")
missing_vars=()

for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        missing_vars+=("$var")
    fi
done

if [ ${#missing_vars[@]} -eq 0 ]; then
    print_success "All required environment variables are set"
else
    print_warning "Missing environment variables (set these in Render dashboard):"
    for var in "${missing_vars[@]}"; do
        echo "  - $var"
    done
fi

echo ""

# Build test
echo "🔨 Step 7: Testing build..."
echo ""

print_info "Building application..."
if CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /tmp/gold-go-test ./cmd/main.go; then
    print_success "Build successful"
    rm -f /tmp/gold-go-test
else
    print_error "Build failed"
    exit 1
fi

echo ""

# Summary
echo "📊 Summary"
echo "=========="
echo ""

if [ "$files_ok" = true ]; then
    print_success "All required files are present"
else
    print_error "Some required files are missing"
    exit 1
fi

echo ""
echo "🎯 Next Steps:"
echo ""
echo "1. Commit and push your changes:"
echo "   git add ."
echo "   git commit -m 'Deploy to Render'"
echo "   git push origin main"
echo ""
echo "2. Go to Render Dashboard:"
echo "   https://dashboard.render.com"
echo ""
echo "3. Create a new Blueprint:"
echo "   - Click 'New +' → 'Blueprint'"
echo "   - Select your repository"
echo "   - Render will detect render.yaml"
echo ""
echo "4. Set APPWRITE_API_KEY:"
echo "   - In service settings → Environment"
echo "   - Add your Appwrite API key securely"
echo ""
echo "5. Deploy:"
echo "   - Click 'Apply' to create services"
echo "   - Monitor deployment in Render dashboard"
echo ""
echo "📚 For detailed instructions, see DEPLOYMENT.md"
echo ""
print_success "Ready to deploy! 🚀"
