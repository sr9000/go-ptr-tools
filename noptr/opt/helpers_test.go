package opt_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sr9000/go-noptr/noptr/opt"
	"github.com/sr9000/go-noptr/pkg"
)

func testMapGetInline[K comparable, V any, M ~map[K]V](m M, k K) opt.Opt[V] {
	if v, ok := m[k]; ok {
		return opt.Of(v)
	} else {
		return opt.Empty[V]()
	}
}

func BenchmarkMapGet(b *testing.B) {
	n := 100_000
	m := make(map[int]int, n)
	for i := 0; i < n; i++ {
		m[i] = i
	}

	b.Run("Bare", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = m[i%(2*n)]
		}
	})

	b.Run("MapGet", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = opt.MapGet(m, i%(2*n))
		}
	})

	b.Run("Inline", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = testMapGetInline(m, i%(2*n))
		}
	})
}

func BenchmarkCastTo(b *testing.B) {
	bar := &testBar{}
	nul := (*testBar)(nil)
	str := "string"
	for i := 0; i < b.N/4; i++ {
		_ = opt.CastTo[testFooer](bar)
		_ = opt.CastTo[testFooer](nul)
		_ = opt.CastTo[testFooer](nil)
		_ = opt.CastTo[testFooer](str)
	}
}

func BenchmarkValidateInterface(b *testing.B) {
	bar := &testBar{}
	nul := (*testBar)(nil)
	str := "string"
	for i := 0; i < b.N/4; i++ {
		_ = opt.ValidateInterface[testFooer](bar)
		_ = opt.ValidateInterface[testFooer](nul)
		_ = opt.ValidateInterface[testFooer](nil)
		_ = opt.ValidateInterface[testFooer](str)
	}
}

func TestWrap(t *testing.T) {
	cases := []struct {
		name    string
		value   int
		cond    any
		isEmpty bool
	}{
		{"TrueCondition", 42, true, false},
		{"FalseCondition", 42, false, true},
		{"NilCondition", 42, nil, false},
		{"NonNilCondition", 42, "non-nil", true},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			result := opt.Wrap(cs.value, cs.cond)
			require.Equal(t, cs.isEmpty, result.Ptr() == nil)
		})
	}
}

func TestMapGet(t *testing.T) {
	cases := []struct {
		name     string
		m        map[int]string
		key      int
		expected *string
	}{
		{"ExistingKey", map[int]string{1: "one", 2: "two"}, 1, pkg.ToPtr("one")},
		{"NonExistingKey", map[int]string{1: "one", 2: "two"}, 3, nil},
		{"NilMap", nil, 1, nil},
		{"EmptyMap", map[int]string{}, 1, nil},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			result := opt.MapGet(cs.m, cs.key)
			require.Equal(t, cs.expected, result.Ptr())
		})
	}
}

func TestCastTo_Int(t *testing.T) {
	cases := []struct {
		name     string
		value    any
		expected *int
	}{
		{"ValidCast", 42, pkg.ToPtr(42)},
		{"InvalidCast", "string", nil},
		{"NilValue", nil, nil},
		{"NilPointer", (*int)(nil), nil},
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
	cases := []struct {
		name     string
		value    any
		expected *[]string
	}{
		{"ValidCast", []string{"foo", "bar"}, pkg.ToPtr([]string{"foo", "bar"})},
		{"InvalidCast", "string", nil},
		{"NilValue", nil, nil},
		{"NilPointer", (*[]string)(nil), nil},
		{"NilSlice", ([]string)(nil), pkg.ToPtr([]string(nil))},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			result := opt.CastTo[[]string](cs.value)
			require.Equal(t, cs.expected, result.Ptr())
		})
	}
}

func TestValidateInterface(t *testing.T) {
	cases := []struct {
		name     string
		value    any
		expected *testFooer
	}{
		{"ValueReceiverInterface", testBar{}, pkg.ToPtr((testFooer)(testBar{}))},
		{"PointerReceiverInterface", &testFoo{a: 1, b: 2}, pkg.ToPtr((testFooer)(&testFoo{a: 1, b: 2}))},
		{"NilInterface", (testFooer)(nil), nil},
		{"NilPointer", (*testBar)(nil), nil},
		{"NonInterfaceType", "string", nil},
		{"NilValue", nil, nil},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			t.Parallel()

			result := opt.ValidateInterface[testFooer](cs.value)
			require.Equal(t, cs.expected, result.Ptr())
		})
	}
}
