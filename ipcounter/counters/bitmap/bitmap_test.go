package bitmap

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		size    uint32
		wantErr bool
	}{
		{"Valid size", 1000, false},
		{"Minimum size", 1, false},
		{"Maximum size", MaxSize, false},
		{"Size zero", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm, err := New(tt.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && bm.Size() != tt.size {
				t.Errorf("New() size = %v, want %v", bm.Size(), tt.size)
			}
		})
	}
}

func TestBitMap_SetBit(t *testing.T) {
	bm, _ := New(100)

	tests := []struct {
		name     string
		position uint32
		value    bool
		want     bool
	}{
		{"Set bit 0", 0, true, true},
		{"Set bit 50", 50, true, true},
		{"Set bit 99", 99, true, true},
		{"Set bit out of range", 100, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bm.SetBit(tt.position, tt.value); got != tt.want {
				t.Errorf("BitMap.SetBit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitMap_GetBit(t *testing.T) {
	bm, _ := New(100)
	bm.SetBit(0, true)
	bm.SetBit(50, true)
	bm.SetBit(99, true)

	tests := []struct {
		name     string
		position uint32
		want     bool
	}{
		{"Get bit 0", 0, true},
		{"Get bit 1", 1, false},
		{"Get bit 50", 50, true},
		{"Get bit 98", 98, false},
		{"Get bit 99", 99, true},
		{"Get bit out of range", 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bm.GetBit(tt.position); got != tt.want {
				t.Errorf("BitMap.GetBit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBitMap_Size(t *testing.T) {
	sizes := []uint32{1, 8, 100, 1000, MaxSize}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("Size %d", size), func(t *testing.T) {
			bm, _ := New(size)
			if got := bm.Size(); got != size {
				t.Errorf("BitMap.Size() = %v, want %v", got, size)
			}
		})
	}
}

func TestBitMap_SetAndGetBit(t *testing.T) {
	bm, _ := New(1000)

	// Set some bits
	bm.SetBit(0, true)
	bm.SetBit(1, false)
	bm.SetBit(500, true)
	bm.SetBit(999, true)

	// Check the bits
	if !bm.GetBit(0) {
		t.Errorf("Expected bit 0 to be set")
	}
	if bm.GetBit(1) {
		t.Errorf("Expected bit 1 to be unset")
	}
	if !bm.GetBit(500) {
		t.Errorf("Expected bit 500 to be set")
	}
	if !bm.GetBit(999) {
		t.Errorf("Expected bit 999 to be set")
	}

	// Check a bit that wasn't set
	if bm.GetBit(750) {
		t.Errorf("Expected bit 750 to be unset")
	}
}

func TestBitMap_Count(t *testing.T) {
	bm, _ := New(1000)

	// Test initial count
	if count := bm.Count(); count != 0 {
		t.Errorf("Initial count should be 0, got %d", count)
	}

	// Set some bits and check count
	bm.SetBit(0, true)
	bm.SetBit(1, true)
	bm.SetBit(500, true)
	if count := bm.Count(); count != 3 {
		t.Errorf("Count should be 3, got %d", count)
	}

	// Unset a bit and check count
	bm.SetBit(1, false)
	if count := bm.Count(); count != 2 {
		t.Errorf("Count should be 2, got %d", count)
	}

	// Set an already set bit and check count
	bm.SetBit(0, true)
	if count := bm.Count(); count != 2 {
		t.Errorf("Count should still be 2, got %d", count)
	}

	// Set many bits and check count
	for i := uint32(0); i < 1000; i++ {
		bm.SetBit(i, true)
	}
	if count := bm.Count(); count != 1000 {
		t.Errorf("Count should be 1000, got %d", count)
	}

	// Unset all bits and check count
	for i := uint32(0); i < 1000; i++ {
		bm.SetBit(i, false)
	}
	if count := bm.Count(); count != 0 {
		t.Errorf("Final count should be 0, got %d", count)
	}
}
