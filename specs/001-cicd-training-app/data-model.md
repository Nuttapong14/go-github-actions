# Data Model: CI/CD Training Application

**Feature Branch**: `001-cicd-training-app`
**Date**: 2025-12-29
**Purpose**: Define entities and their relationships for the training application

## Entity Overview

The data model is intentionally minimal to focus on CI/CD training rather than application complexity. Entities represent deployment-related concepts that trainees will interact with during pipeline exercises.

```
┌─────────────────┐       ┌─────────────────┐
│   Environment   │       │     Release     │
│─────────────────│       │─────────────────│
│ id              │       │ id              │
│ name            │       │ version         │
│ url             │       │ git_sha         │
│ approval_req    │◄──────│ created_at      │
│ is_production   │       │ changelog       │
└─────────────────┘       └─────────────────┘
         │                        │
         │                        │
         ▼                        ▼
┌─────────────────────────────────────────┐
│              Deployment                  │
│─────────────────────────────────────────│
│ id                                       │
│ environment_id (FK)                      │
│ release_id (FK)                          │
│ status (pending/running/success/failed)  │
│ started_at                               │
│ completed_at                             │
│ triggered_by                             │
│ rollback_of_id (FK, nullable)            │
└─────────────────────────────────────────┘
         │
         │
         ▼
┌─────────────────────────────────────────┐
│             HealthCheck                  │
│─────────────────────────────────────────│
│ id                                       │
│ deployment_id (FK)                       │
│ check_type (liveness/readiness)          │
│ status (passing/failing)                 │
│ response_time_ms                         │
│ checked_at                               │
└─────────────────────────────────────────┘
```

## Entity Definitions

### Environment

Represents a deployment target in the promotion pipeline.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK | Unique identifier |
| name | string | unique, not null | Environment name (alpha, beta, nonprod, prod) |
| url | string | not null | Base URL for the environment |
| approval_required | boolean | default: false | Whether manual approval is needed |
| is_production | boolean | default: false | Production environment flag |
| soak_period_hours | integer | default: 0 | Required wait time before next promotion |
| created_at | timestamp | not null | Creation timestamp |
| updated_at | timestamp | not null | Last update timestamp |

**Seed Data**:
```json
[
  {"name": "alpha", "url": "https://alpha.example.com", "approval_required": false, "is_production": false, "soak_period_hours": 0},
  {"name": "beta", "url": "https://beta.example.com", "approval_required": true, "is_production": false, "soak_period_hours": 0},
  {"name": "nonprod", "url": "https://nonprod.example.com", "approval_required": true, "is_production": false, "soak_period_hours": 24},
  {"name": "prod", "url": "https://prod.example.com", "approval_required": true, "is_production": true, "soak_period_hours": 0}
]
```

### Release

Represents a versioned software release with artifacts.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK | Unique identifier |
| version | string | unique, not null | Semantic version (e.g., 1.2.3) |
| git_sha | string | unique, not null | Full git commit SHA |
| git_sha_short | string | not null | Short git SHA (7 chars) |
| changelog | text | nullable | Auto-generated release notes |
| docker_image_backend | string | not null | Backend image URI with tag |
| docker_image_frontend | string | not null | Frontend image URI with tag |
| created_at | timestamp | not null | Creation timestamp |
| created_by | string | not null | User/system that created release |

**Validation Rules**:
- version must follow semver pattern: `^\d+\.\d+\.\d+(-[a-zA-Z0-9.]+)?$`
- git_sha must be 40 hex characters
- docker_image must include registry/repository:tag format

### Deployment

Represents a single deployment event to an environment.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK | Unique identifier |
| environment_id | UUID | FK(Environment) | Target environment |
| release_id | UUID | FK(Release) | Release being deployed |
| status | enum | not null | pending, running, success, failed, rolled_back |
| started_at | timestamp | nullable | Deployment start time |
| completed_at | timestamp | nullable | Deployment completion time |
| triggered_by | string | not null | User or automation that triggered |
| trigger_type | enum | not null | manual, automatic, hotfix |
| rollback_of_id | UUID | FK(Deployment), nullable | If rollback, original deployment |
| error_message | text | nullable | Error details if failed |
| created_at | timestamp | not null | Record creation timestamp |

**Status State Machine**:
```
pending → running → success
                 ↘ failed → rolled_back
```

