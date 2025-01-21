package ptr_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sr9000/go-ptr-tools/ptr"
	"github.com/sr9000/go-ptr-tools/ref"
)

func TestCoalesce(t *testing.T) {
	t.Parallel()

	one := ptr.New(1)
	two := ptr.New(2)
	three := ptr.New(3)

	tests := []struct {
		name     string
		expected *int
		pointers []*int
	}{
		{"single valid", one,
			[]*int{one}},

		{"single nil", nil,
			[]*int{nil}},

		{"two valid", one,
			[]*int{one, two}},

		{"three valid", one,
			[]*int{one, two, three}},

		{"first nil two valid", two,
			[]*int{nil, two, three}},

		{"second nil two valid", one,
			[]*int{one, nil, three}},

		{"third nil two valid", one,
			[]*int{one, two, nil}},

		{"first valid two nil", one,
			[]*int{one, nil, nil}},

		{"last valid two nil", three,
			[]*int{nil, nil, three}},

		{"three nil", nil,
			[]*int{nil, nil, nil}},

		{"no pointers", nil,
			[]*int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := ptr.Coalesce(tt.pointers...)
			require.Same(t, tt.expected, result)
		})
	}
}

func TestFinalize(t *testing.T) {
	t.Parallel()

	one := ptr.New(1)
	two := ptr.New(2)
	three := ptr.New(3)
	final := ref.Guaranteed(ptr.New(42))

	tests := []struct {
		name     string
		expected ref.Ref[int]
		final    ref.Ref[int]
		pointers []*int
	}{
		{"single valid", ref.Guaranteed(one), final,
			[]*int{one}},

		{"single nil", final, final,
			[]*int{nil}},

		{"two valid", ref.Guaranteed(one), final,
			[]*int{one, two}},

		{"three valid", ref.Guaranteed(one), final,
			[]*int{one, two, three}},

		{"first nil two valid", ref.Guaranteed(two), final,
			[]*int{nil, two, three}},

		{"second nil two valid", ref.Guaranteed(one), final,
			[]*int{one, nil, three}},

		{"third nil two valid", ref.Guaranteed(one), final,
			[]*int{one, two, nil}},

		{"first valid two nil", ref.Guaranteed(one), final,
			[]*int{one, nil, nil}},

		{"last valid two nil", ref.Guaranteed(three), final,
			[]*int{nil, nil, three}},

		{"three nil", final, final,
			[]*int{nil, nil, nil}},

		{"no pointers", final, final,
			[]*int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := ptr.Else(tt.final, tt.pointers...)
			require.Same(t, tt.expected.Ptr(), result.Ptr())
		})
	}
}
