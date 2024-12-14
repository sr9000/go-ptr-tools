package opt

import "reflect"

// Wrap creates an optional value based on the given value and condition.
// Returns Opt is `cond` is true or a nil interface (supposed to be nil-error).
func Wrap[T any](v T, cond any) Opt[T] {
	if cond == nil || cond == true {
		return Of(v)
	}

	return Empty[T]()
}

// ParseInterface checks if the given value is a valid interface (exclude nil).
// If the interface or its underlying value is nil, it returns Empty.
func ParseInterface[I any](iface any) Opt[I] {
	if iface == nil {
		return Empty[I]()
	}

	rv := reflect.ValueOf(iface)
	if !rv.IsValid() || (rv.Kind() == reflect.Ptr && rv.IsNil()) {
		return Empty[I]()
	}

	return CastTo[I](iface)
}

// Unwrap2 unwraps either both values or none.
func Unwrap2[T1, T2 any](o1 Opt[T1], o2 Opt[T2]) (v1 T1, v2 T2, ok bool) {
	t1, ok1 := o1.Get()
	t2, ok2 := o2.Get()

	if ok1 && ok2 {
		return t1, t2, true
	}

	return
}

// Unwrap3 unwraps either three values or none.
func Unwrap3[T1, T2, T3 any](o1 Opt[T1], o2 Opt[T2], o3 Opt[T3]) (v1 T1, v2 T2, v3 T3, ok bool) {
	t1, ok1 := o1.Get()
	t2, ok2 := o2.Get()
	t3, ok3 := o3.Get()

	if ok1 && ok2 && ok3 {
		return t1, t2, t3, true
	}

	return
}
