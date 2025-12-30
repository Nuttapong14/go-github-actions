# Quickstart Guide: CI/CD Training Application

**Feature Branch**: `001-cicd-training-app`
**Time to Complete**: ~5 minutes
**Prerequisites**: Docker, Go 1.22+, Node.js 22 LTS, bun

## Prerequisites Checklist

Before starting, verify you have the following installed:

```bash
# Check Docker (required)
docker --version        # Docker 24.0+ recommended
docker compose version  # Docker Compose v2.20+ required

# Check Go (required for backend development)
go version              # Go 1.22+ required

# Check Node.js (required for frontend development)
node --version          # Node.js 22 LTS required
bun --version           # bun 1.0+ required

# Check Git (required)
git --version           # Git 2.40+ recommended
```

### Missing Prerequisites?

| Tool | Installation |
|------|--------------|
| Docker | [Get Docker](https://docs.docker.com/get-docker/) |
| Go | [Download Go](https://go.dev/dl/) |
| Node.js | [Download Node.js](https://nodejs.org/) |
| bun | `npm install -g bun` or [bun.sh](https://bun.sh) |

## Quick Start (5 Minutes)

### 1. Clone and Enter Repository

```bash
git clone <repository-url>
cd go-github-actions
```

### 2. Start the Full Stack

```bash
# Single command to start everything
docker compose -f docker/docker-compose.yml up -d

# Watch the logs (optional)
docker compose -f docker/docker-compose.yml logs -f
```

This starts:
- **Backend** (Go/Fiber): http://localhost:8080
- **Frontend** (Next.js): http://localhost:3000
- **PostgreSQL**: localhost:5432

### 3. Verify Services

```bash
# Check backend health
curl http://localhost:8080/health/live
# Expected: {"status":"ok","timestamp":"..."}

# Check backend readiness
curl http://localhost:8080/health/ready
# Expected: {"status":"ok","timestamp":"...","checks":{"database":"ok"}}

# Check metrics endpoint
curl http://localhost:8080/metrics
# Expected: Prometheus metrics output

# Open frontend dashboard
open http://localhost:3000
```

### 4. Stop Services

```bash
docker compose -f docker/docker-compose.yml down

# To also remove volumes (database data)
docker compose -f docker/docker-compose.yml down -v
```

## Development Workflow

### Backend Development (Hot Reload)

```bash
# Start with hot reload (Air)
cd backend
go run main.go --auto-migrate

# Or use Air directly
air
```

**Endpoints**:
- `GET /health/live` - Liveness probe
- `GET /health/ready` - Readiness probe (checks DB)
- `GET /metrics` - Prometheus metrics
- `GET /api/v1/environments` - List environments
- `GET /api/v1/releases` - List releases
- `GET /api/v1/deployments` - List deployments

### Frontend Development (Hot Reload)

```bash
cd frontend
bun install
bun dev
```

**Pages**:
- `/` - Dashboard home (status summary)
- `/health` - Health endpoint visualization
- `/deployments` - Deployment history

### Running Tests Locally

```bash
# Backend tests with coverage
cd backend
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # View coverage report

# Frontend tests
cd frontend
bun test
bun test:coverage
```

### Running Linting Locally

```bash
# Backend linting
cd backend
golangci-lint run

# Frontend linting
cd frontend
bun lint
```

### Building Container Images

```bash
# Build backend image
docker build -f docker/backend.Dockerfile -t cicd-training-backend:local .

# Build frontend image
docker build -f docker/frontend.Dockerfile -t cicd-training-frontend:local .
```

## Database Operations

### Run Migrations

```bash
# Using golang-migrate CLI
migrate -path ./backend/database/migrations -database "postgres://postgres:postgres@localhost:5432/cicd_training?sslmode=disable" up

# Rollback last migration
migrate -path ./backend/database/migrations -database "postgres://postgres:postgres@localhost:5432/cicd_training?sslmode=disable" down 1
```

### Connect to Database

```bash
# Using Docker
docker compose -f docker/docker-compose.yml exec db psql -U postgres -d cicd_training

# Common queries
SELECT * FROM environments;
SELECT * FROM releases ORDER BY created_at DESC LIMIT 5;
SELECT * FROM deployments WHERE status = 'running';
```

### Reset Database

```bash
docker compose -f docker/docker-compose.yml down -v
docker compose -f docker/docker-compose.yml up -d db
# Wait for DB to be ready, then run migrations
```

## CI/CD Training Exercises

### Exercise 1: Trigger CI Pipeline

1. Create a feature branch:
   ```bash
   git checkout -b feature/my-first-change
   ```

2. Make a small change (e.g., update a comment)

3. Commit and push:
   ```bash
   git add .
   git commit -m "feat: my first CI/CD change"
   git push -u origin feature/my-first-change
   ```

4. Observe the CI pipeline in GitHub Actions

### Exercise 2: Create a Pull Request

1. Go to GitHub and create a PR from your branch
2. Watch the status checks run
3. Review the SonarCloud analysis
4. Merge when all checks pass

### Exercise 3: Observe Deployment to Alpha

1. After merging to `develop`, watch the CD pipeline
2. Verify deployment in alpha environment
3. Check health endpoints in alpha

### Exercise 4: Promote Through Environments

1. Request promotion to beta (requires approval)
2. After beta validation, request nonprod promotion
3. Wait for 24-hour soak period (or simulate)
4. Request production deployment

### Exercise 5: Execute Hotfix

1. Create hotfix branch from main:
   ```bash
   git checkout main
   git pull
   git checkout -b hotfix/critical-fix
   ```

2. Make fix and push
3. Watch expedited deployment pipeline
4. Verify automatic backport to develop

## Troubleshooting

### Docker Issues

**Container won't start**:
```bash
# Check logs
docker compose -f docker/docker-compose.yml logs backend
docker compose -f docker/docker-compose.yml logs frontend
docker compose -f docker/docker-compose.yml logs db

# Rebuild containers
docker compose -f docker/docker-compose.yml build --no-cache
docker compose -f docker/docker-compose.yml up -d
```

**Port already in use**:
```bash
# Find process using port
lsof -i :8080
lsof -i :3000

# Kill process or change ports in docker-compose.yml
```

### Database Issues

**Connection refused**:
```bash
# Check if DB is running
docker compose -f docker/docker-compose.yml ps db

# Check DB health
docker compose -f docker/docker-compose.yml exec db pg_isready

# Restart DB
docker compose -f docker/docker-compose.yml restart db
```

**Migration errors**:
```bash
# Check migration status
migrate -path ./backend/database/migrations -database "$DATABASE_URL" version

# Force to specific version (use with caution)
migrate -path ./backend/database/migrations -database "$DATABASE_URL" force VERSION
```

### Backend Issues

**Module not found**:
```bash
cd backend
go mod tidy
go mod download
```

**Air not working**:
```bash
# Install Air
go install github.com/air-verse/air@latest

# Check Air config
cat .air.toml
```

### Frontend Issues

**Dependencies not installing**:
```bash
cd frontend
rm -rf node_modules bun.lockb
bun install
```

**Build errors**:
```bash
# Clear Next.js cache
rm -rf .next
bun build
```

## Environment Variables

### Backend

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:postgres@localhost:5432/cicd_training?sslmode=disable` |
| `PORT` | HTTP server port | `8080` |
| `LOG_LEVEL` | Logging level | `info` |
| `ENV` | Environment name | `development` |

### Frontend

| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | Backend API URL | `http://localhost:8080` |
| `NODE_ENV` | Node environment | `development` |

## Useful Commands

```bash
# Start everything
make up

# Stop everything
make down

# Run all tests
make test

# Run linting
make lint

# Build images
make build

# View logs
make logs

# Clean up
make clean
```

## Next Steps

After completing the quickstart:

1. **Read the Constitution**: `.specify/memory/constitution.md`
2. **Review the Spec**: `specs/001-cicd-training-app/spec.md`
3. **Explore the Plan**: `specs/001-cicd-training-app/plan.md`
4. **Start Training Exercises**: `docs/exercises/`

## Support

- **Documentation**: `docs/` directory
- **Troubleshooting**: `docs/troubleshooting.md`
- **Runbooks**: `docs/runbooks/`
