-- Create extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create custom enum types
CREATE TYPE deployment_status AS ENUM ('pending', 'running', 'success', 'failed', 'rolled_back');
CREATE TYPE trigger_type AS ENUM ('manual', 'automatic', 'hotfix', 'rollback');
CREATE TYPE health_check_type AS ENUM ('liveness', 'readiness');
CREATE TYPE health_status AS ENUM ('passing', 'failing');

-- Create environments table
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

-- Create releases table
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

-- Create deployments table
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

-- Create health_checks table
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

-- Create indexes for better query performance
CREATE INDEX idx_environment_name ON environments(name);
CREATE INDEX idx_release_version ON releases(version);
CREATE INDEX idx_release_git_sha ON releases(git_sha);
CREATE INDEX idx_release_created_at ON releases(created_at DESC);
CREATE INDEX idx_deployment_environment ON deployments(environment_id);
CREATE INDEX idx_deployment_release ON deployments(release_id);
CREATE INDEX idx_deployment_status ON deployments(status);
CREATE INDEX idx_deployment_started_at ON deployments(started_at DESC);
CREATE UNIQUE INDEX idx_deployment_running ON deployments(environment_id) WHERE status = 'running';
CREATE INDEX idx_healthcheck_deployment ON health_checks(deployment_id);
CREATE INDEX idx_healthcheck_checked_at ON health_checks(checked_at DESC);

-- Seed environments table with default data
INSERT INTO environments (name, url, approval_required, is_production, soak_period_hours) VALUES
    ('alpha', 'https://alpha.example.com', false, false, 0),
    ('beta', 'https://beta.example.com', true, false, 0),
    ('nonprod', 'https://nonprod.example.com', true, false, 24),
    ('prod', 'https://prod.example.com', true, true, 0);
