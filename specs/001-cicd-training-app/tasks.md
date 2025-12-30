# Tasks: CI/CD Training Application

**Input**: Design documents from `/specs/001-cicd-training-app/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: Tests are included as part of quality requirements (FR-006 requires 80% coverage).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Backend**: `backend/` (Go with Fiber v2, GORM)
- **Frontend**: `frontend/` (Next.js 14, Tailwind CSS, shadcn/ui)
- **Infrastructure**: `docker/`, `k8s/`
- **CI/CD**: `.github/workflows/`, `.github/actions/`
- **Documentation**: `docs/`

---

## Phase 1: Setup (Project Initialization)

**Purpose**: Create project structure and initialize dependencies

- [X] T001 Create backend directory structure per plan.md in backend/
- [X] T002 [P] Initialize Go module with go.mod in backend/go.mod
- [X] T003 [P] Create frontend directory structure per plan.md in frontend/
- [X] T004 [P] Initialize Next.js project with package.json in frontend/package.json
- [X] T005 [P] Create docker directory structure in docker/
- [X] T006 [P] Create k8s directory structure with base/ and overlays/ in k8s/
- [X] T007 [P] Create .github/workflows/ and .github/actions/ directories
- [X] T008 [P] Create docs/exercises/ and docs/runbooks/ directories
- [X] T009 [P] Create alerts/ directory for monitoring rules

**Checkpoint**: Project skeleton created - all directories in place

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Backend Foundation

- [X] T010 Install Go dependencies (Fiber v2, GORM, fiberprometheus, uuid) in backend/go.mod
- [X] T011 Create environment configuration loader in backend/config/config.go
- [X] T012 [P] Create database connection manager in backend/database/database.go
- [X] T013 [P] Create correlation ID middleware in backend/middleware/correlation.go
- [X] T014 [P] Create structured JSON logging middleware in backend/middleware/logging.go
- [X] T015 Create SQL migration file 000001_init.up.sql in backend/database/migrations/000001_init.up.sql
- [X] T016 [P] Create SQL rollback file 000001_init.down.sql in backend/database/migrations/000001_init.down.sql
- [X] T017 Create Environment GORM model in backend/models/environment.go
- [X] T018 [P] Create Release GORM model in backend/models/release.go
- [X] T019 [P] Create Deployment GORM model with status enum in backend/models/deployment.go
- [X] T020 [P] Create HealthCheck GORM model in backend/models/health_check.go

### Frontend Foundation

- [X] T021 Install frontend dependencies (Tailwind CSS, shadcn/ui) in frontend/package.json
- [X] T022 Configure Tailwind CSS in frontend/tailwind.config.js
- [X] T023 [P] Configure Next.js settings in frontend/next.config.js
- [X] T024 [P] Create root layout with navigation in frontend/app/layout.tsx
- [X] T025 [P] Install and configure shadcn/ui base components in frontend/components/ui/

### Configuration Files

- [X] T026 [P] Create .golangci.yml for Go linting at repository root
- [X] T027 [P] Create sonar-project.properties for SonarCloud at repository root
- [X] T028 [P] Create Makefile with development commands at repository root
- [X] T029 [P] Create .air.toml for Go hot-reload in backend/.air.toml

**Checkpoint**: Foundation ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Local Development Setup (Priority: P1) 🎯 MVP

**Goal**: Trainees can clone the repo and run `docker compose up` to start the full stack within 3 minutes

**Independent Test**: Clone repo → run `docker compose up` → access frontend at localhost:3000 → see backend health status displayed

### Tests for User Story 1

- [X] T030 [P] [US1] Create health endpoint unit test in backend/tests/health_test.go
- [X] T031 [P] [US1] Create config loading test in backend/tests/config_test.go

### Implementation for User Story 1

- [X] T032 [US1] Create liveness health handler (/health/live) in backend/handlers/health.go
- [X] T033 [US1] Create readiness health handler (/health/ready) in backend/handlers/health.go
- [X] T034 [US1] Setup Prometheus metrics with fiberprometheus in backend/handlers/metrics.go
- [X] T035 [US1] Create main.go with Fiber app setup and routes in backend/main.go
- [X] T036 [P] [US1] Create backend Dockerfile with multi-stage build in docker/backend.Dockerfile
- [X] T037 [P] [US1] Create frontend Dockerfile with Next.js production build in docker/frontend.Dockerfile
- [X] T038 [US1] Create docker-compose.yml with 3 services (backend, frontend, db) in docker/docker-compose.yml
- [X] T039 [P] [US1] Create dashboard home page in frontend/app/page.tsx
- [X] T040 [P] [US1] Create health status page in frontend/app/health/page.tsx
- [X] T041 [P] [US1] Create health-card component in frontend/components/health-card.tsx
- [X] T042 [US1] Add hot-reload support with volume mounts in docker/docker-compose.yml
- [X] T043 [US1] Create prerequisite check script in scripts/check-prereqs.sh

**Checkpoint**: User Story 1 complete - `docker compose up` starts full stack with health dashboard

---

## Phase 4: User Story 2 - Execute CI Pipeline Locally (Priority: P2)

**Goal**: Trainees can run the same lint, test, and build commands locally that run in GitHub Actions

**Independent Test**: Run `make lint` → see linting output → run `make test` → see coverage report → run `make build` → see tagged container images

### Tests for User Story 2

- [X] T044 [P] [US2] Create lint configuration test to verify rules match CI in backend/tests/lint_test.go

### Implementation for User Story 2

- [X] T045 [US2] Add lint command to Makefile (golangci-lint + eslint) in Makefile
- [X] T046 [US2] Add test command with coverage output to Makefile in Makefile
- [X] T047 [US2] Add build command with git SHA tagging to Makefile in Makefile
- [X] T048 [P] [US2] Configure ESLint for frontend in frontend/.eslintrc.js
- [X] T049 [P] [US2] Create frontend component tests setup in frontend/tests/
- [X] T050 [US2] Create local CI simulation script in scripts/local-ci.sh
- [X] T051 [US2] Document local CI commands in docs/exercises/02-ci-pipeline.md

**Checkpoint**: User Story 2 complete - local CI commands mirror GitHub Actions

---

## Phase 5: User Story 3 - Push Code and Observe Pipeline (Priority: P3)

**Goal**: Trainees push changes to GitHub and see CI/CD pipeline execute with distinct stages

**Independent Test**: Create feature branch → push code → see CI workflow trigger within 30 seconds → observe lint/test/build/scan stages

### Implementation for User Story 3

- [X] T052 [US3] Create reusable lint workflow in .github/workflows/reusable-lint.yml
- [X] T053 [P] [US3] Create reusable test workflow in .github/workflows/reusable-test.yml
- [X] T054 [P] [US3] Create reusable build workflow in .github/workflows/reusable-build.yml
- [X] T055 [US3] Create main CI workflow calling reusable workflows in .github/workflows/ci.yml
- [X] T056 [P] [US3] Create Go setup composite action in .github/actions/setup-go/action.yml
- [X] T057 [P] [US3] Create Node.js setup composite action in .github/actions/setup-node/action.yml
- [X] T058 [US3] Add SonarCloud scan job to CI workflow in .github/workflows/ci.yml
- [X] T059 [US3] Configure branch protection rules documentation in docs/github-setup.md
- [X] T060 [US3] Document pipeline observation exercise in docs/exercises/03-deployment.md

**Checkpoint**: User Story 3 complete - CI pipeline runs on push with visible stages

---

## Phase 6: User Story 4 - Deploy to Alpha Environment (Priority: P4)

**Goal**: Code merged to develop automatically deploys to alpha environment with health check validation

**Independent Test**: Merge PR to develop → observe CD workflow trigger → see deployment to alpha → verify health check passes → access alpha URL

### Tests for User Story 4

- [X] T061 [P] [US4] Create deployment model test in backend/tests/deployment_test.go
- [X] T062 [P] [US4] Create API endpoint tests in backend/tests/api_test.go

### Implementation for User Story 4

- [X] T063 [US4] Create deployment API handlers in backend/handlers/api.go
- [X] T064 [US4] Create environment listing endpoint in backend/handlers/api.go
- [X] T065 [US4] Create release management endpoints in backend/handlers/api.go
- [X] T066 [P] [US4] Create deployments page in frontend/app/deployments/page.tsx
- [X] T067 [P] [US4] Create deployment-list component in frontend/components/deployment-list.tsx
- [X] T068 [US4] Create Kustomize base for backend in k8s/base/backend/
- [X] T069 [P] [US4] Create Kustomize base for frontend in k8s/base/frontend/
- [X] T070 [P] [US4] Create base kustomization.yaml in k8s/base/kustomization.yaml
- [X] T071 [US4] Create alpha environment overlay in k8s/overlays/alpha/
- [X] T072 [US4] Create CD workflow for alpha deployment in .github/workflows/cd-alpha.yml
- [X] T073 [US4] Add health check validation step to CD workflow in .github/workflows/cd-alpha.yml
- [X] T074 [US4] Add automatic rollback on health check failure in .github/workflows/cd-alpha.yml

**Checkpoint**: User Story 4 complete - merge to develop triggers alpha deployment with health validation

---

## Phase 7: User Story 5 - Promote Through Environments (Priority: P5)

**Goal**: Trainees can promote releases through beta → nonprod → prod with approval gates and soak periods

**Independent Test**: Request beta promotion → approve → see deployment → wait soak period → request prod → approve → see canary deployment

### Implementation for User Story 5

- [X] T075 [US5] Create beta environment overlay in k8s/overlays/beta/
- [X] T076 [P] [US5] Create nonprod environment overlay in k8s/overlays/nonprod/
- [X] T077 [P] [US5] Create prod environment overlay in k8s/overlays/prod/
- [X] T078 [US5] Create CD workflow for beta (manual trigger, approval) in .github/workflows/cd-beta.yml
- [X] T079 [US5] Create CD workflow for nonprod (manual trigger, approval) in .github/workflows/cd-nonprod.yml
- [X] T080 [US5] Create CD workflow for prod (manual trigger, approval, canary) in .github/workflows/cd-prod.yml
- [X] T081 [US5] Add soak period enforcement (24h) before prod promotion in .github/workflows/cd-prod.yml
- [X] T082 [US5] Add canary deployment logic with automatic rollback in .github/workflows/cd-prod.yml
- [X] T083 [US5] Create promotion documentation in docs/exercises/04-observability.md

**Checkpoint**: User Story 5 complete - full environment promotion with gates and soak periods

---

## Phase 8: User Story 6 - Execute Hotfix Workflow (Priority: P6)

**Goal**: Trainees can deploy emergency fixes directly to production with expedited (not bypassed) gates

**Independent Test**: Create hotfix branch from main → push fix → see expedited workflow → approve → see production deployment → verify backport to develop

### Implementation for User Story 6

- [X] T084 [US6] Create hotfix workflow with expedited path in .github/workflows/hotfix.yml
- [X] T085 [US6] Add automatic backport to develop branch in .github/workflows/hotfix.yml
- [X] T086 [US6] Add 2-minute rollback timeout on health failure in .github/workflows/hotfix.yml
- [X] T087 [US6] Create rollback runbook in docs/runbooks/rollback.md
- [X] T088 [P] [US6] Create hotfix runbook in docs/runbooks/hotfix.md
- [X] T089 [US6] Create post-incident review template in docs/runbooks/post-incident-template.md

**Checkpoint**: User Story 6 complete - hotfix workflow provides expedited production path

---

## Phase 9: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, training materials, and quality improvements

### Documentation

- [X] T090 [P] Create comprehensive README.md with learning path at repository root
- [X] T091 [P] Create local setup exercise in docs/exercises/01-local-setup.md
- [X] T092 [P] Create troubleshooting guide in docs/troubleshooting.md
- [X] T093 Update quickstart.md with final paths in specs/001-cicd-training-app/quickstart.md

### Observability

- [X] T094 [P] Create backend alerting rules in alerts/backend-alerts.yaml
- [X] T095 [P] Create frontend alerting rules in alerts/frontend-alerts.yaml
- [X] T096 Add custom deployment metrics (deployments_total, health_check_duration) in backend/handlers/metrics.go

### Quality & Security

- [ ] T097 Run full test suite and verify 80% coverage
- [ ] T098 Run SonarCloud scan and fix critical issues
- [ ] T099 [P] Security scan container images for vulnerabilities
- [ ] T100 Validate all quickstart.md steps work end-to-end

---

## Dependencies & Execution Order

### Phase Dependencies

```
Phase 1 (Setup) ──────────────────────────────────────┐
                                                       │
