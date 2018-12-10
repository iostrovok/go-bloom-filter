package scalable

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"sync"

	"github.com/iostrovok/go-bloom-filter/bloomfilter"
)

/*
	Implements a space-efficient probabilistic data structure that grows as more items
	are added while maintaining a steady false positive rate.
	more details see here - https://github.com/jaybaird/python-bloomfilter/
*/

const (
	// SmallSetGrowth is a constant for increasing capacity
	SmallSetGrowth = 2 // slower, but takes up less memory
	// LargeSetGrowth is a constant for increasing capacity
	LargeSetGrowth = 4 // faster, but takes up more memory faster
)

// Filter is a structure for scalable bloom filter.
type Filter struct {
	mc sync.RWMutex

	filters []*bloomfilter.BloomFilter

	scale           int
	ratio           float64
	initialCapacity int64
	errorRate       float64
}

// New is constructor. It checks parameters and creates new scalable bloom filter.
func New(initialCapacity int, errorRate float64, modes ...int) (*Filter, error) {

	mode := SmallSetGrowth
	if len(modes) > 0 {
		mode = modes[0]
	}

	sbf := &Filter{
		filters: []*bloomfilter.BloomFilter{},
	}

	err := sbf.Setup(mode, 0.9, int64(initialCapacity), errorRate)
	if err != nil {
		return nil, err
	}

	return sbf, nil
}

// Setup for main parameters. We may change after filter creation.
func (sbf *Filter) Setup(mode int, ratio float64, initialCapacity int64, errorRate float64) error {

	if errorRate <= 0 || 1.0 < errorRate {
		return fmt.Errorf("error Rate must be between 0 and 1")
	}

	if ratio <= 0 || 1.0 < ratio {
		return fmt.Errorf("ratio must be between 0 and 1")
	}

	if mode <= 1 {
		return fmt.Errorf("mode must be more than 0")
	}

	if initialCapacity < 1 {
		return fmt.Errorf("initialCapacity must be > 0")
	}

	sbf.scale = mode
	sbf.ratio = ratio
	sbf.initialCapacity = initialCapacity
	sbf.errorRate = errorRate

	return nil
}

// Add new key. Returns true/false for key and error.
func (sbf *Filter) Add(key []byte, skipChecks ...bool) (bool, error) {

	if sbf.Check(key) {
		return true, nil
	}

	filter, err := sbf.getEmptyFilter()
	if err != nil {
		return false, err
	}

	return filter.Add(key, skipChecks...)
}

// Check key. Returns true/false
func (sbf *Filter) Check(key []byte) bool {
	for i := len(sbf.filters) - 1; i > -1; i-- {
		if sbf.filters[i].Check(key) {
			return true
		}
	}

	return false
}

func (sbf *Filter) getEmptyFilter() (*bloomfilter.BloomFilter, error) {

	if len(sbf.filters) == 0 {

		sbf.mc.Lock()
		defer sbf.mc.Unlock()

		if len(sbf.filters) == 0 {
			filter, err := bloomfilter.New(sbf.initialCapacity, sbf.errorRate*(1.0-sbf.ratio))
			if err != nil {
				return nil, err
			}
			sbf.filters = append(sbf.filters, filter)
		}

		return sbf.filters[len(sbf.filters)-1], nil
	}

	filter := sbf.filters[len(sbf.filters)-1]
	if filter.Count() < filter.Capacity() {
		return filter, nil
	}

	sbf.mc.Lock()
	defer sbf.mc.Unlock()

	if filter.Count() >= filter.Capacity() {
		newFilter, err := bloomfilter.New(filter.Capacity()*int64(sbf.scale), filter.ErrorRate()*sbf.ratio)
		if err != nil {
			return nil, err
		}
		sbf.filters = append(sbf.filters, newFilter)
	}

	return sbf.filters[len(sbf.filters)-1], nil
}

// Merge integrates 2 scalable bloom filters. Filters must have the same parameters
func (sbf *Filter) Merge(sbfNew *Filter) error {
	/*
		You should believe that you have not added keys to target filter before merging.
		This function goal is synchronize local read-only filter and remote one.
	*/
	if float32(sbf.ratio) != float32(sbfNew.ratio) {
		return fmt.Errorf("Wrong length for scale: %f != %f", sbf.ratio, sbfNew.ratio)
	}

	if float32(sbf.errorRate) != float32(sbfNew.errorRate) {
		return fmt.Errorf("Wrong length for errorRate: %f != %f", sbf.errorRate, sbfNew.errorRate)
	}

	if sbf.initialCapacity != sbfNew.initialCapacity {
		return fmt.Errorf("Wrong length for initialCapacity: %d != %d", sbf.initialCapacity, sbfNew.initialCapacity)
	}

	if sbf.scale != sbfNew.scale {
		return fmt.Errorf("Wrong length for scale: %d != %d", sbf.scale, sbfNew.scale)
	}

	for i := 0; i < len(sbf.filters) && i < len(sbfNew.filters); i++ {
		err := sbf.filters[i].Merge(sbfNew.filters[i])
		if err != nil {
			return err
		}
	}

	sbf.mc.Lock()
	defer sbf.mc.Unlock()

	if len(sbf.filters) >= len(sbfNew.filters) {
		return nil
	}

	for i := len(sbf.filters); i < len(sbfNew.filters); i++ {
		sbf.filters = append(sbf.filters, sbfNew.filters[i])
	}

	return nil
}

