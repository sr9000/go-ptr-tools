package ref

type Ref[T any, PT interface{ *T }] struct {
	val PT
}

// Ptr returns the pointer to the value.
func (r Ref[T, PT]) Ptr() PT {
	return r.val
}
