# Exercise 03: Push Code and Observe Pipeline

This exercise teaches you how to push code to GitHub and observe the CI/CD pipeline executing in real-time. You'll learn to interpret pipeline stages, understand job dependencies, and troubleshoot common issues.

## Learning Objectives

By the end of this exercise, you will be able to:

1. Create feature branches following naming conventions
2. Push changes and trigger the CI pipeline
3. Navigate the GitHub Actions interface
4. Interpret pipeline stages and job outputs
5. Understand status checks and their requirements
6. Create pull requests with passing CI checks

## Prerequisites

Before starting this exercise, ensure you have:

- Completed [Exercise 01: Local Development Setup](./01-local-setup.md)
- Completed [Exercise 02: Running CI Pipeline Locally](./02-ci-pipeline.md)
- GitHub repository access with push permissions
- Git configured with your GitHub credentials

## Exercise Steps

### Step 1: Understanding the CI Pipeline Structure

Before pushing code, understand the pipeline stages:

```
┌─────────────────────────────────────────────────────────────────┐
│                      CI Pipeline                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Stage 1: LINT                                                   │
│  ┌──────────────────┐    ┌───────────────────┐                  │
│  │ Backend Lint (Go)│    │ Frontend Lint (TS)│                  │
│  │ - golangci-lint  │    │ - ESLint          │                  │
│  │ - gofmt check    │    │ - TypeScript      │                  │
│  └────────┬─────────┘    └─────────┬─────────┘                  │
│           │                        │                             │
│           └────────────┬───────────┘                             │
│                        ▼                                         │
│  Stage 2: TEST                                                   │
│  ┌──────────────────┐    ┌───────────────────┐                  │
│  │ Backend Tests    │    │ Frontend Tests    │                  │
│  │ - Unit tests     │    │ - Component tests │                  │
│  │ - Race detection │    │ - Coverage        │                  │
│  │ - Coverage: 80%+ │    │                   │                  │
│  └────────┬─────────┘    └─────────┬─────────┘                  │
│           │                        │                             │
│           └────────────┬───────────┘                             │
│                        ▼                                         │
│  Stage 3: BUILD                                                  │
│  ┌──────────────────┐    ┌───────────────────┐                  │
│  │ Backend Build    │    │ Frontend Build    │                  │
│  │ - Go binary      │    │ - Next.js build   │                  │
│  │ - Version embed  │    │ - Static assets   │                  │
│  └────────┬─────────┘    └─────────┬─────────┘                  │
│           │                        │                             │
│           └────────────┬───────────┘                             │
│                        ▼                                         │
│  Stage 4: SCAN                                                   │
│  ┌──────────────────────────────────────────┐                   │
│  │ SonarCloud Analysis                       │                   │
│  │ - Code quality                            │                   │
│  │ - Security vulnerabilities                │                   │
│  │ - Technical debt                          │                   │
│  └────────┬─────────────────────────────────┘                   │
│           ▼                                                      │
│  ┌──────────────────────────────────────────┐                   │
│  │ CI Summary                                │                   │
│  │ - Aggregates all results                  │                   │
│  │ - Final status report                     │                   │
│  └──────────────────────────────────────────┘                   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Step 2: Create a Feature Branch

Create a new branch for your changes:

```bash
# Ensure you're on the latest develop branch
git checkout develop
git pull origin develop

# Create a feature branch
git checkout -b feature/add-readme-badge
```

### Step 3: Make a Simple Change

Add a CI status badge to the README:

```bash
# Edit the README.md file
# Add the following line near the top:
```

```markdown
# CI/CD Training Application

