# Go GitHub Actions Training Project

## Project Overview
This project is a minimal full-stack application designed to teach CI/CD principles using GitHub Actions. It consists of a **Go (Fiber)** backend and a **Next.js** frontend, featuring workflows for local development, quality assurance, deployment, and observability.

The architecture includes:
- **Backend**: Go with Fiber framework.
- **Frontend**: Next.js (React) managed with Bun.
- **Infrastructure**: Docker Compose for local dev, Kubernetes manifests for orchestration.
- **CI/CD**: GitHub Actions workflows for CI, CD (Alpha/Beta/Prod), and Hotfixes.

## Technology Stack
- **Language**: Go (Backend), TypeScript/JavaScript (Frontend)
- **Frameworks**: Fiber (Go), Next.js (Frontend)
- **Package Managers**: Go Modules, Bun
- **Containerization**: Docker, Docker Compose
- **Orchestration**: Kubernetes (k8s)
- **CI/CD**: GitHub Actions

## Building and Running
The project uses a `Makefile` to standardize commands for local development and CI execution.

### Local Development
- **Start Stack (Backend + Frontend + DB)**:
  ```bash
  make up
  ```
  - Backend: http://localhost:8080
  - Frontend: http://localhost:3000
  - Database: localhost:5432

- **Stop Stack**:
  ```bash
  make down
  ```

- **Run Backend/Frontend Individually**:
  ```bash
  make backend.dev   # Run backend with hot-reload (Air)
  make frontend.dev  # Run frontend in dev mode
  ```

### Testing and Linting
- **Run All Tests**:
  ```bash
  make test
  ```
- **Run All Linters**:
  ```bash
  make lint
  ```
- **Run CI Pipeline Locally** (Lint -> Test -> Build):
  ```bash
  make ci
  ```

### Docker
- **Build Images**:
  ```bash
  make docker.build
  ```
- **Tag Images**:
  ```bash
  make docker.tag
  ```

## Directory Structure
- `backend/`: Go source code (Fiber application).
- `frontend/`: Next.js source code.
- `docker/`: Dockerfiles and Compose configurations.
- `k8s/`: Kubernetes manifests.
- `.github/workflows/`: CI/CD pipeline definitions.
- `docs/`: Documentation and learning exercises.
- `alerts/`: Observability alert rules.
- `specs/`: Specifications and quickstart guides.

## Development Conventions
- **Formatting & Linting**:
  - Go: Uses `golangci-lint` with configuration in `.golangci.yml`.
  - Frontend: Uses ESLint.
  - Run `make lint.fix` to auto-fix issues.
- **Testing**:
  - Backend tests run with race detection (`-race`) and coverage analysis.
  - Frontend tests use `bun test`.
- **Version Control**:
  - Docker images are tagged with the Git SHA by default (`short-sha`).
  - Production releases use Git Tags.
- **Dependencies**:
  - Backend dependencies are managed via `go.mod`.
  - Frontend dependencies are managed via `bun` (using `package.json`).
