package util

import "testing"

func BenchmarkPick(b *testing.B) {
	var num int

	for i := 0; i < b.N; i++ {
		num = Pick(true, 1, 2)
		_ = num
	}
}

func BenchmarkIf(b *testing.B) {
	var num int

	for i := 0; i < b.N; i++ {
		if true {
			num = 1
		} else {
			num = 2
		}

		_ = num
	}
}
