package examples_test

import (
	"encoding/json"
	"fmt"
	"github.com/sr9000/go-ptr-tools/ptr"
	"github.com/sr9000/go-ptr-tools/ref"
	"sync"
)

// -----------------------------------------------------------------------------
// 1. Using bool to check validity.

func Example_wrong_usingBoolToCheckValidity() {
	find := func(xs []int, x int) (int, bool) {
		for i, v := range xs {
			if v == x {
				return i, true
			}
		}
		return 0, false
	}

	xs := []int{1, 2, 3}
	if idx, found := find(xs, 2); found {
		fmt.Println("Found at index", idx)
	} else {
		fmt.Println("Not found")
	}
	// Output:
	// Found at index 1
}

func Example_good_usingPointerToCheckValidity() {
	find := func(xs []int, x int) *int {
		for i, v := range xs {
			if v == x {
				return &i
			}
		}
		return nil
	}

	xs := []int{1, 2, 3}
	if idx := find(xs, 2); idx != nil {
		fmt.Println("Found at index", *idx)
	} else {
		fmt.Println("Not found")
	}
	// Output:
	// Found at index 1
}

// -----------------------------------------------------------------------------
// 2. Missing pointer checks.

func Example_wrong_missingPointerChecks() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic")
		}
	}()
	plusOneToPtr := func(n *int) {
		*n++
	}

	var n *int
	plusOneToPtr(n)
	// Output:
	// Recovered from panic
}

func Example_good_usingRefArgs() {
	plusOne := func(r ref.Ref[int]) {
		*r.Ptr()++
	}

	var n int
	plusOne(ref.Guaranteed(&n))
	fmt.Println(n)
	// Output: 1
}

func Example_good_withPointerChecks() {
	plusOne := func(n *int) {
		if n == nil {
			return
		}
		*n++
	}

	var n *int
	plusOne(n)
	// Output:
}

// -----------------------------------------------------------------------------
// 3. ref declarations.

func Example_wrong_refRefDeclarations() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic")
		}
	}()
	var n ref.Ref[int]
	fmt.Println(n.Val())
	// Output:
	// Recovered from panic
}

func Example_good_usingInitialization() {
	n := ref.Literal(42)
	fmt.Println(n.Val())
	// Output: 42
}

// -----------------------------------------------------------------------------
// 4. ref struct fields.

func Example_wrong_refRefStructFields() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic")
		}
	}()
	type S struct {
		N ref.Ref[int]
	}
	s := S{}
	fmt.Println(s.N.Val())
	// Output:
	// Recovered from panic
}

func Example_good_usingDirectTypes() {
	type S struct {
		N int
	}
	s := S{}
	fmt.Println(s.N)
	// Output: 0
}

func Example_good_usingBuilderWithFunctionalOptions() {
	type S struct {
		N        ref.Ref[int]
		isValidN bool
	}
	withN := func(n int) func(s ref.Ref[S]) {
		return func(s ref.Ref[S]) {
			s.Ptr().N = ref.Literal(n)
			s.Ptr().isValidN = true
		}
	}

	builer := func(opts ...func(ref.Ref[S])) (S, error) {
		s := S{}
		for _, opt := range opts {
			opt(ref.Guaranteed(&s))
		}
		if !s.isValidN {
			return S{}, fmt.Errorf("missing N")
		}
		return s, nil
	}

	s, err := builer(withN(42))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s.N.Val())
	// Output: 42
}

func Example_good_usingPointers() {
	type S struct {
		N *int
	}
	s := S{}
	fmt.Println(s.N)
	// Output: <nil>
}

// -----------------------------------------------------------------------------
// 5. using pointer to ref (or storing ref as any).

func Example_wrong_usingPointerOrAny() {
	n := ref.Literal(42)
	foo := func(n *ref.Ref[int]) {
		if n != nil {
			fmt.Println(n.Val())
		}
	}
	foo(&n)

	bar := func(n any) {
		if n, ok := n.(ref.Ref[int]); ok {
			fmt.Println(n.Val())
		}
	}
	bar(n)
	// Output:
	// 42
	// 42
}

func Example_good_usingRef() {
	n := ref.Literal(42)
	foobar := func(n ref.Ref[int]) {
		fmt.Println(n.Val())
	}
	foobar(n)
	// Output: 42
}

// -----------------------------------------------------------------------------
// 6. wrap any value or pointer with ref.

