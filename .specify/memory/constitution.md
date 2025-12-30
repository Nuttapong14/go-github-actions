<!--
===================================================================
SYNC IMPACT REPORT
===================================================================
Version Change: [NEW] -> 1.0.0
Bump Rationale: Initial constitution creation (MAJOR - new document)

Modified Principles: N/A (initial creation)
Added Sections:
  - Core Principles (7 principles)
  - Technology Stack Constraints
  - Branch Strategy
  - Security Requirements
  - Workflow File Standards
  - Learning Objectives
  - Governance

Removed Sections: N/A (initial creation)

Templates Status:
  - .specify/templates/plan-template.md: Compatible (no changes needed)
  - .specify/templates/spec-template.md: Compatible (no changes needed)
  - .specify/templates/tasks-template.md: Compatible (no changes needed)
  - .specify/templates/checklist-template.md: Compatible (no changes needed)

Follow-up TODOs: None
===================================================================
-->

# Go-GitHub-Actions CI/CD Training Constitution

## Core Principles

### I. Pipeline as Code

All CI/CD pipeline definitions MUST be version-controlled in `.github/workflows/` with no manual
console configurations. Every workflow change requires code review and MUST follow the same
branching strategy as application code. Workflows MUST be modular, using reusable workflows and
composite actions to prevent duplication across stages.

**Non-Negotiable Rules:**
- Zero manual pipeline modifications in GitHub Actions UI
- All secrets managed via GitHub Secrets/Environments, never hardcoded
- Workflow changes require PR approval with at least one reviewer
- Composite actions for repeated patterns (build, test, deploy)

**Measurable Criteria:**
- 100% of pipeline logic exists in repository YAML files
- 0 manual workflow runs without code changes
- Reusable workflow coverage: >80% of repeated logic abstracted

### II. Progressive Environment Promotion

Code MUST flow through environments in strict order: `alpha` -> `beta` -> `nonprod` -> `prod`.
No environment may be skipped except through documented hotfix procedures. Each promotion requires
passing all quality gates from the previous environment plus environment-specific validation.

**Non-Negotiable Rules:**
- Alpha: Automatic deployment on merge to `develop` branch
- Beta: Manual approval required after alpha validation passes
- NonProd: Requires beta sign-off + smoke tests pass
- Prod: Requires nonprod soak period (minimum 24h) + change advisory board approval
- Hotfix path: `main` -> `prod` with expedited (not eliminated) gates

**Measurable Criteria:**
- 0 direct deployments to prod without nonprod validation
- Environment promotion time tracked and visible
- Rollback capability demonstrated within 5 minutes at each stage

### III. Quality Gates Before Deployment

Every deployment stage MUST pass mandatory quality gates. No code reaches any environment without
passing linting, unit tests, and security scanning. Higher environments add progressive gates:
integration tests for beta, performance tests for nonprod, and manual approval for prod.

**Gate Progression Matrix:**

| Gate | Alpha | Beta | NonProd | Prod |
|------|-------|------|---------|------|
| Lint + Format | Required | Required | Required | Required |
| Unit Tests (>80% coverage) | Required | Required | Required | Required |
| Integration Tests | Optional | Required | Required | Required |
| SonarQube Scan (0 critical) | Required | Required | Required | Required |
| Container Scan (0 critical CVE) | Required | Required | Required | Required |
| Performance Tests | - | - | Required | Required |
| Manual Approval | - | Required | Required | Required |
| Smoke Tests | Auto | Auto | Auto | Auto + Manual |

**Measurable Criteria:**
- Gate pass rate visible in deployment dashboard
- Mean time to gate failure resolution tracked
- Zero deployments bypassing mandatory gates

### IV. Immutable Artifacts and Semantic Versioning

Build artifacts (Docker images) MUST be immutable once created and tagged. The same artifact that
passes alpha MUST be promoted to beta, nonprod, and prod without rebuilding. Version numbers
follow Semantic Versioning (MAJOR.MINOR.PATCH) with automatic tagging based on conventional commits.

