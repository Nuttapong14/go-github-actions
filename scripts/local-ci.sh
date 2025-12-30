#!/bin/bash
#
# Local CI Simulation Script for CI/CD Training Application
# Runs the same checks that execute in GitHub Actions CI pipeline
#
# Usage: ./scripts/local-ci.sh [options]
#
# Options:
#   --lint-only     Run only linting checks
#   --test-only     Run only tests
#   --build-only    Run only build
#   --skip-lint     Skip linting checks
#   --skip-test     Skip tests
#   --skip-build    Skip build
#   --docker        Also build Docker images
#   --verbose       Enable verbose output
#   --help          Show this help message

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Default options
RUN_LINT=true
RUN_TEST=true
RUN_BUILD=true
RUN_DOCKER=false
VERBOSE=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --lint-only)
            RUN_LINT=true
            RUN_TEST=false
            RUN_BUILD=false
            shift
            ;;
        --test-only)
            RUN_LINT=false
            RUN_TEST=true
            RUN_BUILD=false
            shift
            ;;
        --build-only)
            RUN_LINT=false
            RUN_TEST=false
            RUN_BUILD=true
            shift
            ;;
        --skip-lint)
            RUN_LINT=false
            shift
            ;;
        --skip-test)
            RUN_TEST=false
            shift
            ;;
        --skip-build)
            RUN_BUILD=false
            shift
            ;;
        --docker)
            RUN_DOCKER=true
            shift
            ;;
        --verbose|-v)
            VERBOSE=true
            shift
            ;;
        --help|-h)
            head -24 "$0" | tail -20
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            exit 1
            ;;
    esac
done

# Change to project root
cd "$PROJECT_ROOT"

# Git information
GIT_SHA=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
GIT_TAG=$(git describe --tags --exact-match 2>/dev/null || echo "")

# Timing
START_TIME=$(date +%s)

# Results tracking
LINT_RESULT=0
TEST_RESULT=0
BUILD_RESULT=0
DOCKER_RESULT=0

# Logging functions
log_step() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${CYAN}▶ $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

log_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

log_failure() {
    echo -e "${RED}✗ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

log_info() {
    echo -e "${CYAN}ℹ $1${NC}"
}

# Header
echo ""
echo -e "${CYAN}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║     CI/CD Training Application - Local CI Pipeline       ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "Git SHA:     ${YELLOW}${GIT_SHA}${NC}"
echo -e "Git Branch:  ${YELLOW}${GIT_BRANCH}${NC}"
echo -e "Git Tag:     ${YELLOW}${GIT_TAG:-none}${NC}"
echo -e "Date:        ${YELLOW}$(date '+%Y-%m-%d %H:%M:%S')${NC}"
echo ""

# Step 1: Lint
if [ "$RUN_LINT" = true ]; then
    log_step "Step 1/3: Running Linters"

    # Backend lint
    log_info "Running golangci-lint on backend..."
    if cd backend && golangci-lint run --config ../.golangci.yml ./... 2>&1; then
        log_success "Backend linting passed"
    else
        log_failure "Backend linting failed"
        LINT_RESULT=1
    fi
    cd "$PROJECT_ROOT"

    # Frontend lint
    log_info "Running ESLint on frontend..."
    if cd frontend && bun run lint 2>&1; then
        log_success "Frontend linting passed"
    else
        log_warning "Frontend linting has issues (continuing anyway)"
        # Don't fail on frontend lint for now
    fi
    cd "$PROJECT_ROOT"
else
    log_info "Skipping lint step"
fi

# Step 2: Test
if [ "$RUN_TEST" = true ]; then
    log_step "Step 2/3: Running Tests"

    # Backend tests
    log_info "Running Go tests with coverage..."
    if cd backend && go test -v -race -coverprofile=coverage.out ./... 2>&1; then
        log_success "Backend tests passed"

        # Generate coverage report
        go tool cover -func=coverage.out | tail -1
        go tool cover -html=coverage.out -o coverage.html
        log_success "Coverage report generated: backend/coverage.html"
    else
        log_failure "Backend tests failed"
        TEST_RESULT=1
    fi
    cd "$PROJECT_ROOT"

    # Frontend tests
    log_info "Running frontend tests..."
    if cd frontend && bun test 2>&1; then
        log_success "Frontend tests passed"
    else
        log_warning "Frontend tests have issues (continuing anyway)"
    fi
    cd "$PROJECT_ROOT"
else
    log_info "Skipping test step"
fi

# Step 3: Build
if [ "$RUN_BUILD" = true ]; then
    log_step "Step 3/3: Building Artifacts"

    # Backend build
    log_info "Building backend binary..."
    IMAGE_TAG="${GIT_TAG:-$GIT_SHA}"
    mkdir -p bin
    if cd backend && go build \
        -ldflags "-X main.Version=$IMAGE_TAG -X main.GitSHA=$GIT_SHA -X main.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
        -o ../bin/backend main.go 2>&1; then
        log_success "Backend binary built: bin/backend"
    else
        log_failure "Backend build failed"
        BUILD_RESULT=1
    fi
    cd "$PROJECT_ROOT"

    # Frontend build
    log_info "Building frontend..."
    if cd frontend && bun run build 2>&1; then
        log_success "Frontend built: frontend/.next/"
    else
        log_failure "Frontend build failed"
        BUILD_RESULT=1
    fi
    cd "$PROJECT_ROOT"
else
    log_info "Skipping build step"
fi

# Step 4: Docker (optional)
if [ "$RUN_DOCKER" = true ]; then
    log_step "Step 4: Building Docker Images"

    IMAGE_TAG="${GIT_TAG:-$GIT_SHA}"
    BACKEND_IMAGE="cicd-training-backend:$IMAGE_TAG"
    FRONTEND_IMAGE="cicd-training-frontend:$IMAGE_TAG"

    log_info "Building backend Docker image..."
    if docker build -t "$BACKEND_IMAGE" \
        --build-arg GIT_SHA="$GIT_SHA" \
        --build-arg BUILD_DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
        -f docker/backend.Dockerfile . 2>&1; then
        log_success "Backend image built: $BACKEND_IMAGE"
    else
        log_failure "Backend Docker build failed"
        DOCKER_RESULT=1
    fi

    log_info "Building frontend Docker image..."
    if docker build -t "$FRONTEND_IMAGE" \
        -f docker/frontend.Dockerfile . 2>&1; then
        log_success "Frontend image built: $FRONTEND_IMAGE"
    else
        log_failure "Frontend Docker build failed"
        DOCKER_RESULT=1
    fi
fi

# Calculate duration
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))
MINUTES=$((DURATION / 60))
SECONDS=$((DURATION % 60))

