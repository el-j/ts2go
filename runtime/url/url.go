package url

import (
	"net/url"
	"strings"
)

// URL represents a parsed URL
type URL struct {
	Href     string            // Full URL
	Protocol string            // e.g., "http:", "https:"
	Host     string            // Hostname and port
	Hostname string            // Just hostname
	Port     string            // Just port
	Pathname string            // Path
	Search   string            // Query string with ?
	Query    map[string]string // Parsed query
	Hash     string            // Fragment with #
}

// Parse parses a URL string into a URL object
func Parse(urlStr string) (*URL, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// Parse query string
	query := make(map[string]string)
	for k, v := range u.Query() {
		query[k] = strings.Join(v, ",")
	}

	port := u.Port()
	hostname := u.Hostname()

	// Add protocol suffix if not present
	protocol := u.Scheme
	if protocol != "" && !strings.HasSuffix(protocol, ":") {
		protocol += ":"
	}

	// Format search string
	search := ""
	if u.RawQuery != "" {
		search = "?" + u.RawQuery
	}

	// Format hash string
	hash := ""
	if u.Fragment != "" {
		hash = "#" + u.Fragment
	}

	return &URL{
		Href:     urlStr,
		Protocol: protocol,
		Host:     u.Host,
		Hostname: hostname,
		Port:     port,
		Pathname: u.Path,
		Search:   search,
		Query:    query,
		Hash:     hash,
	}, nil
}

// Format formats a URL object into a string
func Format(u *URL) string {
	var result strings.Builder

	// Protocol
	if u.Protocol != "" {
		result.WriteString(u.Protocol)
		if !strings.HasSuffix(u.Protocol, "//") {
			result.WriteString("//")
		}
	}

	// Host (hostname:port)
	if u.Host != "" {
		result.WriteString(u.Host)
	} else if u.Hostname != "" {
		result.WriteString(u.Hostname)
		if u.Port != "" {
			result.WriteString(":")
			result.WriteString(u.Port)
		}
	}

	// Pathname
	if u.Pathname != "" {
		if !strings.HasPrefix(u.Pathname, "/") && result.Len() > 0 {
			result.WriteString("/")
		}
		result.WriteString(u.Pathname)
	}

	// Search (query string)
	if u.Search != "" {
		result.WriteString(u.Search)
	}

	// Hash (fragment)
	if u.Hash != "" {
		result.WriteString(u.Hash)
	}

	return result.String()
}

// Resolve resolves a relative URL against a base URL
func Resolve(base string, relative string) (string, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	relURL, err := url.Parse(relative)
	if err != nil {
		return "", err
	}

	resolved := baseURL.ResolveReference(relURL)
	return resolved.String(), nil
}

// QueryString provides utilities for working with query strings
type QueryString struct{}

// Parse parses a query string into a map
func (qs *QueryString) Parse(queryStr string) map[string]string {
	// Remove leading ? if present
	queryStr = strings.TrimPrefix(queryStr, "?")

	values, err := url.ParseQuery(queryStr)
	if err != nil {
		return make(map[string]string)
	}

	result := make(map[string]string)
	for k, v := range values {
		result[k] = strings.Join(v, ",")
	}

	return result
}

// Stringify converts a map to a query string
func (qs *QueryString) Stringify(params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return values.Encode()
}

// Escape escapes a string for use in a URL
func (qs *QueryString) Escape(s string) string {
	return url.QueryEscape(s)
}

// Unescape unescapes a URL-encoded string
func (qs *QueryString) Unescape(s string) (string, error) {
	return url.QueryUnescape(s)
}

// Global QueryString instance
var DefaultQueryString = &QueryString{}

// Package-level convenience functions for QueryString

// ParseQuery parses a query string into a map
func ParseQuery(queryStr string) map[string]string {
	return DefaultQueryString.Parse(queryStr)
}

// StringifyQuery converts a map to a query string
func StringifyQuery(params map[string]string) string {
	return DefaultQueryString.Stringify(params)
}

// Escape escapes a string for use in a URL
func Escape(s string) string {
	return DefaultQueryString.Escape(s)
}

// Unescape unescapes a URL-encoded string
func Unescape(s string) (string, error) {
	return DefaultQueryString.Unescape(s)
}
