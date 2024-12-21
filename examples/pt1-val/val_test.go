package pt1val_test

import (
	"fmt"
	"github.com/sr9000/go-noptr/val"
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

var _ = Point{x: 1, y: 2} // suppress unused
