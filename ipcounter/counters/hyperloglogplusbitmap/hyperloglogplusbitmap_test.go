package hyperloglogplusbitmap

import (
	"awesomeProject/ipcounter/utils/murmur3"
	"fmt"
	"math"
	"testing"
)

func TestHyperLogLogPlusIpBitMapAddSparse(t *testing.T) {
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
			if ok := hll.sparseSet.GetBit(tc.hash); ok != tc.expected {
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

func TestHyperLogLogPlusIpBitMapAddNotSparse(t *testing.T) {
	hll, _ := New(8)
	hll.isSparse = false

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

func TestHyperLogLogPlusIpBitMapCount(t *testing.T) {
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

func TestHyperLogLogPlusBitMapThresholdTransition(t *testing.T) {
	precision := uint8(4) // Use low precision for quick threshold reach
	hll, err := New(precision)
	if err != nil {
		t.Fatalf("Failed to create HyperLogLogPlusBitMap: %v", err)
	}
	// Check that the initial state is sparse
	if !hll.isSparse {
		t.Errorf("Expected initial state to be sparse")
	}

	// Add elements until reaching the threshold
	var expectedCount uint64
	for i := uint32(0); i < uint32(hll.sparseSetThreshold)+1; i++ {
		hll.Add(murmur3.Sum32([]byte(fmt.Sprintf("%d", i))))
		expectedCount++
		if i < uint32(hll.sparseSetThreshold) {
			if !hll.isSparse {
				t.Errorf("Expected to remain in sparse mode at iteration %d", i)
			}
		} else {
			if hll.isSparse {
				t.Errorf("Expected to switch to dense mode at iteration %d", i)
			}
		}
	}

	// Check that registers contain non-zero values
	zeroCount := 0
	for _, v := range hll.registers {
		if v == 0 {
			zeroCount++
		}
	}
	if zeroCount == len(hll.registers) {
		t.Errorf("Expected some non-zero registers after transition")
	}
}