func Example_wrong_wrapAnyValueOrPointerWithRefRef() {
	var x any
	foo := func(n ref.Ref[any]) {
		*n.Ptr() = 100
	}
	foo(ref.Guaranteed(&x))
	fmt.Println(x)

	var n int
	bar := func(n ref.Ref[*int]) {
		**n.Ptr() = 500
	}
	bar(ref.Literal(&n))
	fmt.Println(n)
	// Output:
	// 100
	// 500
}

func Example_good_usingRefAgain() {
	var n int
	foobar := func(n ref.Ref[int]) {
		*n.Ptr() = 42
	}
	foobar(ref.Guaranteed(&n))
	fmt.Println(n)
	// Output: 42
}

// -----------------------------------------------------------------------------
// 7. using ref.Ref for optional results.

func Example_wrong_usingRefForOptionalResults() {
	minimum := func(xs []*int) ref.Ref[*int] {
		if len(xs) == 0 {
			return ref.Literal[*int](nil) // no minimum for empty list
		}

		// assume that all pointers are valid
		mn := &xs[0]
		for i := range xs {
			if *xs[i] < **mn {
				mn = &xs[i]
			}
		}

		return ref.Guaranteed(mn) // reference to slice element itself
	}

	xs := []*int{ptr.New(1), ptr.New(2), ptr.New(3)}
	res := minimum(xs)
	if res.Ptr() != nil {
		fmt.Println("Minimum is", **res.Ptr())
	} else {
		fmt.Println("No minimum")
	}
	// Output: Minimum is 1
}

func Example_good_usingPointerForOptionalResults() {
	minimum := func(xs []*int) **int {
		if len(xs) == 0 {
			return nil
		}

		// assume that all pointers are valid
		mn := &xs[0]
		for i := range xs {
			if *xs[i] < **mn {
				mn = &xs[i]
			}
		}

		return mn
	}

	xs := []*int{ptr.New(1), ptr.New(2), ptr.New(3)}
	if res := minimum(xs); res != nil {
		fmt.Println("Minimum is", **res)
	} else {
		fmt.Println("No minimum")
	}
	// Output: Minimum is 1
}

// -----------------------------------------------------------------------------
// 8. extra actions to make a copy.

func Example_wrong_extraActionsToMakeACopy() {
	n, m := 5, 7
	r1 := ref.Guaranteed(&n)
	// ...
	r2 := ref.Guaranteed(r1.Ptr()) // extra actions to make a copy
	r1 = ref.Guaranteed(&m)        // change r1

	fmt.Println(r1.Val(), r2.Val())
	// Output: 7 5
}

func Example_good_usingAssignment() {
	n, m := 5, 7
	r1 := ref.Guaranteed(&n)
	// ...
	r2 := r1                // safe to copy
	r1 = ref.Guaranteed(&m) // reassignment doesn't affect r2

	fmt.Println(r1.Val(), r2.Val())
	// Output: 7 5
}

// -----------------------------------------------------------------------------
// 9. unprotected concurrent access.

func Example_wrong_unprotectedConcurrentAccess() {
	var n int
	var wg sync.WaitGroup

	r := ref.Guaranteed(&n)
	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			*r.Ptr()++
		}()
	}

	wg.Wait()
	fmt.Println(n < 1000)
	// Output: true
}

func Example_good_usingMutex() {
	var n int
	var mu sync.Mutex
	var wg sync.WaitGroup
	r := ref.Guaranteed(&n)

	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			*r.Ptr()++
		}()
	}

	wg.Wait()
	fmt.Println(n)
	// Output: 1000
}

// -----------------------------------------------------------------------------
// 10. marshaling ref.Ref values.

func Example_wrong_marshalingRefValues() {
	type S struct {
		N ref.Ref[int]
	}
	saveJSON := func(s S) {
		data, _ := json.Marshal(s)
		fmt.Println(string(data))
	}
	saveJSON(S{N: ref.Literal(42)})
	// Output: {"N":{}}
}

func Example_good_usingValueDirectly() {
	type S struct {
		N ref.Ref[int]
	}
	saveJSON := func(s S) {
		x := struct {
			N int
		}{
			N: s.N.Val(),
		}
		data, _ := json.Marshal(x)
		fmt.Println(string(data))
	}
	saveJSON(S{N: ref.Literal(42)})
	// Output: {"N":42}
}
