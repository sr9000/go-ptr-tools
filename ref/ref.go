// Package ref provides a reference type Ref that can be used to store an "always valid pointer".
package ref

// Ref must be used as mental anchor for "always valid pointer" guarantees.
type Ref[T any] struct {
	ptr *T
}

// Ptr returns the pointer to the value.
// There is no need to check for nil, because the pointer is guaranteed to be not nil.
func (r Ref[T]) Ptr() *T {
	return r.ptr
}

// Val returns the value itself.
func (r Ref[T]) Val() T {
	return *r.ptr
}
