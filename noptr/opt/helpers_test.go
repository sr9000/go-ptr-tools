package opt_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sr9000/go-noptr/noptr/opt"
	"github.com/sr9000/go-noptr/noptr/ptr"
)

func BenchmarkValidateInterface(b *testing.B) {
	bar := &testBar{}
	nul := (*testBar)(nil)
	str := "string"

	for range b.N / 4 {
		_ = opt.ParseInterface[testFooer](bar)
		_ = opt.ParseInterface[testFooer](nul)
		_ = opt.ParseInterface[testFooer](nil)
		_ = opt.ParseInterface[testFooer](str)
	}
}

func TestWrap(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		value   int
		cond    any
		isEmpty bool
	}{
		{"true", 42, true, false},
		{"false", 42, false, true},
		{"nil", 42, nil, false},
		{"not nil", 42, "non-nil", true},
		{"error", 42, errors.New("error"), true}}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			result := opt.Wrap(cs.value, cs.cond)
			require.Equal(t, cs.isEmpty, result.Ptr() == nil)
		})
	}
}

func TestParseInterface(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		value    any
		expected *testFooer
	}{
		{"value receiver interface", testBar{}, ptr.Of((testFooer)(testBar{}))},
		{"pointer receiver interface", &testFoo{a: 1, b: 2}, ptr.Of((testFooer)(&testFoo{a: 1, b: 2}))},
		{"nil interface", (testFooer)(nil), nil},
		{"nil pointer", (*testBar)(nil), nil},
		{"not interface type", "string", nil},
		{"nil value", nil, nil},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			result := opt.ParseInterface[testFooer](cs.value)
			require.Equal(t, cs.expected, result.Ptr())
		})
	}
}

func TestUnwrap2(t *testing.T) {
	t.Parallel()

	cases := []struct {
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

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			o1 := opt.Of(cs.v1).NotZero()
			o2 := opt.Of(cs.v2).NotZero()

			v1, v2, ok := opt.Unwrap2(o1, o2)
			require.Equal(t, cs.ok, ok)
			if ok {
				require.Equal(t, cs.v1, v1)
				require.Equal(t, cs.v2, v2)
			}
		})
	}
}

func TestUnwrap3(t *testing.T) {
	t.Parallel()

	cases := []struct {
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

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			o1 := opt.Of(cs.v1).NotZero()
			o2 := opt.Of(cs.v2).NotZero()
			o3 := opt.Of(cs.v3).NotZero()

			v1, v2, v3, ok := opt.Unwrap3(o1, o2, o3)
			require.Equal(t, cs.ok, ok)
			if ok {
				require.Equal(t, cs.v1, v1)
				require.Equal(t, cs.v2, v2)
				require.Equal(t, cs.v3, v3)
			}
		})
	}
}
