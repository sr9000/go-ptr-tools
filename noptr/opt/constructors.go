package opt

// Empty creates an empty optional value.
func Empty[T any]() Opt[T] {
	return Opt[T]{ptr: nil}
}

// Of creates an optional value from the given value.
func Of[T any](v T) Opt[T] {
	return Opt[T]{ptr: &v}
}

// OfPtr creates an optional value based on the given pointer.
func OfPtr[T any](ptr *T) Opt[T] {
	if ptr == nil {
		return Empty[T]()
	}

	return Opt[T]{ptr: ptr}
}

// OfMap takes a value from a map by key.
// If the key is not found, it returns Empty.
func OfMap[K comparable, V any, M ~map[K]V](m M, k K) Opt[V] {
	v, ok := m[k]

	return Wrap(v, ok)
}
