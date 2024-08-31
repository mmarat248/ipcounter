package ipcounter

import (
	"testing"
)

func BenchmarkSeqBitMapCounter(b *testing.B) {
	mp, _ := NewIPBitMap()
	ic := NewIPCounter(mp, false, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ic.CountIPFromFile("./ipsbig")
	}
}

func BenchmarkParallelBitMapCounter(b *testing.B) {
	mp, _ := NewIPBitMap()
	ic := NewIPCounter(mp, true, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ic.CountIPFromFile("./ipsbig")
	}
}

func BenchmarkParallelHyperLogLog(b *testing.B) {
	mp, _ := NewHyperLogLog(14)
	ic := NewIPCounter(mp, true, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ic.CountIPFromFile("./ipsbig")
	}
}

func BenchmarkParallelHyperLogLogPlus(b *testing.B) {
	mp, _ := NewHyperLogLogPlus(14)
	ic := NewIPCounter(mp, true, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ic.CountIPFromFile("./ipsbig")
	}
}

func BenchmarkParallelHyperLogLogPlusBitMpa(b *testing.B) {
	mp, _ := NewHyperLogLogPlusBitMap(14)
	ic := NewIPCounter(mp, true, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ic.CountIPFromFile("./ipsbig")
	}
}
