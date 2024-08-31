package ipcounter

import (
	"fmt"
	"math/rand"
	"testing"
)

func generateRandomHashes(n int) []uint32 {
	hashes := make([]uint32, n)
	for i := 0; i < n; i++ {
		hashes[i] = rand.Uint32()
	}
	return hashes
}

func BenchmarkBitMapCount(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000, 1000000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("HyperLogLogBitMap-Size-%d", size), func(b *testing.B) {
			bm, _ := NewIPBitMap()
			hashes := generateRandomHashes(size)
			for _, hash := range hashes {
				bm.Add(hash)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bm.Count()
			}
		})
	}
}

func BenchmarkHyperLogLogCount(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000, 1000000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("HyperLogLogCount-Size-%d", size), func(b *testing.B) {
			hll, _ := NewHyperLogLog(14)
			hashes := generateRandomHashes(size)
			for _, hash := range hashes {
				hll.Add(hash)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				hll.Count()
			}
		})
	}
}

func BenchmarkHyperLogLogPlusCount(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000, 1000000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("HyperLogLogPlusCount-Size-%d", size), func(b *testing.B) {
			hll, _ := NewHyperLogLogPlus(14)
			hashes := generateRandomHashes(size)
			for _, hash := range hashes {
				hll.Add(hash)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				hll.Count()
			}
		})
	}
}

func BenchmarkHyperLogLogPlusBitMapCount(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000, 1000000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("HyperLogLogPlusCount-Size-%d", size), func(b *testing.B) {
			hll, _ := NewHyperLogLogPlusBitMap(14)
			hashes := generateRandomHashes(size)
			for _, hash := range hashes {
				hll.Add(hash)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				hll.Count()
			}
		})
	}
}
