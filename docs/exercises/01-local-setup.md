# Exercise 1: Local Setup

## Goal
Start the full stack locally and verify health endpoints.

## Prerequisites
- Docker Engine and Docker Compose v2
- Go 1.22+
- Node.js 22+
- bun

## Steps

1. Clone and enter the repo:
   ```bash
   git clone <repository-url>
   cd go-github-actions
   ```

2. Start the stack:
   ```bash
   docker compose -f docker/docker-compose.yml up -d
   ```

3. Verify backend health:
   ```bash
   curl http://localhost:8080/health/live
   curl http://localhost:8080/health/ready
   ```

4. Verify frontend:
   ```bash
   curl http://localhost:3000
   ```

5. Watch logs (optional):
   ```bash
   docker compose -f docker/docker-compose.yml logs -f
   ```

6. Stop the stack:
   ```bash
   docker compose -f docker/docker-compose.yml down
   ```

## Success Criteria
- Backend health endpoints return 200.
- Frontend responds on port 3000.
- Stack starts with a single command.

## Troubleshooting
If something fails, check `docs/troubleshooting.md`.
