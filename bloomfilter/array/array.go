package array

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"sync"
	// "unsafe"
)

// const sizeOneByte = uint64(unsafe.Sizeof(uint8(0)) * 8)
const sizeOneByte = uint64(8)

func init() {
	if strconv.IntSize != 64 {
		panic("Program works on 64 bits system only")
	}
}

// Array is main structure
type Array struct {
	mc sync.RWMutex

	bArray      []byte `json:"array"`
	Length      uint64 `json:"length"`
	SizeOneByte uint64
}

// New is constructor
func New(length uint64) *Array {

	l := length / sizeOneByte
	if l*sizeOneByte < length {
		l++
	}

	out := &Array{
		SizeOneByte: sizeOneByte,
		bArray:      make([]byte, l, l),
		Length:      length,
	}

	return out
}

// Set adds new point to array
func (b *Array) Set(i uint64) {
	j := int(i / sizeOneByte)
	k := uint8(1 << (i % sizeOneByte))

	b.mc.Lock()
	defer b.mc.Unlock()

	b.bArray[j] = b.bArray[j] | k
}

// Get return true if point is found in array
func (b *Array) Get(i uint64) bool {
	j := int(i / sizeOneByte)
	k := uint8(1 << (i % sizeOneByte))

	b.mc.RLock()
	defer b.mc.RUnlock()

	return b.bArray[j]&k != 0
}

// Merge adds values from outside array into current
func (b *Array) Merge(a *Array) error {

	b.mc.Lock()
	defer b.mc.Unlock()

	for i := range a.bArray {
		b.bArray[i] |= a.bArray[i]
	}

	return nil
}

// Compare checks characteristics of arrays
func (b *Array) Compare(a *Array) error {
	if a.SizeOneByte != b.SizeOneByte {
		return fmt.Errorf("Wrong SizeOneByte for Array: %d != %d", a.SizeOneByte, b.SizeOneByte)
	}

	if a.Length != b.Length {
		return fmt.Errorf("Wrong length for Array: %d != %d", a.Length, b.Length)
	}

	if len(b.bArray) != len(b.bArray) {
		return fmt.Errorf("Wrong length for byteArray: %d != %d", len(b.bArray), len(a.bArray))
	}

	return nil
}

// ToBytes save internal byte array to buffer
func (b *Array) ToBytes(binBuf *bytes.Buffer) error {
	b.mc.RLock()
	defer b.mc.RUnlock()

	_, err := binBuf.Write(b.bArray)

	return err
}

// Read read internal byte array from buffer
func (b *Array) Read(reader *bufio.Reader, length int64) error {

	b.mc.Lock()
	defer b.mc.Unlock()

	readByLength := false
	if length == 0 {
		readByLength = true
		length = int64(b.Length / sizeOneByte)
		if b.Length*sizeOneByte > uint64(length) {
			length++
		}
	}

	maxLength := int64(1024)
	array := make([]byte, maxLength, maxLength)

	var err error
	b.bArray = []byte{}

	var n int
	total := 0
	for err == nil {

		if maxLength > length-int64(total) {
			array = make([]byte, length-int64(total), length-int64(total))
		}

		n, err = reader.Read(array)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			return nil
		}

		b.bArray = append(b.bArray, array[:n]...)
		if err == io.EOF {
			break
		}
		total += n

		if readByLength && int64(total) >= length {
			break
		}
	}

	return nil
}
