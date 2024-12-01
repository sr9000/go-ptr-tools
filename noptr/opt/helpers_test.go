package opt_test

import (
	"github.com/sr9000/go-noptr/noptr/opt"
	"testing"
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

// todo implement tests
