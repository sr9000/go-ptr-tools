package opt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sr9000/go-noptr/noptr/opt"
	"github.com/sr9000/go-noptr/noptr/ptr"
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
		{"int", 42,
			true, false, false},

		{"zero int", zeroInt,
			false, false, false},

		{"string", "hello",
			true, false, true},

		{"empty string", "",
			false, false, false},

		{"zero string", zeroString,
			false, false, false},

		{"struct", testFoo{a: 1, b: 2},
			true, false, false},

		{"zero struct", zeroStruct,
			false, false, false},

		{"int ptr", ptr.Of(42),
			true, true, false},

		{"zero int ptr", new(int),
			true, true, false}, // int is zero BUT pointer is not zero (not nil)

		{"nil int ptr", nilIntPointer,
			false, false, false},

		{"zero struct ptr", new(testFoo),
			true, true, false}, // testFoo is zero BUT pointer is not zero (not nil)

		{"nil struct ptr", nilStructPointer,
			false, false, false},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			x := opt.Of(cs.value)
			assert.Equal(t, cs.notZero, nil != x.NotZero().Ptr(), "NotZero")
			assert.Equal(t, cs.notNil, nil != x.NotNil().Ptr(), "NotNil")
			assert.Equal(t, cs.notEmpty, nil != x.NotEmpty().Ptr(), "NotEmpty")
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
		{"slice", []int{1, 2, 3},
			true, true, true},

		{"empty slice", []int{},
			true, true, false},

		{"zero slice", zeroSlice,
			false, false, false},

		{"map", map[string]int{"a": 1},
			true, true, true},

		{"empty map", map[string]int{},
			true, true, false},

		{"zero map", zeroMap,
			false, false, false},

		{"interface", testFooer(&testBar{}),
			true, true, false},

		{"nil ptr interface", testFooer((*testBar)(nil)),
			false, false, false},

		{"zero interface", zeroInterface,
			false, false, false},

		{"empty chan", make(chan int),
			true, true, false},

		{"empty buff chan", make(chan int, 10),
			true, true, false},

		{"filled chan", filledChan,
			true, true, true},

		{"zero chan", zeroChan,
			false, false, false},

		{"func", func() {},
			true, true, false},

		{"zero func", zeroFunc,
			false, false, false},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			x := opt.Of(cs.value)
			assert.Equal(t, cs.notZero, nil != x.NotZero().Ptr(), "NotZero")
			assert.Equal(t, cs.notNil, nil != x.NotNil().Ptr(), "NotNil")
			assert.Equal(t, cs.notEmpty, nil != x.NotEmpty().Ptr(), "NotEmpty")
		})
	}
}
