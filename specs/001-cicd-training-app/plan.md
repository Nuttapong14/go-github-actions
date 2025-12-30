# Implementation Plan: CI/CD Training Application

**Branch**: `001-cicd-training-app` | **Date**: 2025-12-29 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-cicd-training-app/spec.md`

## Summary

Build a minimal CI/CD training application demonstrating GitHub Actions deployment automation. The application consists of a Go backend API with health/metrics endpoints and a minimal Next.js frontend dashboard. The focus is on teaching DevOps practices (pipeline design, environment promotion, observability) rather than application complexity.

**Technical Approach**: Flat Go project structure with Fiber v2, minimal Next.js 14 App Router frontend, Docker Compose for local development, and Kubernetes manifests for production deployment. GitHub Actions workflows implement the full CI/CD pipeline with progressive environment promotion (alpha → beta → nonprod → prod).

## Technical Context

**Language/Version**: Go 1.22+, Node.js 22 LTS (Next.js 14)
**Primary Dependencies**: Fiber v2 (Go HTTP), GORM (ORM), Next.js 14, Tailwind CSS, shadcn/ui
**Storage**: PostgreSQL 16
**Testing**: go test (race detector), Jest + React Testing Library
**Target Platform**: Linux containers (Docker/Kubernetes)
**Project Type**: Web application (backend + frontend)
**Performance Goals**: Health endpoints <100ms, local CI <3min, alpha deploy <10min
**Constraints**: Single-command local setup, 80% test coverage, zero critical vulnerabilities
**Scale/Scope**: Training environment, <10 concurrent trainees

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Requirement | Status | Implementation |
|-----------|-------------|--------|----------------|
| I. Pipeline as Code | All workflows in `.github/workflows/` | PASS | ci.yml, cd-*.yml, reusable-*.yml |
| II. Progressive Environment Promotion | alpha → beta → nonprod → prod | PASS | 4 deployment workflows with gates |
| III. Quality Gates | Lint, test, scan before deploy | PASS | CI pipeline stages with fail-fast |
| IV. Immutable Artifacts | SHA-tagged Docker images | PASS | Build once, promote everywhere |
| V. Infrastructure as Code | docker-compose.yml + k8s/ | PASS | Docker Compose (dev), Kustomize (prod) |
| VI. Observability by Default | /health/*, /metrics, JSON logs | PASS | fiberprometheus + slog/json |
| VII. Fail Fast, Recover Faster | <5min rollback, actionable errors | PASS | Blue-green deploy, rollback scripts |

**Gate Status**: All constitutional principles addressed. Proceed with Phase 0.

## Project Structure

### Documentation (this feature)

```text
specs/001-cicd-training-app/
├── plan.md              # This file
├── research.md          # Phase 0 output - technology decisions
├── data-model.md        # Phase 1 output - entity definitions
├── quickstart.md        # Phase 1 output - getting started guide
├── contracts/           # Phase 1 output - API contracts
│   └── openapi.yaml     # Backend API specification
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
# Monorepo structure for Go backend + Next.js frontend

backend/
├── main.go                  # Application entry point
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
├── config/
│   └── config.go            # Environment configuration
├── handlers/
│   ├── health.go            # Health check handlers
│   ├── metrics.go           # Prometheus metrics setup
│   └── api.go               # Demo API handlers
├── middleware/
│   ├── logging.go           # Structured JSON logging
│   └── correlation.go       # Correlation ID injection
├── models/
│   └── deployment.go        # GORM models
├── database/
│   ├── migrations/          # SQL migration files
│   │   ├── 000001_init.up.sql
│   │   └── 000001_init.down.sql
│   └── database.go          # Connection management
└── tests/
    ├── health_test.go       # Health endpoint tests
    └── api_test.go          # API integration tests

frontend/
├── package.json             # Node.js dependencies
├── next.config.js           # Next.js configuration
├── tailwind.config.js       # Tailwind CSS configuration
├── app/
│   ├── layout.tsx           # Root layout
│   ├── page.tsx             # Dashboard home
│   ├── health/
│   │   └── page.tsx         # Health status page
│   └── deployments/
│       └── page.tsx         # Deployment history page
├── components/
│   ├── ui/                  # shadcn/ui components
│   ├── health-card.tsx      # Health status display
│   └── deployment-list.tsx  # Deployment history
└── tests/
    └── components/          # Component tests

# Infrastructure
docker/
├── backend.Dockerfile       # Go multi-stage build
├── frontend.Dockerfile      # Next.js production build
└── docker-compose.yml       # Local development stack

k8s/
├── base/
│   ├── backend/
│   │   ├── deployment.yaml
│   │   └── service.yaml
│   ├── frontend/
│   │   ├── deployment.yaml
│   │   └── service.yaml
│   └── kustomization.yaml
└── overlays/
    ├── alpha/
    ├── beta/
    ├── nonprod/
    └── prod/

# CI/CD Workflows
.github/
├── workflows/
│   ├── ci.yml               # Main CI pipeline
│   ├── cd-alpha.yml         # Alpha deployment
│   ├── cd-beta.yml          # Beta deployment
│   ├── cd-nonprod.yml       # NonProd deployment
│   ├── cd-prod.yml          # Production deployment
│   ├── hotfix.yml           # Emergency hotfix
│   ├── reusable-lint.yml    # Reusable lint workflow
│   ├── reusable-test.yml    # Reusable test workflow
│   └── reusable-build.yml   # Reusable build workflow
└── actions/
    ├── setup-go/
    │   └── action.yml       # Go environment setup
    └── setup-node/
        └── action.yml       # Node.js environment setup

# Documentation
docs/
├── exercises/               # Training exercises
│   ├── 01-local-setup.md
│   ├── 02-ci-pipeline.md
│   ├── 03-deployment.md
│   └── 04-observability.md
├── runbooks/
│   ├── rollback.md          # Rollback procedure
│   └── hotfix.md            # Hotfix procedure
└── troubleshooting.md       # Common issues

# Configuration
alerts/
├── backend-alerts.yaml      # Backend alerting rules
└── frontend-alerts.yaml     # Frontend alerting rules

# Root files
README.md                    # Main documentation
Makefile                     # Development commands
.golangci.yml                # Go linting configuration
sonar-project.properties     # SonarCloud configuration
```

**Structure Decision**: Web application structure with separate `backend/` and `frontend/` directories. This provides clear separation for training purposes while maintaining a single repository for GitOps simplicity.

## Complexity Tracking

> No constitution violations requiring justification. Structure follows minimal viable patterns.

| Decision | Justification | Simpler Alternative |
|----------|---------------|---------------------|
| Monorepo | Single repo for all code + infrastructure | Multi-repo adds complexity without training value |
| Kustomize overlays | Environment-specific configs without branching | Helm charts over-engineered for training scope |
| PostgreSQL | Required by constitution, demonstrates migrations | SQLite simpler but lacks production patterns |