**Validation Rules**:
- Only one deployment per environment can be `running` at a time
- rollback_of_id only valid when trigger_type = 'rollback'
- completed_at required when status in (success, failed, rolled_back)

### HealthCheck

Represents health check results after deployment.

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| id | UUID | PK | Unique identifier |
| deployment_id | UUID | FK(Deployment) | Associated deployment |
| check_type | enum | not null | liveness, readiness |
| status | enum | not null | passing, failing |
| response_time_ms | integer | nullable | Response time in milliseconds |
| status_code | integer | nullable | HTTP status code received |
| error_message | text | nullable | Error details if failing |
| checked_at | timestamp | not null | When check was performed |

**Check Types**:
- `liveness`: `/health/live` - Is the service running?
- `readiness`: `/health/ready` - Is the service ready to accept traffic?

## Indexes

### Environment
- `idx_environment_name` on `name` (unique)

### Release
- `idx_release_version` on `version` (unique)
- `idx_release_git_sha` on `git_sha` (unique)
- `idx_release_created_at` on `created_at` (descending)

### Deployment
- `idx_deployment_environment` on `environment_id`
- `idx_deployment_release` on `release_id`
- `idx_deployment_status` on `status`
- `idx_deployment_started_at` on `started_at` (descending)
- `idx_deployment_running` on `environment_id` where `status = 'running'` (partial, unique)

### HealthCheck
- `idx_healthcheck_deployment` on `deployment_id`
- `idx_healthcheck_checked_at` on `checked_at` (descending)

## Migration Files

### 000001_init.up.sql
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE deployment_status AS ENUM ('pending', 'running', 'success', 'failed', 'rolled_back');
CREATE TYPE trigger_type AS ENUM ('manual', 'automatic', 'hotfix', 'rollback');
CREATE TYPE health_check_type AS ENUM ('liveness', 'readiness');
CREATE TYPE health_status AS ENUM ('passing', 'failing');

CREATE TABLE environments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    url VARCHAR(255) NOT NULL,
    approval_required BOOLEAN DEFAULT false,
    is_production BOOLEAN DEFAULT false,
    soak_period_hours INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE releases (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    version VARCHAR(50) UNIQUE NOT NULL,
    git_sha CHAR(40) UNIQUE NOT NULL,
    git_sha_short CHAR(7) NOT NULL,
    changelog TEXT,
    docker_image_backend VARCHAR(255) NOT NULL,
    docker_image_frontend VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by VARCHAR(100) NOT NULL
);

CREATE TABLE deployments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    environment_id UUID NOT NULL REFERENCES environments(id),
    release_id UUID NOT NULL REFERENCES releases(id),
    status deployment_status NOT NULL DEFAULT 'pending',
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    triggered_by VARCHAR(100) NOT NULL,
    trigger_type trigger_type NOT NULL,
    rollback_of_id UUID REFERENCES deployments(id),
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE health_checks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    deployment_id UUID NOT NULL REFERENCES deployments(id),
    check_type health_check_type NOT NULL,
    status health_status NOT NULL,
    response_time_ms INTEGER,
    status_code INTEGER,
    error_message TEXT,
    checked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_release_created_at ON releases(created_at DESC);
CREATE INDEX idx_deployment_environment ON deployments(environment_id);
CREATE INDEX idx_deployment_release ON deployments(release_id);
CREATE INDEX idx_deployment_status ON deployments(status);
CREATE INDEX idx_deployment_started_at ON deployments(started_at DESC);
CREATE UNIQUE INDEX idx_deployment_running ON deployments(environment_id) WHERE status = 'running';
CREATE INDEX idx_healthcheck_deployment ON health_checks(deployment_id);
CREATE INDEX idx_healthcheck_checked_at ON health_checks(checked_at DESC);

-- Seed environments
INSERT INTO environments (name, url, approval_required, is_production, soak_period_hours) VALUES
    ('alpha', 'https://alpha.example.com', false, false, 0),
    ('beta', 'https://beta.example.com', true, false, 0),
    ('nonprod', 'https://nonprod.example.com', true, false, 24),
    ('prod', 'https://prod.example.com', true, true, 0);
