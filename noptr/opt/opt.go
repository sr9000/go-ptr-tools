// Package opt provides an optional type Opt that can be used to store possibly missing values.
package opt

// Opt optional value.
type Opt[T any] struct {
	ptr *T
}

// Ptr returns the actual pointer to the value.
func (o Opt[T]) Ptr() *T {
	return o.ptr
}

// Get returns the value itself.
func (o Opt[T]) Get() (v T, ok bool) {
	if o.ptr == nil {
		return
	}

	return *o.ptr, true
}
