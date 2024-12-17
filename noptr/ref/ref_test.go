package ref_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sr9000/go-noptr/noptr/ref"
)

func BenchmarkToFromPtr(b *testing.B) {
	b.Run("Ref", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = *ref.Of(42).Ptr()
		}
	})

	b.Run("Bare Ptr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := 42
			ptr := &x
			_ = *ptr
		}
	})
}

func TestOf(t *testing.T) {
	t.Parallel()

	x := 42
	rx := ref.Of(x)
	require.NotNil(t, rx.Ptr())
	require.NotSame(t, &x, rx.Ptr())
	require.Equal(t, x, rx.Val())
}

func TestOfPtr_NotNil(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		x := 42
		rx, err := ref.OfPtr(&x)
		require.NoError(t, err)
		require.NotNil(t, rx.Ptr())
		require.Same(t, &x, rx.Ptr())
		require.Equal(t, x, rx.Val())
	})

	t.Run("nil cause err", func(t *testing.T) {
		t.Parallel()

		_, err := ref.OfPtr((*int)(nil))
		require.ErrorIs(t, err, ref.ErrPtrMustNotBeNil)
	})
}

func TestMust(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		x := 42
		rx := ref.Must(&x)
		require.NotNil(t, rx.Ptr())
		require.Same(t, &x, rx.Ptr())
		require.Equal(t, x, rx.Val())
	})

	t.Run("nil ptr", func(t *testing.T) {
		t.Parallel()

		require.PanicsWithError(t, ref.ErrPtrMustNotBeNil.Error(),
			func() {
				ref.Must[int](nil)
			})
	})
}
