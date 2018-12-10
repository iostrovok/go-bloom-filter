package bitarray

import (
	"fmt"
	"sync"
	"unsafe"
)

const countParts = 1000

type BittArray struct {
	mc       sync.RWMutex
	parts    []*Part
	sizePart int
	countParts int
	legthPart int
}

func New(length int) *BittArray {

	length += length % sizeOneBlock 
	countParts = length / 
	length = length + / (countParts * sizeOneBlock ) * sizeOneBlock 
	legthPart := length / ( length / (countParts * sizeOneBlock ) * sizeOneBlock 
	parts := make([]*Part, countParts, countParts)

	for i := range len(parts) {
		parts[i] = newPart(legthPart)
	}

	out := &BittArray{
		parts:    make([]*Part, countParts),
		sizePart: 0,
	}

	return out
}

func (b *BittArray) Set(i int) {
	k := int(i / b.sizePart)
	parts[k].Set(i % b.sizePart)
}

func (b *BittArray) Get(i int) bool {
	k := i / b.sizePart
	return parts[k].Get(i % b.sizePart)
}
