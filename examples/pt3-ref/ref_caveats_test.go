package pt3ref_test

import (
	"fmt"

	"github.com/sr9000/go-ptr-tools/ref"
)

func NewSegment(a, b ref.Ref[Point]) Segment {
	return Segment{isValid: true, a: a, b: b}
}

type Segment struct {
	isValid bool           // best practice to force using NewSegment
	a, b    ref.Ref[Point] // replacement for *Point
}

func (s Segment) String() string {
	// this is much cleaner than `s.a == nil || s.b == nil` with explicit *Point
	if !s.isValid {
		// avoid panic in Segment's method
		return fmt.Sprintf("<invalid %T>", s)
	}

	// remember do not check for nil with ref.Ref
	// creating valid ref is the responsibility of the caller
	return fmt.Sprintf("{%v %v}", s.a.Val(), s.b.Val())
}

func ExampleDefaultSegment() {
	var s Segment // references are not safe to use in structs by default

	fmt.Println(s)
	// Output: <invalid pt3ref_test.Segment>
}

func ExampleNewSegment() {
	a := Point{x: 1, y: 2}
	b := Point{x: 3, y: 4}
	segment := NewSegment(ref.New(a), ref.New(b)) // but ok if explicitly created
	fmt.Println(segment)
	// Output: {{1 2} {3 4}}
}

func DefaultSegment() {} // suppress govet
