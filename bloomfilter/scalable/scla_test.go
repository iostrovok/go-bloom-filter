package scalable

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/iostrovok/go-bloom-filter/bloomfilter/fortesting"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type scalTestSuite struct{}

var _ = Suite(&scalTestSuite{})

func (s *scalTestSuite) TestSetup(c *C) {

	filter, err := New(100, 0.0001, 4)
	c.Assert(err, IsNil)

	mode := 2
	ratio := float64(.9)
	initialCapacity := int64(100)
	errorRate := float64(.05)
	c.Assert(filter.Setup(mode, ratio, initialCapacity, errorRate), IsNil)

	mode = 1
	c.Assert(filter.Setup(mode, ratio, initialCapacity, errorRate), NotNil)

	mode = 2
	ratio = float64(0.0)
	c.Assert(filter.Setup(mode, ratio, initialCapacity, errorRate), NotNil)
	ratio = float64(10.0)
	c.Assert(filter.Setup(mode, ratio, initialCapacity, errorRate), NotNil)

	ratio = float64(.9)
	errorRate = float64(1.1)
	c.Assert(filter.Setup(mode, ratio, initialCapacity, errorRate), NotNil)

	errorRate = float64(0)
	c.Assert(filter.Setup(mode, ratio, initialCapacity, errorRate), NotNil)

	errorRate = float64(.05)
	initialCapacity = int64(0)
	c.Assert(filter.Setup(mode, ratio, initialCapacity, errorRate), NotNil)
}

func (s *scalTestSuite) TestSimple(c *C) {

	filter, err := New(100, 0.0001)
	c.Assert(err, IsNil)

	testArray := fortesting.ArrayForTesting()

	for _, s := range testArray {
		res, err := filter.Add([]byte(s))
		c.Assert(err, IsNil)
		c.Assert(res, Equals, false)
	}

	for _, s := range testArray {
		res := filter.Check([]byte(s))
		c.Assert(res, Equals, true)
	}

	for _, s := range testArray {
		res := filter.Check([]byte(s + "eee4"))
		c.Assert(res, Equals, false)
	}
}

func (s *scalTestSuite) TestFromFile(c *C) {

	filter, err := FromFile(fortesting.Dir() + "/test_scal.bin")
	c.Assert(err, IsNil)
	c.Assert(filter, NotNil)

	c.Assert(filter.Count(), Equals, int64(1930))
	c.Assert(filter.Capacity(), Equals, int64(3100))

	testArray := fortesting.ArrayForTesting()

	for _, s := range testArray {
		res := filter.Check([]byte(s))
		c.Assert(res, Equals, true)
	}

	for _, s := range testArray {
		res := filter.Check([]byte(s + "eee-33"))
		c.Assert(res, Equals, false)
	}
}

func (s *scalTestSuite) TestCapacity(c *C) {

	testArray := fortesting.ArrayForTesting()

	filter, err := New(len(testArray), 0.01)
	c.Assert(err, IsNil)
	c.Assert(filter, NotNil)

	for _, s := range testArray {
		res, err := filter.Add([]byte(s))
		c.Assert(err, IsNil)
		c.Assert(res, Equals, false)
	}

	for _, s := range testArray {
		res := filter.Check([]byte(s))
		c.Assert(res, Equals, true)
	}

	countBad := 0
	for _, s := range testArray {
		if filter.Check([]byte(s + "eee-2-1")) {
			countBad++
		}
		if filter.Check([]byte(s + "eee-2-2")) {
			countBad++
		}
		if filter.Check([]byte(s + "eee-2-3")) {
			countBad++
		}
	}

	check := (float32(countBad) / float32(len(testArray))) < 0.01

	c.Assert(check, Equals, true)
}

