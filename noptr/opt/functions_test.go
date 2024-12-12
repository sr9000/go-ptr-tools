package opt_test

import (
	"github.com/sr9000/go-noptr/noptr/opt"
	"github.com/sr9000/go-noptr/noptr/ptr"
	"github.com/stretchr/testify/require"
	"testing"
)

func BenchmarkCastTo(b *testing.B) {
	bar := &testBar{}
	nul := (*testBar)(nil)
	str := "string"

	for range b.N / 4 {
		_ = opt.CastTo[testFooer](bar)
		_ = opt.CastTo[testFooer](nul)
		_ = opt.CastTo[testFooer](nil)
		_ = opt.CastTo[testFooer](str)
	}
}

func TestCastTo_Int(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		value    any
		expected *int
	}{
		{"valid cast", 42, ptr.Of(42)},
		{"invalid cast", "string", nil},
		{"nil value", nil, nil},
		{"nil pointer", (*int)(nil), nil},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			result := opt.CastTo[int](cs.value)
			require.Equal(t, cs.expected, result.Ptr())
		})
	}
}

func TestCastTo_SliceOfString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		value    any
		expected *[]string
	}{
		{"valid cast", []string{"foo", "bar"}, ptr.Of([]string{"foo", "bar"})},
		{"invalid cast", "string", nil},
		{"nil value", nil, nil},
		{"nil pointer", (*[]string)(nil), nil},
		{"nil slice", ([]string)(nil), ptr.Of([]string(nil))},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			result := opt.CastTo[[]string](cs.value)
			require.Equal(t, cs.expected, result.Ptr())
		})
	}
}
