// Package ref provides a reference type Ref that can be used to store an "always valid pointer".
package ref

// Ref must be used as metal anchor for "always valid pointer" guarantees.
type Ref[T any] struct {
	ptr *T
}

// Ptr returns the pointer to the value.
//
// PLEASE don't write `ref.Ptr() == nil` checks because it's against the purpose of this type.
// Use this kind of hack only with a to-do promise to fix this technical debt in the near future.
func (r Ref[T]) Ptr() *T {
	return r.ptr
}
