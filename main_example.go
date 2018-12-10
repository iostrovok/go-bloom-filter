package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/iostrovok/go-bloom-filter/bloomfilter"
	"github.com/iostrovok/go-bloom-filter/bloomfilter/scalable"
)

var r *rand.Rand
var a1, a2 [][]byte
var countScalable, countTestKeys, falseAlarmCountTestKeys int
var errorRate float64

type filterType interface {
	Check([]byte) bool
	Add([]byte, ...bool) (bool, error)
}

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))

	countScalable = 10
	countTestKeys = 100000
	falseAlarmCountTestKeys = 100 * countTestKeys
	errorRate = 0.001

	fmt.Println("start generate random arrays..........")
	a1 = testArray("a-a", countTestKeys)
	a2 = testArray("b-b", falseAlarmCountTestKeys)
	fmt.Println("..........finish generate random arrays\n")
}

func filnaTest(filter filterType) {
	falseAlarm := 0
	falseCheck := 0
	fmt.Printf("start checking existed keys (%d)...\n", countTestKeys)
	for _, key := range a1 {
		if !filter.Check(key) {
			falseCheck++
		}
	}

	fmt.Printf("start checking non-existent keys (%d)...\n", falseAlarmCountTestKeys)
	for _, key := range a2 {
		if filter.Check(key) {
			falseAlarm++
		}
	}

	errorRateNow := float64(falseAlarm) / float64(falseAlarmCountTestKeys)
	fmt.Printf("falseCheck: %d from %d, falseAlarm: %d from %d\n", falseCheck, countTestKeys, falseAlarm, falseAlarmCountTestKeys)
	fmt.Printf("current error rate: %0.10f (expected: %0.10f)\n", errorRateNow, errorRate)
}

func addKeys(filter filterType) {
	fmt.Printf("Add %d keys\n", countTestKeys)
	for _, key := range a1 {
		_, err := filter.Add(key)
		if err != nil {
			panic(err)
		}
	}
}

func testArray(prefix string, count int) [][]byte {
	out := [][]byte{}
	t := 1
	a := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	for {
		if count%10000 == 0 {
			t = 1 + r.Int()%3
		}
		rv := fmt.Sprintf("%20.20f", r.Float64())
		for _, i := range a {
			i1 := prefix + "-" + rv + "-" + strings.Repeat(i, t)
			for _, j := range a {
				j1 := i1 + "-" + strings.Repeat(j, t)
				for _, k := range a {
					k1 := j1 + "-" + strings.Repeat(k, t)
					for _, l := range a {
						out = append(out, []byte(k1+strings.Repeat(l, t)))
						count--
						if count < 1 {
							return out
						}
					}
				}
			}
		}
	}
	return [][]byte{}
}

func simpleBloomFilter() {

	fmt.Println("Start simpleBloomFilter..........")
	fmt.Printf("create new bloom filter with capacity %d, errorRate: %0.3f\n", countTestKeys, errorRate)
	filter, err := bloomfilter.New(int64(countTestKeys), errorRate)
	if err != nil {
		panic(err)
	}

	addKeys(filter)
	filnaTest(filter)
	fmt.Println("..........Finish simpleBloomFilter\n\n")
}

func scalableBloomFilter() {

	fmt.Println("Start scalableBloomFilter..........")
	fmt.Printf("create new scalable bloom filter with capacity %d, errorRate: %0.3f\n", countTestKeys/countScalable, errorRate)
	filter, err := scalable.New(countTestKeys/countScalable, errorRate)
	if err != nil {
		panic(err)
	}

	addKeys(filter)
	filnaTest(filter)
	fmt.Println("..........Finish scalableBloomFilter\n\n")
}

func bloomFilterAndFile() {

	fmt.Println("Start bloomFilterAndFile..........")
	fmt.Printf("create new bloom filter with capacity %d, errorRate: %0.3f\n", countTestKeys, errorRate)
	filter, err := bloomfilter.New(int64(countTestKeys), errorRate)
	if err != nil {
		panic(err)
	}

	addKeys(filter)

	tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
	tmpFile.Close()
	fileName := tmpFile.Name()
	// Remember to clean up the file afterwards
	defer os.Remove(fileName)

	fmt.Println("Save bloom filter to file")
	filter.ToFile(fileName)

	fmt.Println("Read new bloom filter from file")
	filter2, err := bloomfilter.FromFile(fileName)
	if err != nil {
		panic(err)
	}

	filnaTest(filter2)
	fmt.Println("..........Finish bloomFilterAndFile\n\n")
}

func scalableBloomFilterAndFile() {

	fmt.Println("Start scalableBloomFilterAndFile..........")
	fmt.Printf("create new scalable bloom filter with capacity %d, errorRate: %0.3f\n", countTestKeys/countScalable, errorRate)
	filter, err := scalable.New(countTestKeys/countScalable, errorRate)
	if err != nil {
		panic(err)
	}

	addKeys(filter)

	tmpFile, err := ioutil.TempFile(os.TempDir(), "scalable-")
	tmpFile.Close()
	fileName := tmpFile.Name()

	// Remember to clean up the file afterwards
	defer os.Remove(fileName)

	fmt.Println("Save scalable bloom filter to file")
	filter.ToFile(fileName)

	fmt.Println("Read new scalable bloom filter from file")
	filter2, err := scalable.FromFile(fileName)
	if err != nil {
		panic(err)
	}

	filnaTest(filter2)
	fmt.Println("..........Finish scalableBloomFilterAndFile\n\n")
}

func main() {
	simpleBloomFilter()
	scalableBloomFilter()
	bloomFilterAndFile()
	scalableBloomFilterAndFile()
}
