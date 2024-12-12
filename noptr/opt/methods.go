package opt

func (o Opt[T]) Trigger(funcs ...func(T)) {
	if v, ok := o.Get(); ok {
		for _, f := range funcs {
			f(v)
		}
	}
}

func (o Opt[T]) Apply(funcs ...func(*T)) {
	if ptr := o.Ptr(); ptr != nil {
		for _, f := range funcs {
			f(ptr)
		}
	}
}

func (o Opt[T]) ApplyEx(funcs ...func(*T) error) error {
	if ptr := o.Ptr(); ptr != nil {
		for _, f := range funcs {
			if err := f(ptr); err != nil {
				return err
			}
		}
	}

	return nil
}

func (o Opt[T]) Else(val T) T {
	if v, ok := o.Get(); ok {
		return v
	}

	return val
}

func (o Opt[T]) ElseGet(f func() T) T {
	if v, ok := o.Get(); ok {
		return v
	}

	return f()
}

func (o Opt[T]) ElseGetEx(f func() (T, error)) (T, error) {
	if v, ok := o.Get(); ok {
		return v, nil
	}

	return f()
}
