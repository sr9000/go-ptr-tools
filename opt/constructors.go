package opt

func New[T any](val T, ok bool) Opt[T] {
	return Opt[T]{val: val, ok: ok}
}

func Literal[T any](val T) Opt[T] {
	return Opt[T]{val: val, ok: true}
}

func FromErr[T any](val T, err error) Opt[T] {
	if err == nil {
		return Opt[T]{val: val, ok: true}
	}

	var empty Opt[T] // zero opt is valid empty opt

	return empty
}

func FromPtr[T any](ptr *T) Opt[T] {
	if ptr != nil {
		return Opt[T]{val: *ptr, ok: true}
	}

	var empty Opt[T] // zero opt is valid empty opt

	return empty
}

func FromZero[T comparable](v T) Opt[T] {
	var zero T
	if v != zero {
		return Opt[T]{val: v, ok: true}
	}

	var empty Opt[T] // zero opt is valid empty opt

	return empty
}
