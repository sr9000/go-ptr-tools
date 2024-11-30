package ref

func Of[T any, PT interface{ *T }](val T) Ref[T, PT] {
	return Ref[T, PT]{&val}
}

func OfPtr[T any, PT interface{ *T }](val PT) (Ref[T, PT], error) {
	if val == nil {
		return Ref[T, PT]{val}, ErrPtrMustNotBeNil
	}

	return Ref[T, PT]{val}, nil
}
