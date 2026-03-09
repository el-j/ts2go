package http

import (
	"io"
	"net/http"
	"strings"
)

// Server represents an HTTP server
type Server struct {
	server *http.Server
	mux    *http.ServeMux
}

// Request represents an HTTP request
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
	Query   map[string]string
}

// Response represents an HTTP response
type Response struct {
	writer http.ResponseWriter
}

// RequestHandler is a function that handles HTTP requests
type RequestHandler func(req *Request, res *Response)

// CreateServer creates a new HTTP server
func CreateServer(handler RequestHandler) *Server {
	mux := http.NewServeMux()

	// Wrap the handler to convert between Node.js-style and Go-style
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Convert Go request to our Request type
		headers := make(map[string]string)
		for k, v := range r.Header {
			headers[k] = strings.Join(v, ", ")
		}

		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		query := make(map[string]string)
		for k, v := range r.URL.Query() {
			query[k] = strings.Join(v, ", ")
		}

		req := &Request{
			Method:  r.Method,
			URL:     r.URL.String(),
			Headers: headers,
			Body:    string(body),
			Query:   query,
		}

		res := &Response{writer: w}

		// Call the user's handler
		handler(req, res)
	})

	return &Server{
		server: &http.Server{
			Handler: mux,
		},
		mux: mux,
	}
}

// Listen starts the server on the specified address
func (s *Server) Listen(addr string) error {
	s.server.Addr = addr
	return s.server.ListenAndServe()
}

// Close shuts down the server
func (s *Server) Close() error {
	return s.server.Close()
}

// WriteHead sets the status code and headers
func (r *Response) WriteHead(statusCode int, headers map[string]string) {
	for k, v := range headers {
		r.writer.Header().Set(k, v)
	}
	r.writer.WriteHeader(statusCode)
}

// Write writes data to the response
func (r *Response) Write(data string) error {
	_, err := r.writer.Write([]byte(data))
	return err
}

// End ends the response
func (r *Response) End(data string) error {
	if data != "" {
		_, err := r.writer.Write([]byte(data))
		return err
	}
	return nil
}

// SetHeader sets a response header
func (r *Response) SetHeader(key, value string) {
	r.writer.Header().Set(key, value)
}

// Get makes a GET request to the specified URL
func Get(url string) (*http.Response, error) {
	return http.Get(url)
}

// Post makes a POST request to the specified URL
func Post(url string, contentType string, body string) (*http.Response, error) {
	return http.Post(url, contentType, strings.NewReader(body))
}

// RequestOptions represents options for making HTTP requests
type RequestOptions struct {
	Method      string
	URL         string
	Headers     map[string]string
	Body        string
	ContentType string
}

// MakeRequest makes an HTTP request with the specified options
func MakeRequest(options RequestOptions) (*http.Response, error) {
	var bodyReader io.Reader
	if options.Body != "" {
		bodyReader = strings.NewReader(options.Body)
	}

	req, err := http.NewRequest(options.Method, options.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set headers
	if options.ContentType != "" {
		req.Header.Set("Content-Type", options.ContentType)
	}

	for k, v := range options.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	return client.Do(req)
}

// StatusText returns the text for an HTTP status code
func StatusText(code int) string {
	return http.StatusText(code)
}

// Common status codes
const (
	StatusOK                  = http.StatusOK                  // 200
	StatusCreated             = http.StatusCreated             // 201
	StatusAccepted            = http.StatusAccepted            // 202
	StatusNoContent           = http.StatusNoContent           // 204
	StatusMovedPermanently    = http.StatusMovedPermanently    // 301
	StatusFound               = http.StatusFound               // 302
	StatusNotModified         = http.StatusNotModified         // 304
	StatusBadRequest          = http.StatusBadRequest          // 400
	StatusUnauthorized        = http.StatusUnauthorized        // 401
	StatusForbidden           = http.StatusForbidden           // 403
	StatusNotFound            = http.StatusNotFound            // 404
	StatusMethodNotAllowed    = http.StatusMethodNotAllowed    // 405
	StatusInternalServerError = http.StatusInternalServerError // 500
	StatusNotImplemented      = http.StatusNotImplemented      // 501
	StatusBadGateway          = http.StatusBadGateway          // 502
	StatusServiceUnavailable  = http.StatusServiceUnavailable  // 503
)
