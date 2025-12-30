# Makefile for CI/CD Training Application
# Provides commands for local development, testing, and building
# Mirrors GitHub Actions CI/CD pipeline for local execution

.PHONY: help up down logs build lint test clean install ci docker.build docker.tag

# Default target
.DEFAULT_GOAL := help

# Variables
DOCKER_COMPOSE := docker compose -f docker/docker-compose.yml
GO := go
BUN := bun
BACKEND_DIR := backend
FRONTEND_DIR := frontend
DOCKER_DIR := docker

# Git information for tagging
GIT_SHA := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
GIT_TAG := $(shell git describe --tags --exact-match 2>/dev/null || echo "")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Docker image names
DOCKER_REGISTRY := ghcr.io/nuttapong14
BACKEND_IMAGE := $(DOCKER_REGISTRY)/cicd-training-backend
FRONTEND_IMAGE := $(DOCKER_REGISTRY)/cicd-training-frontend

# Image tag logic: use git tag if available, otherwise sha
IMAGE_TAG := $(if $(GIT_TAG),$(GIT_TAG),$(GIT_SHA))

# Environment variables
export GOFLAGS := -mod=mod
export NODE_ENV := development

## help: Display this help message
help:
	@echo "CI/CD Training Application - Development Commands"
	@echo ""
	@echo "Build Info:"
	@echo "  Git SHA:    $(GIT_SHA)"
	@echo "  Git Branch: $(GIT_BRANCH)"
	@echo "  Git Tag:    $(if $(GIT_TAG),$(GIT_TAG),none)"
	@echo "  Image Tag:  $(IMAGE_TAG)"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make up           - Start all services with Docker Compose"
	@echo "  make down         - Stop all services"
	@echo "  make logs         - View logs from all services"
	@echo "  make rebuild      - Rebuild and restart all services"
	@echo ""
	@echo "Backend Commands:"
	@echo "  make backend.dev  - Run backend with hot-reload (Air)"
	@echo "  make backend.run  - Run backend without hot-reload"
	@echo "  make backend.test - Run backend tests with coverage"
	@echo "  make backend.lint - Run backend linting (golangci-lint)"
	@echo "  make backend.build- Build backend binary"
	@echo ""
	@echo "Frontend Commands:"
	@echo "  make frontend.dev - Run frontend in development mode"
	@echo "  make frontend.build- Build frontend for production"
	@echo "  make frontend.test- Run frontend tests"
	@echo "  make frontend.lint- Run frontend linting (ESLint)"
	@echo ""
	@echo "CI/CD Commands (mirrors GitHub Actions):"
	@echo "  make ci           - Run full CI pipeline locally (lint → test → build)"
	@echo "  make lint         - Run all linters (backend + frontend)"
	@echo "  make test         - Run all tests with coverage"
	@echo "  make build        - Build all artifacts with git SHA tagging"
	@echo "  make docker.build - Build Docker images with git SHA tags"
	@echo "  make docker.tag   - Tag images for release (requires GIT_TAG)"
	@echo "  make clean        - Clean build artifacts and dependencies"
	@echo "  make install      - Install all dependencies"
	@echo ""
	@echo "Database Commands:"
	@echo "  make db.migrate   - Run database migrations"
	@echo "  make db.rollback  - Rollback last migration"
	@echo "  make db.reset     - Reset database (drop and recreate)"
	@echo "  make db.connect   - Connect to PostgreSQL with psql"

## up: Start all services with Docker Compose
up:
	@echo "Starting all services..."
	$(DOCKER_COMPOSE) up -d
	@echo "Services started!"
	@echo "Backend:  http://localhost:8080"
	@echo "Frontend: http://localhost:3000"
	@echo "Database: localhost:5432"

## down: Stop all services
down:
	@echo "Stopping all services..."
	$(DOCKER_COMPOSE) down

## logs: View logs from all services
logs:
	$(DOCKER_COMPOSE) logs -f

## rebuild: Rebuild and restart all services
rebuild:
	@echo "Rebuilding all services..."
	$(DOCKER_COMPOSE) build --no-cache
	$(DOCKER_COMPOSE) up -d

