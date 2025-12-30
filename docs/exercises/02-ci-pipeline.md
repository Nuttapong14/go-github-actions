# Exercise 02: Running CI Pipeline Locally

This exercise teaches you how to run the same CI/CD pipeline checks locally that execute in GitHub Actions. By running these checks before pushing code, you can catch issues early and avoid failed CI builds.

## Learning Objectives

By the end of this exercise, you will be able to:

1. Run linting checks for both backend and frontend code
2. Execute tests with coverage reporting
3. Build production artifacts with git SHA tagging
4. Use the local CI simulation script to mirror GitHub Actions
5. Understand the relationship between local development and CI/CD

## Prerequisites

Before starting this exercise, ensure you have:

- Completed [Exercise 01: Local Development Setup](./01-local-setup.md)
- All prerequisites installed (run `./scripts/check-prereqs.sh` to verify)
- The development environment running (`make up`)

## Exercise Steps

### Step 1: Understanding the Makefile

The project Makefile provides commands that mirror the GitHub Actions CI pipeline. View available commands:

```bash
make help
```

You'll see output showing:
- Build information (Git SHA, branch, tag)
- Docker commands for local development
- Backend and frontend specific commands
- **CI/CD commands that mirror GitHub Actions**

### Step 2: Running Linting Checks

Linting ensures code follows consistent style and catches potential bugs early.

#### Backend Linting (golangci-lint)

Run Go linting:

```bash
make backend.lint
```

This uses the configuration in `.golangci.yml` which enables:
- `errcheck` - Unchecked error detection
- `gosec` - Security vulnerability detection
- `govet` - Suspicious construct detection
- `staticcheck` - Static analysis
- `gofmt` / `goimports` - Code formatting

#### Frontend Linting (ESLint)

Run TypeScript/React linting:

```bash
make frontend.lint
```

This uses the configuration in `frontend/.eslintrc.js` which enforces:
- TypeScript best practices
- React hooks rules
- Import ordering
- Next.js specific rules

#### Run All Linters

To run both backend and frontend linting:

```bash
make lint
```

### Step 3: Running Tests with Coverage

Tests validate that your code works as expected.

#### Backend Tests

Run Go tests with race detection and coverage:

```bash
make backend.test
```

This will:
1. Run all Go tests with `-race` flag for race condition detection
2. Generate `backend/coverage.out` with coverage data
3. Create `backend/coverage.html` for visual coverage report

View the coverage report:

```bash
# On macOS
open backend/coverage.html

# On Linux
xdg-open backend/coverage.html
```

#### Frontend Tests

Run frontend component tests:

```bash
make frontend.test
```

#### Run All Tests

```bash
make test
```

### Step 4: Building Production Artifacts

Building creates the production-ready binaries and assets.

#### Backend Build

Build the Go binary with version information embedded:

```bash
make backend.build
```

This creates `bin/backend` with embedded:
- Version (from git tag or SHA)
- Git SHA
- Build date

Verify the build:

```bash
./bin/backend --version  # If version flag is implemented
ls -la bin/backend
```

#### Frontend Build

Build the Next.js production bundle:

```bash
make frontend.build
```

This creates the optimized production build in `frontend/.next/`.

#### Build All

```bash
make build
```

### Step 5: Running the Full CI Pipeline Locally

The `local-ci.sh` script simulates the complete GitHub Actions CI pipeline:

```bash
./scripts/local-ci.sh
```

This runs:
1. **Lint** - All linting checks
2. **Test** - All tests with coverage
3. **Build** - Production artifacts

#### Script Options

Run only specific stages:

```bash
# Only linting
./scripts/local-ci.sh --lint-only

# Only tests
./scripts/local-ci.sh --test-only

# Only build
./scripts/local-ci.sh --build-only

# Skip specific stages
./scripts/local-ci.sh --skip-lint
./scripts/local-ci.sh --skip-test
```

Build Docker images too:

```bash
./scripts/local-ci.sh --docker
```

### Step 6: Using the CI Make Target

For a quick CI check, use:

```bash
make ci
```

This runs lint → test → build in sequence and displays a summary.

### Step 7: Docker Image Building

Build Docker images with git SHA tagging:

```bash
make docker.build
```

This creates:
- `ghcr.io/nuttapong14/cicd-training-backend:<git-sha>`
- `ghcr.io/nuttapong14/cicd-training-frontend:<git-sha>`

Tag images for release:

```bash
make docker.tag
```

## Understanding Git SHA Tagging

Every build includes the git SHA for traceability:

```
Git SHA:     abc1234
Git Branch:  feature/my-feature
Git Tag:     v1.0.0 (if exists)
Image Tag:   v1.0.0 (or abc1234 if no tag)
```

This allows you to:
- Know exactly which code is in each build
- Trace production issues back to specific commits
- Maintain reproducible builds

## Comparison: Local vs GitHub Actions

| Check | Local Command | GitHub Actions |
|-------|---------------|----------------|
| Backend Lint | `make backend.lint` | `golangci-lint run` |
| Frontend Lint | `make frontend.lint` | `bun run lint` |
| Backend Test | `make backend.test` | `go test -race -coverprofile` |
| Frontend Test | `make frontend.test` | `bun test` |
| Backend Build | `make backend.build` | `go build -ldflags` |
| Frontend Build | `make frontend.build` | `bun run build` |
| Docker Build | `make docker.build` | `docker build` |

## Troubleshooting

### golangci-lint not found

Install it:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

Or use Homebrew on macOS:

```bash
brew install golangci-lint
```

### ESLint errors about missing dependencies

Install frontend dependencies:

```bash
cd frontend && bun install
```

### Tests failing with import errors

Ensure Go modules are downloaded:

```bash
cd backend && go mod download
```

### Build failing with "permission denied"

Make scripts executable:

```bash
chmod +x scripts/*.sh
```

## Practice Exercises

1. **Fix a Lint Error**: Introduce a lint error (e.g., unused variable) and run `make lint` to see it caught.

2. **Add a Test**: Create a new test in `backend/tests/` and verify it runs with `make backend.test`.

3. **Check Coverage**: Run tests and view the coverage report. Identify areas with low coverage.

4. **Pre-Push Hook**: Create a git pre-push hook that runs `make ci` before pushing.

## Summary

In this exercise, you learned:

- How to run the same CI checks locally that run in GitHub Actions
- The importance of linting for code quality
- How tests and coverage work in the project
- How builds are tagged with git SHA for traceability
- Using the local CI script to validate changes before pushing

## Next Steps

Proceed to [Exercise 03: Push and Observe Pipeline](./03-deployment.md) to learn how to push code and watch the GitHub Actions CI/CD pipeline execute.
