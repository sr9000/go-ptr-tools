package ref

type Ref[T any, PT interface{ *T }] struct {
	val PT
}