**Non-Negotiable Rules:**
- Docker images tagged with Git SHA + semantic version
- No `latest` tag usage in deployment manifests
- Artifacts stored in container registry with retention policy
- Version bumps automated via conventional commit prefixes:
  - `feat:` -> MINOR bump
  - `fix:` -> PATCH bump
  - `BREAKING CHANGE:` -> MAJOR bump
- Release notes auto-generated from commit messages

**Measurable Criteria:**
- 100% of deployments use SHA-tagged images
- Version history traceable in release notes
- Artifact provenance verifiable via container registry

### V. Infrastructure as Code with GitOps

All infrastructure (Docker Compose for dev, Kubernetes manifests for prod) MUST be declarative,
version-controlled, and applied via GitOps principles. Environment-specific configuration uses
overlays/variants, not branching. Changes to infrastructure follow the same PR workflow as
application code.

**Non-Negotiable Rules:**
- Development: `docker-compose.yml` with all services defined
- Production: Kubernetes manifests in `k8s/` directory with Kustomize overlays
- No `kubectl apply` or `docker` commands in CI without manifest files
- Environment variables via ConfigMaps/Secrets, not inline
- Database migrations version-controlled and applied via CI

**Measurable Criteria:**
- 100% of infrastructure defined in repository files
- Environment parity score: dev mirrors prod architecture
- Infrastructure change audit trail in Git history

### VI. Observability by Default

Every deployed service MUST expose health endpoints, structured logs, and metrics from day one.
CI pipelines MUST validate observability requirements before deployment. Alerting rules are
defined as code alongside the services they monitor.

**Non-Negotiable Rules:**
- Health endpoints: `/health/live` (liveness), `/health/ready` (readiness)
- Structured logging: JSON format with correlation IDs
- Metrics: Prometheus-compatible `/metrics` endpoint
- Minimum metrics: request count, latency histogram, error rate
- Deployment includes smoke test hitting health endpoints
- Alert definitions in `alerts/` directory, deployed with application

**Measurable Criteria:**
- 100% of services expose required endpoints
- Log aggregation functional within 5 minutes of deployment
- Alert coverage for critical paths documented

### VII. Fail Fast, Recover Faster

Pipelines MUST fail immediately on first error with clear, actionable messages. Every deployment
MUST have a documented and tested rollback procedure. Recovery time (MTTR) is more important
than prevention time (MTBF) for learning environments.

**Non-Negotiable Rules:**
- Pipeline stages use `fail-fast: true` for parallel jobs
- Error messages include: what failed, why it failed, how to fix it
- Rollback procedure documented and executable in <5 minutes
- Canary/blue-green deployment patterns for production
- Circuit breaker patterns in inter-service communication
- Post-incident review template for learning from failures

**Measurable Criteria:**
- Pipeline failure provides actionable error within 30 seconds
- Rollback tested monthly, executable in documented timeframe
- Post-incident reviews completed within 48 hours of incidents

## Technology Stack Constraints

### Backend
- **Language**: Go 1.22+ with standard library preferred
- **Framework**: Fiber v2 for HTTP routing
- **ORM**: GORM for database access
- **Database**: PostgreSQL 16
- **Linting**: golangci-lint with project configuration
- **Testing**: go test with race detector enabled

### Frontend
- **Framework**: Next.js 14 with App Router
- **Styling**: Tailwind CSS
- **Components**: shadcn/ui component library
- **Build System**: Turborepo for monorepo management
- **Linting**: ESLint + Prettier

### Infrastructure
- **Local Development**: Docker Compose
- **Production**: Kubernetes with Kustomize
- **CI/CD**: GitHub Actions
- **Container Registry**: GitHub Container Registry (ghcr.io)
- **Code Quality**: SonarCloud

### Prerequisites
- Docker & Docker Compose
- Go 1.22+
- Node.js (Latest LTS)
- bun (package manager)

## Branch Strategy

