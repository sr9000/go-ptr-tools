package ref_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sr9000/go-ptr-tools/ref"
)

func BenchmarkToFromPtr(b *testing.B) {
	b.Run("ref", func(b *testing.B) {
		for i := range b.N / 2 {
			r := ref.Guaranteed(&i)
			*r.Ptr()++
			*r.Ptr()--
		}
	})

	b.Run("bare ptr", func(b *testing.B) {
		for i := range b.N / 2 {
			ptr := &i
			*ptr++
			*ptr--
		}
	})
}

func TestFromPtr(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		x := 42
		rx, err := ref.New(&x)
		require.NoError(t, err)
		require.NotNil(t, rx.Ptr())
		require.Same(t, &x, rx.Ptr())
		require.Equal(t, x, rx.Val())
	})

	t.Run("nil cause err", func(t *testing.T) {
		t.Parallel()

		_, err := ref.New[int](nil)
		require.ErrorIs(t, err, ref.ErrPtrMustBeNotNil)
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

func TestLiteral(t *testing.T) {
	t.Parallel()

	t.Run("literal", func(t *testing.T) {
		t.Parallel()

		rx := ref.Literal(42)
		require.NotNil(t, rx.Ptr())
		require.Equal(t, 42, rx.Val())
	})

	t.Run("value", func(t *testing.T) {
		t.Parallel()

		x := 42
		rx := ref.Literal(x)
		require.NotNil(t, rx.Ptr())
		require.NotSame(t, &x, rx.Ptr())
		require.Equal(t, x, rx.Val())
	})
}
