package opt_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sr9000/go-noptr/noptr/opt"
)

func TestPtr(t *testing.T) {
	t.Parallel()

	t.Run("value 42", func(t *testing.T) {
		t.Parallel()

		value := 42
		o := opt.Of(value)
		require.Equal(t, &value, o.Ptr())
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		o := opt.Empty[int]()
		require.Nil(t, o.Ptr())
	})
}

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("value 42", func(t *testing.T) {
		t.Parallel()

		value := 42
		o := opt.Of(value)
		v, ok := o.Get()
		require.True(t, ok)
		require.Equal(t, value, v)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		o := opt.Empty[int]()
		v, ok := o.Get()
		require.False(t, ok)
		require.Zero(t, v)
	})
}
