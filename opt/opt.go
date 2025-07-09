// Package opt provides a generic optional type Opt
// that can be used to represent a value that may or may not be present.
package opt

// Opt is a generic optional type that can hold a value of type T or be empty.
type Opt[T any] struct {
	val T
	ok  bool
}

// IsPresent returns true if the Opt contains a value.
func (o Opt[T]) IsPresent() bool {
	return o.ok
}

// IsMissing returns true if the Opt does not contain a value.
func (o Opt[T]) IsMissing() bool {
	return !o.ok
}

// Get returns the value and a boolean indicating if the value is present.
func (o Opt[T]) Get() (T, bool) {
	return o.val, o.ok
}

// Ptr returns a pointer to the value if it is present, or nil if it is not.
func (o Opt[T]) Ptr() *T {
	if !o.ok {
		return nil
	}

	return &o.val
}
