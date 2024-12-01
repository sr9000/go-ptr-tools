package opt

import "reflect"

type Filter func(reflect.Value) bool

var (
	NotZero Filter = func(rv reflect.Value) bool {
		return !rv.IsZero()
	}
	NotNil Filter = func(rv reflect.Value) bool {
		switch rv.Kind() {
		default:
			return false
		case reflect.Ptr, reflect.Interface,
			reflect.Slice, reflect.Map,
			reflect.Chan, reflect.Func:
			return !rv.IsNil()
		}
	}
	NotEmpty Filter = func(rv reflect.Value) bool {
		switch rv.Kind() {
		default:
			return false
		case reflect.String, reflect.Array,
			reflect.Slice, reflect.Map,
			reflect.Chan:
			return rv.Len() > 0
		}
	}
)

// OfAny creates an optional value based on the given value and filter functions.
// If any filter function returns false, it returns Empty.
func OfAny[T any](v T, filter Filter) Opt[T] {
	rv := reflect.ValueOf(v)
	if rv.IsValid() && filter(rv) {
		return Of(v)
	}

	return Empty[T]()
}
