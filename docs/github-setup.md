# GitHub Repository Setup Guide

This guide documents how to configure your GitHub repository for optimal CI/CD pipeline execution, including branch protection rules, secrets management, and required integrations.

## Table of Contents

1. [Repository Settings](#repository-settings)
2. [Branch Protection Rules](#branch-protection-rules)
3. [Secrets Configuration](#secrets-configuration)
4. [SonarCloud Integration](#sonarcloud-integration)
5. [GitHub Container Registry](#github-container-registry)
6. [Notifications and Webhooks](#notifications-and-webhooks)

## Repository Settings

### General Settings

Navigate to **Settings → General** and configure:

1. **Default Branch**: Set to `main`
2. **Features**:
   - ✅ Issues
   - ✅ Projects
   - ✅ Preserve this repository
3. **Pull Requests**:
   - ✅ Allow merge commits
   - ✅ Allow squash merging (recommended default)
   - ✅ Allow rebase merging
   - ✅ Automatically delete head branches

### Actions Permissions

Navigate to **Settings → Actions → General**:

1. **Actions permissions**: Allow all actions and reusable workflows
2. **Workflow permissions**:
   - ✅ Read and write permissions
   - ✅ Allow GitHub Actions to create and approve pull requests
3. **Fork pull request workflows**:
   - ⚠️ Require approval for first-time contributors

## Branch Protection Rules

### Main Branch Protection

Navigate to **Settings → Branches → Add branch protection rule**

**Branch name pattern**: `main`

#### Required Settings

```yaml
# Basic Protection
✅ Require a pull request before merging
  ✅ Require approvals: 1
  ✅ Dismiss stale pull request approvals when new commits are pushed
  ✅ Require review from Code Owners (if CODEOWNERS file exists)

# Status Checks
✅ Require status checks to pass before merging
  ✅ Require branches to be up to date before merging

  Required checks:
  - CI Pipeline / Lint / Backend Lint (Go)
  - CI Pipeline / Lint / Frontend Lint (TypeScript)
  - CI Pipeline / Test / Backend Tests (Go)
  - CI Pipeline / Test / Frontend Tests (TypeScript)
  - CI Pipeline / Build / Build Backend (Go)
  - CI Pipeline / Build / Build Frontend (Next.js)
  - CI Pipeline / SonarCloud Scan
  - CI Pipeline / CI Summary

# Additional Protection
✅ Require conversation resolution before merging
✅ Do not allow bypassing the above settings
✅ Restrict who can push to matching branches
  - Only allow administrators
```

#### Recommended Additional Settings

```yaml
❌ Require signed commits (optional, for security-sensitive repos)
❌ Require linear history (optional, depends on team preference)
✅ Include administrators (enforce rules for everyone)
✅ Allow force pushes → Disabled
✅ Allow deletions → Disabled
```

### Develop Branch Protection

**Branch name pattern**: `develop`

```yaml
# Similar to main, but with relaxed settings for development
✅ Require a pull request before merging
  ✅ Require approvals: 1

✅ Require status checks to pass before merging
  Required checks:
  - CI Pipeline / Lint / Backend Lint (Go)
  - CI Pipeline / Lint / Frontend Lint (TypeScript)
  - CI Pipeline / Test / Backend Tests (Go)
  - CI Pipeline / Test / Frontend Tests (TypeScript)
  - CI Pipeline / CI Summary

✅ Do not allow bypassing the above settings
```

### Feature Branch Pattern

**Branch name pattern**: `feature/**`

```yaml
# Lighter protection for feature branches
❌ Require a pull request before merging (direct push allowed)
✅ Require status checks to pass before merging
  Required checks:
  - CI Pipeline / Lint / Backend Lint (Go)
  - CI Pipeline / Lint / Frontend Lint (TypeScript)
```

## Secrets Configuration

Navigate to **Settings → Secrets and variables → Actions**

### Required Secrets

| Secret Name | Description | How to Obtain |
|-------------|-------------|---------------|
| `SONAR_TOKEN` | SonarCloud authentication token | See [SonarCloud Integration](#sonarcloud-integration) |

### Repository Variables

Navigate to **Settings → Secrets and variables → Actions → Variables**

| Variable Name | Value | Description |
|---------------|-------|-------------|
| `GO_VERSION` | `1.23` | Go version for CI |
| `NODE_VERSION` | `22` | Node.js version for CI |
| `COVERAGE_THRESHOLD` | `80` | Minimum coverage percentage |

## SonarCloud Integration

### Initial Setup

1. **Create SonarCloud Account**:
   - Go to [sonarcloud.io](https://sonarcloud.io)
   - Sign up with your GitHub account
   - Import your organization

2. **Add Project**:
   - Click "Analyze new project"
   - Select `go-github-actions` repository
   - Choose "GitHub Actions" as analysis method

3. **Get Token**:
   - Go to My Account → Security
   - Generate new token with name: `go-github-actions-ci`
   - Copy the token

4. **Add Secret to GitHub**:
   - Go to repository Settings → Secrets → Actions
   - Add new secret: `SONAR_TOKEN` with the token value

### Quality Gate Configuration

In SonarCloud, configure Quality Gate:

1. Go to **Quality Gates** → Create new or use default
2. Recommended conditions:

| Metric | Operator | Value |
|--------|----------|-------|
| Coverage | is less than | 80% |
| Duplicated Lines (%) | is greater than | 3% |
| Maintainability Rating | is worse than | A |
| Reliability Rating | is worse than | A |
| Security Rating | is worse than | A |
| Security Hotspots Reviewed | is less than | 100% |

### Project Configuration

The `sonar-project.properties` file in the repository root configures:

```properties
sonar.projectKey=go-github-actions
sonar.organization=devops-training

# Source locations
sonar.sources=backend,frontend
sonar.tests=backend/tests,frontend/tests

# Coverage reports
sonar.go.coverage.reportPaths=backend/coverage.out
sonar.javascript.lcov.reportPaths=frontend/coverage/lcov.info
```

## GitHub Container Registry

### Enable GHCR

1. Ensure GitHub Packages is enabled for your organization
2. The CI workflow uses `GITHUB_TOKEN` for authentication (automatic)

### Package Visibility

Navigate to **Packages** → Select package → **Package settings**:

- **Visibility**: Internal or Public (as appropriate)
- **Manage Actions access**: Add repository with write access

### Pull Images

```bash
# Login to GHCR
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# Pull images
docker pull ghcr.io/OWNER/cicd-training-backend:latest
docker pull ghcr.io/OWNER/cicd-training-frontend:latest
```

## Notifications and Webhooks

### Email Notifications

Configure in **Settings → Notifications**:

- ✅ Receive notifications for failed workflows only
- ✅ Receive notifications when workflows require intervention

### Slack Integration (Optional)

1. Create Slack webhook URL
2. Add as secret: `SLACK_WEBHOOK_URL`
3. Add notification step to workflows

### Status Badges

Add to README.md:

```markdown
![CI Pipeline](https://github.com/OWNER/go-github-actions/actions/workflows/ci.yml/badge.svg)
![SonarCloud](https://sonarcloud.io/api/project_badges/measure?project=go-github-actions&metric=alert_status)
```

## Verification Checklist

After setup, verify:

- [ ] Branch protection rules are active on `main` and `develop`
- [ ] `SONAR_TOKEN` secret is configured
- [ ] CI workflow runs on push to protected branches
- [ ] Pull requests trigger CI checks
- [ ] Status checks appear in PR merge requirements
- [ ] SonarCloud receives analysis results
- [ ] Build artifacts are uploaded successfully

## Troubleshooting

### CI Checks Not Running

1. Verify workflow file syntax with `yamllint`
2. Check Actions tab for workflow errors
3. Ensure branch matches trigger patterns

### SonarCloud Analysis Fails

1. Verify `SONAR_TOKEN` is set correctly
2. Check SonarCloud project configuration
3. Ensure coverage files are uploaded as artifacts

### Branch Protection Bypass

1. Verify "Include administrators" is checked
2. Check for rules with higher priority
3. Ensure required status checks are correctly named

## Next Steps

After configuring your repository:

1. Make a small change on a feature branch
2. Push and observe the CI pipeline
3. Create a pull request to `develop`
4. Verify all status checks pass
5. Merge and observe the develop branch CI

See [Exercise 03: Push and Observe Pipeline](./exercises/03-deployment.md) for hands-on practice.
