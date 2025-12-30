# Research: CI/CD Training Application

**Feature Branch**: `001-cicd-training-app`
**Date**: 2025-12-29
**Purpose**: Resolve technical decisions for minimal CI/CD training application

## Research Summary

This document captures technology decisions made during Phase 0 research. All decisions align with the project constitution while optimizing for training clarity and minimal complexity.

---

## 1. Go Project Layout

**Decision**: Flat package structure without `cmd/` and `internal/` directories

**Rationale**:
- Single binary application with no CLI subcommands
- Training focus benefits from direct code navigation
- Over-engineering obscures learning objectives
- Constitution requires simplicity; flat layout is sufficient

**Alternatives Considered**:
| Layout | Pros | Cons | Verdict |
|--------|------|------|---------|
| Standard (cmd/internal/) | Enterprise pattern | Over-engineered for scope | Rejected |
| Flat (handlers/, models/) | Clear, minimal | Less isolation | **Selected** |
| Hexagonal | Clean architecture | High learning curve | Rejected |

**Implementation**:
```text
backend/
├── main.go           # Entry point with Fiber setup
├── config/           # Environment configuration
├── handlers/         # HTTP handlers (health, metrics, api)
├── middleware/       # Request middleware (logging, correlation)
├── models/           # GORM models
└── database/         # Connection and migrations
```

---

## 2. Next.js Minimal Setup

**Decision**: Next.js 14 App Router with 3 pages (home, health, deployments)

**Rationale**:
- App Router is the modern Next.js standard (required by constitution)
- Minimal pages demonstrate patterns without feature bloat
- Dashboard-only scope matches "minimal" user requirement
- shadcn/ui provides accessible components without custom CSS

**Alternatives Considered**:
| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Pages Router | Simpler mental model | Legacy pattern | Rejected |
| App Router (full) | Complete feature set | Over-scoped | Rejected |
| App Router (minimal) | Modern, focused | Limited features | **Selected** |
| Static HTML | Simplest | No interactivity | Rejected |

**Implementation**:
```text
frontend/app/
├── layout.tsx        # Root layout with nav
├── page.tsx          # Dashboard home (status summary)
├── health/page.tsx   # Health endpoint visualization
└── deployments/page.tsx  # Deployment history list
```

**Page Purposes**:
- **Home**: Overview dashboard with backend status, latest deployment
- **Health**: Real-time health check status (/health/live, /health/ready)
- **Deployments**: History of deployments with version, status, timestamp

---

## 3. Docker Compose Patterns

**Decision**: Three-service stack with hot-reload support

**Rationale**:
- Single `docker compose up` command meets SC-001 (5-minute setup)
- Air for Go hot-reload, Next.js dev server for frontend
- PostgreSQL with health check ensures DB ready before app start
- Volume mounts enable code changes without rebuild

**Alternatives Considered**:
| Pattern | Pros | Cons | Verdict |
|---------|------|------|---------|
| All-in-one container | Simple | No separation | Rejected |
| Separate services | Realistic | Complex orchestration | Partial |
| Dev services with hot-reload | Fast iteration | More config | **Selected** |

**Implementation**:
```yaml
services:
  backend:
    build: ./docker/backend.Dockerfile
    volumes:
      - ./backend:/app
    environment:
      - DATABASE_URL=postgres://...
    depends_on:
      db:
        condition: service_healthy

  frontend:
    build: ./docker/frontend.Dockerfile
    volumes:
      - ./frontend:/app
    environment:
      - NEXT_PUBLIC_API_URL=http://localhost:8080

  db:
    image: postgres:16
    healthcheck:
      test: ["CMD-SHELL", "pg_isready"]
    volumes:
      - pgdata:/var/lib/postgresql/data
```

---

## 4. GitHub Actions Workflow Organization

**Decision**: Modular workflows with reusable components and composite actions

**Rationale**:
- Constitution Principle I requires >80% reusable workflow coverage
- Separate CD workflows per environment enable approval gates
- Composite actions reduce duplication in setup steps
- Clear file naming (ci.yml, cd-*.yml, reusable-*.yml) follows standards

**Alternatives Considered**:
| Organization | Pros | Cons | Verdict |
|--------------|------|------|---------|
| Single workflow | Simple | Monolithic, hard to maintain | Rejected |
| Matrix jobs | DRY | Complex conditions | Partial |
| Reusable workflows | Modular, clear | More files | **Selected** |

**Implementation**:
```text
.github/workflows/
├── ci.yml                 # Triggered on: push/PR to any branch
│   └── jobs: lint → test → build → scan
├── cd-alpha.yml           # Triggered on: merge to develop
├── cd-beta.yml            # Triggered on: workflow_dispatch (manual)
├── cd-nonprod.yml         # Triggered on: workflow_dispatch (manual)
├── cd-prod.yml            # Triggered on: workflow_dispatch (manual)
├── hotfix.yml             # Triggered on: push to hotfix/*
├── reusable-lint.yml      # Called by: ci.yml
├── reusable-test.yml      # Called by: ci.yml
└── reusable-build.yml     # Called by: ci.yml, cd-*.yml
```

