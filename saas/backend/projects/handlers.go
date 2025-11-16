package projects

import (
	"net/http"
	"strings"

	"github.com/el-j/ts2go/saas/backend/db/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateProjectRequest represents the project creation request
type CreateProjectRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=200"`
	Slug        string  `json:"slug" binding:"required,min=1,max=200"`
	Description string  `json:"description"`
	Visibility  string  `json:"visibility" binding:"required,oneof=private team public"`
	TeamID      *string `json:"team_id"`
}

// UpdateProjectRequest represents the project update request
type UpdateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Visibility  *string `json:"visibility" binding:"omitempty,oneof=private team public"`
}

// Handler handles project endpoints
type Handler struct {
	repo *Repository
}

// NewHandler creates a new project handler
func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// Create creates a new project
// @Summary Create a new project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateProjectRequest true "Project details"
// @Success 201 {object} models.Project
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /api/v1/projects [post]
func (h *Handler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate slug format
	if !isValidSlug(req.Slug) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slug format. Use lowercase letters, numbers, and hyphens only"})
		return
	}

	// Check if slug already exists for this user
	existing, _ := h.repo.GetBySlug(c.Request.Context(), userID.(uuid.UUID), req.Slug)
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Project with this slug already exists"})
		return
	}

	// Parse team ID if provided
	var teamID *uuid.UUID
	if req.TeamID != nil && *req.TeamID != "" {
		tid, err := uuid.Parse(*req.TeamID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
			return
		}
		teamID = &tid
	}

	project, err := h.repo.Create(
		c.Request.Context(),
		userID.(uuid.UUID),
		req.Name,
		req.Slug,
		req.Description,
		req.Visibility,
		teamID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// List lists all projects for the authenticated user
// @Summary List user's projects
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param archived query boolean false "Include archived projects"
// @Success 200 {array} models.Project
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/projects [get]
func (h *Handler) List(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	includeArchived := c.Query("archived") == "true"

	projects, err := h.repo.ListByUser(c.Request.Context(), userID.(uuid.UUID), includeArchived)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list projects"})
		return
	}

	if projects == nil {
		projects = make([]*models.Project, 0)
	}

	c.JSON(http.StatusOK, projects)
}

// Get retrieves a single project
// @Summary Get project by ID
// @Tags projects
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {object} models.Project
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/projects/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Check access
	hasAccess, err := h.repo.CheckAccess(c.Request.Context(), projectID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check access"})
		return
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	project, err := h.repo.GetByID(c.Request.Context(), projectID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get project"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// Update updates a project
// @Summary Update project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Param request body UpdateProjectRequest true "Update details"
// @Success 200 {object} models.Project
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/projects/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check access
	hasAccess, err := h.repo.CheckAccess(c.Request.Context(), projectID, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check access"})
		return
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	err = h.repo.Update(c.Request.Context(), projectID, req.Name, req.Description, req.Visibility)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	// Return updated project
	project, _ := h.repo.GetByID(c.Request.Context(), projectID)
	c.JSON(http.StatusOK, project)
}

// Delete deletes (archives) a project
// @Summary Delete project
// @Tags projects
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 204
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/projects/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	err = h.repo.Delete(c.Request.Context(), projectID, userID.(uuid.UUID))
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "unauthorized") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found or access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.Status(http.StatusNoContent)
}

// isValidSlug checks if a slug is valid
func isValidSlug(slug string) bool {
	if len(slug) == 0 || len(slug) > 200 {
		return false
	}

	for _, c := range slug {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}

	return true
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
