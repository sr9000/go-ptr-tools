package ref_test

import (
	"github.com/sr9000/go-noptr/ref"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkToFromPtr(b *testing.B) {
	b.Run("ref", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = *ref.New(42).Ptr()
		}
	})

	b.Run("bare ptr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			x := 42
			ptr := &x
			_ = *ptr
		}
	})
}

func TestNew(t *testing.T) {
	t.Parallel()

	x := 42
	rx := ref.New(x)
	require.NotNil(t, rx.Ptr())
	require.NotSame(t, &x, rx.Ptr())
	require.Equal(t, x, rx.Val())
}

func TestFrom(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		x := 42
		rx, err := ref.From(&x)
		require.NoError(t, err)
		require.NotNil(t, rx.Ptr())
		require.Same(t, &x, rx.Ptr())
		require.Equal(t, x, rx.Val())
	})

	t.Run("nil cause err", func(t *testing.T) {
		t.Parallel()

		_, err := ref.From[int](nil)
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

func TestGuaranteed(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		x := 42
		rx := ref.Guaranteed(&x)
		require.NotNil(t, rx.Ptr())
		require.Same(t, &x, rx.Ptr())
		require.Equal(t, x, rx.Val())
	})

	t.Run("nil ptr", func(t *testing.T) {
		t.Parallel()

		require.NotPanics(t, func() {
			rnil := ref.Guaranteed[int](nil)
			require.Nil(t, rnil.Ptr())
		})
	})
}
