package ringbuffer

import "testing"

func Benchmark_Mod1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = i % 1024
	}
}

func Benchmark_Mod2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = i % 7
	}
}

func Benchmark_And(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = i & 1023
	}
}