![CI Pipeline](https://github.com/YOUR_USERNAME/go-github-actions/actions/workflows/ci.yml/badge.svg)

[Rest of README content...]
```

### Step 4: Run Local CI Before Pushing

Always verify locally before pushing:

```bash
# Run the full local CI
make ci

# Or use the simulation script
./scripts/local-ci.sh
```

Expected output:
```
╔══════════════════════════════════════════════════════════════╗
║              ✓ CI Pipeline Passed Successfully               ║
╚══════════════════════════════════════════════════════════════╝

Artifacts:
  • Backend binary: bin/backend
  • Frontend build: frontend/.next/
  • Coverage report: backend/coverage.html
```

### Step 5: Commit and Push

Commit your changes with a descriptive message:

```bash
# Stage changes
git add README.md

# Commit with conventional commit message
git commit -m "docs: add CI pipeline status badge to README"

# Push to GitHub
git push -u origin feature/add-readme-badge
```

### Step 6: Observe the Pipeline in GitHub

1. **Navigate to Actions Tab**:
   - Go to your repository on GitHub
   - Click the "Actions" tab
   - You should see "CI Pipeline" workflow running

2. **Watch the Pipeline Progress**:
   - Click on the running workflow
   - Observe the job graph showing dependencies
   - Watch jobs transition: queued → in progress → completed

3. **View Job Logs**:
   - Click on any job to see detailed logs
   - Expand steps to see command output
   - Look for timing information

### Step 7: Understanding the Actions Interface

#### Workflow Run View

```
Workflow Run: CI Pipeline #42
Triggered by: push to feature/add-readme-badge
Duration: 3m 45s
Status: ✓ Success

Jobs:
├── ✓ Lint (2m 12s)
│   ├── ✓ Backend Lint (Go)
│   └── ✓ Frontend Lint (TypeScript)
├── ✓ Test (1m 45s)
│   ├── ✓ Backend Tests (Go)
│   └── ✓ Frontend Tests (TypeScript)
├── ✓ Build (1m 30s)
│   ├── ✓ Build Backend (Go)
│   └── ✓ Build Frontend (Next.js)
├── ✓ SonarCloud Scan (45s)
└── ✓ CI Summary
```

#### Job Details

Click on a job to see:
- Step-by-step execution
- Command outputs
- Timing for each step
- Annotations (errors/warnings)
- Artifacts produced

### Step 8: Create a Pull Request

1. **Click "Compare & pull request"** (appears after push)

2. **Fill in PR details**:
   - Title: `docs: add CI pipeline status badge to README`
   - Description:
     ```markdown
     ## Summary
     - Adds CI pipeline status badge to README for visibility

     ## Test Plan
     - [x] Local CI passes
     - [x] Badge URL is correct
     ```

3. **Observe Status Checks**:
   - Required checks will appear
   - All must pass before merging
   - Click "Details" to see specific job output

### Step 9: Review Status Check Requirements

The PR should show:

```
✓ All checks have passed
8 successful checks

CI Pipeline / Lint / Backend Lint (Go) — Successful in 45s
CI Pipeline / Lint / Frontend Lint (TypeScript) — Successful in 38s
CI Pipeline / Test / Backend Tests (Go) — Successful in 52s
CI Pipeline / Test / Frontend Tests (TypeScript) — Successful in 28s
CI Pipeline / Build / Build Backend (Go) — Successful in 35s
CI Pipeline / Build / Build Frontend (Next.js) — Successful in 42s
CI Pipeline / SonarCloud Scan — Successful in 45s
CI Pipeline / CI Summary — Successful in 5s
```

### Step 10: Merge the Pull Request

Once all checks pass:

1. Click "Squash and merge"
2. Confirm the merge
3. Delete the feature branch

## Practice Exercises

### Exercise A: Trigger a Lint Failure

1. Create a new branch: `feature/lint-failure-test`
2. Add an unused variable to a Go file:
   ```go
   func SomeFunction() {
       unusedVar := "this will fail lint"
       // ...
   }
   ```
3. Push and observe the lint job fail
4. Fix the issue and push again
5. Observe the pipeline succeed

### Exercise B: Trigger a Test Failure

1. Create branch: `feature/test-failure-test`
2. Break a test assertion intentionally
3. Push and observe the test job fail
4. Check the test output in Actions
5. Fix and verify success

### Exercise C: Observe SonarCloud Analysis

1. Push code with a code smell (e.g., duplicate code)
2. Check SonarCloud dashboard after pipeline completes
3. Review the analysis findings
4. Fix issues and observe quality gate improvement

## Troubleshooting

### Pipeline Doesn't Start

**Symptoms**: No workflow run appears in Actions tab

**Solutions**:
1. Check that workflow file exists in `.github/workflows/ci.yml`
2. Verify branch name matches trigger patterns
3. Check for YAML syntax errors
4. Ensure Actions is enabled in repository settings

### Job Fails Immediately

**Symptoms**: Job shows as failed without running steps

**Solutions**:
1. Check for invalid YAML syntax
2. Verify action versions exist (e.g., `actions/checkout@v4`)
3. Look for missing required inputs
4. Check repository secrets are configured

### Status Checks Not Appearing in PR

**Symptoms**: PR shows no status checks

**Solutions**:
1. Verify branch protection rules are configured
2. Check workflow trigger includes `pull_request`
3. Ensure status check names match exactly
4. Wait a few seconds for checks to register

### SonarCloud Analysis Fails

**Symptoms**: SonarCloud job fails

**Solutions**:
1. Verify `SONAR_TOKEN` secret is set
2. Check `sonar-project.properties` configuration
3. Ensure coverage files are generated
4. Verify SonarCloud project is set up

## Key Takeaways

1. **Always run CI locally first** - Use `make ci` before pushing
2. **Pipeline stages have dependencies** - Test requires Lint to pass
3. **Status checks protect branches** - All must pass to merge
4. **Artifacts are preserved** - Download from workflow runs
5. **Logs are searchable** - Use browser search in job logs
6. **Failures provide context** - Read error messages carefully

## Performance Tips

- **Use caching effectively** - Dependencies are cached between runs
- **Parallelize when possible** - Independent jobs run concurrently
- **Skip unnecessary work** - Use path filters for targeted changes
- **Monitor duration trends** - Investigate sudden increases

## Next Steps

After completing this exercise, proceed to:
- [Exercise 04: Deploy to Alpha Environment](./04-alpha-deployment.md) - Learn about automated deployments
- Review the [GitHub Actions Documentation](https://docs.github.com/en/actions)
- Explore [SonarCloud Dashboard](https://sonarcloud.io) features

## Quick Reference

| Action | Command |
|--------|---------|
| Create branch | `git checkout -b feature/name` |
| Run local CI | `make ci` |
| Push branch | `git push -u origin branch-name` |
| View Actions | GitHub → Actions tab |
| View logs | Click job → Expand steps |
| Download artifacts | Workflow run → Artifacts section |
