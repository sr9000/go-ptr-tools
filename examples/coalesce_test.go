package examples_test

import (
	"fmt"

	"github.com/sr9000/go-ptr-tools/ptr"
	"github.com/sr9000/go-ptr-tools/ref"
)

type Point struct {
	X, Y    int
	IsValid bool
}

func XCoordFromArray(coords []int) *int {
	if len(coords) == 0 {
		return nil
	}

	return &coords[0]
}

func XCoordFromStruct(p Point) *int {
	if !p.IsValid {
		return nil
	}

	return &p.X
}

func PrintPointer[T any, PT *T](name string, ptr PT) {
	name += ":"
	if ptr == nil {
		fmt.Println(name, "<nil>") //nolint
	} else {
		fmt.Println(name, *ptr) //nolint
	}
}

func Example_chooseXCoord() {
	coords := []int{1, 2}
	p := Point{X: 3, Y: 4, IsValid: true}

	x1 := XCoordFromArray(coords)
	x2 := XCoordFromStruct(p)

	x := ptr.Coalesce(x1, x2)

	PrintPointer("x1", x1)
	PrintPointer("x2", x2)
	PrintPointer("x", x)
	// Output:
	// x1: 1
	// x2: 3
	// x: 1
}

func Example_noXCoord() {
	var coords []int

	var p Point

	x := ptr.Else(ref.Literal(5), XCoordFromArray(coords), XCoordFromStruct(p))

	PrintPointer("x1", XCoordFromArray(coords))
	PrintPointer("x2", XCoordFromStruct(p))
	fmt.Println("x:", x.Val())
	// Output:
	// x1: <nil>
	// x2: <nil>
	// x: 5
}
