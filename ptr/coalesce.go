package ptr

import "github.com/sr9000/go-ptr-tools/ref"

// Coalesce returns the first not nil pointer.
func Coalesce[T any](pointers ...*T) *T {
	for _, ptr := range pointers {
		if ptr != nil {
			return ptr
		}
	}

	return nil
}

// Finalize returns the first not nil pointer as a Ref or the final Ref.
func Finalize[T any](final ref.Ref[T], pointers ...*T) ref.Ref[T] {
	for _, ptr := range pointers {
		if ptr != nil {
			return ref.Guaranteed(ptr)
		}
	}

	return final
}
