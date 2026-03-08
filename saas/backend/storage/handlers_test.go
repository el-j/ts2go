package storage

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/el-j/ts2go/saas/backend/projects"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestDownloadFile_OwnershipCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storageClient := &Client{} // Not actually tested because auth fails first
	repo := NewRepository(db)
	projectsRepo := projects.NewRepository(db)
	h := NewHandler(storageClient, repo, projectsRepo)

	fileID := uuid.New()
	projectID := uuid.New()
	callerID := uuid.New()

	// 1. Mock GetFile
	rows := sqlmock.NewRows([]string{
		"id", "project_id", "transpilation_id", "file_type",
		"original_name", "storage_path", "mime_type", "size_bytes",
		"checksum", "created_at",
	}).AddRow(
		fileID, projectID, nil, "input",
		"test.ts", "path/to/test.ts", "text/plain", 100,
		"checksum", time.Now(),
	)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, project_id, transpilation_id, file_type, original_name, storage_path, mime_type, size_bytes, checksum, created_at FROM files WHERE id = $1")).
		WithArgs(fileID).
		WillReturnRows(rows)

	// 2. Mock CheckAccess returning 0 (no access)
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*)")).
		WithArgs(projectID, callerID).
		WillReturnRows(countRows)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/api/v1/files/"+fileID.String()+"/download", nil)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: fileID.String()}}
	c.Set("user_id", callerID.String())

	h.DownloadFile(c)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected 403 Forbidden, got %d", w.Code)
	}
}

func TestGetPresignedURL_OwnershipCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storageClient := &Client{}
	repo := NewRepository(db)
	projectsRepo := projects.NewRepository(db)
	h := NewHandler(storageClient, repo, projectsRepo)

	fileID := uuid.New()
	projectID := uuid.New()
	callerID := uuid.New()

	// 1. Mock GetFile
	rows := sqlmock.NewRows([]string{
		"id", "project_id", "transpilation_id", "file_type",
		"original_name", "storage_path", "mime_type", "size_bytes",
		"checksum", "created_at",
	}).AddRow(
		fileID, projectID, nil, "input",
		"test.ts", "path/to/test.ts", "text/plain", 100,
		"checksum", time.Now(),
	)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, project_id, transpilation_id, file_type, original_name, storage_path, mime_type, size_bytes, checksum, created_at FROM files WHERE id = $1")).
		WithArgs(fileID).
		WillReturnRows(rows)

	// 2. Mock CheckAccess returning 0 (no access)
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*)")).
		WithArgs(projectID, callerID).
		WillReturnRows(countRows)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("GET", "/api/v1/files/"+fileID.String()+"/url", nil)
	c.Request = req
	c.Params = []gin.Param{{Key: "id", Value: fileID.String()}}
	c.Set("user_id", callerID.String())

	h.GetPresignedURL(c)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected 403 Forbidden, got %d", w.Code)
	}
}