# Summary
echo ""
echo -e "${CYAN}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${CYAN}║                    Pipeline Summary                       ║${NC}"
echo -e "${CYAN}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""

# Results table
echo "┌─────────────┬─────────────┐"
echo "│ Stage       │ Result      │"
echo "├─────────────┼─────────────┤"

if [ "$RUN_LINT" = true ]; then
    if [ $LINT_RESULT -eq 0 ]; then
        echo -e "│ Lint        │ ${GREEN}✓ Passed${NC}    │"
    else
        echo -e "│ Lint        │ ${RED}✗ Failed${NC}    │"
    fi
fi

if [ "$RUN_TEST" = true ]; then
    if [ $TEST_RESULT -eq 0 ]; then
        echo -e "│ Test        │ ${GREEN}✓ Passed${NC}    │"
    else
        echo -e "│ Test        │ ${RED}✗ Failed${NC}    │"
    fi
fi

if [ "$RUN_BUILD" = true ]; then
    if [ $BUILD_RESULT -eq 0 ]; then
        echo -e "│ Build       │ ${GREEN}✓ Passed${NC}    │"
    else
        echo -e "│ Build       │ ${RED}✗ Failed${NC}    │"
    fi
fi

if [ "$RUN_DOCKER" = true ]; then
    if [ $DOCKER_RESULT -eq 0 ]; then
        echo -e "│ Docker      │ ${GREEN}✓ Passed${NC}    │"
    else
        echo -e "│ Docker      │ ${RED}✗ Failed${NC}    │"
    fi
fi

echo "└─────────────┴─────────────┘"
echo ""
echo -e "Duration: ${YELLOW}${MINUTES}m ${SECONDS}s${NC}"
echo ""

# Overall result
OVERALL_RESULT=$((LINT_RESULT + TEST_RESULT + BUILD_RESULT + DOCKER_RESULT))

if [ $OVERALL_RESULT -eq 0 ]; then
    echo -e "${GREEN}╔══════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║              ✓ CI Pipeline Passed Successfully           ║${NC}"
    echo -e "${GREEN}╚══════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo "Artifacts:"
    [ -f bin/backend ] && echo "  • Backend binary: bin/backend"
    [ -d frontend/.next ] && echo "  • Frontend build: frontend/.next/"
    [ -f backend/coverage.out ] && echo "  • Coverage data: backend/coverage.out"
    [ -f backend/coverage.html ] && echo "  • Coverage report: backend/coverage.html"
    echo ""
    exit 0
else
    echo -e "${RED}╔══════════════════════════════════════════════════════════╗${NC}"
    echo -e "${RED}║              ✗ CI Pipeline Failed                        ║${NC}"
    echo -e "${RED}╚══════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo "Please review the errors above and fix them before pushing."
    echo ""
    exit 1
fi
