package examples_test

import "fmt"

type Aint = *int         // alias
type Pint *int           // new type (wrapper)
type Sint struct{ *int } // struct embedding

func foo[T any](msg string, x *T) {
	fmt.Printf("foo %s: %T\n", msg, x)
}

func bar[T any, PT *T](msg string, x PT) {
	fmt.Printf("bar %s: %T\n", msg, x)
}

func Example_tPtDifference() {
	x := 42 //nolint
	foo("pointer", &x)
	bar("pointer", &x)

	px := Pint(&x)
	foo("wrapper", px)
	bar("wrapper", (*int)(px)) // requires explicit conversion

	ax := Aint(px)
	foo("alias", ax)
	bar("alias", ax)

	sx := Sint{&x}
	foo("struct", sx.int) // requires field access
	bar("struct", sx.int) // requires field access

	// Output:
	// foo pointer: *int
	// bar pointer: *int
	// foo wrapper: *int
	// bar wrapper: *int
	// foo alias: *int
	// bar alias: *int
	// foo struct: *int
	// bar struct: *int
}
