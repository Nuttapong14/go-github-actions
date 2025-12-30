#!/bin/bash
#
# Prerequisite Check Script for CI/CD Training Application
# Verifies all required tools are installed and properly configured

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
PASSED=0
FAILED=0
WARNINGS=0

# Check function
check() {
    local name="$1"
    local command="$2"
    local min_version="${3:-}"
    local required="${4:-true}"

    printf "Checking %-20s " "$name..."

    if ! command -v "$command" &> /dev/null; then
        if [ "$required" = "true" ]; then
            echo -e "${RED}FAILED${NC} - Not found"
            FAILED=$((FAILED + 1))

        else
            echo -e "${YELLOW}WARNING${NC} - Not found (optional)"
            WARNINGS=$((WARNINGS + 1))
        fi
        return 1
    fi

    local version
    case "$command" in
        docker)
            version=$(docker --version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)
            ;;
        go)
            version=$(go version 2>/dev/null | grep -oE 'go[0-9]+\.[0-9]+' | sed 's/go//')
            ;;
        node)
            version=$(node --version 2>/dev/null | sed 's/v//')
            ;;
        bun)
            version=$(bun --version 2>/dev/null)
            ;;
        git)
            version=$(git --version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')
            ;;
        *)
            version=$($command --version 2>/dev/null | head -1)
            ;;
    esac

    echo -e "${GREEN}OK${NC} - Version: $version"
    PASSED=$((PASSED + 1))
    return 0
}

# Header
echo "=============================================="
echo "CI/CD Training Application - Prerequisite Check"
echo "=============================================="
echo ""

# Core Requirements
echo "--- Core Requirements ---"
check "Docker" "docker" "24.0" true
check "Docker Compose" "docker" "" true

# Verify Docker Compose v2
if command -v docker &> /dev/null; then
    if docker compose version &> /dev/null; then
        compose_version=$(docker compose version --short 2>/dev/null)
        printf "Checking %-20s " "Docker Compose v2..."
        echo -e "${GREEN}OK${NC} - Version: $compose_version"
        PASSED=$((PASSED + 1))
    else
        printf "Checking %-20s " "Docker Compose v2..."
        echo -e "${RED}FAILED${NC} - Docker Compose v2 not available"
        FAILED=$((FAILED + 1))
    fi
fi

check "Git" "git" "2.40" true

echo ""
echo "--- Development Tools ---"
check "Go" "go" "1.22" true
check "Node.js" "node" "22" true
check "bun" "bun" "1.0" true

echo ""
echo "--- Optional Tools ---"
check "golangci-lint" "golangci-lint" "" false
check "Air (hot-reload)" "air" "" false
check "migrate (db)" "migrate" "" false

# Docker Daemon Check
echo ""
echo "--- Docker Daemon ---"
printf "Checking %-20s " "Docker daemon..."
if docker info &> /dev/null; then
    echo -e "${GREEN}OK${NC} - Running"
    PASSED=$((PASSED + 1))
else
    echo -e "${RED}FAILED${NC} - Not running or permission denied"
    FAILED=$((FAILED + 1))
fi

# Port Availability Check
echo ""
echo "--- Port Availability ---"
check_port() {
    local port="$1"
    local service="$2"
    printf "Checking %-20s " "Port $port ($service)..."
    if lsof -Pi :$port -sTCP:LISTEN -t &> /dev/null || netstat -tuln 2>/dev/null | grep -q ":$port "; then
        echo -e "${YELLOW}WARNING${NC} - Port in use"
        WARNINGS=$((WARNINGS + 1))
    else
        echo -e "${GREEN}OK${NC} - Available"
        PASSED=$((PASSED + 1))
    fi
}

check_port 3000 "Frontend"
check_port 8080 "Backend"
check_port 5432 "PostgreSQL"

# Summary
echo ""
echo "=============================================="
echo "Summary"
echo "=============================================="
echo -e "Passed:   ${GREEN}$PASSED${NC}"
echo -e "Failed:   ${RED}$FAILED${NC}"
echo -e "Warnings: ${YELLOW}$WARNINGS${NC}"
echo ""

if [ $FAILED -gt 0 ]; then
    echo -e "${RED}Some prerequisites are missing!${NC}"
    echo "Please install the missing tools before proceeding."
    exit 1
else
    echo -e "${GREEN}All required prerequisites are satisfied!${NC}"
    if [ $WARNINGS -gt 0 ]; then
        echo -e "${YELLOW}Some optional tools are missing. Consider installing them for the best experience.${NC}"
    fi
    echo ""
    echo "You can now start the application with:"
    echo "  cd docker && docker compose up -d"
    echo ""
    exit 0
fi
