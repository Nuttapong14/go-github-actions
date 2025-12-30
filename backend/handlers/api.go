package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/Nuttapong14/go-github-actions/backend/models"
)

// Repository interface for data access
type Repository interface {
	// Environment operations
	ListEnvironments() ([]models.Environment, error)
	GetEnvironment(id uuid.UUID) (*models.Environment, error)
	GetEnvironmentByName(name string) (*models.Environment, error)

	// Release operations
	ListReleases() ([]models.Release, error)
	GetRelease(id uuid.UUID) (*models.Release, error)
	CreateRelease(release *models.Release) error

	// Deployment operations
	ListDeployments(environmentID *uuid.UUID) ([]models.Deployment, error)
	GetDeployment(id uuid.UUID) (*models.Deployment, error)
	CreateDeployment(deployment *models.Deployment) error
	UpdateDeployment(deployment *models.Deployment) error

	// Health check operations
	CreateHealthCheck(check *models.HealthCheck) error
	GetLatestHealthChecks(deploymentID uuid.UUID) ([]models.HealthCheck, error)
}

// APIHandler holds the repository and provides HTTP handlers
type APIHandler struct {
	repo Repository
}

// NewAPIHandler creates a new API handler with the given repository
func NewAPIHandler(repo Repository) *APIHandler {
	return &APIHandler{repo: repo}
}

// APIResponse wraps API responses in a consistent format
type APIResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error *APIError   `json:"error,omitempty"`
}