Phase 2 (Foundational) ◄──────────────────────────────┘
         │
         │  ⚠️ BLOCKS ALL USER STORIES
         ▼
┌────────┴────────┬────────────────┬────────────────┬────────────────┬────────────────┐
│                 │                │                │                │                │
▼                 ▼                ▼                ▼                ▼                ▼
Phase 3 (US1)   Phase 4 (US2)   Phase 5 (US3)   Phase 6 (US4)   Phase 7 (US5)   Phase 8 (US6)
MVP             CI Local        CI Cloud        Deploy Alpha    Promotions      Hotfix
│                 │                │                │                │                │
└────────┬────────┴────────────────┴────────────────┴────────────────┴────────────────┘
         │
         ▼
Phase 9 (Polish)
```

### User Story Dependencies

| Story | Depends On | Can Parallelize With |
|-------|------------|----------------------|
| US1 (P1) | Phase 2 only | US2, US3 (partial) |
| US2 (P2) | Phase 2 only | US1, US3 |
| US3 (P3) | Phase 2 only | US1, US2 |
| US4 (P4) | US1 (health endpoints) | US5, US6 (after US1 complete) |
| US5 (P5) | US4 (alpha deployment) | US6 |
| US6 (P6) | US4 (deployment infrastructure) | US5 |

### Within Each User Story

1. Tests → written first (if included), must fail
2. Models → before services
3. Services/Handlers → before endpoints
4. Infrastructure → before deployment workflows
5. Documentation → after implementation

### Parallel Opportunities

**Phase 2 Parallel Groups**:
```
Group A (Backend): T012, T013, T014 (middleware)
Group B (Models): T017, T018, T019, T020 (all models)
Group C (Frontend): T022, T023, T024, T025
Group D (Config): T026, T027, T028, T029
```

**User Story Parallel Groups**:
```
US1 Dockerfiles: T036, T037 (different files)
US1 Frontend Pages: T039, T040, T041 (different components)
US3 Workflows: T052, T053, T054 (reusable workflows)
US4 Kustomize: T068, T069, T070 (base resources)
US5 Overlays: T075, T076, T077 (environment overlays)
```

---

## Parallel Example: User Story 1

```bash
# Launch tests in parallel:
Task: "Create health endpoint unit test in backend/tests/health_test.go"
Task: "Create config loading test in backend/tests/config_test.go"

