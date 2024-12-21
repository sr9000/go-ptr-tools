// Package ptr is an impostor of the "go-ptr-tools" library :P
// It provides a pointer factory functions.
package ptr

// New is a helper routine that allocates a new any value
// to store v and returns a pointer to it.
func New[T any](v T) *T {
	return &v
}

// Nil is a helper routine that returns a nil pointer of any type.
func Nil[T any]() (nilPtr *T) { return }
