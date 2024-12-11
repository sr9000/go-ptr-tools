package opt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sr9000/go-noptr/noptr/opt"
)

func TestOfAny(t *testing.T) {
	x := 42
	emptyChan := make(chan int)
	emptyBuffChan := make(chan int, 10)
	filledChan := make(chan int, 10)
	filledChan <- 42

	cases := []struct {
		name                      string
		value                     any
		notZero, notNil, notEmpty bool
	}{
		{"Int", 42,
			true, false, false},

		{"ZeroInt", *new(int),
			false, false, false},

		{"String", "hello",
			true, false, true},

		{"EmptyString", "",
			false, false, false},

		{"ZeroString", *new(string),
			false, false, false},

		{"Slice", []int{1, 2, 3},
			true, true, true},

		{"EmptySlice", []int{},
			true, true, false},

		{"ZeroSlice", *new([]int),
			false, false, false},

		{"Map", map[string]int{"a": 1},
			true, true, true},

		{"EmptyMap", map[string]int{},
			true, true, false},

		{"ZeroMap", *new(map[string]int),
			false, false, false},

		{"Struct", testFoo{a: 1, b: 2},
			true, false, false},

		{"ZeroStruct", *new(testFoo),
			false, false, false},

		{"Interface", testFooer(&testBar{}),
			true, true, false},

		{"Interface of nil ptr", testFooer((*testBar)(nil)),
			false, false, false},

		{"ZeroInterface", *new(testFooer),
			false, false, false},

		{"IntPointer", &x,
			true, true, false},

		{"ZeroIntPointer", new(int),
			true, true, false}, // int is zero BUT pointer is not zero (not nil)

		{"NilIntPointer", *new(*int),
			false, false, false},

		{"ZeroStructPointer", new(testFoo),
			true, true, false}, // testFoo is zero BUT pointer is not zero (not nil)

		{"NilStructPointer", *new(*testFoo),
			false, false, false},

		{"EmptyChan", emptyChan,
			true, true, false},

		{"EmptyBuffChan", emptyBuffChan,
			true, true, false},

		{"FilledChan", filledChan,
			true, true, true},

		{"ZeroChan", *new(chan int),
			false, false, false},

		{"Func", func() {},
			true, true, false},

		{"ZeroFunc", *new(func()),
			false, false, false},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, cs.notZero, nil != opt.OfAny(cs.value, opt.NotZero).Ptr(), "NotZero")
			assert.Equal(t, cs.notNil, nil != opt.OfAny(cs.value, opt.NotNil).Ptr(), "NotNil")
			assert.Equal(t, cs.notEmpty, nil != opt.OfAny(cs.value, opt.NotEmpty).Ptr(), "NotEmpty")
		})
	}
}
