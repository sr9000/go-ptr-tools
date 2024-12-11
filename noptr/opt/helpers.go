package opt

import "reflect"

// Wrap creates an optional value based on the given value and condition.
// Returns Empty if the condition is false or a non-nil error.
func Wrap[T any](v T, cond any) Opt[T] {
	if cond == nil || cond == true {
		return Of(v)
	}

	return Empty[T]()
}

// MapGet takes a value from a map by key.
// If the key is not found, it returns Empty.
func MapGet[K comparable, V any, M ~map[K]V](m M, k K) Opt[V] {
	v, ok := m[k]

	return Wrap(v, ok)
}

// CastTo tries to cast the given value to the target type.
// If the cast isn't possible, it returns Empty.
func CastTo[T any](v any) Opt[T] {
	r, ok := v.(T)

	return Wrap(r, ok)
}

// ValidateInterface checks if the given value is a valid interface.
// If the interface or its underlying value is nil, it returns Empty.
func ValidateInterface[I any](iface any) Opt[I] {
	if iface == nil {
		return Empty[I]()
	}

	rv := reflect.ValueOf(iface)
	if !rv.IsValid() || (rv.Kind() == reflect.Ptr && rv.IsNil()) {
		return Empty[I]()
	}

	return CastTo[I](iface)
}
