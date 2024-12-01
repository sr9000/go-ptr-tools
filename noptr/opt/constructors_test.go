package opt_test

import (
	"github.com/sr9000/go-noptr/noptr/opt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmpty(t *testing.T) {
	t.Run("EmptyInt", func(t *testing.T) {
		t.Parallel()
		optVal := opt.Empty[int]()
		require.Nil(t, optVal.Ptr())
	})

	t.Run("EmptyString", func(t *testing.T) {
		t.Parallel()
		optValStr := opt.Empty[string]()
		require.Nil(t, optValStr.Ptr())
	})
}

func TestOf(t *testing.T) {
	t.Run("Int", func(t *testing.T) {
		t.Parallel()
		val := 42
		optVal := opt.Of(val)
		require.NotNil(t, optVal.Ptr())
		require.Equal(t, val, *optVal.Ptr())
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		strVal := "hello"
		optValStr := opt.Of(strVal)
		require.NotNil(t, optValStr.Ptr())
		require.Equal(t, strVal, *optValStr.Ptr())
	})
}

func TestOfPtr(t *testing.T) {
	t.Run("Int", func(t *testing.T) {
		t.Parallel()
		val := 42
		optVal := opt.OfPtr(&val)
		require.NotNil(t, optVal.Ptr())
		require.Equal(t, &val, optVal.Ptr())
		require.Equal(t, val, *optVal.Ptr())
	})

	t.Run("NilInt", func(t *testing.T) {
		t.Parallel()
		nilOptVal := opt.OfPtr[int](nil)
		require.Nil(t, nilOptVal.Ptr())
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		strVal := "hello"
		optValStr := opt.OfPtr(&strVal)
		require.NotNil(t, optValStr.Ptr())
		require.Equal(t, &strVal, optValStr.Ptr())
		require.Equal(t, strVal, *optValStr.Ptr())
	})

	t.Run("NilString", func(t *testing.T) {
		t.Parallel()
		nilOptValStr := opt.OfPtr[string](nil)
		require.Nil(t, nilOptValStr.Ptr())
	})
}
