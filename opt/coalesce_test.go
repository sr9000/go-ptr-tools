package opt_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sr9000/go-ptr-tools/opt"
)

func TestOptCoalesce(t *testing.T) {
	t.Parallel()

	one := 1
	two := 2
	three := 3

	tests := []struct {
		name     string
		expected opt.Opt[int]
		opts     []opt.Opt[int]
	}{
		{"single valid", opt.Literal(one),
			[]opt.Opt[int]{opt.Literal(one)}},

		{"single nil", opt.Opt[int]{},
			[]opt.Opt[int]{{}}},

		{"two valid", opt.Literal(one),
			[]opt.Opt[int]{opt.Literal(one), opt.Literal(two)}},

		{"three valid", opt.Literal(one),
			[]opt.Opt[int]{opt.Literal(one), opt.Literal(two), opt.Literal(three)}},

		{"first nil two valid", opt.Literal(two),
			[]opt.Opt[int]{{}, opt.Literal(two), opt.Literal(three)}},

		{"second nil two valid", opt.Literal(one),
			[]opt.Opt[int]{opt.Literal(one), {}, opt.Literal(three)}},

		{"third nil two valid", opt.Literal(one),
			[]opt.Opt[int]{opt.Literal(one), opt.Literal(two), {}}},

		{"first valid two nil", opt.Literal(one),
			[]opt.Opt[int]{opt.Literal(one), {}, {}}},

		{"last valid two nil", opt.Literal(three),
			[]opt.Opt[int]{{}, {}, opt.Literal(three)}},

		{"three nil", opt.Opt[int]{},
			[]opt.Opt[int]{{}, {}, {}}},

		{"no opts", opt.Opt[int]{},
			[]opt.Opt[int]{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := opt.Coalesce(tt.opts...)
			expVal, expOk := tt.expected.Get()
			resVal, resOk := result.Get()
			require.Equal(t, expOk, resOk)
			require.Equal(t, expVal, resVal)
		})
	}
}

func TestFinalize(t *testing.T) {
	t.Parallel()

	one := opt.Literal(1)
	two := opt.Literal(2)
	three := opt.Literal(3)
	final := 42

	tests := []struct {
		name     string
		expected int
		final    int
		opts     []opt.Opt[int]
	}{
		{"single valid", 1, final,
			[]opt.Opt[int]{one}},

		{"single nil", final, final,
			[]opt.Opt[int]{{}}},

		{"two valid", 1, final,
			[]opt.Opt[int]{one, two}},

		{"three valid", 1, final,
			[]opt.Opt[int]{one, two, three}},

		{"first nil two valid", 2, final,
			[]opt.Opt[int]{{}, two, three}},

		{"second nil two valid", 1, final,
			[]opt.Opt[int]{one, {}, three}},

		{"third nil two valid", 1, final,
			[]opt.Opt[int]{one, two, {}}},

		{"first valid two nil", 1, final,
			[]opt.Opt[int]{one, {}, {}}},

		{"last valid two nil", 3, final,
			[]opt.Opt[int]{{}, {}, three}},

		{"three nil", final, final,
			[]opt.Opt[int]{{}, {}, {}}},

		{"no opts", final, final,
			[]opt.Opt[int]{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := opt.Else(tt.final, tt.opts...)
			require.Equal(t, tt.expected, result)
		})
	}
}
