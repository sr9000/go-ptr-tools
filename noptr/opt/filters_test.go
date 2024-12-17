package opt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sr9000/go-noptr/noptr/opt"
	"github.com/sr9000/go-noptr/noptr/ptr"
	"github.com/sr9000/go-noptr/noptr/val"
)

func TestFilters_Empty(t *testing.T) {
	t.Parallel()

	x := opt.Empty[int]()
	assert.Nil(t, x.Ptr(), "Ptr")
	assert.Nil(t, x.NotZero().Ptr(), "NotZero")
	assert.Nil(t, x.NotNil().Ptr(), "NotNil")
	assert.Nil(t, x.NotEmpty().Ptr(), "NotEmpty")
}

func TestFilters_Primitives(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                      string
		value                     any
		notZero, notNil, notEmpty bool
	}{
		{"int", 42,
			true, false, false},

		{"zero int", val.Zero[int](),
			false, false, false},

		{"string", "hello",
			true, false, true},

		{"empty string", "",
			false, false, false},

		{"zero string", val.Zero[string](),
			false, false, false},

		{"struct", testFoo{a: 1, b: 2},
			true, false, false},

		{"zero struct", val.Zero[testFoo](),
			false, false, false},

		{"int ptr", ptr.Of(42),
			true, true, false},

		{"zero int ptr", new(int),
			true, true, false}, // int is zero BUT pointer is not zero (not nil)

		{"nil int ptr", ptr.Nil[int](),
			false, false, false},

		{"zero struct ptr", new(testFoo),
			true, true, false}, // testFoo is zero BUT pointer is not zero (not nil)

		{"nil struct ptr", ptr.Nil[testFoo](),
			false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			x := opt.Of(tt.value)
			assert.Equal(t, tt.notZero, nil != x.NotZero().Ptr(), "NotZero")
			assert.Equal(t, tt.notNil, nil != x.NotNil().Ptr(), "NotNil")
			assert.Equal(t, tt.notEmpty, nil != x.NotEmpty().Ptr(), "NotEmpty")
		})
	}
}

func TestFilters_Nillable(t *testing.T) {
	t.Parallel()

	filledChan := make(chan int, 10)
	filledChan <- 42

	tests := []struct {
		name                      string
		value                     any
		notZero, notNil, notEmpty bool
	}{
		{"slice", []int{1, 2, 3},
			true, true, true},

		{"empty slice", []int{},
			true, true, false},

		{"zero slice", val.Zero[[]int](),
			false, false, false},

		{"map", map[string]int{"a": 1},
			true, true, true},

		{"empty map", map[string]int{},
			true, true, false},

		{"zero map", val.Zero[map[string]int](),
			false, false, false},

		{"interface", testFooer(new(testBar)),
			true, true, false},

		{"nil ptr interface", testFooer(ptr.Nil[testBar]()),
			false, false, false},

		{"zero interface", val.Zero[testFooer](),
			false, false, false},

		{"empty chan", make(chan int),
			true, true, false},

		{"empty buff chan", make(chan int, 10),
			true, true, false},

		{"filled chan", filledChan,
			true, true, true},

		{"zero chan", val.Zero[chan int](),
			false, false, false},

		{"func", func() {},
			true, true, false},

		{"zero func", val.Zero[func()](),
			false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			x := opt.Of(tt.value)
			assert.Equal(t, tt.notZero, nil != x.NotZero().Ptr(), "NotZero")
			assert.Equal(t, tt.notNil, nil != x.NotNil().Ptr(), "NotNil")
			assert.Equal(t, tt.notEmpty, nil != x.NotEmpty().Ptr(), "NotEmpty")
		})
	}
}
