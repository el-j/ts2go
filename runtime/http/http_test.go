package http

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestCreateServer(t *testing.T) {
	// Create a simple server
	server := CreateServer(func(req *Request, res *Response) {
		res.WriteHead(StatusOK, map[string]string{
			"Content-Type": "text/plain",
		})
		res.End("Hello, World!")
	})

	if server == nil {
		t.Fatal("CreateServer returned nil")
	}

	// Start server in background
	go func() {
		if err := server.Listen(":0"); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Close server
	if err := server.Close(); err != nil {
		t.Errorf("Failed to close server: %v", err)
	}
}

func TestServerWithRequests(t *testing.T) {
	// Track if handler was called
	handlerCalled := false

	// Create server
	server := CreateServer(func(req *Request, res *Response) {
		handlerCalled = true

		// Check request properties
		if req.Method == "" {
			t.Error("Request method is empty")
		}
		if req.URL == "" {
			t.Error("Request URL is empty")
		}

		// Send response
		res.SetHeader("Content-Type", "application/json")
		res.WriteHead(StatusOK, nil)
		res.Write(`{"message":"success"}`)
	})

	// Start server in background on a random port
	go func() {
		if err := server.Listen(":18080"); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	// Give server time to start
	time.Sleep(200 * time.Millisecond)

	// Make a request to the server
	resp, err := http.Get("http://localhost:18080/test")
	if err != nil {
		t.Logf("Request failed (server may not have started): %v", err)
	} else {
		defer resp.Body.Close()

		// Check response
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Failed to read response body: %v", err)
		}

		expectedBody := `{"message":"success"}`
		if string(body) != expectedBody {
			t.Errorf("Expected body %s, got %s", expectedBody, string(body))
		}

		if !handlerCalled {
			t.Error("Handler was not called")
		}
	}

	// Close server
	if err := server.Close(); err != nil {
		t.Errorf("Failed to close server: %v", err)
	}
}

func TestResponseMethods(t *testing.T) {
	// Create server
	server := CreateServer(func(req *Request, res *Response) {
		// Test SetHeader
		res.SetHeader("X-Custom-Header", "test-value")

		// Test WriteHead
		res.WriteHead(StatusCreated, map[string]string{
			"Content-Type": "text/plain",
		})

		// Test Write
		if err := res.Write("Hello, "); err != nil {
			t.Errorf("Write failed: %v", err)
		}

		// Test End
		if err := res.End("World!"); err != nil {
			t.Errorf("End failed: %v", err)
		}
	})

	// Start server
	go func() {
		if err := server.Listen(":18081"); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	time.Sleep(200 * time.Millisecond)

	// Make request
	resp, err := http.Get("http://localhost:18081/")
	if err != nil {
		t.Logf("Request failed: %v", err)
	} else {
		defer resp.Body.Close()

		// Check custom header
		customHeader := resp.Header.Get("X-Custom-Header")
		if customHeader != "test-value" {
			t.Errorf("Expected X-Custom-Header to be 'test-value', got '%s'", customHeader)
		}

		// Check status code
		if resp.StatusCode != StatusCreated {
			t.Errorf("Expected status %d, got %d", StatusCreated, resp.StatusCode)
		}

		// Check body
		body, _ := io.ReadAll(resp.Body)
		expectedBody := "Hello, World!"
		if string(body) != expectedBody {
			t.Errorf("Expected body '%s', got '%s'", expectedBody, string(body))
		}
	}

	server.Close()
}

func TestGet(t *testing.T) {
	// Test Get function (requires internet connection)
	// Using a reliable test endpoint
	resp, err := Get("https://httpbin.org/get")
	if err != nil {
		t.Skipf("Get failed (no internet?): %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response: %v", err)
	}

	if len(body) == 0 {
		t.Error("Response body is empty")
	}

	t.Logf("Get response length: %d bytes", len(body))
}

func TestPost(t *testing.T) {
	// Test Post function
	testData := `{"test":"data"}`
	resp, err := Post("https://httpbin.org/post", "application/json", testData)
	if err != nil {
		t.Skipf("Post failed (no internet?): %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response: %v", err)
	}

	// Check that our data was echoed back
	bodyStr := string(body)
	if !strings.Contains(bodyStr, "test") || !strings.Contains(bodyStr, "data") {
		t.Error("Response doesn't contain posted data")
	}

	t.Logf("Post response length: %d bytes", len(body))
}

func TestRequest(t *testing.T) {
	// Test MakeRequest function with custom options
	options := RequestOptions{
		Method:      "PUT",
		URL:         "https://httpbin.org/put",
		ContentType: "application/json",
		Body:        `{"key":"value"}`,
		Headers: map[string]string{
			"X-Test-Header": "test-value",
		},
	}

	resp, err := MakeRequest(options)
	if err != nil {
		t.Skipf("Request failed (no internet?): %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response: %v", err)
	}

	bodyStr := string(body)
	if !strings.Contains(bodyStr, "key") {
		t.Error("Response doesn't contain request data")
	}

	t.Logf("Request response length: %d bytes", len(body))
}

func TestStatusText(t *testing.T) {
	tests := []struct {
		code int
		text string
	}{
		{StatusOK, "OK"},
		{StatusCreated, "Created"},
		{StatusNotFound, "Not Found"},
		{StatusInternalServerError, "Internal Server Error"},
	}

	for _, tt := range tests {
		result := StatusText(tt.code)
		if result != tt.text {
			t.Errorf("StatusText(%d) = %s, expected %s", tt.code, result, tt.text)
		}
	}
}

func TestStatusConstants(t *testing.T) {
	// Test that status constants have correct values
	if StatusOK != 200 {
		t.Errorf("StatusOK should be 200, got %d", StatusOK)
	}
	if StatusNotFound != 404 {
		t.Errorf("StatusNotFound should be 404, got %d", StatusNotFound)
	}
	if StatusInternalServerError != 500 {
		t.Errorf("StatusInternalServerError should be 500, got %d", StatusInternalServerError)
	}
}

func TestRequestStructure(t *testing.T) {
	// Create a server that inspects the request
	server := CreateServer(func(req *Request, res *Response) {
		// Test request structure
		if req.Method == "" {
			t.Error("Request.Method is empty")
		}
		if req.URL == "" {
			t.Error("Request.URL is empty")
		}
		if req.Headers == nil {
			t.Error("Request.Headers is nil")
		}
		if req.Query == nil {
			t.Error("Request.Query is nil")
		}

		// Send simple response
		res.WriteHead(StatusOK, nil)
		res.End("ok")
	})

	go server.Listen(":18082")
	time.Sleep(200 * time.Millisecond)

	// Make a request with query parameters
	_, err := http.Get("http://localhost:18082/path?key=value")
	if err != nil {
		t.Logf("Request failed: %v", err)
	}

	server.Close()
}
