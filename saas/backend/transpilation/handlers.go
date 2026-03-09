package transpilation

import (
	"net/http"

	"github.com/el-j/ts2go/saas/backend/errors"
	"github.com/el-j/ts2go/saas/backend/queue"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles transpilation requests
type Handler struct {
	queue *queue.Queue
}

// NewHandler creates a new transpilation handler
func NewHandler(queue *queue.Queue) *Handler {
	return &Handler{
		queue: queue,
	}
}

// TranspileRequest represents a transpilation request
type TranspileRequest struct {
	ProjectID uuid.UUID              `json:"project_id" binding:"required"`
	FileIDs   []string               `json:"file_ids" binding:"required,min=1"`
	Priority  int                    `json:"priority"`
	Settings  map[string]interface{} `json:"settings,omitempty"`
}

// TranspileResponse represents the transpilation response
type TranspileResponse struct {
	JobID     string `json:"job_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	ProjectID string `json:"project_id"`
}

// Transpile creates a new transpilation job
// @Summary Create transpilation job
// @Description Submits files for TypeScript to Go transpilation
// @Tags transpilation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TranspileRequest true "Transpilation request"
// @Success 202 {object} TranspileResponse
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/v1/transpile [post]
func (h *Handler) Transpile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Error(errors.ErrUnauthorized)
		return
	}

	var req TranspileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.ErrValidation.WithDetails(err.Error()))
		return
	}

	// Validate priority
	priority := queue.PriorityNormal
	if req.Priority >= int(queue.PriorityLow) && req.Priority <= int(queue.PriorityUrgent) {
		priority = queue.JobPriority(req.Priority)
	}

	// Parse file IDs to file metadata (simplified - in production, fetch actual metadata)
	inputFiles := make([]queue.FileMetadata, len(req.FileIDs))
	for i, fileID := range req.FileIDs {
		inputFiles[i] = queue.FileMetadata{
			Name: fileID, // Placeholder - should fetch actual file name
			Path: fileID,
		}
	}

	// Create job
	userUUID, _ := uuid.Parse(userID.(string))
	job := &queue.TranspilationJob{
		ID:         uuid.New(),
		UserID:     userUUID,
		ProjectID:  &req.ProjectID,
		InputFiles: inputFiles,
		Settings:   req.Settings,
		Priority:   priority,
	}

	// Enqueue job
	if err := h.queue.Enqueue(c.Request.Context(), job); err != nil {
		c.Error(errors.ErrInternal.WithDetails("Failed to enqueue job"))
		return
	}

	c.JSON(http.StatusAccepted, TranspileResponse{
		JobID:     job.ID.String(),
		Status:    "queued",
		Message:   "Transpilation job created successfully",
		ProjectID: req.ProjectID.String(),
	})
}

// GetJobStatus retrieves the status of a transpilation job
// @Summary Get job status
// @Description Get the current status of a transpilation job
// @Tags transpilation
// @Produce json
// @Security BearerAuth
// @Param job_id path string true "Job ID"
// @Success 200 {object} queue.TranspilationJob
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/v1/transpile/{job_id} [get]
func (h *Handler) GetJobStatus(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.Error(errors.ErrUnauthorized)
		return
	}

	jobID := c.Param("job_id")
	if jobID == "" {
		c.Error(errors.ErrValidation.WithDetails("Job ID is required"))
		return
	}

	jobUUID, err := uuid.Parse(jobID)
	if err != nil {
		c.Error(errors.ErrValidation.WithDetails("Invalid job ID"))
		return
	}

	job, err := h.queue.GetJob(c.Request.Context(), jobUUID)
	if err != nil {
		c.Error(errors.ErrNotFound.WithDetails("Job not found"))
		return
	}

	userID, _ := c.Get("user_id")
	userUUID, _ := uuid.Parse(userID.(string))
	if job.UserID != userUUID {
		c.Error(errors.ErrForbidden.WithDetails("You do not have access to this job"))
		return
	}

	c.JSON(http.StatusOK, job)
}

// ListUserJobs lists all jobs for the current user
// @Summary List user jobs
// @Description Get all transpilation jobs for the authenticated user
// @Tags transpilation
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(50)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} errors.ErrorResponse
// @Router /api/v1/transpile [get]
func (h *Handler) ListUserJobs(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Error(errors.ErrUnauthorized)
		return
	}

	userUUID, _ := uuid.Parse(userID.(string))

	// Get limit from query param
	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := uuid.Parse(limitStr); err == nil {
			limit = int(l.ID())
		}
	}

	jobs, err := h.queue.ListUserJobs(c.Request.Context(), userUUID, limit)
	if err != nil {
		c.Error(errors.ErrInternal.WithDetails("Failed to list jobs"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs":  jobs,
		"total": len(jobs),
	})
}

// CancelJob cancels a pending or processing job
// @Summary Cancel job
// @Description Cancel a transpilation job
// @Tags transpilation
// @Security BearerAuth
// @Param job_id path string true "Job ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} errors.ErrorResponse
// @Router /api/v1/transpile/{job_id}/cancel [post]
func (h *Handler) CancelJob(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.Error(errors.ErrUnauthorized)
		return
	}

	jobID := c.Param("job_id")
	if jobID == "" {
		c.Error(errors.ErrValidation.WithDetails("Job ID is required"))
		return
	}

	jobUUID, err := uuid.Parse(jobID)
	if err != nil {
		c.Error(errors.ErrValidation.WithDetails("Invalid job ID"))
		return
	}

	// Get job to verify ownership
	job, err := h.queue.GetJob(c.Request.Context(), jobUUID)
	if err != nil {
		c.Error(errors.ErrNotFound.WithDetails("Job not found"))
		return
	}

	userID, _ := c.Get("user_id")
	userUUID, _ := uuid.Parse(userID.(string))
	if job.UserID != userUUID {
		c.Error(errors.ErrForbidden.WithDetails("You do not have access to this job"))
		return
	}

	// Only cancel if pending or processing
	if job.Status != "pending" && job.Status != "processing" {
		c.Error(errors.ErrBadRequest.WithDetails("Job cannot be cancelled in current state"))
		return
	}

	// Mark as failed with cancellation message
	err = h.queue.FailJob(c.Request.Context(), jobUUID, "Cancelled by user")
	if err != nil {
		c.Error(errors.ErrInternal.WithDetails("Failed to cancel job"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Job cancelled successfully",
		"job_id":  jobID,
	})
}
