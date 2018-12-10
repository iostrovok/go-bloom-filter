package bloomfilter

import (
	"encoding/binary"
	"math"
)

func packInt(val int) []byte {
	buf := make([]byte, 4, 4)
	binary.LittleEndian.PutUint16(buf, uint16(val))
	return buf
}

func packFloat(val float64) []byte {
	buf := make([]byte, 8, 8)
	binary.LittleEndian.PutUint64(buf, math.Float64bits(val))
	return buf
}

func unpack(chunkSize int, b []byte) []uint64 {

	out := make([]uint64, len(b)/chunkSize, len(b)/chunkSize)

	k := 0
	for i := 0; i < len(b); i += chunkSize {

		switch chunkSize {
		case 8:
			out[k] = binary.LittleEndian.Uint64(b[i : i+chunkSize])
		case 4:
			out[k] = uint64(binary.LittleEndian.Uint32(b[i : i+chunkSize]))
		case 2:
			out[k] = uint64(binary.LittleEndian.Uint16(b[i : i+chunkSize]))
		}

		k++
	}

	return out
}
