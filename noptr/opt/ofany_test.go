package opt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sr9000/go-noptr/noptr/internal"
	"github.com/sr9000/go-noptr/noptr/opt"
)

func TestOfAny_Primitives(t *testing.T) {
	t.Parallel()

	var (
		zeroInt          int
		zeroString       string
		zeroStruct       testFoo
		nilIntPointer    *int
		nilStructPointer *testFoo
	)

	cases := []struct {
		name                      string
		value                     any
		notZero, notNil, notEmpty bool
	}{
		{"Int", 42,
			true, false, false},

		{"ZeroInt", zeroInt,
			false, false, false},

		{"String", "hello",
			true, false, true},

		{"EmptyString", "",
			false, false, false},

		{"ZeroString", zeroString,
			false, false, false},

		{"Struct", testFoo{a: 1, b: 2},
			true, false, false},

		{"ZeroStruct", zeroStruct,
			false, false, false},

		{"IntPointer", internal.ToPtr(42),
			true, true, false},

		{"ZeroIntPointer", new(int),
			true, true, false}, // int is zero BUT pointer is not zero (not nil)

		{"NilIntPointer", nilIntPointer,
			false, false, false},

		{"ZeroStructPointer", new(testFoo),
			true, true, false}, // testFoo is zero BUT pointer is not zero (not nil)

		{"NilStructPointer", nilStructPointer,
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

func TestOfAny_Nillable(t *testing.T) {
	t.Parallel()

	var (
		zeroSlice     []int
		zeroMap       map[string]int
		zeroInterface testFooer
		zeroChan      chan int
		zeroFunc      func()
	)

	filledChan := make(chan int, 10)
	filledChan <- 42

	cases := []struct {
		name                      string
		value                     any
		notZero, notNil, notEmpty bool
	}{
		{"Slice", []int{1, 2, 3},
			true, true, true},

		{"EmptySlice", []int{},
			true, true, false},

		{"ZeroSlice", zeroSlice,
			false, false, false},

		{"Map", map[string]int{"a": 1},
			true, true, true},

		{"EmptyMap", map[string]int{},
			true, true, false},

		{"ZeroMap", zeroMap,
			false, false, false},

		{"Interface", testFooer(&testBar{}),
			true, true, false},

		{"Interface of nil ptr", testFooer((*testBar)(nil)),
			false, false, false},

		{"ZeroInterface", zeroInterface,
			false, false, false},

		{"EmptyChan", make(chan int),
			true, true, false},

		{"EmptyBuffChan", make(chan int, 10),
			true, true, false},

		{"FilledChan", filledChan,
			true, true, true},

		{"ZeroChan", zeroChan,
			false, false, false},

		{"Func", func() {},
			true, true, false},

		{"ZeroFunc", zeroFunc,
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
