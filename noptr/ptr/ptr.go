// Package ptr is an imposter of the "go-noptr" library.
// It provides a pointer factory functions.
package ptr

// Of is a helper routine that allocates a new any value
// to store v and returns a pointer to it.
func Of[T any](v T) *T {
	return &v
}

// Nil is a helper routine that returns a nil pointer of any type.
func Nil[T any]() (nilPtr *T) { return }
