package opt

// Coalesce returns the first non-empty Opt.
// If all Opt values are empty, it returns empty Opt.
func Coalesce[T any](opts ...Opt[T]) (r Opt[T]) {
	for _, o := range opts {
		if o.IsPresent() {
			return o
		}
	}

	return
}

// Else returns the value from the first non-empty Opt.
// If all Opt values are empty, it returns the final value.
func Else[T any](final T, opts ...Opt[T]) T {
	for _, o := range opts {
		if val, ok := o.Get(); ok {
			return val
		}
	}

	return final
}
