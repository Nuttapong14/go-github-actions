# Troubleshooting

## Docker Compose File Not Found
If `docker compose up` fails, make sure to reference the compose file:

```bash
docker compose -f docker/docker-compose.yml up -d
```

## Ports Already In Use
Check what is using the port:

```bash
lsof -i :8080
lsof -i :3000
```

Stop the process or change the port mappings in `docker/docker-compose.yml`.

## Database Not Ready
The backend waits for the DB, but you can verify readiness:

```bash
docker compose -f docker/docker-compose.yml ps db
docker compose -f docker/docker-compose.yml exec db pg_isready -U postgres
```

## Backend Fails to Connect to DB
Confirm environment variables in `docker/docker-compose.yml` and container logs:

```bash
docker compose -f docker/docker-compose.yml logs backend
```

## Frontend Build or Dev Errors
Clear caches and reinstall deps:

```bash
cd frontend
rm -rf node_modules .next bun.lockb
bun install
```

## Go Module Issues
Refresh dependencies:

```bash
cd backend
go mod tidy
go mod download
```

## GitHub Actions Failures
- Missing GHCR permissions: ensure workflow has `packages: write`.
- Missing secrets: verify `SONAR_TOKEN` and environment approvals.
- Kustomize step fails: check `kustomize` version and manifests.

## Kustomize Build Errors
Validate overlays locally:

```bash
kustomize build k8s/overlays/alpha
kustomize build k8s/overlays/beta
kustomize build k8s/overlays/nonprod
kustomize build k8s/overlays/prod
```
