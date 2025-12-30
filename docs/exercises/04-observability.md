# Exercise 4: Promotion and Observability

## Goal
Promote a release through beta → nonprod → prod while observing health checks, metrics, and rollback behavior.

## Prerequisites
- Alpha deployment completed and images exist in GHCR.
- GitHub environments configured with required reviewers for `beta`, `nonprod`, and `prod`.
- Access to GitHub Actions for this repository.

## Step 1: Identify the Image Tag
1. Open the **CD Alpha** workflow summary.
2. Note the `Image Tag` (usually a short SHA).

## Step 2: Promote to Beta
1. Trigger **CD Beta** via `workflow_dispatch`.
2. Provide the image tag from alpha.
3. Approve the environment gate if required.
4. Verify the workflow reports successful health checks.

## Step 3: Promote to NonProd
1. Trigger **CD NonProd** via `workflow_dispatch`.
2. Use the same image tag promoted to beta.
3. Approve the environment gate if required.
4. Verify health checks pass.

## Step 4: Soak Period Enforcement
Production promotion requires a 24h soak after the last successful nonprod deployment.

- The **CD Prod** workflow checks the timestamp of the last successful `cd-nonprod` run.
- Optionally, provide `nonprod_deployed_at` (ISO 8601) to set the soak start time.

If the soak period is not met, the workflow fails with a clear message.

## Step 5: Promote to Prod with Canary
1. Trigger **CD Prod** via `workflow_dispatch`.
2. Provide the same image tag.
3. Approve the environment gate.
4. Watch the canary rollout and health check steps.
5. If health checks pass, the workflow scales to full replicas.

## Step 6: Observe Health and Metrics
- Backend liveness: `GET /health/live`
- Backend readiness: `GET /health/ready`
- Metrics: `GET /metrics`

Use these endpoints to confirm:
- The rollout completed successfully.
- Health checks are stable.
- Metrics reflect requests and deployment activity.

## Step 7: Rollback Drill
To observe rollback behavior:
1. Trigger a promotion with an invalid image tag.
2. Watch the health check failure.
3. Confirm the rollback job runs and reports completion.

## Success Criteria
- Beta and nonprod promotions complete with approvals and health checks.
- Prod promotion enforces the 24h soak period.
- Canary rollout succeeds before full promotion.
- Rollback triggers automatically on failure.
