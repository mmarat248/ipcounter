package hyperloglogplus

import (
	"fmt"
	"math"
	"math/bits"
)

const two32 = 1 << 32

type HyperLogLogPlus struct {
	registers          []uint8 // Array of registers
	precision          uint8   // Precision (number of bits for addressing registers)
	numRegisters       uint32  // Number of registers (2^precision)
	IsSparse           bool
	SparseSet          map[uint32]bool
	SparseSetThreshold uint32
}

func New(precision uint8) (*HyperLogLogPlus, error) {
	if precision < 4 || precision > 16 {
		return nil, fmt.Errorf("invalid precision: %d, must be between 4 and 16", precision)
	}

	numRegisters := uint32(1 << precision)
	return &HyperLogLogPlus{
		registers:          make([]uint8, numRegisters),
		numRegisters:       numRegisters,
		precision:          precision,
		IsSparse:           true,
		SparseSet:          make(map[uint32]bool),
		SparseSetThreshold: uint32(float64(numRegisters) * 0.75),
	}, nil
}

func (h *HyperLogLogPlus) Add(hash uint32) {
	if h.IsSparse {
		h.SparseSet[hash] = true

		if uint32(len(h.SparseSet)) > h.SparseSetThreshold {
			h.IsSparse = false
			for k := range h.SparseSet {
				h.Add(k)
			}
			h.SparseSet = nil
		}
	} else {
		// Extract register address from the most significant bits of the hash
		registerMask := uint32(((1 << 32) - 1) << (32 - h.precision))
		registerIndex := (hash & registerMask) >> (32 - h.precision)

		// Clear the bits used for the register address
		remainingHash := hash ^ (hash & registerMask)

		// Count the number of leading zeros + 1
		leadingZeros := countTrailingRightZeros(remainingHash) + 1
		// Update the register if the new value is greater
		if leadingZeros > h.registers[registerIndex] {
			h.registers[registerIndex] = leadingZeros
		}
	}

}

func (h *HyperLogLogPlus) Count() uint64 {
	if h.IsSparse {
		return uint64(len(h.SparseSet))
	}

	estimate := calculateRawEstimate(h.registers)
	if estimate <= 2.5*float64(h.numRegisters) {
		// Use linear counting for small values
		zeroRegisters := countZeroRegisters(h.registers)
		if zeroRegisters != 0 {
			return uint64(linearCounting(h.numRegisters, zeroRegisters))
		}
		return uint64(estimate)
	} else if estimate < two32/30 {
		// Use raw estimate for medium values
		return uint64(estimate)
	}

	// Use correction for large values
	return uint64(-two32 * math.Log(1-estimate/two32))
}

func calculateRawEstimate(registers []uint8) float64 {
	sum := 0.0
	for _, val := range registers {
		sum += 1.0 / float64(uint64(1)<<val)
	}

	m := uint32(len(registers))
	return alpha(m) * float64(m) * float64(m) / sum
}

func linearCounting(m uint32, v uint32) float64 {
	return float64(m) * math.Log(float64(m)/float64(v))
}

func countZeroRegisters(registers []uint8) uint32 {
	var count uint32
	for _, v := range registers {
		if v == 0 {
			count++
		}
	}
	return count
}

func alpha(m uint32) float64 {
	switch m {
	case 16:
		return 0.673
	case 32:
		return 0.697
	case 64:
		return 0.709
	default:
		return 0.7213 / (1 + 1.079/float64(m))
	}
}

func countTrailingRightZeros(value uint32) uint8 {
	return uint8(bits.TrailingZeros32(value))
}
