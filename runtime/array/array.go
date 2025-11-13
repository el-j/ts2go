package array

// Array provides JavaScript-like array methods for Go slices
// This runtime library enables transpiled TypeScript code to use
// familiar array methods like filter, map, reduce, etc.

// Filter returns a new slice containing only elements that pass the test
// implemented by the provided function.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map creates a new slice with the results of calling a provided function
// on every element in the slice.
func Map[T any, R any](slice []T, mapper func(T) R) []R {
	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = mapper(item)
	}
	return result
}

// Reduce executes a reducer function on each element of the slice,
// resulting in a single output value.
func Reduce[T any, R any](slice []T, reducer func(R, T) R, initialValue R) R {
	acc := initialValue
	for _, item := range slice {
		acc = reducer(acc, item)
	}
	return acc
}

// Find returns the first element in the slice that satisfies the provided
// testing function. Returns the zero value if no element is found.
func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
	for _, item := range slice {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// FindIndex returns the index of the first element in the slice that
// satisfies the provided testing function. Returns -1 if not found.
func FindIndex[T any](slice []T, predicate func(T) bool) int {
	for i, item := range slice {
		if predicate(item) {
			return i
		}
	}
	return -1
}

// Some tests whether at least one element in the slice passes the test
// implemented by the provided function.
func Some[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}

// Every tests whether all elements in the slice pass the test implemented
// by the provided function.
func Every[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// Push adds one or more elements to the end of a slice and returns the new length.
func Push[T any](slice *[]T, elements ...T) int {
	*slice = append(*slice, elements...)
	return len(*slice)
}

// Pop removes the last element from a slice and returns that element.
func Pop[T any](slice *[]T) (T, bool) {
	var zero T
	if len(*slice) == 0 {
		return zero, false
	}
	lastIndex := len(*slice) - 1
	element := (*slice)[lastIndex]
	*slice = (*slice)[:lastIndex]
	return element, true
}

// Shift removes the first element from a slice and returns that element.
func Shift[T any](slice *[]T) (T, bool) {
	var zero T
	if len(*slice) == 0 {
		return zero, false
	}
	element := (*slice)[0]
	*slice = (*slice)[1:]
	return element, true
}

// Unshift adds one or more elements to the beginning of a slice and
// returns the new length.
func Unshift[T any](slice *[]T, elements ...T) int {
	*slice = append(elements, *slice...)
	return len(*slice)
}

// Includes determines whether a slice includes a certain element.
func Includes[T comparable](slice []T, searchElement T) bool {
	for _, item := range slice {
		if item == searchElement {
			return true
		}
	}
	return false
}

// IndexOf returns the first index at which a given element can be found
// in the slice, or -1 if it is not present.
func IndexOf[T comparable](slice []T, searchElement T) int {
	for i, item := range slice {
		if item == searchElement {
			return i
		}
	}
	return -1
}

// LastIndexOf returns the last index at which a given element can be found
// in the slice, or -1 if it is not present.
func LastIndexOf[T comparable](slice []T, searchElement T) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == searchElement {
			return i
		}
	}
	return -1
}

// Reverse reverses a slice in place.
func Reverse[T any](slice []T) []T {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

// Slice returns a shallow copy of a portion of a slice.
func Slice[T any](slice []T, start, end int) []T {
	if start < 0 {
		start = 0
	}
	if end > len(slice) {
		end = len(slice)
	}
	if start > end {
		start = end
	}
	return slice[start:end]
}

// Concat merges two or more slices.
func Concat[T any](slices ...[]T) []T {
	totalLen := 0
	for _, s := range slices {
		totalLen += len(s)
	}
	result := make([]T, 0, totalLen)
	for _, s := range slices {
		result = append(result, s...)
	}
	return result
}

// Join concatenates all elements of a slice into a string.
func Join(slice []string, separator string) string {
	if len(slice) == 0 {
		return ""
	}
	result := slice[0]
	for i := 1; i < len(slice); i++ {
		result += separator + slice[i]
	}
	return result
}

// ForEach executes a provided function once for each slice element.
func ForEach[T any](slice []T, callback func(T, int)) {
	for i, item := range slice {
		callback(item, i)
	}
}

// Sort sorts the elements of a slice in place using the provided comparison function.
// The comparison function should return true if the first argument is less than the second.
func Sort[T any](slice []T, compare func(T, T) bool) []T {
	// Simple bubble sort for now (can be optimized later)
	n := len(slice)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if !compare(slice[j], slice[j+1]) {
				slice[j], slice[j+1] = slice[j+1], slice[j]
			}
		}
	}
	return slice
}
