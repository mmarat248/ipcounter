package ipcounter

import (
	"awesomeProject/ipcounter/counters/bitmap"
	"awesomeProject/ipcounter/counters/hyperloglog"
	"awesomeProject/ipcounter/counters/hyperloglogplus"
	"awesomeProject/ipcounter/counters/hyperloglogplusbitmap"
)

type IPMap interface {
	Add(ip uint32)
	Count() uint64
}

func NewIPBitMap() (IPMap, error) {
	bm, err := bitmap.New(bitmap.MaxSize)
	if err != nil {
		return nil, err
	}
	return &IPBitMap{
		bm: bm,
	}, nil
}

func NewSet() (IPMap, error) {
	return &IPSet{
		set: make(map[uint32]struct{}, bitmap.MaxSize),
	}, nil
}

func NewHyperLogLog(precision uint8) (IPMap, error) {
	return hyperloglog.New(precision)
}

func NewHyperLogLogPlus(precision uint8) (IPMap, error) {
	return hyperloglogplus.New(precision)
}

func NewHyperLogLogPlusBitMap(precision uint8) (IPMap, error) {
	return hyperloglogplusbitmap.New(precision)
}
