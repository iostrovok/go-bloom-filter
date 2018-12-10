package bloomfilter

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"unsafe"

	"github.com/iostrovok/go-bloom-filter/bloomfilter/array"
)

var log2Const float64
var capacityError error

func init() {
	log2Const = math.Log(2) * math.Log(2)
	capacityError = fmt.Errorf("BloomFilter is at capacity")
}

// --------------------------------------------------------

func sum(mark string, b []byte) []byte {

	switch mark {
	case "sha512":
		t := sha512.Sum512(b)
		return t[:]
	case "sha384":
		t := sha512.Sum384(b)
		return t[:]
	case "sha256":
		t := sha256.Sum256(b)
		return t[:]
	case "sha1":
		t := sha1.Sum(b)
		return t[:]
	case "md5":
		t := md5.Sum(b)
		return t[:]
	default:
		// nothing. We will never be here.
	}

	return []byte{}
}

// BloomFilter is a structure for scalable bloom filter.
type BloomFilter struct {
	errorRate    float64
	numSlices    int
	bitsPerSlice uint64
	capacity     int64
	numBits      uint64
	count        int64
	chunkSize    int
	hashfnname   string

	saltFunctions []*hashFN

	bitarray *array.Array
}

// New is constructor. It checks parameters and creates new bloom filter.
func New(capacity int64, errorRates ...float64) (*BloomFilter, error) {

	errorRate := float64(0.001)
	if len(errorRates) > 0 {
		errorRate = errorRates[0]
	}

	if errorRate <= 0 || 1.0 < errorRate {
		return nil, fmt.Errorf("error Rate must be between 0 and 1")
	}

	if capacity < 1 {
		return nil, fmt.Errorf("capacity must be > 0")
	}

	numSlices := int(math.Ceil(math.Log2(float64(1.0) / errorRate)))
	bitsPerSlice := uint64(math.Ceil((float64(capacity) * math.Abs(math.Log(errorRate))) / (float64(numSlices) * log2Const)))

	bf := &BloomFilter{}
	bf.setup(errorRate, bitsPerSlice, numSlices, capacity, int64(0))
	return bf, nil
}

func (bf *BloomFilter) setup(errorRates float64, bitsPerSlice uint64, numSlices int, capacity, count int64) {

	bf.errorRate = errorRates
	bf.numSlices = numSlices
	bf.bitsPerSlice = bitsPerSlice
	bf.capacity = capacity
	bf.count = count
	bf.numBits = uint64(numSlices) * bitsPerSlice
	bf.makeSalts()

	bf.bitarray = array.New(bf.numBits)
}

// Add new key. Returns true/false for key and error.
func (bf *BloomFilter) Add(key []byte, skipChecks ...bool) (bool, error) {

	if bf.count > bf.capacity {
		return false, capacityError
	}

	skipCheck := false
	if len(skipChecks) > 0 {
		skipCheck = skipChecks[0]
	}

	foundAllBits := true
	offset := uint64(0)
	hashes := bf.makeSaltIterator(key)
	k, find := hashes.next()
	for find {
		if !skipCheck && foundAllBits {
			if !bf.bitarray.Get(offset + k) {
				foundAllBits = false
			}
		}

		bf.bitarray.Set(offset + k)
		offset += bf.bitsPerSlice

		k, find = hashes.next()
	}

	if skipCheck || !foundAllBits {
		bf.count++
		return false, nil
	}

	return true, nil
}

// Check key. Returns true/false
func (bf *BloomFilter) Check(key []byte) bool {
	offset := uint64(0)
	hashes := bf.makeSaltIterator(key)
	k, find := hashes.next()
	for find {
		if !bf.bitarray.Get(offset + k) {
			return false
		}
		offset += bf.bitsPerSlice
		k, find = hashes.next()
	}
	return true
}

func (bf *BloomFilter) makeSaltIterator(key []byte) *saltIterator {
	iterator := &saltIterator{
		numSlices:        bf.numSlices,
		saltFunctions:    bf.saltFunctions,
		numSaltFunctions: len(bf.saltFunctions),
		key:              key,
		tmpBody:          []uint64{},
		bitsPerSlice:     bf.bitsPerSlice,
		chunkSize:        bf.chunkSize,
	}

	return iterator
}

