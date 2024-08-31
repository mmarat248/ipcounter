package ipcounter

import (
	"awesomeProject/ipcounter/counters/bitmap"
)

type IPBitMap struct {
	bm *bitmap.BitMap
}

func (m *IPBitMap) Add(ip uint32) {
	m.bm.SetBit(ip, true)
}

func (m *IPBitMap) Count() uint64 {
	return m.bm.Count()
}
