package opt_test

type testFooer interface {
	Foo()
}

type testFoo struct {
	a, b int
}

func (*testFoo) Foo() {}

type testBar struct{}

func (testBar) Foo() {}
