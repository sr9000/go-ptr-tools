package pt3ref_test

import (
	"fmt"
	ref2 "github.com/sr9000/go-noptr/ref"
)

type Point struct {
	x, y float64
}

func ExampleRefNew() {
	func(r ref2.Ref[int]) {
		fmt.Println(r.Val()) // no need to check for the nil
	}(ref2.New(42))
	// Output: 42
}

func ExampleRefFrom() {
	var err error

	value := 42

	valid, err := ref2.From(&value)
	if err != nil {
		panic(err)
	}

	fmt.Println(valid.Val(), err)

	_, err = ref2.From[int](nil)
	fmt.Println(err)

	// Output:
	// 42 <nil>
	// ptr must not be nil
}

func ExampleModifyingValue() {
	point := Point{x: 1, y: 2}

	func(r ref2.Ref[Point]) {
		r.Ptr().x = 10
		r.Ptr().y = 20
	}(ref2.Must(&point))

	fmt.Println(point)
	// Output: {10 20}
}

func RefNew()         {} // suppress govet
func RefFrom()        {} // suppress govet
func ModifyingValue() {} // suppress govet