**Quality Gate Implementation**:
```yaml
# ci.yml structure
jobs:
  lint:
    uses: ./.github/workflows/reusable-lint.yml
  test:
    needs: lint
    uses: ./.github/workflows/reusable-test.yml
  build:
    needs: test
    uses: ./.github/workflows/reusable-build.yml
  scan:
    needs: build
    # SonarCloud analysis
```

---

## 5. Prometheus Metrics in Go

**Decision**: fiberprometheus middleware with custom application metrics

**Rationale**:
- Native integration with Fiber v2 framework
- Automatic HTTP metrics (requests, latency, errors)
- Custom metrics demonstrate observability patterns
- Constitution requires `/metrics` endpoint with specific metrics

**Alternatives Considered**:
| Library | Pros | Cons | Verdict |
|---------|------|------|---------|
| prometheus/client_golang | Official | Manual integration | Partial |
| fiberprometheus | Fiber-native | Less control | **Selected** |
| OpenTelemetry | Modern, traces | Over-scoped | Future consideration |

**Implementation**:
```go
// handlers/metrics.go
import (
    "github.com/ansrivas/fiberprometheus/v2"
    "github.com/prometheus/client_golang/prometheus"
)

var (
    deploymentsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "deployments_total",
            Help: "Total number of deployments",
        },
        []string{"environment", "status"},
    )
    healthCheckDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "health_check_duration_seconds",
            Help:    "Health check response time",
            Buckets: prometheus.DefBuckets,
        },
        []string{"check_type"},
    )
)

func SetupMetrics(app *fiber.App) {
    prometheus.MustRegister(deploymentsTotal, healthCheckDuration)
    prom := fiberprometheus.New("cicd_training")
    prom.RegisterAt(app, "/metrics")
    app.Use(prom.Middleware)
}
```

---

## 6. Database Migration Strategy

**Decision**: golang-migrate for CI/CD, GORM AutoMigrate for local development

**Rationale**:
- golang-migrate provides version-controlled, reversible migrations
- SQL files enable code review for schema changes
- AutoMigrate speeds up local iteration during development
- Constitution requires version-controlled migrations applied via CI

**Alternatives Considered**:
| Tool | Pros | Cons | Verdict |
|------|------|------|---------|
| GORM AutoMigrate only | Simple | No rollback, no version control | Local dev only |
| golang-migrate | Version control, rollback | More setup | **Selected (CI/CD)** |
| Atlas | Modern, declarative | Learning curve | Future consideration |
| Manual SQL | Full control | Error-prone | Rejected |

**Implementation**:
```text
backend/database/migrations/
├── 000001_init.up.sql      # Create tables
├── 000001_init.down.sql    # Drop tables
├── 000002_add_indexes.up.sql
└── 000002_add_indexes.down.sql
```

**Migration Commands**:
```bash
# CI/CD (production-like)
migrate -path ./database/migrations -database $DATABASE_URL up

# Local development (rapid iteration)
go run main.go --auto-migrate
```

---

## 7. SonarCloud Integration

**Decision**: Parallel job in CI pipeline with multi-language project scan

**Rationale**:
- SonarCloud provides unified quality metrics for Go + JavaScript
- Parallel execution doesn't block deployment on non-critical issues
- Constitution requires code quality scanning with configurable thresholds
- 0 critical/blocker issues required; major issues allowed with limits

**Alternatives Considered**:
| Tool | Pros | Cons | Verdict |
|------|------|------|---------|
| golangci-lint only | Fast, Go-native | No unified dashboard | Partial (linting) |
| CodeClimate | Simple setup | Less Go support | Rejected |
| SonarCloud | Comprehensive | Requires token setup | **Selected** |

**Implementation**:
```yaml
# ci.yml - scan job
scan:
  name: SonarCloud Analysis
  needs: build
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
      with:
        fetch-depth: 0
    - uses: SonarSource/sonarcloud-github-action@master
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}
```

**Configuration** (`sonar-project.properties`):
```properties
sonar.projectKey=go-github-actions
sonar.organization=<org>
sonar.sources=backend,frontend
sonar.exclusions=**/*_test.go,**/node_modules/**
sonar.go.coverage.reportPaths=backend/coverage.out
sonar.javascript.lcov.reportPaths=frontend/coverage/lcov.info
```

---

## Decision Summary

| Topic | Decision | Constitution Alignment |
|-------|----------|----------------------|
| Go Layout | Flat structure | Simplicity principle |
| Next.js Setup | Minimal App Router | Tech stack constraint |
| Docker Compose | 3 services + hot-reload | IaC principle |
| GitHub Actions | Modular + reusable | Pipeline as Code (>80% reuse) |
| Prometheus | fiberprometheus | Observability (required metrics) |
| Migrations | golang-migrate + AutoMigrate | IaC (version-controlled) |
| SonarCloud | Parallel CI job | Quality Gates (0 critical) |

---

## Open Items

None. All NEEDS CLARIFICATION items resolved through research.

## Next Steps

1. Create `data-model.md` with entity definitions
2. Create `contracts/openapi.yaml` with API specification
3. Create `quickstart.md` with setup instructions
4. Generate `tasks.md` via `/speckit.tasks` command
