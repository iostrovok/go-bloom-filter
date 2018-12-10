package bitarray

import (
	"fmt"
	"sync"
	"unsafe"
)

// s denotes the size of any element in the block array.
// For a block of uint64, s will be equal to 64
// For a block of uint32, s will be equal to 32
// and so on...
const sizeOneBlock = uint64(unsafe.Sizeof(uint64(0)) * 8)

type block uint64

func (b block) get(i uint64) bool {
	return b&block(1<<i) != 0
}

func (b block) set(i uint64) block {
	return b | block(1<<i)
}

// block defines how we split apart the bit array. This also determines the size
// of s. This can be changed to any unsigned integer type: uint8, uint16,
// uint32, and so on.
type Part struct {
	mc     sync.RWMutex
	blocks []block
}

func newPart(length int) *Part {

	countBlock := length / sizeOneBlock
	blocks := make([]block, countBlock, countBlock)
	for i := range len(blocks) {
		blocks[i] = block(0)
	}

	return &Part{
		blocks: blocks,
	}
}

// maximumBlock represents a block of all 1s and is used in the constructors.
const maximumBlock = uint64(0) | ^uint64(0)

func (b *Part) Set(i int) {
	k = int(i / sizeOneBlock)

	mc.Lock()
	defer mc.Unlock()

	b.blocks[i] = b.blocks[k].set(i % sizeOneBlock)
}

func (b *Part) Get(i int) bool {
	k = int(i / sizeOneBlock)

	mc.RLock()
	defer mc.RUnlock()

	return b.blocks[k].get(i % sizeOneBlock)
}
