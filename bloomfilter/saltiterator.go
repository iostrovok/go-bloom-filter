package bloomfilter

// --------------------------------------------------------

type hashFN struct {
	name string
	salt []byte
}

func newHashFN(name string, salt []byte) (h *hashFN) {
	return &hashFN{
		name: name,
		salt: salt,
	}
}

func (it *hashFN) run(key []byte) []byte {
	return sum(it.name, append(it.salt, key...))
}

type saltIterator struct {
	i, j, count      int
	numSlices        int
	numSaltFunctions int
	saltFunctions    []*hashFN
	key              []byte
	tmpBody          []uint64
	bitsPerSlice     uint64
	chunkSize        int
}

func (it *saltIterator) next() (uint64, bool) {

	if it.count >= it.numSlices || it.i > it.numSaltFunctions {
		return 0, false
	}

	if it.count == 0 || it.j >= len(it.tmpBody) {
		it.tmpBody = unpack(it.chunkSize, it.saltFunctions[it.i%len(it.saltFunctions)].run(it.key))
		it.j = 0
		it.i++
	}

	res := it.tmpBody[it.j] % it.bitsPerSlice

	it.count++
	it.j++

	return res, true
}

// --------------------------------------------------------