## backend.dev: Run backend with hot-reload (Air)
backend.dev:
	@echo "Starting backend with hot-reload..."
	cd $(BACKEND_DIR) && air

## backend.run: Run backend without hot-reload
backend.run:
	@echo "Starting backend..."
	cd $(BACKEND_DIR) && $(GO) run main.go

## backend.test: Run backend tests with coverage
backend.test:
	@echo "Running backend tests..."
	cd $(BACKEND_DIR) && $(GO) test -v -race -coverprofile=coverage.out ./...
	@echo "Coverage report:"
	cd $(BACKEND_DIR) && $(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: backend/coverage.html"

## backend.lint: Run backend linting (golangci-lint)
backend.lint:
	@echo "Running backend linting..."
	@echo "Using config: .golangci.yml"
	cd $(BACKEND_DIR) && golangci-lint run --config ../.golangci.yml ./...

## backend.build: Build backend binary with git SHA tagging
backend.build:
	@echo "Building backend..."
	@echo "Git SHA: $(GIT_SHA)"
	@echo "Build Date: $(BUILD_DATE)"
	@mkdir -p bin
	cd $(BACKEND_DIR) && $(GO) build \
		-ldflags "-X main.Version=$(IMAGE_TAG) -X main.GitSHA=$(GIT_SHA) -X main.BuildDate=$(BUILD_DATE)" \
		-o ../bin/backend main.go
	@echo "Backend binary built: bin/backend (version: $(IMAGE_TAG))"

## frontend.dev: Run frontend in development mode
frontend.dev:
	@echo "Starting frontend in development mode..."
	cd $(FRONTEND_DIR) && $(BUN) dev

## frontend.build: Build frontend for production
frontend.build:
	@echo "Building frontend for production..."
	cd $(FRONTEND_DIR) && $(BUN) run build
	@echo "Frontend built: frontend/.next/"

## frontend.test: Run frontend tests
frontend.test:
	@echo "Running frontend tests..."
	cd $(FRONTEND_DIR) && $(BUN) test

## frontend.lint: Run frontend linting (ESLint)
frontend.lint:
	@echo "Running frontend linting..."
	cd $(FRONTEND_DIR) && $(BUN) run lint

## lint: Run all linters (backend + frontend)
lint: backend.lint frontend.lint

## test: Run all tests (backend + frontend)
test: backend.test frontend.test

## build: Build all artifacts
build: backend.build frontend.build

## clean: Clean build artifacts and dependencies
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf $(BACKEND_DIR)/coverage.out $(BACKEND_DIR)/coverage.html
	rm -rf $(FRONTEND_DIR)/.next/
	rm -rf $(FRONTEND_DIR)/node_modules/
	rm -rf $(BACKEND_DIR)/tmp/
	@echo "Clean complete!"

## install: Install all dependencies
install:
	@echo "Installing backend dependencies..."
	cd $(BACKEND_DIR) && $(GO) mod download
	@echo "Installing frontend dependencies..."
	cd $(FRONTEND_DIR) && $(BUN) install
	@echo "Installing development tools..."
	$(GO) install github.com/air-verse/air@latest
	@echo "Dependencies installed!"

## db.migrate: Run database migrations
db.migrate:
	@echo "Running database migrations..."
	$(DOCKER_COMPOSE) exec -T db psql -U postgres -d cicd_training -f /docker-entrypoint-initdb.d/000001_init.up.sql || true

## db.rollback: Rollback last migration
db.rollback:
	@echo "Rolling back last migration..."
	$(DOCKER_COMPOSE) exec -T db psql -U postgres -d cicd_training -f /docker-entrypoint-initdb.d/000001_init.down.sql || true

## db.reset: Reset database (drop and recreate)
db.reset:
	@echo "Resetting database..."
	$(DOCKER_COMPOSE) down -v
	$(DOCKER_COMPOSE) up -d db
	sleep 5
	$(MAKE) db.migrate

## db.connect: Connect to PostgreSQL with psql
db.connect:
	@echo "Connecting to database..."
	$(DOCKER_COMPOSE) exec db psql -U postgres -d cicd_training

# =============================================================================
# CI/CD Pipeline Commands (mirrors GitHub Actions)
# =============================================================================

## ci: Run the full CI pipeline locally (same as GitHub Actions)
ci: ci.info lint test build
	@echo ""
	@echo "=============================================="
	@echo "CI Pipeline Complete!"
	@echo "=============================================="
	@echo "All checks passed successfully."
	@echo ""
	@echo "Artifacts:"
	@echo "  - Backend binary: bin/backend"
	@echo "  - Frontend build: frontend/.next/"
	@echo "  - Coverage report: backend/coverage.html"
	@echo ""

## ci.info: Display CI build information
ci.info:
	@echo "=============================================="
	@echo "CI/CD Pipeline - Local Execution"
	@echo "=============================================="
	@echo "Git SHA:     $(GIT_SHA)"
	@echo "Git Branch:  $(GIT_BRANCH)"
	@echo "Git Tag:     $(if $(GIT_TAG),$(GIT_TAG),none)"
	@echo "Image Tag:   $(IMAGE_TAG)"
	@echo "Build Date:  $(BUILD_DATE)"
	@echo "=============================================="
	@echo ""

## docker.build: Build Docker images with git SHA tagging
docker.build:
	@echo "Building Docker images..."
	@echo "Backend image: $(BACKEND_IMAGE):$(IMAGE_TAG)"
	@echo "Frontend image: $(FRONTEND_IMAGE):$(IMAGE_TAG)"
	docker build -t $(BACKEND_IMAGE):$(IMAGE_TAG) \
		--build-arg GIT_SHA=$(GIT_SHA) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-f $(DOCKER_DIR)/backend.Dockerfile .
	docker build -t $(FRONTEND_IMAGE):$(IMAGE_TAG) \
		-f $(DOCKER_DIR)/frontend.Dockerfile .
	@echo ""
	@echo "Docker images built:"
	@echo "  $(BACKEND_IMAGE):$(IMAGE_TAG)"
	@echo "  $(FRONTEND_IMAGE):$(IMAGE_TAG)"

## docker.tag: Tag Docker images for release (uses git tag or 'latest')
docker.tag:
	@echo "Tagging Docker images..."
	docker tag $(BACKEND_IMAGE):$(IMAGE_TAG) $(BACKEND_IMAGE):latest
	docker tag $(FRONTEND_IMAGE):$(IMAGE_TAG) $(FRONTEND_IMAGE):latest
ifdef GIT_TAG
	docker tag $(BACKEND_IMAGE):$(IMAGE_TAG) $(BACKEND_IMAGE):$(GIT_TAG)
	docker tag $(FRONTEND_IMAGE):$(IMAGE_TAG) $(FRONTEND_IMAGE):$(GIT_TAG)
	@echo "Tagged with version: $(GIT_TAG)"
endif
	@echo "Tagged with: latest"

## docker.push: Push Docker images to registry (requires authentication)
docker.push:
	@echo "Pushing Docker images to $(DOCKER_REGISTRY)..."
	docker push $(BACKEND_IMAGE):$(IMAGE_TAG)
	docker push $(FRONTEND_IMAGE):$(IMAGE_TAG)
	docker push $(BACKEND_IMAGE):latest
	docker push $(FRONTEND_IMAGE):latest
ifdef GIT_TAG
	docker push $(BACKEND_IMAGE):$(GIT_TAG)
	docker push $(FRONTEND_IMAGE):$(GIT_TAG)
endif
	@echo "Images pushed successfully!"

## coverage.report: Generate and display coverage report
coverage.report:
	@echo "Generating coverage report..."
	cd $(BACKEND_DIR) && $(GO) tool cover -func=coverage.out
	@echo ""
	@echo "HTML report available at: backend/coverage.html"

## lint.fix: Run linters with auto-fix where possible
lint.fix:
	@echo "Running linters with auto-fix..."
	cd $(BACKEND_DIR) && golangci-lint run --config ../.golangci.yml --fix ./...
	cd $(FRONTEND_DIR) && $(BUN) run lint --fix || true
	@echo "Lint fixes applied!"
