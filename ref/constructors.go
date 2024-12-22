package ref

// New returns a new Ref based on the given pointer.
// It returns an error if the pointer is nil.
func New[T any](ptr *T) (Ref[T], error) {
	if ptr == nil {
		return Ref[T]{}, ErrPtrMustBeNotNil
	}

	return Ref[T]{ptr}, nil
}

// Guaranteed returns a new Ref based on the given pointer that is guaranteed to be not nil.
// It is useful when you know that the pointer is not nil, but the compiler does not.
func Guaranteed[T any](notNilPtr *T) Ref[T] {
	return Ref[T]{notNilPtr}
}
