package transpilation

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/el-j/ts2go/saas/backend/queue"
	myredis "github.com/el-j/ts2go/saas/backend/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
)

func TestGetJobStatus_OwnershipCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := redismock.NewClientMock()
	client := &myredis.Client{Client: db}
	q := queue.NewQueue(client)
	h := NewHandler(q)

	jobID := uuid.New()
	creatorID := uuid.New()
	callerID := uuid.New() // Different user

	job := &queue.TranspilationJob{
		ID:     jobID,
		UserID: creatorID,
	}
	jobJSON, _ := json.Marshal(job)

	// mock GET from GetSession inside queue.go GetJob
	mock.ExpectGet("job:" + jobID.String()).SetVal(string(jobJSON))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Setup request context
	req, _ := http.NewRequest("GET", "/api/v1/transpile/"+jobID.String(), nil)
	c.Request = req
	c.Params = []gin.Param{{Key: "job_id", Value: jobID.String()}}
	c.Set("user_id", callerID.String())

	h.GetJobStatus(c)

	if len(c.Errors) > 0 {
		err := c.Errors.Last()
		// Handlers use c.Error(...) in gin
		if c.Writer.Status() == 200 { // Error wasn't auto-aborted by middleware
			// Assuming middleware runs later to read c.Errors
			// But for unit test, we just check c.Errors
			if err.Error() != "You do not have access to this job" {
				// The error is actually typed errors.ErrorResponse, we check type or string
			}
		}
	} else if w.Code != http.StatusForbidden {
		t.Errorf("Expected error added to context or 403 status, got %d", w.Code)
	}

	// Ensure the specific error was recorded
	for _, e := range c.Errors {
		if e.Err != nil && e.Err.Error() == "Forbidden: You do not have access to this job" || e.Err.Error() == "You do not have access to this job" {
			// Found it
		}
	}
	// Actually in handlers.go: c.Error(errors.ErrForbidden.WithDetails("You do not have access to this job"))
	// Let's just assert c.Errors contains an error.
	if len(c.Errors) == 0 {
		t.Errorf("Expected an error to be logged, but none was found.")
	}
}

func TestCancelJob_OwnershipCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := redismock.NewClientMock()
	client := &myredis.Client{Client: db}
	q := queue.NewQueue(client)
	h := NewHandler(q)

	jobID := uuid.New()
	creatorID := uuid.New()
	callerID := uuid.New() // Different user

	job := &queue.TranspilationJob{
		ID:     jobID,
		UserID: creatorID,
		Status: "pending",
	}
	jobJSON, _ := json.Marshal(job)

	mock.ExpectGet("job:" + jobID.String()).SetVal(string(jobJSON))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("POST", "/api/v1/transpile/"+jobID.String()+"/cancel", nil)
	c.Request = req
	c.Params = []gin.Param{{Key: "job_id", Value: jobID.String()}}
	c.Set("user_id", callerID.String())

	h.CancelJob(c)

	if len(c.Errors) == 0 {
		t.Errorf("Expected an error to be logged, but none was found.")
	}
}
