package mapper

import (
	"fmt"
	"regexp"
	"strings"
)

// CallTransformer transforms TypeScript/JavaScript API calls to Go equivalents
type CallTransformer struct {
	db         *MappingDatabase
	classifier *Classifier
}

// NewCallTransformer creates a new call transformer
func NewCallTransformer(db *MappingDatabase) *CallTransformer {
	return &CallTransformer{
		db:         db,
		classifier: NewClassifier(db),
	}
}

// TransformResult contains the result of transforming an API call
type TransformResult struct {
	Original    string  // Original code (e.g., "axios.get(url)")
	Transformed string  // Transformed code (e.g., "resty.R().Get(url)")
	PackageUsed string  // Package name used
	Confidence  float64 // Confidence in transformation
	Notes       string  // Additional notes
	Error       error   // Error if transformation failed
}

// TransformCall transforms a single API call expression
func (t *CallTransformer) TransformCall(code string) *TransformResult {
	result := &TransformResult{
		Original:   code,
		Confidence: 0.0,
	}

	// Extract package.method pattern
	pattern := regexp.MustCompile(`^(\w+)\.(\w+)\s*\(`)
	matches := pattern.FindStringSubmatch(code)

	if len(matches) < 3 {
		result.Error = fmt.Errorf("unable to parse call expression: %s", code)
		return result
	}

	packageName := matches[1]
	methodName := matches[2]

	// Try to find mapping
	mapping, err := t.db.GetMapping(packageName)
	if err != nil {
		// Check if it might be an alias
		result.Error = fmt.Errorf("no mapping found for package: %s", packageName)
		result.Confidence = 0.0
		return result
	}

	// Look for API mapping
	apiKey := fmt.Sprintf("%s.%s", packageName, methodName)
	goAPI, ok := mapping.APIMappings[apiKey]

	if !ok {
		// No specific API mapping, try generic transformation
		result.Transformed = t.genericTransform(code, packageName, methodName, mapping)
		result.Confidence = 0.5
		result.Notes = "No specific API mapping found, used generic transformation"
	} else {
		// Use specific API mapping
		result.Transformed = t.applyAPIMapping(code, goAPI)
		result.Confidence = 0.9
	}

	result.PackageUsed = mapping.Go
	return result
}

// TransformMultipleCalls transforms multiple API calls in code
func (t *CallTransformer) TransformMultipleCalls(code string, patterns []string) map[string]*TransformResult {
	results := make(map[string]*TransformResult)

	for _, pattern := range patterns {
		result := t.TransformCall(pattern)
		results[pattern] = result
	}

	return results
}

// applyAPIMapping applies a specific API mapping from the database
func (t *CallTransformer) applyAPIMapping(code, goAPI string) string {
	// Extract the arguments from the original call
	args := t.extractArguments(code)

	// If goAPI doesn't contain parentheses, it's just a method name
	if !strings.Contains(goAPI, "(") {
		return fmt.Sprintf("%s(%s)", goAPI, args)
	}

	// If goAPI ends with empty parens, add arguments there
	if strings.HasSuffix(goAPI, "()") {
		return goAPI[:len(goAPI)-2] + fmt.Sprintf("(%s)", args)
	}

	// Otherwise add arguments to the end
	return fmt.Sprintf("%s(%s)", goAPI, args)
}

// genericTransform provides a generic transformation when no specific mapping exists
func (t *CallTransformer) genericTransform(code, packageName, methodName string, mapping *Mapping) string {
	args := t.extractArguments(code)

	// Get the Go package alias
	alias := t.getPackageAlias(mapping.Go)

	// Capitalize method name (Go convention)
	goMethodName := capitalizeFirst(methodName)

	return fmt.Sprintf("%s.%s(%s)", alias, goMethodName, args)
}

// extractArguments extracts arguments from a function call
func (t *CallTransformer) extractArguments(code string) string {
	// Find the opening parenthesis
	start := strings.Index(code, "(")
	if start == -1 {
		return ""
	}

	// Find the matching closing parenthesis
	depth := 0
	for i := start; i < len(code); i++ {
		switch code[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return strings.TrimSpace(code[start+1 : i])
			}
		}
	}

	return ""
}

