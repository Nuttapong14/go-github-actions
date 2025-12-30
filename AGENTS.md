# Repository Guidelines

## Project Structure & Module Organization
- `backend/` Go Fiber API, GORM models, middleware, handlers, and `*_test.go` tests.
- `frontend/` Next.js App Router app, UI components, and `tests/` for component tests.
- `docker/` Dockerfiles and `docker-compose.yml` for the local stack.
- `k8s/` Kustomize base plus overlays (`alpha/`, `beta/`, `nonprod/`, `prod/`).
- `.github/workflows/` CI/CD workflows and `.github/actions/` composite actions.
- `docs/` exercises and runbooks; `alerts/` for Prometheus rules; `specs/` for design docs.

## Build, Test, and Development Commands
- `make up` / `make down` / `make logs`: start, stop, and inspect the local stack (Docker Compose).
- `make backend.dev` / `make frontend.dev`: hot reload backend (Air) and frontend.
- `make lint`, `make test`, `make build`: local equivalents of CI stages.
- `make ci`: full local pipeline (lint -> test -> build).
- Direct compose example: `docker compose -f docker/docker-compose.yml up -d`.

## Coding Style & Naming Conventions
- Go: run `gofmt` and `golangci-lint` (`make backend.lint`); keep idiomatic naming; tests are `*_test.go`.
- Frontend: `next lint` via `bun lint`; tests use `*.test.ts` / `*.test.tsx`.
- Indentation: gofmt uses tabs; TypeScript/TSX uses 2 spaces.

## Testing Guidelines
- Backend: `go test -v -race -coverprofile=coverage.out ./...` (`make backend.test`).
- Frontend: `bun test` and `bun test:coverage`.
- Target at least 80% coverage per spec; add tests for new endpoints and UI components.

## Commit & Pull Request Guidelines
- Git history only includes an initial commit, so there is no established message convention yet.
- Suggested: short, imperative subject lines; consider Conventional Commits for consistency.
- PRs should include a clear summary, test evidence (`make ci` or relevant subset), and screenshots for UI changes. Link related specs/issues when applicable.

## Security & Configuration Tips
- Do not commit secrets; use `.env*` locally and Kubernetes secrets for deployments.
- Keep Docker contexts lean (`.dockerignore` is in place). For emergency changes, follow `docs/runbooks/`.
