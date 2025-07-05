package test_test

import (
	"github.com/sr9000/go-ptr-tools/opt"
	"testing"
)

func parsePrimeBool(number int) (int, bool) {
	if number < 2 {
		return 0, false
	}

	for i := 2; i*i <= number; i++ {
		if number%i == 0 {
			return 0, false
		}
	}

	return number, true
}

func parsePrimeOpt(number int) (res opt.Opt[int]) {
	if number < 2 {
		return // return empty opt
	}

	for i := 2; i*i <= number; i++ {
		if number%i == 0 {
			return // return empty opt
		}
	}

	return opt.Literal(number)
}

func parsePrimePtr(number int) *int {
	if number < 2 {
		return nil
	}

	for i := 2; i*i <= number; i++ {
		if number%i == 0 {
			return nil
		}
	}

	return &number
}

func parseAny(number int) any {
	if number < 2 {
		return nil
	}

	for i := 2; i*i <= number; i++ {
		if number%i == 0 {
			return nil
		}
	}

	return number
}

const testIters = 100_000

func boolLoop() []int {
	res := make([]int, testIters)

	for i := range testIters {
		if n, ok := parsePrimeBool(i); ok {
			res[i] = n
		}
	}

	return res
}

func ptrLoop() []*int {
	res := make([]*int, testIters)
	for i := range testIters {
		res[i] = parsePrimePtr(i)
	}

	return res
}

func optLoop() []opt.Opt[int] {
	res := make([]opt.Opt[int], testIters)

	for i := range testIters {
		res[i] = parsePrimeOpt(i)
	}

	return res
}

func anyLoop() []any {
	res := make([]any, testIters)

	for i := range testIters {
		res[i] = parseAny(i)
	}
	return res
}

func BenchmarkBool(b *testing.B) {
	for range b.N {
		arr := boolLoop()
		s := 0

		for _, n := range arr {
			if n != 0 {
				s += n
			}
		}
	}
}

func BenchmarkPtr(b *testing.B) {
	for range b.N {
		arr := ptrLoop()
		s := 0

		for _, n := range arr {
			if n != nil {
				s += *n
			}
		}
	}
}

func BenchmarkOpt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		arr := optLoop()
		s := 0

		for _, o := range arr {
			if n, ok := o.Get(); ok {
				s += n
			}
		}
	}
}

func BenchmarkAny(b *testing.B) {
	for i := 0; i < b.N; i++ {
		arr := anyLoop()
		s := 0

		for _, n := range arr {
			if v, ok := n.(int); ok {
				s += v
			}
		}
	}
}
