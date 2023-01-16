package benchmark_test

import (
	"testing"
)

func X(i int) bool {
	return i == 0
}

func Y(i any) bool {
	return i == 0
}

func Z(i any) bool {
	return i.(int) == 0
}

func BenchmarkA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		X(i)
	}
}

func BenchmarkB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Y(i)
	}
}

func BenchmarkC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Z(i)
	}
}