### Branches
- `main`: Production-ready code, protected branch
- `develop`: Integration branch, auto-deploys to alpha
- `feature/*`: Feature development branches
- `hotfix/*`: Emergency production fixes
- `release/*`: Release preparation branches

### Merge Rules
- `feature/*` -> `develop`: Squash merge, requires PR approval
- `develop` -> `main`: Merge commit, requires release checklist
- `hotfix/*` -> `main`: Fast-forward if possible, immediate deploy path

### Protection Rules
- `main`: Require 2 approvals, require status checks, no force push
- `develop`: Require 1 approval, require status checks

## Security Requirements

### Secrets Management
- All secrets via GitHub Secrets or external vault
- No secrets in code, configs, or logs
- Secret rotation documented and automated where possible
- Environment-scoped secrets (alpha secrets != prod secrets)

### Dependency Security
- Dependabot enabled for all package managers
- Weekly dependency vulnerability scan
- Critical vulnerabilities block deployment

### Access Control
- Principle of least privilege for service accounts
- Audit log retention: 90 days minimum
- GitHub environments for deployment approvals

## Workflow File Standards

### Naming Convention
- `ci.yml`: Main CI pipeline (lint, test, build, scan)
- `cd-{environment}.yml`: Environment-specific deployment
- `reusable-{action}.yml`: Reusable workflows in `.github/workflows/`
- `hotfix.yml`: Emergency hotfix deployment workflow

### Required Workflow Attributes
- `name`: Descriptive, includes trigger context
- `on`: Explicit triggers, no implicit defaults
- `concurrency`: Prevent parallel runs where inappropriate
- `timeout-minutes`: All jobs have explicit timeouts
- `env`: Environment variables at workflow level when shared

### CI Pipeline Stages
1. `lint-test`: Linting and static analysis (Go + JS)
2. `unit-test`: Unit testing with coverage
3. `build`: Build Docker images for backend and frontend
4. `sonar-scan`: SonarCloud code quality analysis
5. `deploy_alpha`: Alpha environment deployment
6. `deploy_beta`: Beta environment deployment
7. `deploy_nonprod`: Non-production deployment
8. `deploy_prod`: Production deployment
9. `hotfix`: Emergency hotfix workflow

## Learning Objectives

Students completing this training MUST demonstrate:

### Pipeline Skills
- Create multi-stage GitHub Actions workflow from scratch
- Implement quality gates with conditional execution
- Configure environment-specific deployments
- Set up reusable workflows and composite actions

### Infrastructure Skills
- Write Docker Compose for multi-service local development
- Create Kubernetes manifests with Kustomize overlays
- Implement health checks and readiness probes
- Configure secrets management without hardcoding

### Observability Skills
- Implement structured logging with correlation IDs
- Expose Prometheus metrics from application code
- Create alerting rules for SLI-based monitoring
- Debug deployment issues using logs and metrics

### Recovery Skills
- Execute rollback procedure under time pressure
- Identify root cause from pipeline failure logs
- Complete post-incident review documentation
- Implement fix and verify through full pipeline

## Governance

### Amendment Process
1. **Proposal**: Submit PR to `constitution.md` with rationale
2. **Review Period**: 48-hour comment period for team input
3. **Approval**: Requires 2 approvals from maintainers
4. **Migration**: If breaking, include migration guide
5. **Communication**: Announce changes in team channel

### Version Policy
- **MAJOR**: Principle removed or fundamentally changed
- **MINOR**: New principle added or significant clarification
- **PATCH**: Typo fixes, minor wording improvements

### Compliance Verification
- **PR Review**: Reviewers verify constitution compliance
- **Automated Checks**: CI validates measurable criteria
- **Monthly Audit**: Review adherence to non-automated principles
- **Training Sessions**: New contributors briefed on constitution

### Exception Process
Exceptions to constitution principles require:
1. Written justification in PR description
2. Explicit approval from 2 maintainers
3. Time-bounded scope (must revisit in 30 days)
4. Documentation of exception and rationale

**Version**: 1.0.0 | **Ratified**: 2025-12-29 | **Last Amended**: 2025-12-29
