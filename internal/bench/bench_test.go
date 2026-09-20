package bench

import "testing"

func BenchmarkPlaceholder(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {

	}
}
