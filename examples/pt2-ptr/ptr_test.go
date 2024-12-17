package pt2ptr_test

import (
	"fmt"

	"github.com/sr9000/go-noptr/noptr/ptr"
)

func ExamplePtrOf() {
	fmt.Println(*ptr.Of(42))
	fmt.Println(*ptr.Of("hello"))

	// Output:
	// 42
	// hello
}

func ExamplePtrNil() {
	fmt.Printf("%#v\n", ptr.Nil[int]())
	fmt.Printf("%#v\n", ptr.Nil[string]())
	fmt.Printf("%#v\n", ptr.Nil[[]int]())

	// Output:
	// (*int)(nil)
	// (*string)(nil)
	// (*[]int)(nil)
}

func PtrOf()  {} // suppress govet
func PtrNil() {} // suppress govet