func (bf *BloomFilter) makeSalts() {

	if bf.bitsPerSlice >= (1 << 31) {
		bf.chunkSize = 8
	} else if bf.bitsPerSlice >= (1 << 15) {
		bf.chunkSize = 4
	} else {
		bf.chunkSize = 2
	}

	totalHashBits := 8 * bf.numSlices * bf.chunkSize
	fmtLength := md5.Size
	bf.hashfnname = "md5"

	if totalHashBits > 384 {
		bf.hashfnname = "sha512"
		fmtLength = sha512.Size
	} else if totalHashBits > 256 {
		bf.hashfnname = "sha384"
		fmtLength = sha512.Size384
	} else if totalHashBits > 160 {
		bf.hashfnname = "sha256"
		fmtLength = sha256.Size
	} else if totalHashBits > 128 {
		bf.hashfnname = "sha1"
		fmtLength = sha1.Size
	}

	numSalts := bf.numSlices / fmtLength
	if bf.numSlices > numSalts*fmtLength {
		numSalts++
	}

	bf.saltFunctions = []*hashFN{}
	for i := 0; i < numSalts; i++ {
		s := sum(bf.hashfnname, packInt(i))
		bf.saltFunctions = append(bf.saltFunctions, newHashFN(bf.hashfnname, s))
	}
}

// Count is a "getter". Returns all number of added keys.
func (bf *BloomFilter) Count() int64 {
	return bf.count
}

// Capacity is a "getter". Returns full Capacity
func (bf *BloomFilter) Capacity() int64 {
	return bf.capacity
}

// ErrorRate is a "getter". Returns error rate for current filter
func (bf *BloomFilter) ErrorRate() float64 {
	return bf.errorRate
}

// Merge integrates 2 filters. Filters must have the same parameters
func (bf *BloomFilter) Merge(bfNew *BloomFilter) error {

	if err := bf.compare(bfNew); err != nil {
		return err
	}

	return bf.bitarray.Merge(bfNew.bitarray)
}

func (bf *BloomFilter) compare(bfNew *BloomFilter) error {

	if bf.numBits != bfNew.numBits {
		return fmt.Errorf("Wrong length for numBits: %d != %d", bf.numBits, bfNew.numBits)
	}

	if bf.bitsPerSlice != bfNew.bitsPerSlice {
		return fmt.Errorf("Wrong length for bitsPerSlice: %d != %d", bf.bitsPerSlice, bfNew.bitsPerSlice)
	}

	if bf.capacity != bfNew.capacity {
		return fmt.Errorf("Wrong length for capacity: %d != %d", bf.capacity, bfNew.capacity)
	}

	if bf.chunkSize != bfNew.chunkSize {
		return fmt.Errorf("Wrong length for chunkSize: %d != %d", bf.chunkSize, bfNew.chunkSize)
	}

	if bf.numSlices != bfNew.numSlices {
		return fmt.Errorf("Wrong length for numSlices: %d != %d", bf.numSlices, bfNew.numSlices)
	}

	if float32(bf.errorRate) != float32(bfNew.errorRate) {
		return fmt.Errorf("Wrong length for errorRate: %f != %f", bf.errorRate, bfNew.errorRate)
	}

	return bf.bitarray.Compare(bfNew.bitarray)
}

// ToFile saves bloom filter to file by file name.
func (bf *BloomFilter) ToFile(fileName string) error {
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		return err
	}

	binBuf := bytes.NewBuffer([]byte{})
	if err := bf.ToBytes(binBuf); err != nil {
		return err
	}

	_, err = file.Write(binBuf.Bytes())
	return err
}

// ToBytes returns binary image of bloom filter
func (bf *BloomFilter) ToBytes(binBuf *bytes.Buffer) error {

	binary.Write(binBuf, binary.LittleEndian, bf.errorRate)
	binary.Write(binBuf, binary.LittleEndian, uint64(bf.numSlices))
	binary.Write(binBuf, binary.LittleEndian, uint64(bf.bitsPerSlice))
	binary.Write(binBuf, binary.LittleEndian, uint64(bf.capacity))
	binary.Write(binBuf, binary.LittleEndian, uint64(bf.count))

	return bf.bitarray.ToBytes(binBuf)
}

// FromFile creates new bloom filter from file
func FromFile(fileName string, lengths ...int64) (*BloomFilter, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	length := int64(0)
	if len(lengths) > 1 {
		length = lengths[1]
	}

	return FromReader(bufio.NewReader(file), length)
}

// FromReader creates new bloom filter from bufio.Reader
func FromReader(reader *bufio.Reader, length int64) (*BloomFilter, error) {

	var header struct {
		ErrorRate    float64
		NumSlices    int64
		BitsPerSlice int64
		Capacity     int64
		Count        int32
	}

	const headerLen = int64(unsafe.Sizeof(header))

	b := make([]byte, headerLen, headerLen)
	_, err := reader.Read(b)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(b)
	if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	bf := &BloomFilter{}
	bf.setup(header.ErrorRate, uint64(header.BitsPerSlice), int(header.NumSlices), int64(header.Capacity), int64(header.Count))

	if length > 0 {
		length = length - headerLen
	}
	bf.bitarray.Read(reader, length)

	return bf, nil
}
