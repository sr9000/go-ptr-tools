package pt1val_test

import (
	"fmt"

	"github.com/sr9000/go-noptr/noptr/val"
)

type Point struct {
	x, y float64
}

func Example() {
	fmt.Println(val.Zero[int]())
	fmt.Println(val.Zero[string]())
	fmt.Println(val.Zero[[]int]())
	fmt.Println(val.Zero[any]())
	fmt.Println(val.Zero[Point]())

	// Output:
	// 0
	//
	// []
	// <nil>
	// {0 0}
}

var (
	_ = val.Zero[Point]().x // suppress unused
	_ = val.Zero[Point]().y // suppress unused
)
