package pt2ptr_test

import (
	"fmt"

	"github.com/sr9000/go-ptr-tools/ptr"
)

func ExamplePtrNew() {
	fmt.Println(*ptr.New(42))
	fmt.Println(*ptr.New("hello"))

	// Output:
	// 42
	// hello
}

func PtrNew() {} // suppress govet
