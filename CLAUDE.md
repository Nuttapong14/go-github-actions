# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CI/CD Training Application - A Go + Next.js stack for teaching GitHub Actions CI/CD concepts. Demonstrates quality gates, environment promotion (alpha → beta → nonprod → prod), observability, and hotfix workflows.

## Development Commands

```bash
# Full stack with Docker Compose
make up                  # Start backend + frontend + PostgreSQL
make down                # Stop all services
make logs                # Tail container logs

# Local CI pipeline (mirrors GitHub Actions)
make ci                  # lint → test → build (full pipeline)
make lint                # Backend (golangci-lint) + Frontend (ESLint)
make test                # Backend + Frontend with coverage
make build               # Binary + Next.js production build

# Backend (Go)
make backend.dev         # Hot-reload with Air
make backend.test        # Tests with coverage report
make backend.lint        # golangci-lint using .golangci.yml
make backend.build       # Build binary with git SHA version

# Frontend (Bun/Next.js)
make frontend.dev        # Development server
make frontend.test       # Bun test
make frontend.lint       # ESLint
make frontend.build      # Production build

# Database
make db.connect          # psql shell to PostgreSQL
make db.reset            # Drop and recreate database

# Docker
make docker.build        # Build images with git SHA tags
make docker.push         # Push to ghcr.io/nuttapong14
```

## Architecture

### Backend (Go 1.23 + Fiber v2)
```
backend/
├── main.go              # Entry point, Fiber app setup, route registration
├── config/              # Environment-based configuration (Config struct)
├── database/            # GORM PostgreSQL connection
├── handlers/            # HTTP handlers (health, metrics, API)
├── middleware/          # Correlation ID, structured logging
├── models/              # GORM models (Environment, Release, Deployment)
└── tests/               # Go tests (*_test.go)
```

Key patterns:
- Configuration via environment variables (DB_HOST, DB_PORT, etc.)
- Health endpoints: `/health/live`, `/health/ready`
- Prometheus metrics at `/metrics`
- API routes under `/api/v1`
- Graceful shutdown with signal handling

### Frontend (Next.js 14 + App Router)
```
frontend/
├── app/                 # Next.js App Router pages
│   ├── layout.tsx       # Root layout
│   ├── page.tsx         # Home page
│   ├── health/          # Health monitoring page
│   └── deployments/     # Deployments dashboard
├── components/
│   ├── ui/              # shadcn/ui primitives (button, card, badge, table)
│   └── *.tsx            # Feature components (health-card, deployment-list)
└── package.json         # Bun scripts: dev, build, test, lint
```

### CI/CD (GitHub Actions)
```
.github/workflows/
├── ci.yml               # Main pipeline: lint → test → build → SonarCloud
├── reusable-*.yml       # Modular reusable workflows
├── cd-alpha.yml         # Deploy to alpha environment
├── cd-beta.yml          # Promote to beta
├── cd-nonprod.yml       # Promote to non-prod
├── cd-prod.yml          # Production deployment
└── hotfix.yml           # Emergency hotfix workflow
```

### Kubernetes (Kustomize)
```
k8s/
├── base/                # Base manifests
└── overlays/            # Environment-specific patches
```

## Environment Variables

Backend configuration:
- `ENV` - Environment (development/production)
- `PORT` - HTTP server port (default: 8080)
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`
- `LOG_LEVEL` - Logging level (default: info)

Frontend:
- `NEXT_PUBLIC_API_URL` - Backend API URL

## Testing

```bash
# Run single Go test
cd backend && go test -v -run TestFunctionName ./...

# Run tests with race detection
cd backend && go test -v -race ./...

# Frontend specific test file
cd frontend && bun test path/to/file.test.ts
```

## Code Quality

- Go linting: `.golangci.yml` with errcheck, gosec, govet, staticcheck, gocyclo, revive
- Frontend: ESLint with Next.js config
- Coverage threshold: 80%
- SonarCloud integration for code quality scanning
