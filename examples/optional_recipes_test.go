package examples_test

import (
	"fmt"
	"strconv"

	"github.com/sr9000/go-ptr-tools/ptr"
	"github.com/sr9000/go-ptr-tools/ref"
)

func Map[T any, R any](opt *T, f func(T) R) *R {
	if opt == nil {
		return nil
	}

	return ptr.New(f(*opt))
}

// 1. Optional.Of() value from a variable.
func Example_ofVariable() {
	var n int //nolint

	opt := &n
	if opt != nil { //nolint
		fmt.Println("Got non-nil pointer")
	}
	// Output:
	// Got non-nil pointer
}

// 2. Optional.Of() literal value.
func Example_ofLiteral() {
	opt := ptr.New(42)
	if opt != nil && *opt == 42 {
		fmt.Println("Got int 42")
	}
	// Output:
	// Got int 42
}

// 3. Optional.Empty().
func Example_empty() {
	optEmpty := (*int)(nil)
	if optEmpty == nil {
		fmt.Println("Got nil pointer")
	}
	// Output:
	// Got nil pointer
}

// 4. Optional.IsPresent().
func Example_isPresent() {
	opt := ptr.New(42)
	if opt != nil {
		fmt.Println("Got non-nil pointer")
	}
	// Output:
	// Got non-nil pointer
}

// 5. Optional.IsEmpty().
func Example_isEmpty() {
	optEmpty := (*int)(nil)
	if optEmpty == nil {
		fmt.Println("Got nil pointer")
	}
	// Output:
	// Got nil pointer
}

// 6. Optional.Get().
func Example_get() {
	var val int

	opt := ptr.New(42)
	if opt != nil {
		val = *opt
	}

	if val == 42 {
		fmt.Println("Got value 42")
	}
	// Output:
	// Got value 42
}

// 7. Optional.OrElse().
func Example_orElse() {
	optEmpty := (*int)(nil)
	refVal := ptr.Else(ref.Literal(42), optEmpty)

	if refVal.Ptr() != nil && refVal.Val() == 42 {
		fmt.Println("Got pointer to 42")
	}
	// Output:
	// Got pointer to 42
}

// 8. Optional.Coalesce().
func Example_coalesce() {
	optA := ptr.New(1)
	optB := ptr.New(2)
	single := ptr.Coalesce(optA, optB)

	if single != nil && *single == 1 {
		fmt.Println("Got pointer to 1")
	}
	// Output:
	// Got pointer to 1
}

// 9. Optional.Map().
func Example_mapWithFunc() {
	opt := ptr.New(42)
	mapped := Map(opt, func(v int) int { return v * 2 })

	if mapped != nil && *mapped == 84 {
		fmt.Println("Got pointer to 84")
	}
	// Output:
	// Got pointer to 84
}

// 9. Optional.Map() inside a function.
func Example_mapInPlace() {
	var res *string

	optEmpty := (*int)(nil)
	if optEmpty != nil {
		res = ptr.New(strconv.Itoa(*optEmpty))
	}

	if res == nil {
		fmt.Println("Got nil pointer")
	}
	// Output:
	// Got nil pointer
}