// Capacity is a "getter". Returns full Capacity
func (sbf *Filter) Capacity() int64 {
	// Returns the total capacity for all filters in this SBF
	res := int64(0)
	for _, f := range sbf.filters {
		res += int64(f.Capacity())
	}
	return res
}

// Count is a "getter". Returns all number of added keys.
func (sbf *Filter) Count() int64 {
	// Returns the total number of elements stored in this SBF
	res := int64(0)
	for _, f := range sbf.filters {
		res += int64(f.Count())
	}
	return res
}

// ToFile saves scalable bloom filter to file by file name.
func (sbf *Filter) ToFile(fileName string) error {
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		return err
	}

	_, err = file.Write(sbf.ToBytes())
	return err
}

// ToBytes returns binary image of scalable bloom filter
func (sbf *Filter) ToBytes() []byte {

	sbf.mc.RLock()
	defer sbf.mc.RUnlock()

	binBuf := bytes.NewBuffer([]byte{})
	binary.Write(binBuf, binary.LittleEndian, uint32(sbf.scale))
	binary.Write(binBuf, binary.LittleEndian, sbf.ratio)
	binary.Write(binBuf, binary.LittleEndian, int64(sbf.initialCapacity))
	binary.Write(binBuf, binary.LittleEndian, float64(sbf.errorRate))
	binary.Write(binBuf, binary.LittleEndian, int32(len(sbf.filters)))

	if len(sbf.filters) == 0 {
		return binBuf.Bytes()
	}

	// Then each filter directly, with a header describing
	// their lengths.
	headerPos := binBuf.Len()

	filterSizes := make([]uint64, len(sbf.filters), len(sbf.filters))
	for _, i := range filterSizes {
		binary.Write(binBuf, binary.LittleEndian, i)
	}

	var last, begin int
	for i, filter := range sbf.filters {
		begin = binBuf.Len()
		filter.ToBytes(binBuf)
		last = binBuf.Len()
		filterSizes[i] = uint64(last - begin)
	}

	return saveInt64List(headerPos, binBuf, filterSizes)
}

func saveInt64List(pos int, binBuf *bytes.Buffer, mylist []uint64) []byte {
	var tmpBinBuf bytes.Buffer
	for _, i := range mylist {
		binary.Write(&tmpBinBuf, binary.LittleEndian, i)
	}

	// copy to out list
	buf := tmpBinBuf.Bytes()
	out := binBuf.Bytes()
	for i := range buf {
		out[i+pos] = buf[i]
	}

	return out
}

// FromFile creates new scalable bloom filter from file
func FromFile(fileName string) (*Filter, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	return FromReader(reader)
}

// FromReader creates new scalable bloom filter from bufio.Reader
func FromReader(reader *bufio.Reader) (*Filter, error) {

	var header struct {
		Scale           int32
		Ratio           float64
		InitialCapacity int64
		ErrorRate       float64
		CountFilters    int32
	}

	// const headerLen = int64(unsafe.Sizeof(header))
	headerLen := 32
	b := make([]byte, headerLen, headerLen)
	_, err := reader.Read(b)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(b)
	if err := binary.Read(r, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	sbf, err := New(int(header.InitialCapacity), float64(header.ErrorRate), int(header.Scale))
	if err != nil {
		return nil, err
	}

	sbf.filters = make([]*bloomfilter.BloomFilter, header.CountFilters, header.CountFilters)

	if header.CountFilters == 0 {
		return sbf, nil
	}

	filterSizes := []uint64{}
	filterSizes, err = readArrayOfUint64(reader, int(header.CountFilters))
	if err != nil {
		return nil, err
	}

	for i, length := range filterSizes {

		sbf.filters[i], err = bloomfilter.FromReader(reader, int64(length))

		if err != nil {
			return nil, err
		}
	}

	return sbf, nil
}

func readArrayOfUint64(reader *bufio.Reader, count int) ([]uint64, error) {

	out := make([]uint64, count, count)
	for i := 0; i < count; i++ {
		array := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
		_, err := reader.Read(array)
		if err != nil {
			return nil, err
		}

		out[i] = binary.LittleEndian.Uint64(array)
	}
	return out, nil
}
