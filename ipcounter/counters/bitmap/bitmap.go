package bitmap

import (
	"fmt"
)

const MaxSize uint32 = 1<<32 - 1

type BitMap struct {
	data        []byte
	cardinality uint32
	count       uint64
}

// New creates a new BitMap.
func New(size uint32) (*BitMap, error) {
	if size < 1 || size > MaxSize {
		return nil, fmt.Errorf("invalid size: %d, must be between 1 and %d", size, MaxSize)
	}
	arraySize := (size) >> 3
	if size&7 > 0 {
		arraySize++
	}

	return &BitMap{
		data:        make([]byte, arraySize), // "size" >> 3 = number of bytes, "+ 7" = get extra byte
		cardinality: size,
	}, nil
}

func (b *BitMap) GetBit(position uint32) bool {
	if position >= b.cardinality {
		return false
	}

	index, bit := position>>3, position&7
	return (b.data[index]>>(bit))&0x01 != 0
}

func (b *BitMap) SetBit(position uint32, value bool) bool {
	if position >= b.cardinality {
		return false
	}

	index, bit := position>>3, position&7 // = position / 8 , position % 8

	oldValue := b.data[index]
	if value {
		b.data[index] |= 1 << bit
	} else {
		b.data[index] ^= 1 << bit
	}

	if oldValue != b.data[index] {
		if value {
			b.count++
		} else {
			b.count--
		}
	}
	return true
}

func (b *BitMap) Size() uint32 {
	return b.cardinality
}

func (b *BitMap) Count() uint64 {
	return b.count
}

// BitIterator is an iterator for a BitMap that allows iterating over the set bits.
// It skips over zero bytes for efficiency.
type BitIterator struct {
	bitmap    *BitMap
	byteIndex uint32
	bitIndex  uint8
}

func (b *BitMap) Iterator() *BitIterator {
	return &BitIterator{
		bitmap: b,
	}
}

func (it *BitIterator) HasNext() bool {
	return it.byteIndex < uint32(len(it.bitmap.data))
}

func (it *BitIterator) Next() (uint32, bool) {
	for it.byteIndex < uint32(len(it.bitmap.data)) {
		if it.bitmap.data[it.byteIndex] == 0 {
			it.byteIndex++
			it.bitIndex = 0
			continue
		}

		for it.bitIndex < 8 {
			if (it.bitmap.data[it.byteIndex] & (1 << it.bitIndex)) != 0 {
				position := (it.byteIndex << 3) + uint32(it.bitIndex)
				it.bitIndex++
				if it.bitIndex == 8 {
					it.byteIndex++
					it.bitIndex = 0
				}
				return position, true
			}
			it.bitIndex++
		}

		it.byteIndex++
		it.bitIndex = 0
	}
	return 0, false
}
