package opt_test

import "errors"

var errTest = errors.New("test error")

type testFooer interface {
	Foo()
}

type testFoo struct {
	a, b int
}

func (*testFoo) Foo() {}

type testBar struct{}

func (testBar) Foo() {}
