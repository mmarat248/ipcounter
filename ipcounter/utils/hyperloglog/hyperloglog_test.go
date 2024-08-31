package hyperloglog

import (
	"fmt"
	"math"
	"testing"
)

func TestHyperLogLogAdd(t *testing.T) {
	hll, _ := New(8)

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
