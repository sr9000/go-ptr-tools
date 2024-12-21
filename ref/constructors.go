package ref

import (
	"github.com/sr9000/go-ptr-tools/val"
)

// New returns a new Ref with the given value.
func New[T any](val T) Ref[T] {
	return Ref[T]{&val}
}

// From returns a new Ref based on the given pointer.
func From[T any](ptr *T) (Ref[T], error) {
	if ptr == nil {
		return val.Zero[Ref[T]](), ErrPtrMustNotBeNil
	}

	return Ref[T]{ptr}, nil
}

// Must returns a new Ref based on the given pointer.
// It panics if the pointer is nil.
func Must[T any](ptr *T) Ref[T] {
	ref, err := From(ptr)
	if err != nil {
		panic(err)
	}

	return ref
}

// Guaranteed returns a new Ref based on the given pointer that is guaranteed to be not nil.
// It is useful when you know that the pointer is not nil, but the compiler does not.
func Guaranteed[T any](notNilPtr *T) Ref[T] {
	return Ref[T]{notNilPtr}
}
