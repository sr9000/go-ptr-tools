package opt

// CastTo tries to cast the given value to the target type.
// If the cast isn't possible, it returns Empty.
func CastTo[T any](v any) Opt[T] {
	r, ok := v.(T)

	return Wrap(r, ok)
}

// Coalesce returns the first non-empty optional value.
func Coalesce[T any](opts ...Opt[T]) Opt[T] {
	for _, o := range opts {
		if o.Ptr() != nil {
			return o
		}
	}

	return Empty[T]()
}

// FMap applies the function inside the Opt monad.
func FMap[T, R any](o Opt[T], f func(T) R) Opt[R] {
	if o.Ptr() == nil {
		return Empty[R]()
	}

	return Of(f(*o.Ptr()))
}

// FMapEx applies the function inside the Opt monad and signals about any error.
func FMapEx[T, R any](o Opt[T], f func(T) (R, error)) (Opt[R], error) {
	if o.Ptr() == nil {
		return Empty[R](), nil
	}

	v, err := f(*o.Ptr())
	if err != nil {
		return Empty[R](), err
	}

	return Of(v), nil
}
