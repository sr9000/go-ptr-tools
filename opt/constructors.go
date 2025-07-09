package opt

// Of creates a new Opt with the given value and ok set to true.
func Of[T any](val T) Opt[T] {
	return Opt[T]{val: val, ok: true}
}

// FromOk creates a new Opt with the given value and ok status.
func FromOk[T any](val T, ok bool) Opt[T] {
	return Opt[T]{val: val, ok: ok}
}

// FromErr creates a new Opt from a value and an error.
//   - If the error is nil, it returns an Opt with the value and ok set to true.
//   - If the error is not nil, it returns an empty Opt.
func FromErr[T any](val T, err error) (o Opt[T]) {
	if err == nil {
		return Opt[T]{val: val, ok: true}
	}

	return // zero opt is valid empty opt
}

// FromPtr creates a new Opt from a pointer.
//   - If the pointer is not nil, it returns an Opt with the value pointed to and ok set to true.
//   - If the pointer is nil, it returns an empty Opt.
func FromPtr[T any](ptr *T) (o Opt[T]) {
	if ptr != nil {
		return Opt[T]{val: *ptr, ok: true}
	}

	return // zero opt is valid empty opt
}

// FromZero creates a new Opt from a zeroable (emptiable) value.
//   - If the value is not equal to its zero value, it returns an Opt with the value and ok set to true.
//   - If the value is equal to its zero value, it returns an empty Opt.
func FromZero[T comparable](v T) (o Opt[T]) {
	var zero T
	if v != zero {
		return Opt[T]{val: v, ok: true}
	}

	return // zero opt is valid empty opt
}
