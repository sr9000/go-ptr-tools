package opt_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sr9000/go-noptr/noptr/opt"
	"github.com/sr9000/go-noptr/noptr/ptr"
)

func testOfMapInline[K comparable, V any, M ~map[K]V](m M, k K) opt.Opt[V] {
	if v, ok := m[k]; ok {
		return opt.Of(v)
	}

	return opt.Empty[V]()
}

func BenchmarkOfMap(b *testing.B) {
	mapSize := 100_000
	mapOfInts := make(map[int]int, mapSize)

	for i := range mapSize {
		mapOfInts[i] = i
	}

	b.Run("bare", func(b *testing.B) {
		for i := range b.N {
			_, ignored := mapOfInts[i%(2*mapSize)]
			_ = ignored
		}
	})

	b.Run("OfMap", func(b *testing.B) {
		for i := range b.N {
			_ = opt.OfMap(mapOfInts, i%(2*mapSize))
		}
	})

	b.Run("inline", func(b *testing.B) {
		for i := range b.N {
			_ = testOfMapInline(mapOfInts, i%(2*mapSize))
		}
	})
}

func TestEmpty(t *testing.T) {
	t.Parallel()

	t.Run("empty int", func(t *testing.T) {
		t.Parallel()

		optVal := opt.Empty[int]()
		require.Nil(t, optVal.Ptr())
	})

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()

		optValStr := opt.Empty[string]()
		require.Nil(t, optValStr.Ptr())
	})
}

func TestOf(t *testing.T) {
	t.Parallel()

	t.Run("int", func(t *testing.T) {
		t.Parallel()

		val := 42
		optVal := opt.Of(val)
		require.NotNil(t, optVal.Ptr())
		require.Equal(t, val, *optVal.Ptr())
	})

	t.Run("string", func(t *testing.T) {
		t.Parallel()

		strVal := "hello"
		optValStr := opt.Of(strVal)
		require.NotNil(t, optValStr.Ptr())
		require.Equal(t, strVal, *optValStr.Ptr())
	})
}

func TestOfPtr(t *testing.T) {
	t.Parallel()

	t.Run("int", func(t *testing.T) {
		t.Parallel()

		val := 42
		optVal := opt.OfPtr(&val)
		require.NotNil(t, optVal.Ptr())
		require.Equal(t, &val, optVal.Ptr())
		require.Equal(t, val, *optVal.Ptr())
	})

	t.Run("nil int ptr", func(t *testing.T) {
		t.Parallel()

		nilOptVal := opt.OfPtr[int](nil)
		require.Nil(t, nilOptVal.Ptr())
	})

	t.Run("string", func(t *testing.T) {
		t.Parallel()

		strVal := "hello"
		optValStr := opt.OfPtr(&strVal)
		require.NotNil(t, optValStr.Ptr())
		require.Equal(t, &strVal, optValStr.Ptr())
		require.Equal(t, strVal, *optValStr.Ptr())
	})

	t.Run("nil string ptr", func(t *testing.T) {
		t.Parallel()

		nilOptValStr := opt.OfPtr[string](nil)
		require.Nil(t, nilOptValStr.Ptr())
	})
}

func TestOfMap(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		m        map[int]string
		key      int
		expected *string
	}{
		{"existing key", map[int]string{1: "one", 2: "two"}, 1, ptr.Of("one")},
		{"non existing key", map[int]string{1: "one", 2: "two"}, 3, nil},
		{"nil map", nil, 1, nil},
		{"empty map", map[int]string{}, 1, nil},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			result := opt.OfMap(cs.m, cs.key)
			require.Equal(t, cs.expected, result.Ptr())
		})
	}
}
