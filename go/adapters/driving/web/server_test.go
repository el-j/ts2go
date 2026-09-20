package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createTestServer(t *testing.T) *Server {
	t.Helper()
	server, err := NewServer(":0")
	if err != nil {
		t.Fatalf("failed to create test server: %v", err)
	}
	return server
}

func TestHealthEndpoint(t *testing.T) {
	server := createTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	server.mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if resp["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got %v", resp["status"])
	}
}

func TestRootEndpoint(t *testing.T) {
	server := createTestServer(t)

	// Valid root
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	server.mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// 404 for unknown path under /
	req404 := httptest.NewRequest(http.MethodGet, "/unknown-page", nil)
	w404 := httptest.NewRecorder()
	server.mux.ServeHTTP(w404, req404)

	if w404.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w404.Code)
	}
}

func TestCORS(t *testing.T) {
	server := createTestServer(t)
	handler := server.enableCORS(server.mux)

	// Test OPTIONS preflight
	req := httptest.NewRequest(http.MethodOptions, "/health", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200 for OPTIONS, got %d", w.Code)
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("expected CORS header to be *")
	}
}

func TestSettingsEndpoints(t *testing.T) {
	server := createTestServer(t)

	// GET settings
	reqGet := httptest.NewRequest(http.MethodGet, "/api/settings", nil)
	wGet := httptest.NewRecorder()
	server.mux.ServeHTTP(wGet, reqGet)

	if wGet.Code != http.StatusOK {
		t.Errorf("expected status 200 for GET /api/settings, got %d", wGet.Code)
	}

	// POST settings
	body := []byte(`{"theme":"dark","autoSave":true,"autoSaveDelay":500}`)
	reqPost := httptest.NewRequest(http.MethodPost, "/api/settings", bytes.NewReader(body))
	wPost := httptest.NewRecorder()
	server.mux.ServeHTTP(wPost, reqPost)

	if wPost.Code != http.StatusOK {
		t.Errorf("expected status 200 for POST /api/settings, got %d: %s", wPost.Code, wPost.Body.String())
	}
}

func TestStateEndpoints(t *testing.T) {
	server := createTestServer(t)

	// Method Not Allowed
	reqPatch := httptest.NewRequest(http.MethodPatch, "/api/state?project=/tmp/test", nil)
	wPatch := httptest.NewRecorder()
	server.mux.ServeHTTP(wPatch, reqPatch)

	if wPatch.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405 for PATCH, got %d", wPatch.Code)
	}

	// GET without project parameter
	reqNoPath := httptest.NewRequest(http.MethodGet, "/api/state", nil)
	wNoPath := httptest.NewRecorder()
	server.mux.ServeHTTP(wNoPath, reqNoPath)

	if wNoPath.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for GET without project param, got %d", wNoPath.Code)
	}

	// GET with project param
	reqGet := httptest.NewRequest(http.MethodGet, "/api/state?project=/tmp/test-project", nil)
	wGet := httptest.NewRecorder()
	server.mux.ServeHTTP(wGet, reqGet)

	if wGet.Code != http.StatusOK {
		t.Errorf("expected status 200 for GET with project, got %d", wGet.Code)
	}

	// POST with project param
	body := []byte(`{"projectPath":"/tmp/test-project"}`)
	reqPost := httptest.NewRequest(http.MethodPost, "/api/state?project=/tmp/test-project", bytes.NewReader(body))
	wPost := httptest.NewRecorder()
	server.mux.ServeHTTP(wPost, reqPost)

	if wPost.Code != http.StatusOK {
		t.Errorf("expected status 200 for POST with project, got %d: %s", wPost.Code, wPost.Body.String())
	}
}

func TestTranspileEndpoints(t *testing.T) {
	server := createTestServer(t)

	// Invalid Method (GET on /api/transpile)
	reqGet := httptest.NewRequest(http.MethodGet, "/api/transpile", nil)
	wGet := httptest.NewRecorder()
	server.mux.ServeHTTP(wGet, reqGet)

	if wGet.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", wGet.Code)
	}

	// Invalid JSON body
	reqBadJSON := httptest.NewRequest(http.MethodPost, "/api/transpile", bytes.NewReader([]byte("{invalid-json")))
	wBadJSON := httptest.NewRecorder()
	server.mux.ServeHTTP(wBadJSON, reqBadJSON)

	if wBadJSON.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for bad JSON, got %d", wBadJSON.Code)
	}
}

func TestBuildAndTestEndpoints(t *testing.T) {
	server := createTestServer(t)

	// GET on /api/build -> 405
	reqBuildGet := httptest.NewRequest(http.MethodGet, "/api/build", nil)
	wBuildGet := httptest.NewRecorder()
	server.mux.ServeHTTP(wBuildGet, reqBuildGet)
	if wBuildGet.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", wBuildGet.Code)
	}

	// GET on /api/test -> 405
	reqTestGet := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	wTestGet := httptest.NewRecorder()
	server.mux.ServeHTTP(wTestGet, reqTestGet)
	if wTestGet.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", wTestGet.Code)
	}

	// Bad JSON on /api/build
	reqBad := httptest.NewRequest(http.MethodPost, "/api/build", bytes.NewReader([]byte("{invalid")))
	wBad := httptest.NewRecorder()
	server.mux.ServeHTTP(wBad, reqBad)
	if wBad.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", wBad.Code)
	}
}
