# CI/CD Training Application

A minimal Go + Next.js application for teaching CI/CD with GitHub Actions. It demonstrates local development, quality gates, environment promotion, observability, and hotfix workflows.

## Quickstart

```bash
# Start the full stack (backend, frontend, database)
docker compose -f docker/docker-compose.yml up -d

# View logs (optional)
docker compose -f docker/docker-compose.yml logs -f
```

- Backend: http://localhost:8080
- Frontend: http://localhost:3000

For a detailed walkthrough, see `specs/001-cicd-training-app/quickstart.md`.

## Learning Path

1. Local setup and hot reload: `docs/exercises/01-local-setup.md`
2. Local CI commands: `docs/exercises/02-ci-pipeline.md`
3. Observe the CI pipeline: `docs/exercises/03-deployment.md`
4. Promote releases and observe health: `docs/exercises/04-observability.md`

## Repository Layout

- Backend (Go/Fiber): `backend/`
- Frontend (Next.js): `frontend/`
- Docker: `docker/`
- Kubernetes: `k8s/`
- GitHub Actions: `.github/workflows/`
- Docs and exercises: `docs/`
- Alerts: `alerts/`

## Local Development Commands

```bash
make up        # Start stack
make down      # Stop stack
make lint      # Lint backend + frontend
make test      # Test backend + frontend
make build     # Build backend binary + frontend
make ci        # Run lint -> test -> build
```

## CI/CD Overview

- CI: `.github/workflows/ci.yml`
- Deploy Alpha: `.github/workflows/cd-alpha.yml`
- Promote Beta/NonProd/Prod: `.github/workflows/cd-beta.yml`, `cd-nonprod.yml`, `cd-prod.yml`
- Hotfix: `.github/workflows/hotfix.yml`

## Observability

- `/health/live`
- `/health/ready`
- `/metrics`

Alert rules live in `alerts/`.

## Support

- Troubleshooting: `docs/troubleshooting.md`
- Runbooks: `docs/runbooks/`
