// Package ptrfunc provides several helper routines for working with pointers.
package ptr

// New returns pointer onto a value, works with primitive types (int, string, etc.).
func New[T any](v T) *T {
	return &v
}

// FromOk creates a new pointer from a value and an ok flag.
// If ok equals true, it returns a pointer to val, else nil.
func FromOk[T any](val T, ok bool) *T {
	if ok {
		return &val
	}

	return nil
}

// FromErr creates a new pointer from a value and an error.
//   - If the error is nil, it returns a pointer to val.
//   - If the error is not nil, it returns nil.
func FromErr[T any](val T, err error) *T {
	if err == nil {
		return &val
	}

	return nil
}

// FromZero creates a new pointer from a zeroable (emptiable) value.
//   - If the value is not equal to its zero value, it returns  a pointer to val.
//   - If the value is equal to its zero value, it returns nil.
func FromZero[T comparable](v T) *T {
	var zero T
	if v != zero {
		return &v
	}

	return nil
}
