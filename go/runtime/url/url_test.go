package url

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected *URL
	}{
		{
			input: "https://example.com:8080/path/to/resource?key=value&foo=bar#section",
			expected: &URL{
				Protocol: "https:",
				Hostname: "example.com",
				Port:     "8080",
				Host:     "example.com:8080",
				Pathname: "/path/to/resource",
				Search:   "?key=value&foo=bar",
				Hash:     "#section",
			},
		},
		{
			input: "http://example.com/",
			expected: &URL{
				Protocol: "http:",
				Hostname: "example.com",
				Port:     "",
				Host:     "example.com",
				Pathname: "/",
				Search:   "",
				Hash:     "",
			},
		},
		{
			input: "https://example.com/path",
			expected: &URL{
				Protocol: "https:",
				Hostname: "example.com",
				Pathname: "/path",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse failed: %v", err)
			}

			if result.Protocol != tt.expected.Protocol {
				t.Errorf("Protocol: got %s, want %s", result.Protocol, tt.expected.Protocol)
			}
			if result.Hostname != tt.expected.Hostname {
				t.Errorf("Hostname: got %s, want %s", result.Hostname, tt.expected.Hostname)
			}
			if result.Port != tt.expected.Port {
				t.Errorf("Port: got %s, want %s", result.Port, tt.expected.Port)
			}
			if result.Pathname != tt.expected.Pathname {
				t.Errorf("Pathname: got %s, want %s", result.Pathname, tt.expected.Pathname)
			}
			if result.Search != tt.expected.Search {
				t.Errorf("Search: got %s, want %s", result.Search, tt.expected.Search)
			}
			if result.Hash != tt.expected.Hash {
				t.Errorf("Hash: got %s, want %s", result.Hash, tt.expected.Hash)
			}
		})
	}
}

func TestParseQuery(t *testing.T) {
	urlStr := "https://example.com/path?key=value&foo=bar&baz=qux"
	result, err := Parse(urlStr)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(result.Query) != 3 {
		t.Errorf("Expected 3 query params, got %d", len(result.Query))
	}

	if result.Query["key"] != "value" {
		t.Errorf("Expected key=value, got key=%s", result.Query["key"])
	}
	if result.Query["foo"] != "bar" {
		t.Errorf("Expected foo=bar, got foo=%s", result.Query["foo"])
	}
	if result.Query["baz"] != "qux" {
		t.Errorf("Expected baz=qux, got baz=%s", result.Query["baz"])
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    *URL
		expected string
	}{
		{
			name: "Full URL",
			input: &URL{
				Protocol: "https:",
				Hostname: "example.com",
				Port:     "8080",
				Pathname: "/path",
				Search:   "?key=value",
				Hash:     "#section",
			},
			expected: "https://example.com:8080/path?key=value#section",
		},
		{
			name: "Simple URL",
			input: &URL{
				Protocol: "http:",
				Hostname: "example.com",
				Pathname: "/",
			},
			expected: "http://example.com/",
		},
		{
			name: "URL with Host",
			input: &URL{
				Protocol: "https:",
				Host:     "example.com:443",
				Pathname: "/path",
			},
			expected: "https://example.com:443/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Format(tt.input)
			if result != tt.expected {
				t.Errorf("Format: got %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	tests := []struct {
		base     string
		relative string
		expected string
	}{
		{
			base:     "https://example.com/path/to/page",
			relative: "other.html",
			expected: "https://example.com/path/to/other.html",
		},
		{
			base:     "https://example.com/path/to/page",
			relative: "/absolute/path",
			expected: "https://example.com/absolute/path",
		},
		{
			base:     "https://example.com/path/to/page",
			relative: "../parent",
			expected: "https://example.com/path/parent",
		},
		{
			base:     "https://example.com/",
			relative: "https://other.com/page",
			expected: "https://other.com/page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.relative, func(t *testing.T) {
			result, err := Resolve(tt.base, tt.relative)
			if err != nil {
				t.Fatalf("Resolve failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Resolve: got %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestQueryString_Parse(t *testing.T) {
	qs := &QueryString{}

	tests := []struct {
		input    string
		expected map[string]string
	}{
		{
			input: "key=value&foo=bar&baz=qux",
			expected: map[string]string{
				"key": "value",
				"foo": "bar",
				"baz": "qux",
			},
		},
		{
			input: "?name=John&age=30",
			expected: map[string]string{
				"name": "John",
				"age":  "30",
			},
		},
		{
			input:    "",
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := qs.Parse(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Parse: got %d params, want %d", len(result), len(tt.expected))
			}

			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("Parse: got %s=%s, want %s=%s", k, result[k], k, v)
				}
			}
		})
	}
}

func TestQueryString_Stringify(t *testing.T) {
	qs := &QueryString{}

	params := map[string]string{
		"key": "value",
		"foo": "bar",
	}

	result := qs.Stringify(params)

	// Check that both params are present (order may vary)
	if !contains(result, "key=value") {
		t.Errorf("Stringify: result doesn't contain 'key=value': %s", result)
	}
	if !contains(result, "foo=bar") {
		t.Errorf("Stringify: result doesn't contain 'foo=bar': %s", result)
	}
}

func TestQueryString_Escape(t *testing.T) {
	qs := &QueryString{}

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "hello world",
			expected: "hello+world",
		},
		{
			input:    "key=value&foo=bar",
			expected: "key%3Dvalue%26foo%3Dbar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := qs.Escape(tt.input)
			if result != tt.expected {
				t.Errorf("Escape: got %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestQueryString_Unescape(t *testing.T) {
	qs := &QueryString{}

	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "hello+world",
			expected: "hello world",
		},
		{
			input:    "key%3Dvalue%26foo%3Dbar",
			expected: "key=value&foo=bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := qs.Unescape(tt.input)
			if err != nil {
				t.Fatalf("Unescape failed: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Unescape: got %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestPackageFunctions(t *testing.T) {
	// Test ParseQuery
	query := ParseQuery("key=value&foo=bar")
	if len(query) != 2 {
		t.Errorf("ParseQuery: got %d params, want 2", len(query))
	}

	// Test StringifyQuery
	params := map[string]string{"a": "1", "b": "2"}
	result := StringifyQuery(params)
	if !contains(result, "a=1") || !contains(result, "b=2") {
		t.Errorf("StringifyQuery: unexpected result %s", result)
	}

	// Test Escape
	escaped := Escape("hello world")
	if escaped != "hello+world" {
		t.Errorf("Escape: got %s, want hello+world", escaped)
	}

	// Test Unescape
	unescaped, err := Unescape("hello+world")
	if err != nil {
		t.Errorf("Unescape failed: %v", err)
	}
	if unescaped != "hello world" {
		t.Errorf("Unescape: got %s, want 'hello world'", unescaped)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