```

### 000001_init.down.sql
```sql
DROP TABLE IF EXISTS health_checks;
DROP TABLE IF EXISTS deployments;
DROP TABLE IF EXISTS releases;
DROP TABLE IF EXISTS environments;
DROP TYPE IF EXISTS health_status;
DROP TYPE IF EXISTS health_check_type;
DROP TYPE IF EXISTS trigger_type;
DROP TYPE IF EXISTS deployment_status;
```

## GORM Models

```go
// models/environment.go
type Environment struct {
    ID               uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    Name             string    `gorm:"uniqueIndex;not null;size:50"`
    URL              string    `gorm:"not null;size:255"`
    ApprovalRequired bool      `gorm:"default:false"`
    IsProduction     bool      `gorm:"default:false"`
    SoakPeriodHours  int       `gorm:"default:0"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
    Deployments      []Deployment `gorm:"foreignKey:EnvironmentID"`
}

// models/release.go
type Release struct {
    ID                   uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    Version              string    `gorm:"uniqueIndex;not null;size:50"`
    GitSHA               string    `gorm:"column:git_sha;uniqueIndex;not null;size:40"`
    GitSHAShort          string    `gorm:"column:git_sha_short;not null;size:7"`
    Changelog            *string   `gorm:"type:text"`
    DockerImageBackend   string    `gorm:"not null;size:255"`
    DockerImageFrontend  string    `gorm:"not null;size:255"`
    CreatedAt            time.Time
    CreatedBy            string    `gorm:"not null;size:100"`
    Deployments          []Deployment `gorm:"foreignKey:ReleaseID"`
}

// models/deployment.go
type DeploymentStatus string
type TriggerType string

const (
    StatusPending    DeploymentStatus = "pending"
    StatusRunning    DeploymentStatus = "running"
    StatusSuccess    DeploymentStatus = "success"
    StatusFailed     DeploymentStatus = "failed"
    StatusRolledBack DeploymentStatus = "rolled_back"
)

const (
    TriggerManual    TriggerType = "manual"
    TriggerAutomatic TriggerType = "automatic"
    TriggerHotfix    TriggerType = "hotfix"
    TriggerRollback  TriggerType = "rollback"
)

type Deployment struct {
    ID            uuid.UUID        `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    EnvironmentID uuid.UUID        `gorm:"type:uuid;not null;index"`
    ReleaseID     uuid.UUID        `gorm:"type:uuid;not null;index"`
    Status        DeploymentStatus `gorm:"type:deployment_status;not null;default:'pending';index"`
    StartedAt     *time.Time
    CompletedAt   *time.Time
    TriggeredBy   string           `gorm:"not null;size:100"`
    TriggerType   TriggerType      `gorm:"type:trigger_type;not null"`
    RollbackOfID  *uuid.UUID       `gorm:"type:uuid"`
    ErrorMessage  *string          `gorm:"type:text"`
    CreatedAt     time.Time

    Environment   Environment  `gorm:"foreignKey:EnvironmentID"`
    Release       Release      `gorm:"foreignKey:ReleaseID"`
    RollbackOf    *Deployment  `gorm:"foreignKey:RollbackOfID"`
    HealthChecks  []HealthCheck `gorm:"foreignKey:DeploymentID"`
}

// models/health_check.go
type HealthCheckType string
type HealthStatus string

const (
    CheckTypeLiveness  HealthCheckType = "liveness"
    CheckTypeReadiness HealthCheckType = "readiness"
)

const (
    HealthPassing HealthStatus = "passing"
    HealthFailing HealthStatus = "failing"
)

type HealthCheck struct {
    ID             uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    DeploymentID   uuid.UUID       `gorm:"type:uuid;not null;index"`
    CheckType      HealthCheckType `gorm:"type:health_check_type;not null"`
    Status         HealthStatus    `gorm:"type:health_status;not null"`
    ResponseTimeMs *int
    StatusCode     *int
    ErrorMessage   *string         `gorm:"type:text"`
    CheckedAt      time.Time       `gorm:"default:now()"`

    Deployment     Deployment      `gorm:"foreignKey:DeploymentID"`
}
```

## Relationships Summary

| From | To | Type | Description |
|------|-----|------|-------------|
| Deployment | Environment | Many-to-One | Each deployment targets one environment |
| Deployment | Release | Many-to-One | Each deployment deploys one release |
| Deployment | Deployment | Self-ref (optional) | Rollback references original deployment |
| HealthCheck | Deployment | Many-to-One | Multiple checks per deployment |
