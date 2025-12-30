# Hotfix Runbook

## Purpose
Deploy an emergency fix directly to production with an expedited workflow and required approvals.

## Trigger
- Push to a `hotfix/*` branch to trigger the hotfix workflow.
- Or use `workflow_dispatch` on the hotfix workflow.

## Steps
1. Create a hotfix branch from `main`:
   ```bash
   git checkout main
   git pull
   git checkout -b hotfix/<short-description>
   ```
2. Apply the fix and run tests locally:
   ```bash
   make test
   ```
3. Push the hotfix branch:
   ```bash
   git push -u origin hotfix/<short-description>
   ```
4. Monitor the **Hotfix** workflow in GitHub Actions.
5. Approve the `prod` environment gate when prompted.
6. Verify deployment health:
   - `GET /health/live`
   - `GET /health/ready`
   - `GET /metrics`
7. Confirm the backport PR to `develop` and merge after review.

## Rollback
If health checks fail, the workflow triggers rollback automatically with a 2-minute timeout.
For manual rollback, follow `docs/runbooks/rollback.md`.

## Notes
- Hotfixes bypass beta and nonprod promotions, but do not bypass prod approvals.
- Always merge the backport PR to keep `develop` in sync.
