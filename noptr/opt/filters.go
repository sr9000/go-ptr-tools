package opt

import "reflect"

// NotZero returns if its value is not zero.
func (o Opt[T]) NotZero() Opt[T] {
	if v, ok := o.Get(); ok {
		rval := reflect.ValueOf(v)
		if !rval.IsValid() || rval.IsZero() {
			return Empty[T]()
		}

		return o
	}

	return Empty[T]()
}

// NotNil returns if its value is not nil.
func (o Opt[T]) NotNil() Opt[T] {
	if v, ok := o.Get(); ok {
		rval := reflect.ValueOf(v)
		if !rval.IsValid() {
			return Empty[T]()
		}

		switch rval.Kind() {
		default:
			return Empty[T]()
		case reflect.Ptr, reflect.Interface,
			reflect.Slice, reflect.Map,
			reflect.Chan, reflect.Func:
			if !rval.IsNil() {
				return o
			}
		}
	}

	return Empty[T]()
}

// NotEmpty returns if its collection len is greater than 0.
func (o Opt[T]) NotEmpty() Opt[T] {
	if v, ok := o.Get(); ok {
		rval := reflect.ValueOf(v)
		if !rval.IsValid() {
			return Empty[T]()
		}

		switch rval.Kind() {
		default:
			return Empty[T]()
		case reflect.String, reflect.Array,
			reflect.Slice, reflect.Map,
			reflect.Chan:
			if rval.Len() > 0 {
				return o
			}
		}
	}

	return Empty[T]()
}
