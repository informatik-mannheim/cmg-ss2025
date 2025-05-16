package utils

import "slices"

// .map function from JavaScript Arrays in Go
func Map[T any, R any](input []T, mapper func(T) R) []R {
	output := make([]R, len(input))
	for i, v := range input {
		output[i] = mapper(v)
	}
	return output
}

// .filter function from JavaScript Arrays in Go
func Filter[T any](input []T, predicate func(T) bool) []T {
	output := make([]T, 0)
	for _, v := range input {
		if predicate(v) {
			output = append(output, v)
		}
	}
	return output
}

// .some function from JavaScript Arrays in Go
func Some[T any](input []T, predicate func(T) bool) bool {
	return slices.ContainsFunc(input, predicate)
}

// .find function from JavaScript Arrays in Go
func Find[T any](input []T, predicate func(T) bool) (T, bool) {
	for _, v := range input {
		if predicate(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}
