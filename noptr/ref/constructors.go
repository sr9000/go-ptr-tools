package ref

func invalid[T any]() (invalid Ref[T]) {
	return
}

// Of returns a new Ref with the given value.
func Of[T any](val T) Ref[T] {
	return Ref[T]{&val}
}

// OfPtr returns a new Ref based on the given pointer.
func OfPtr[T any](ptr *T) (Ref[T], error) {
	if ptr == nil {
		return invalid[T](), ErrPtrMustNotBeNil
	}

	return Ref[T]{ptr}, nil
}
