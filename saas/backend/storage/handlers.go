package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/el-j/ts2go/saas/backend/projects"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	MaxFileSize = 100 * 1024 * 1024 // 100 MB
	MaxFiles    = 50                // Max files per upload
)

var allowedExtensions = map[string]bool{
	".ts":   true,
	".tsx":  true,
	".js":   true,
	".jsx":  true,
	".json": true,
	".txt":  true,
	".md":   true,
}

// Handler handles file storage endpoints
type Handler struct {
	storage      *Client
	repo         *Repository
	projectsRepo *projects.Repository
}

// NewHandler creates a new storage handler
func NewHandler(storage *Client, repo *Repository, projectsRepo *projects.Repository) *Handler {
	return &Handler{
		storage:      storage,
		repo:         repo,
		projectsRepo: projectsRepo,
	}
}

// UploadFile handles file upload
// @Summary Upload files to a project
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param project_id formData string true "Project ID"
// @Param files formData file true "Files to upload"
// @Success 200 {object} UploadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/files/upload [post]
func (h *Handler) UploadFile(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	projectIDStr := c.PostForm("project_id")
	if projectIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project_id"})
		return
	}

	// Get multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files provided"})
		return
	}

	if len(files) > MaxFiles {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Too many files (max %d)", MaxFiles)})
		return
	}

	var uploadedFiles []UploadedFile
	var errors []string

	for _, fileHeader := range files {
		// Validate file size
		if fileHeader.Size > MaxFileSize {
			errors = append(errors, fmt.Sprintf("%s: file too large (max %d MB)", fileHeader.Filename, MaxFileSize/(1024*1024)))
			continue
		}

		// Validate file extension
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if !allowedExtensions[ext] {
			errors = append(errors, fmt.Sprintf("%s: unsupported file type", fileHeader.Filename))
			continue
		}

		// Open file
		file, err := fileHeader.Open()
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: failed to open file", fileHeader.Filename))
			continue
		}

		// Calculate checksum
		hash := sha256.New()
		if _, err := io.Copy(hash, file); err != nil {
			file.Close()
			errors = append(errors, fmt.Sprintf("%s: failed to calculate checksum", fileHeader.Filename))
			continue
		}
		checksum := hex.EncodeToString(hash.Sum(nil))

		// Reset file pointer
		file.Seek(0, 0)

		// Generate storage path
		storagePath := fmt.Sprintf("projects/%s/%s/%s", projectID, uuid.New(), fileHeader.Filename)

		// Upload to storage
		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		_, err = h.storage.UploadFile(c.Request.Context(), storagePath, file, fileHeader.Size, contentType)
		file.Close()

		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: upload failed", fileHeader.Filename))
			continue
		}

		// Save metadata to database
		fileRecord, err := h.repo.CreateFile(
			c.Request.Context(),
			&projectID,
			nil,
			"input",
			fileHeader.Filename,
			storagePath,
			contentType,
			fileHeader.Size,
			checksum,
		)

		if err != nil {
			// Clean up uploaded file
			h.storage.DeleteFile(c.Request.Context(), storagePath)
			errors = append(errors, fmt.Sprintf("%s: failed to save metadata", fileHeader.Filename))
			continue
		}

		uploadedFiles = append(uploadedFiles, UploadedFile{
			ID:           fileRecord.ID,
			OriginalName: fileHeader.Filename,
			Size:         fileHeader.Size,
			ContentType:  contentType,
			Checksum:     checksum,
			StoragePath:  storagePath,
		})
	}

	response := UploadResponse{
		Success: len(uploadedFiles),
		Failed:  len(errors),
		Files:   uploadedFiles,
		Errors:  errors,
	}

	if len(errors) > 0 && len(uploadedFiles) == 0 {
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// DownloadFile handles file download
// @Summary Download a file
// @Tags files
// @Produce application/octet-stream
// @Security BearerAuth
// @Param id path string true "File ID"
// @Success 200 {file} binary
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/files/{id}/download [get]
func (h *Handler) DownloadFile(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	// Get file metadata
	fileRecord, err := h.repo.GetFile(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	userID, _ := c.Get("user_id")
	userUUID, _ := uuid.Parse(userID.(string))

	hasAccess, err := h.projectsRepo.CheckAccess(c.Request.Context(), fileRecord.ProjectID, userUUID)
	if err != nil || !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to project files"})
		return
	}

	// Get file from storage
	reader, err := h.storage.DownloadFile(c.Request.Context(), fileRecord.StoragePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file"})
		return
	}
	defer reader.Close()

	// Set headers
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileRecord.OriginalName))
	c.Header("Content-Type", *fileRecord.MimeType)

	// Stream file
	io.Copy(c.Writer, reader)
}

// GetPresignedURL generates a presigned URL for download
// @Summary Get presigned download URL
// @Tags files
// @Produce json
// @Security BearerAuth
// @Param id path string true "File ID"
// @Param expiry query int false "Expiry in seconds (default 3600)"
// @Success 200 {object} PresignedURLResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/files/{id}/url [get]
func (h *Handler) GetPresignedURL(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	// Get file metadata
	fileRecord, err := h.repo.GetFile(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	userID, _ := c.Get("user_id")
	userUUID, _ := uuid.Parse(userID.(string))

	hasAccess, err := h.projectsRepo.CheckAccess(c.Request.Context(), fileRecord.ProjectID, userUUID)
	if err != nil || !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to project files"})
		return
	}

	// Generate presigned URL (default 1 hour)
	url, err := h.storage.GetPresignedURL(c.Request.Context(), fileRecord.StoragePath, 3600)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate URL"})
		return
	}

	c.JSON(http.StatusOK, PresignedURLResponse{
		URL:       url,
		ExpiresIn: 3600,
	})
}

// DeleteFile deletes a file
// @Summary Delete a file
// @Tags files
// @Security BearerAuth
// @Param id path string true "File ID"
// @Success 204
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/files/{id} [delete]
func (h *Handler) DeleteFile(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	// Get file metadata
	fileRecord, err := h.repo.GetFile(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	userID, _ := c.Get("user_id")
	userUUID, _ := uuid.Parse(userID.(string))

	hasAccess, err := h.projectsRepo.CheckAccess(c.Request.Context(), fileRecord.ProjectID, userUUID)
	if err != nil || !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to project files"})
		return
	}

	// Delete from storage
	if err := h.storage.DeleteFile(c.Request.Context(), fileRecord.StoragePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from storage"})
		return
	}

	// Delete metadata
	if err := h.repo.DeleteFile(c.Request.Context(), fileID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file metadata"})
		return
	}

	c.Status(http.StatusNoContent)
}

// Response types
type UploadResponse struct {
	Success int            `json:"success"`
	Failed  int            `json:"failed"`
	Files   []UploadedFile `json:"files"`
	Errors  []string       `json:"errors,omitempty"`
}

type UploadedFile struct {
	ID           uuid.UUID `json:"id"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	Checksum     string    `json:"checksum"`
	StoragePath  string    `json:"-"`
}

type PresignedURLResponse struct {
	URL       string `json:"url"`
	ExpiresIn int    `json:"expires_in"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
