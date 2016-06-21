package main

import (
	"fmt"
)

const BitmapSize = 1 << 32 // OS Pagesize

type Bitmap struct {
	data    []byte
	bitsize uint64
	maxpos  uint64
}

func NewBitmap() *Bitmap {
	return NewBitmapSize(BitmapSize)
}

func NewBitmapSize(size int) *Bitmap {
	if size == 0 || size > BitmapSize {
		size = BitmapSize
	} else if remainder := size % 8; remainder != 0 {
		size += 8 - remainder
	}
	return &Bitmap{data: make([]byte, size>>3), bitsize: uint64(size - 1)}
}

func (this *Bitmap) SetBit(offset uint64, value uint8) bool {
	index, pos := offset/8, offset%8
	if this.bitsize < offset {
		return false
	}
	if value == 0 {
		this.data[index] &^= 0x01 << pos
	} else {
		this.data[index] |= 0x01 << pos
		if this.maxpos < offset {
			this.maxpos = offset
		}
	}
	return true
}

func (this *Bitmap) GetBit(offset uint64) uint8 {
	index, pos := offset/8, offset%8
	if this.bitsize < offset {
		return 0
	}
	return (this.data[index] >> pos) & 0x01
}

func (this *Bitmap) Maxpos() uint64 {
	return this.maxpos
}

func (this *Bitmap) String() string {
	var maxTotal, bitTotal uint64 = 100, this.maxpos + 1
	if this.maxpos > maxTotal {
		bitTotal = maxTotal
	}
	numSlice := make([]uint64, 0, bitTotal)
	var offset uint64
	for offset = 0; offset < bitTotal; offset++ {
		if this.GetBit(offset) == 1 {
			numSlice = append(numSlice, offset)
		}
	}
	return fmt.Sprintf("%v", numSlice)
}
