package ref

import "testing"

func BenchmarkToFromPtr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Of(42).Ptr()
	}
}
