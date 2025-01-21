package examples_test

import (
	"fmt"
	"sync"

	"github.com/sr9000/go-ptr-tools/ptr"
)

func CoalesceGetters[T any](getters ...func() *T) *T {
	var wg sync.WaitGroup

	wg.Add(len(getters))
	results := make([]*T, len(getters))

	for i, g := range getters {
		go func() {
			defer wg.Done()

			results[i] = g()
		}()
	}

	wg.Wait()

	return ptr.Coalesce(results...)
}

func Example_coalesceGetters() {
	for range 1000 {
		var first, second, third *int

		first = ptr.New(1)
		third = ptr.New(3)

		getters := []func() *int{
			func() *int { return first },
			func() *int { return second },
			func() *int { return third },
		}

		result := CoalesceGetters(getters...)
		if result == nil || *result != 1 {
			fmt.Println("fail")
		}
	}
	// Output:
}