// getPackageAlias extracts an alias from a Go package path
func (t *CallTransformer) getPackageAlias(goPackage string) string {
	parts := strings.Split(goPackage, "/")
	if len(parts) == 0 {
		return ""
	}

	lastPart := parts[len(parts)-1]

	// Remove version suffix if present (e.g., v2, v3)
	if strings.HasPrefix(lastPart, "v") && len(lastPart) > 1 {
		if len(parts) > 1 {
			lastPart = parts[len(parts)-2]
		}
	}

	return lastPart
}

// PatternMatcher helps identify API patterns in code
type PatternMatcher struct {
	patterns []*regexp.Regexp
}

// NewPatternMatcher creates a new pattern matcher
func NewPatternMatcher() *PatternMatcher {
	patterns := []*regexp.Regexp{
		// Method calls: package.method(args)
		regexp.MustCompile(`(\w+)\.(\w+)\s*\([^)]*\)`),

		// Chained calls: package.method().method()
		regexp.MustCompile(`(\w+)\.(\w+)\s*\(\)[^.]*\.(\w+)\s*\([^)]*\)`),

		// Property access with method: package.property.method()
		regexp.MustCompile(`(\w+)\.(\w+)\.(\w+)\s*\([^)]*\)`),
	}

	return &PatternMatcher{patterns: patterns}
}

// FindAPICalls finds all API calls in a code snippet
func (pm *PatternMatcher) FindAPICalls(code string) []string {
	calls := make([]string, 0)
	seen := make(map[string]bool)

	for _, pattern := range pm.patterns {
		matches := pattern.FindAllString(code, -1)
		for _, match := range matches {
			if !seen[match] {
				calls = append(calls, match)
				seen[match] = true
			}
		}
	}

	return calls
}

// TransformationPlan represents a complete transformation plan for a code snippet
type TransformationPlan struct {
	OriginalCode    string
	Transformations []*TransformResult
	Summary         string
	SuccessRate     float64
}

// CreateTransformationPlan creates a complete plan for transforming code
func (t *CallTransformer) CreateTransformationPlan(code string) *TransformationPlan {
	matcher := NewPatternMatcher()
	calls := matcher.FindAPICalls(code)

	transformations := make([]*TransformResult, 0)
	successCount := 0

	for _, call := range calls {
		result := t.TransformCall(call)
		transformations = append(transformations, result)

		if result.Error == nil && result.Confidence > 0.7 {
			successCount++
		}
	}

	successRate := 0.0
	if len(transformations) > 0 {
		successRate = float64(successCount) / float64(len(transformations))
	}

	plan := &TransformationPlan{
		OriginalCode:    code,
		Transformations: transformations,
		SuccessRate:     successRate,
	}

	plan.Summary = fmt.Sprintf(
		"Transformation Plan\n"+
			"Total API Calls: %d\n"+
			"Successful: %d\n"+
			"Success Rate: %.1f%%\n",
		len(transformations),
		successCount,
		successRate*100,
	)

	return plan
}

// ApplyTransformations applies all transformations to code
func (tp *TransformationPlan) ApplyTransformations() string {
	result := tp.OriginalCode

	// Apply transformations in reverse order to preserve positions
	for i := len(tp.Transformations) - 1; i >= 0; i-- {
		trans := tp.Transformations[i]
		if trans.Error == nil && trans.Transformed != "" {
			result = strings.Replace(result, trans.Original, trans.Transformed, 1)
		}
	}

	return result
}

// Helper types for advanced transformations

// ContextualTransformer transforms calls based on surrounding context
type ContextualTransformer struct {
	transformer *CallTransformer
	context     map[string]string // Variable name -> package mapping
}

// NewContextualTransformer creates a transformer that uses variable context
func NewContextualTransformer(db *MappingDatabase) *ContextualTransformer {
	return &ContextualTransformer{
		transformer: NewCallTransformer(db),
		context:     make(map[string]string),
	}
}

// SetVariableContext sets the package for a variable name
func (ct *ContextualTransformer) SetVariableContext(varName, packageName string) {
	ct.context[varName] = packageName
}

// TransformWithContext transforms a call using variable context
// Example: client.get(url) where client is an axios instance
func (ct *ContextualTransformer) TransformWithContext(code, varName string) *TransformResult {
	packageName, ok := ct.context[varName]
	if !ok {
		return &TransformResult{
			Original: code,
			Error:    fmt.Errorf("no context found for variable: %s", varName),
		}
	}

	// Replace variable name with package name temporarily for transformation
	tempCode := strings.Replace(code, varName+".", packageName+".", 1)
	result := ct.transformer.TransformCall(tempCode)

	// Keep original in result
	result.Original = code

	return result
}
