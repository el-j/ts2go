package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/el-j/ts2go/saas/backend/auth"
	"github.com/gin-gonic/gin"
)

func TestRequireAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// We don't need the repository for RequireAuth as it only uses ValidateToken on the service
	svc := auth.NewService("secret", time.Hour, time.Hour)
	m := NewMiddleware(svc, nil)

	tests := []struct {
		name       string
		authHeader string
		wantStatus int
	}{
		{
			name:       "Missing header",
			authHeader: "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Invalid format",
			authHeader: "BearerTokenWithoutSpace",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Invalid token",
			authHeader: "Bearer invalid.token.string",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := http.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			c.Request = req

			m.RequireAuth()(c)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}