// APIError represents an error response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ListResponse wraps list responses with pagination metadata
type ListResponse struct {
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination contains pagination metadata
type Pagination struct {
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

// CreateReleaseRequest represents the request body for creating a release
type CreateReleaseRequest struct {
	Version             string  `json:"version" validate:"required"`
	GitSHA              string  `json:"git_sha" validate:"required,len=40"`
	Changelog           *string `json:"changelog,omitempty"`
	DockerImageBackend  string  `json:"docker_image_backend" validate:"required"`
	DockerImageFrontend string  `json:"docker_image_frontend" validate:"required"`
	CreatedBy           string  `json:"created_by" validate:"required"`
}

// CreateDeploymentRequest represents the request body for creating a deployment
type CreateDeploymentRequest struct {
	EnvironmentID string `json:"environment_id" validate:"required,uuid"`
	ReleaseID     string `json:"release_id" validate:"required,uuid"`
	TriggeredBy   string `json:"triggered_by" validate:"required"`
	TriggerType   string `json:"trigger_type" validate:"required,oneof=manual automatic hotfix rollback"`
}

// RollbackRequest represents the request body for rolling back a deployment
type RollbackRequest struct {
	TriggeredBy string `json:"triggered_by" validate:"required"`
}

// EnvironmentResponse represents an environment in API responses
type EnvironmentResponse struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	URL              string    `json:"url"`
	ApprovalRequired bool      `json:"approval_required"`
	IsProduction     bool      `json:"is_production"`
	SoakPeriodHours  int       `json:"soak_period_hours"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ReleaseResponse represents a release in API responses
type ReleaseResponse struct {
	ID                  string    `json:"id"`
	Version             string    `json:"version"`
	GitSHA              string    `json:"git_sha"`
	GitSHAShort         string    `json:"git_sha_short"`
	Changelog           *string   `json:"changelog,omitempty"`
	DockerImageBackend  string    `json:"docker_image_backend"`
	DockerImageFrontend string    `json:"docker_image_frontend"`
	CreatedBy           string    `json:"created_by"`
	CreatedAt           time.Time `json:"created_at"`
}

// DeploymentResponse represents a deployment in API responses
type DeploymentResponse struct {
	ID            string                `json:"id"`
	EnvironmentID string                `json:"environment_id"`
	ReleaseID     string                `json:"release_id"`
	Status        string                `json:"status"`
	StartedAt     *time.Time            `json:"started_at,omitempty"`
	CompletedAt   *time.Time            `json:"completed_at,omitempty"`
	TriggeredBy   string                `json:"triggered_by"`
	TriggerType   string                `json:"trigger_type"`
	RollbackOfID  *string               `json:"rollback_of_id,omitempty"`
	ErrorMessage  *string               `json:"error_message,omitempty"`
	CreatedAt     time.Time             `json:"created_at"`
	Environment   *EnvironmentResponse  `json:"environment,omitempty"`
	Release       *ReleaseResponse      `json:"release,omitempty"`
	HealthChecks  []HealthCheckResponse `json:"health_checks,omitempty"`
}

// HealthCheckResponse represents a health check in API responses
type HealthCheckResponse struct {
	ID           string    `json:"id"`
	DeploymentID string    `json:"deployment_id"`
	CheckType    string    `json:"check_type"`
	Status       string    `json:"status"`
	ResponseTime int       `json:"response_time_ms"`
	StatusCode   *int      `json:"status_code,omitempty"`
	ErrorMessage *string   `json:"error_message,omitempty"`
	CheckedAt    time.Time `json:"checked_at"`
}

// --- Environment Handlers ---

// ListEnvironments handles GET /api/v1/environments
func (h *APIHandler) ListEnvironments(c *fiber.Ctx) error {
	environments, err := h.repo.ListEnvironments()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Error: &APIError{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to retrieve environments",
			},
		})
	}

	response := make([]EnvironmentResponse, len(environments))
	for i, env := range environments {
		response[i] = toEnvironmentResponse(&env)
	}

	return c.JSON(ListResponse{
		Data: response,
		Pagination: &Pagination{
			Total:   len(response),
			Page:    1,
			PerPage: len(response),
		},
	})
}

// GetEnvironment handles GET /api/v1/environments/:id
func (h *APIHandler) GetEnvironment(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid environment ID format",
			},
		})
	}

	env, err := h.repo.GetEnvironment(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(APIResponse{
				Error: &APIError{
					Code:    "NOT_FOUND",
					Message: "Environment not found",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Error: &APIError{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to retrieve environment",
			},
		})
	}

	return c.JSON(APIResponse{
		Data: toEnvironmentResponse(env),
	})
}

// --- Release Handlers ---

// ListReleases handles GET /api/v1/releases
func (h *APIHandler) ListReleases(c *fiber.Ctx) error {
	releases, err := h.repo.ListReleases()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Error: &APIError{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to retrieve releases",
			},
		})
	}

	response := make([]ReleaseResponse, len(releases))
	for i, rel := range releases {
		response[i] = toReleaseResponse(&rel)
	}

	return c.JSON(ListResponse{
		Data: response,
		Pagination: &Pagination{
			Total:   len(response),
			Page:    1,
			PerPage: len(response),
		},
	})
}

// GetRelease handles GET /api/v1/releases/:id
func (h *APIHandler) GetRelease(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid release ID format",
			},
		})
	}

	release, err := h.repo.GetRelease(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(APIResponse{
				Error: &APIError{
					Code:    "NOT_FOUND",
					Message: "Release not found",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Error: &APIError{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to retrieve release",
			},
		})
	}

	return c.JSON(APIResponse{
		Data: toReleaseResponse(release),
	})
}

// CreateRelease handles POST /api/v1/releases
func (h *APIHandler) CreateRelease(c *fiber.Ctx) error {
	var req CreateReleaseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Error: &APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request body",
			},
		})
	}

	// Validate required fields
	if req.Version == "" || req.GitSHA == "" || req.DockerImageBackend == "" ||
		req.DockerImageFrontend == "" || req.CreatedBy == "" {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Error: &APIError{
				Code:    "VALIDATION_ERROR",
				Message: "Missing required fields",
			},
		})
	}

	if len(req.GitSHA) != 40 {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Error: &APIError{
				Code:    "VALIDATION_ERROR",
				Message: "git_sha must be 40 characters",
			},
		})
	}

	release := &models.Release{
		ID:                  uuid.New(),
		Version:             req.Version,
		GitSHA:              req.GitSHA,
		GitSHAShort:         req.GitSHA[:7],
		Changelog:           req.Changelog,
		DockerImageBackend:  req.DockerImageBackend,
		DockerImageFrontend: req.DockerImageFrontend,
		CreatedBy:           req.CreatedBy,
		CreatedAt:           time.Now().UTC(),
	}

	if err := h.repo.CreateRelease(release); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Error: &APIError{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to create release",
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Data: toReleaseResponse(release),
	})
}

// --- Deployment Handlers ---

// ListDeployments handles GET /api/v1/deployments
func (h *APIHandler) ListDeployments(c *fiber.Ctx) error {
	var envID *uuid.UUID

	// Parse optional environment_id query parameter
	envIDStr := c.Query("environment_id")
	if envIDStr != "" {
		id, err := uuid.Parse(envIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
				Error: &APIError{
					Code:    "INVALID_ID",
					Message: "Invalid environment_id format",
				},
			})
		}
		envID = &id
	}

	deployments, err := h.repo.ListDeployments(envID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Error: &APIError{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to retrieve deployments",
			},
		})
	}

	response := make([]DeploymentResponse, len(deployments))
	for i, dep := range deployments {
		response[i] = toDeploymentResponse(&dep, false)
	}

	return c.JSON(ListResponse{
		Data: response,
		Pagination: &Pagination{
			Total:   len(response),
			Page:    1,
			PerPage: len(response),
		},
	})
}

// GetDeployment handles GET /api/v1/deployments/:id
func (h *APIHandler) GetDeployment(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid deployment ID format",
			},
		})
	}

	deployment, err := h.repo.GetDeployment(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(APIResponse{
				Error: &APIError{
					Code:    "NOT_FOUND",
					Message: "Deployment not found",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Error: &APIError{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to retrieve deployment",
			},
		})
	}

	// Include related data for single deployment
	return c.JSON(APIResponse{
		Data: toDeploymentResponse(deployment, true),
	})
}

// CreateDeployment handles POST /api/v1/deployments
func (h *APIHandler) CreateDeployment(c *fiber.Ctx) error {
	var req CreateDeploymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Error: &APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request body",
			},
		})
	}

	// Parse UUIDs
	envID, err := uuid.Parse(req.EnvironmentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid environment_id format",
			},
		})
	}

	releaseID, err := uuid.Parse(req.ReleaseID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid release_id format",
			},
		})
	}

	// Validate trigger type
	triggerType := models.TriggerType(req.TriggerType)
	if !isValidTriggerType(triggerType) {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Error: &APIError{
				Code:    "VALIDATION_ERROR",
				Message: "Invalid trigger_type",
			},
		})
	}

	now := time.Now().UTC()
	deployment := &models.Deployment{
		ID:            uuid.New(),
		EnvironmentID: envID,
		ReleaseID:     releaseID,
		Status:        models.StatusPending,
		StartedAt:     &now,
		TriggeredBy:   req.TriggeredBy,
		TriggerType:   triggerType,
		CreatedAt:     now,
	}

	if err := h.repo.CreateDeployment(deployment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Error: &APIError{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to create deployment",
			},
		})
	}

	// Record deployment metric
	RecordDeployment("unknown", string(deployment.Status))

	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Data: toDeploymentResponse(deployment, false),
	})
}

// RollbackDeployment handles POST /api/v1/deployments/:id/rollback
func (h *APIHandler) RollbackDeployment(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid deployment ID format",
			},
		})
	}

	// Get the deployment to rollback
	deployment, err := h.repo.GetDeployment(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(APIResponse{
				Error: &APIError{
					Code:    "NOT_FOUND",
					Message: "Deployment not found",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Error: &APIError{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to retrieve deployment",
			},
		})
	}

	var req RollbackRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
			Error: &APIError{
				Code:    "INVALID_REQUEST",
				Message: "Invalid request body",
			},
		})
	}

	// Create a new deployment that is a rollback of the original
	now := time.Now().UTC()
	rollbackDeployment := &models.Deployment{
		ID:            uuid.New(),
		EnvironmentID: deployment.EnvironmentID,
		ReleaseID:     deployment.ReleaseID,
		Status:        models.StatusPending,
		StartedAt:     &now,
		TriggeredBy:   req.TriggeredBy,
		TriggerType:   models.TriggerRollback,
		RollbackOfID:  &id,
		CreatedAt:     now,
	}

	if err := h.repo.CreateDeployment(rollbackDeployment); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
			Error: &APIError{
				Code:    "INTERNAL_ERROR",
				Message: "Failed to create rollback deployment",
			},
		})
	}

	// Record deployment metric
	RecordDeployment("unknown", "rollback")

	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Data: toDeploymentResponse(rollbackDeployment, false),
	})
}

// --- Helper Functions ---

func toEnvironmentResponse(env *models.Environment) EnvironmentResponse {
	return EnvironmentResponse{
		ID:               env.ID.String(),
		Name:             env.Name,
		URL:              env.URL,
		ApprovalRequired: env.ApprovalRequired,
		IsProduction:     env.IsProduction,
		SoakPeriodHours:  env.SoakPeriodHours,
		CreatedAt:        env.CreatedAt,
		UpdatedAt:        env.UpdatedAt,
	}
}

func toReleaseResponse(rel *models.Release) ReleaseResponse {
	return ReleaseResponse{
		ID:                  rel.ID.String(),
		Version:             rel.Version,
		GitSHA:              rel.GitSHA,
		GitSHAShort:         rel.GitSHAShort,
		Changelog:           rel.Changelog,
		DockerImageBackend:  rel.DockerImageBackend,
		DockerImageFrontend: rel.DockerImageFrontend,
		CreatedBy:           rel.CreatedBy,
		CreatedAt:           rel.CreatedAt,
	}
}

func toDeploymentResponse(dep *models.Deployment, includeRelated bool) DeploymentResponse {
	resp := DeploymentResponse{
		ID:            dep.ID.String(),
		EnvironmentID: dep.EnvironmentID.String(),
		ReleaseID:     dep.ReleaseID.String(),
		Status:        string(dep.Status),
		StartedAt:     dep.StartedAt,
		CompletedAt:   dep.CompletedAt,
		TriggeredBy:   dep.TriggeredBy,
		TriggerType:   string(dep.TriggerType),
		ErrorMessage:  dep.ErrorMessage,
		CreatedAt:     dep.CreatedAt,
	}

	if dep.RollbackOfID != nil {
		rollbackStr := dep.RollbackOfID.String()
		resp.RollbackOfID = &rollbackStr
	}

	if includeRelated {
		// Check if Environment was preloaded by checking if ID is not nil
		if dep.Environment.ID != uuid.Nil {
			envResp := toEnvironmentResponse(&dep.Environment)
			resp.Environment = &envResp
		}
		// Check if Release was preloaded by checking if ID is not nil
		if dep.Release.ID != uuid.Nil {
			relResp := toReleaseResponse(&dep.Release)
			resp.Release = &relResp
		}
		// Health checks would be populated here if loaded
	}

	return resp
}

func isValidTriggerType(t models.TriggerType) bool {
	switch t {
	case models.TriggerManual, models.TriggerAutomatic, models.TriggerHotfix, models.TriggerRollback:
		return true
	default:
		return false
	}
}

// RegisterAPIRoutes registers the API routes with the Fiber app
func RegisterAPIRoutes(app *fiber.App, handler *APIHandler) {
	api := app.Group("/api/v1")

	// Environment routes
	api.Get("/environments", handler.ListEnvironments)
	api.Get("/environments/:id", handler.GetEnvironment)

	// Release routes
	api.Get("/releases", handler.ListReleases)
	api.Get("/releases/:id", handler.GetRelease)
	api.Post("/releases", handler.CreateRelease)

	// Deployment routes
	api.Get("/deployments", handler.ListDeployments)
	api.Get("/deployments/:id", handler.GetDeployment)
	api.Post("/deployments", handler.CreateDeployment)
	api.Post("/deployments/:id/rollback", handler.RollbackDeployment)
}
