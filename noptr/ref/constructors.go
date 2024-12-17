package ref

import "github.com/sr9000/go-noptr/noptr/val"

// Of returns a new Ref with the given value.
func Of[T any](val T) Ref[T] {
	return Ref[T]{&val}
}

// OfPtr returns a new Ref based on the given pointer.
func OfPtr[T any](ptr *T) (Ref[T], error) {
	if ptr == nil {
		return val.Zero[Ref[T]](), ErrPtrMustNotBeNil
	}

	return Ref[T]{ptr}, nil
}

// Must returns a new Ref based on the given pointer.
// It panics if the pointer is nil.
func Must[T any](ptr *T) Ref[T] {
	ref, err := OfPtr(ptr)
	if err != nil {
		panic(err)
	}

	return ref
}
