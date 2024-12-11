package opt

func Coalesce[T any](opts ...Opt[T]) Opt[T] {
	for _, o := range opts {
		if o.Ptr() != nil {
			return o
		}
	}

	return Empty[T]()
}

func Map[T, R any](o Opt[T], f func(T) R) Opt[R] {
	if o.Ptr() == nil {
		return Empty[R]()
	}

	return Of(f(*o.Ptr()))
}

func MapEx[T, R any](o Opt[T], f func(T) (R, error)) (Opt[R], error) {
	if o.Ptr() == nil {
		return Empty[R](), nil
	}

	v, err := f(*o.Ptr())
	if err != nil {
		return Empty[R](), err
	}

	return Of(v), nil
}

func Unwrap2[T1, T2 any](o1 Opt[T1], o2 Opt[T2]) (v1 T1, v2 T2, ok bool) {
	t1, ok1 := o1.Get()
	t2, ok2 := o2.Get()

	if ok1 && ok2 {
		return t1, t2, true
	}

	return
}

func Unwrap3[T1, T2, T3 any](o1 Opt[T1], o2 Opt[T2], o3 Opt[T3]) (v1 T1, v2 T2, v3 T3, ok bool) {
	t1, ok1 := o1.Get()
	t2, ok2 := o2.Get()
	t3, ok3 := o3.Get()

	if ok1 && ok2 && ok3 {
		return t1, t2, t3, true
	}

	return
}
