package pt2ptr_test

import (
	"fmt"
	"github.com/sr9000/go-noptr/ptr"
)

func ExamplePtrNew() {
	fmt.Println(*ptr.New(42))
	fmt.Println(*ptr.New("hello"))

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

func PtrNew() {} // suppress govet
func PtrNil() {} // suppress govet
