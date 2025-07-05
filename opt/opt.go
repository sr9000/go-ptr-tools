package opt

type Opt[T any] struct {
	val T
	ok  bool
}

func (o Opt[T]) IsPresent() bool {
	return o.ok
}

func (o Opt[T]) IsEmpty() bool {
	return !o.ok
}

func (o Opt[T]) Get() (T, bool) {
	return o.val, o.ok
}

func (o Opt[T]) Ptr() *T {
	if !o.ok {
		return nil
	}

	return &o.val
}
