package hyperloglogplus

import (
	"fmt"
	"math"
	"testing"
)

func TestHyperLogLogPlusAddSparse(t *testing.T) {
	hll, _ := New(8)

	testCases := []struct {
		hash     uint32
		expected bool
	}{
		{uint32(math.Pow(2, 25) + 2), true},
		{uint32(math.Pow(2, 26) + 8), true},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Hash%d", tc.hash), func(t *testing.T) {
			hll.Add(tc.hash)
			if _, ok := hll.SparseSet[tc.hash]; ok != tc.expected {
				t.Errorf("Expected sparseSet to contain %d: %v, got: %v", tc.hash, tc.expected, ok)
			}
		})
	}

	zeroCount := 0
	for _, v := range hll.registers {
		if v == 0 {
			zeroCount++
		}
	}
	expectedZeroCount := int(math.Pow(2, 8))
	if zeroCount != expectedZeroCount {
		t.Errorf("Expected %d zero registers, got: %d", expectedZeroCount, zeroCount)
	}
}

func TestHyperLogLogPlusAddNotSparse(t *testing.T) {
	hll, _ := New(8)
	hll.IsSparse = false

	testCases := []struct {
		hash             uint32
		expectedRegister int
		expectedValue    uint8
	}{
		{uint32(math.Pow(2, 25) + 2), 2, 2},
		{uint32(math.Pow(2, 26) + 8), 4, 4},
		{uint32(math.Pow(2, 27) + 16), 8, 5},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Hash%d", tc.hash), func(t *testing.T) {
			hll.Add(tc.hash)
			if hll.registers[tc.expectedRegister] != tc.expectedValue {
				t.Errorf("Expected register[%d] to be %d, got: %d", tc.expectedRegister, tc.expectedValue, hll.registers[tc.expectedRegister])
			}
		})
	}
}

func TestHyperLogLogPlusCount(t *testing.T) {
	hll, _ := New(8)

	testCases := []struct {
		hash     uint32
		expected uint64
	}{
		{uint32(math.Pow(2, 25) + 2), 1},
		{uint32(math.Pow(2, 26) + 4), 2},
		{uint32(math.Pow(2, 27) + 8), 3},
		{uint32(math.Pow(2, 27) + 8), 3}, // Duplicate, should not increase count
	}

	for i, tc := range testCases {
		hll.Add(tc.hash)
		count := hll.Count()
		if count != tc.expected {
			t.Errorf("After adding %d elements, expected count: %d, got: %d", i+1, tc.expected, count)
		}
	}
}
