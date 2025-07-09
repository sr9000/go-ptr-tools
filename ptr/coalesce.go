package ptr

import "github.com/sr9000/go-ptr-tools/ref"

// Coalesce returns the first not nil pointer.
// If all pointers are nil, it returns nil.
func Coalesce[T any](pointers ...*T) *T {
	for _, ptr := range pointers {
		if ptr != nil {
			return ptr
		}
	}

	return nil
}

// Else returns the first not nil pointer as a ref.Ref.
// If all pointers are nil, it returns the final ref.Ref.
func Else[T any](final ref.Ref[T], pointers ...*T) ref.Ref[T] {
	for _, ptr := range pointers {
		if ptr != nil {
			return ref.Guaranteed(ptr)
		}
	}

	return final
}
