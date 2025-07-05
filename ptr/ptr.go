// Package ptrfunc provides several helper routines for working with pointers.
package ptr

// New returns pointer onto a value, works with primitive types (int, string, etc.).
func New[T any](v T) *T {
	return &v
}