# Launch Dockerfiles in parallel:
Task: "Create backend Dockerfile in docker/backend.Dockerfile"
Task: "Create frontend Dockerfile in docker/frontend.Dockerfile"

# Launch frontend pages in parallel:
Task: "Create dashboard home page in frontend/app/page.tsx"
Task: "Create health status page in frontend/app/health/page.tsx"
Task: "Create health-card component in frontend/components/health-card.tsx"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T009)
2. Complete Phase 2: Foundational (T010-T029) - **CRITICAL**
3. Complete Phase 3: User Story 1 (T030-T043)
4. **STOP and VALIDATE**: Run `docker compose up`, verify health dashboard
5. Deploy/demo MVP - trainees can start learning immediately

### Incremental Delivery

| Milestone | Stories | Training Capability |
|-----------|---------|---------------------|
| MVP | US1 | Local development setup |
| +CI | US1 + US2 | Local CI execution |
| +Pipeline | US1-US3 | Full CI/CD observation |
| +Deploy | US1-US4 | Alpha deployments |
| +Promote | US1-US5 | Full environment promotion |
| Complete | US1-US6 | Hotfix workflow |

### Parallel Team Strategy

With multiple developers after Phase 2 completion:

```
Developer A: US1 (MVP) → US4 (Deployments)
Developer B: US2 (Local CI) → US5 (Promotions)
Developer C: US3 (Cloud CI) → US6 (Hotfix)
```

---

## Notes

- [P] tasks = different files, no dependencies on incomplete tasks
- [Story] label maps task to specific user story (US1-US6)
- Each user story is independently completable and testable
- 80% test coverage required per FR-006
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- All workflows must use >80% reusable components per Constitution Principle I

## Summary

| Metric | Count |
|--------|-------|
| **Total Tasks** | 100 |
| **Setup Phase** | 9 |
| **Foundational Phase** | 20 |
| **US1 Tasks** | 14 |
| **US2 Tasks** | 8 |
| **US3 Tasks** | 9 |
| **US4 Tasks** | 14 |
| **US5 Tasks** | 9 |
| **US6 Tasks** | 6 |
| **Polish Tasks** | 11 |
| **Parallel Opportunities** | 47 tasks marked [P] |