func (s *scalTestSuite) TestToFile(c *C) {

	filter, err := New(100, 0.0001)
	c.Assert(err, IsNil)

	testArray := fortesting.ArrayForTesting()

	for _, s := range testArray {
		res, err := filter.Add([]byte(s))
		c.Assert(err, IsNil)
		c.Assert(res, Equals, false)
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
	c.Assert(err, IsNil)
	tmpFile.Close()

	fileName := tmpFile.Name()

	// Remember to clean up the file afterwards
	defer os.Remove(fileName)

	filter.ToFile(fileName)
	fortesting.CheckFiles(c, fileName, fortesting.Dir()+"/test_scal.bin")
}

func (s *scalTestSuite) TestToBytes(c *C) {

	filter, err := New(100, 0.0001)
	c.Assert(err, IsNil)

	testArray := fortesting.ArrayForTesting()

	for _, s := range testArray {
		res, err := filter.Add([]byte(s))
		c.Assert(err, IsNil)
		c.Assert(res, Equals, false)
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
	c.Assert(err, IsNil)

	fileName := tmpFile.Name()

	// Remember to clean up the file afterwards
	defer os.Remove(fileName)

	buf := filter.ToBytes()
	c.Assert(buf, NotNil)

	_, err = tmpFile.Write(buf)
	c.Assert(err, IsNil)
	c.Assert(tmpFile.Close(), IsNil)

	fortesting.CheckFiles(c, fileName, fortesting.Dir()+"/test_scal.bin")
}

func (s *scalTestSuite) TestMerge1(c *C) {

	filter, err := New(10, 0.0001)
	c.Assert(err, IsNil)
	c.Assert(filter, NotNil)

	filter2, err := New(10, 0.0001)
	c.Assert(err, IsNil)
	c.Assert(filter2, NotNil)

	testArray := fortesting.ArrayForTesting()

	for i, s := range testArray {
		if i < 100 {
			res, err := filter.Add([]byte(s))
			c.Assert(err, IsNil)
			c.Assert(res, Equals, false)
		}

		res, err := filter2.Add([]byte(s))
		c.Assert(err, IsNil)
		c.Assert(res, Equals, false)
	}

	err = filter.Merge(filter2)
	c.Assert(err, IsNil)

	for _, s := range testArray {
		res := filter.Check([]byte(s))
		c.Assert(res, Equals, true)
	}

	for _, s := range testArray {
		// fmt.Printf("%s\n", s+"-eee-25")
		res := filter.Check([]byte(s + "-eee-25"))
		c.Assert(res, Equals, false)
	}
}

func (s *scalTestSuite) TestMerge2(c *C) {

	filterA, err := New(20, 0.0001)
	c.Assert(err, IsNil)
	c.Assert(filterA, NotNil)

	filterB, err := New(20, 0.0001)
	c.Assert(err, IsNil)
	c.Assert(filterB, NotNil)

	filterC, err := New(20, 0.0001)
	c.Assert(err, IsNil)
	c.Assert(filterC, NotNil)

	testArray := fortesting.ArrayForTesting()

	for i, s := range testArray {

		if i < 100 {
			res, err := filterA.Add([]byte(s))
			c.Assert(err, IsNil)
			c.Assert(res, Equals, false)

			res, err = filterB.Add([]byte(s))
			c.Assert(err, IsNil)
			c.Assert(res, Equals, false)

			res, err = filterC.Add([]byte(s))
			c.Assert(err, IsNil)
			c.Assert(res, Equals, false)

			continue
		}

		if 100 < i && i < 500 {
			res, err := filterB.Add([]byte(s))
			c.Assert(err, IsNil)
			c.Assert(res, Equals, false)

			res, err = filterC.Add([]byte(s))
			c.Assert(err, IsNil)
			c.Assert(res, Equals, false)
			continue
		}

		res, err := filterC.Add([]byte(s))
		c.Assert(err, IsNil)
		c.Assert(res, Equals, false)
	}

	err = filterA.Merge(filterB)
	c.Assert(err, IsNil)

	err = filterA.Merge(filterC)
	c.Assert(err, IsNil)

	for _, s := range testArray {
		res := filterA.Check([]byte(s))
		c.Assert(res, Equals, true)
	}

	for _, s := range testArray {
		res := filterA.Check([]byte(s + "eee-23"))
		c.Assert(res, Equals, false)
	}
}

func (s *scalTestSuite) TestMerge3(c *C) {

	testArray := fortesting.ArrayForTesting()

	filterA, err := New(len(testArray), 0.0001)
	c.Assert(err, IsNil)
	c.Assert(filterA, NotNil)

	filterB, err := New(len(testArray), 0.0001)
	c.Assert(err, IsNil)
	c.Assert(filterB, NotNil)

	for _, s := range testArray {
		res, err := filterB.Add([]byte(s))
		c.Assert(err, IsNil)
		c.Assert(res, Equals, false)
	}

	err = filterA.Merge(filterB)
	c.Assert(err, IsNil)

	for _, s := range testArray {
		res := filterA.Check([]byte(s))
		c.Assert(res, Equals, true)
	}

	for _, s := range testArray {
		res := filterA.Check([]byte(s + "eee-24"))
		c.Assert(res, Equals, false)
	}
}

func (s *scalTestSuite) TestMergeError(c *C) {

	testArray := fortesting.ArrayForTesting()

	filterA, err := New(len(testArray), 0.0001)
	c.Assert(err, IsNil)

	filterB, err := New(len(testArray), 0.0001, 4)
	c.Assert(err, IsNil)

	filterC, err := New(len(testArray), 0.0001, 6)
	c.Assert(err, IsNil)

	for _, s := range testArray {
		res, err := filterA.Add([]byte(s))
		c.Assert(err, IsNil)
		c.Assert(res, Equals, false)

		res, err = filterB.Add([]byte(s))
		c.Assert(err, IsNil)
		c.Assert(res, Equals, false)

		res, err = filterC.Add([]byte(s))
		c.Assert(err, IsNil)
		c.Assert(res, Equals, false)
	}

	err = filterA.Merge(filterB)
	c.Assert(err, NotNil)

	err = filterA.Merge(filterC)
	c.Assert(err, NotNil)

	err = filterB.Merge(filterC)
	c.Assert(err, NotNil)

	err = filterB.Merge(filterA)
	c.Assert(err, NotNil)

	err = filterC.Merge(filterA)
	c.Assert(err, NotNil)

	err = filterC.Merge(filterB)
	c.Assert(err, NotNil)

}

func (s *scalTestSuite) TestMergeErrorB(c *C) {

	l := 3000

	filterA, err := New(l, 0.0001)
	c.Assert(err, IsNil)

	filterB, err := New(l, 0.001)
	c.Assert(err, IsNil)

	filterC, err := New(l+10, 0.0001)
	c.Assert(err, IsNil)

	err = filterA.Merge(filterB)
	c.Assert(err, NotNil)

	err = filterA.Merge(filterC)
	c.Assert(err, NotNil)

	err = filterB.Merge(filterC)
	c.Assert(err, NotNil)

	err = filterB.Merge(filterA)
	c.Assert(err, NotNil)

	err = filterC.Merge(filterA)
	c.Assert(err, NotNil)

	err = filterC.Merge(filterB)
	c.Assert(err, NotNil)

}
