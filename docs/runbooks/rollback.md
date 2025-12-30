# Rollback Runbook

## Purpose
Restore a previously known-good release when a deployment fails or degrades production behavior.

## When to Use
- Health checks fail after a deployment.
- Error rate or latency regressions exceed thresholds.
- Critical user-facing functionality is broken.

## Preconditions
- Access to the cluster and namespace.
- Identify the target environment (alpha, beta, nonprod, prod).
- Identify the previous image tag from the workflow summary.

## Steps
1. Confirm the environment and namespace.
2. Roll back backend and frontend deployments:
   ```bash
   kubectl rollout undo deployment/backend -n <namespace>
   kubectl rollout undo deployment/frontend -n <namespace>
   ```
3. Monitor rollout status (2-minute max recommended):
   ```bash
   kubectl rollout status deployment/backend -n <namespace> --timeout=120s
   kubectl rollout status deployment/frontend -n <namespace> --timeout=120s
   ```
4. Validate service health:
   ```bash
   curl https://api.<env>.example.com/health/live
   curl https://api.<env>.example.com/health/ready
   curl https://<env>.example.com/
   ```

## Alternative: Re-apply Kustomize with Previous Tag
1. Update the image tag to the previous known-good value:
   ```bash
   cd k8s/overlays/<env>
   kustomize edit set image backend=<backend-image>:<previous-tag>
   kustomize edit set image frontend=<frontend-image>:<previous-tag>
   kubectl apply -k .
   ```
2. Re-run the health checks above.

## Communication
- Notify the incident channel with the rollback reason and status.
- Link the workflow run and commit SHA in the notification.

## Follow-Up
- Open or update a post-incident review document.
- Identify root cause and prevent recurrence.
