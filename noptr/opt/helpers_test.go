package opt_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sr9000/go-noptr/noptr/opt"
	"github.com/sr9000/go-noptr/noptr/ptr"
	"github.com/sr9000/go-noptr/noptr/val"
)

func BenchmarkValidateInterface(b *testing.B) {
	bar := &testBar{}
	nilPtr := ptr.Nil[testBar]()
	str := "string"

	for range b.N / 4 {
		_ = opt.ParseInterface[testFooer](bar)
		_ = opt.ParseInterface[testFooer](nilPtr)
		_ = opt.ParseInterface[testFooer](nil)
		_ = opt.ParseInterface[testFooer](str)
	}
}

func TestWrap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   int
		cond    any
		isEmpty bool
	}{
		{"true", 42, true, false},
		{"false", 42, false, true},
		{"nil", 42, nil, false},
		{"string", 42, "hello", true},
		{"error", 42, errTest, true}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := opt.Wrap(tt.value, tt.cond)
			require.Equal(t, tt.isEmpty, result.Ptr() == nil)
		})
	}
}

func TestParseInterface(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    any
		expected *testFooer
	}{
		{"value receiver", testBar{}, ptr.Of((testFooer)(testBar{}))},
		{"pointer to value receiver", new(testBar), ptr.Of((testFooer)(new(testBar)))},
		{"pointer receiver", &testFoo{a: 1, b: 2}, ptr.Of((testFooer)(&testFoo{a: 1, b: 2}))},
		{"value of pointer receiver", testFoo{a: 1, b: 2}, nil},
		{"zero interface", val.Zero[testFooer](), nil},
		{"nil pointer", ptr.Nil[testBar](), nil},
		{"not an interface", "string", nil},
		{"nil value", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := opt.ParseInterface[testFooer](tt.value)
			require.Equal(t, tt.expected, result.Ptr())
		})
	}
}

func TestUnwrap2(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v1   int
		v2   string
		ok   bool
	}{
		{"ok", 42, "hello", true},
		{"first empty", 0, "hello", false},
		{"second empty", 42, "", false},
		{"both empty", 0, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			o1 := opt.Of(tt.v1).NotZero()
			o2 := opt.Of(tt.v2).NotZero()

			v1, v2, ok := opt.Unwrap2(o1, o2)
			require.Equal(t, tt.ok, ok)

			if ok {
				require.Equal(t, tt.v1, v1)
				require.Equal(t, tt.v2, v2)
			}
		})
	}
}

func TestUnwrap3(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		v1   int
		v2   string
		v3   float64
		ok   bool
	}{
		{"ok", 42, "hello", 3.14, true},
		{"first empty", 0, "hello", 3.14, false},
		{"second empty", 42, "", 3.14, false},
		{"third empty", 42, "hello", 0, false},
		{"all empty", 0, "", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			o1 := opt.Of(tt.v1).NotZero()
			o2 := opt.Of(tt.v2).NotZero()
			o3 := opt.Of(tt.v3).NotZero()

			v1, v2, v3, ok := opt.Unwrap3(o1, o2, o3)
			require.Equal(t, tt.ok, ok)

			if ok {
				require.Equal(t, tt.v1, v1)
				require.Equal(t, tt.v2, v2)
				require.InEpsilon(t, tt.v3, v3, 1.0e-6)
			}
		})
	}
}
