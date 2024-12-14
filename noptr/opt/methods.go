package opt

// Trigger calls the given functions with the value if it is present.
func (o Opt[T]) Trigger(funcs ...func(T)) {
	if v, ok := o.Get(); ok {
		for _, f := range funcs {
			f(v)
		}
	}
}

// Apply applies the given functions as patches to the value if it is present.
func (o Opt[T]) Apply(funcs ...func(*T)) {
	if ptr := o.Ptr(); ptr != nil {
		for _, f := range funcs {
			f(ptr)
		}
	}
}

// ApplyEx applies the given functions as patches to the value if it is present.
// If any of the functions returns an error, it stops the execution and returns the error.
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

// Else returns the value from the optional else the given value.
func (o Opt[T]) Else(val T) T {
	if v, ok := o.Get(); ok {
		return v
	}

	return val
}

// ElseGet returns the value from the optional else the result of a given function.
func (o Opt[T]) ElseGet(f func() T) T {
	if v, ok := o.Get(); ok {
		return v
	}

	return f()
}

// ElseGetEx returns the value from the optional else the result of a given function and an error if it occurs.
func (o Opt[T]) ElseGetEx(f func() (T, error)) (T, error) {
	if v, ok := o.Get(); ok {
		return v, nil
	}

	return f()
}
