package bloomfilter

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/iostrovok/go-bloom-filter/bloomfilter/fortesting"
	. "gopkg.in/check.v1"
)

type filterTestSuite struct{}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&filterTestSuite{})

func (s *filterTestSuite) TestNew(c *C) {

	filter, err := New(10000*10000, 10.001)
	c.Assert(err, NotNil)
	c.Assert(filter, IsNil)

	filter, err = New(0, .001)
	c.Assert(err, NotNil)
	c.Assert(filter, IsNil)
}

func (s *filterTestSuite) TestOverload(c *C) {

	filter, err := New(100, .0001)
	c.Assert(err, IsNil)
	c.Assert(filter, NotNil)

	testArray := fortesting.ArrayForTesting()

	for i, s := range testArray {
		_, err := filter.Add([]byte(s))
		if i <= 100 {
			c.Assert(err, IsNil)
		} else {
			c.Assert(err, NotNil)

		}
	}

}

func (s *filterTestSuite) TestNew2(c *C) {

	filter, err := New(10000*10000, 0.001)
	c.Assert(err, IsNil)
	c.Assert(filter, NotNil)

	testArray := fortesting.ArrayForTesting()

	for _, s := range testArray {
		res, err := filter.Add([]byte(s))
		c.Assert(err, IsNil)
		c.Assert(res, Equals, false)
	}

	for _, s := range testArray {
		res, err := filter.Add([]byte(s), true)
		c.Assert(err, IsNil)
		c.Assert(res, Equals, false)

		res, err = filter.Add([]byte(s), false)
		c.Assert(err, IsNil)
		c.Assert(res, Equals, true)
	}

	for _, s := range testArray {
		res := filter.Check([]byte(s))
		c.Assert(res, Equals, true)
	}

	for _, s := range testArray {
		res := filter.Check([]byte(s + "eee1"))
		c.Assert(res, Equals, false)
	}
}

func (s *filterTestSuite) TestToFile(c *C) {

	filter, err := New(10000, 0.001)
	c.Assert(err, IsNil)
	c.Assert(filter, NotNil)

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
	fortesting.CheckFiles(c, fileName, fortesting.Dir()+"/test_simple.bin")
}

func (s *filterTestSuite) TestToBytes(c *C) {

	filter, err := New(10000, 0.001)
	c.Assert(err, IsNil)
	c.Assert(filter, NotNil)

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

	binBuf := bytes.NewBuffer([]byte{})
	err = filter.ToBytes(binBuf)
	c.Assert(err, IsNil)
	_, err = tmpFile.Write(binBuf.Bytes())
	c.Assert(err, IsNil)
	c.Assert(tmpFile.Close(), IsNil)

	fortesting.CheckFiles(c, fileName, fortesting.Dir()+"/test_simple.bin")
}

func (s *filterTestSuite) TestFromFile(c *C) {

	filter, err := FromFile(fortesting.Dir() + "/test_simple.bin")
	c.Assert(err, IsNil)
	c.Assert(filter, NotNil)

	c.Assert(filter.Count(), Equals, int64(1930))
	c.Assert(filter.Capacity(), Equals, int64(10000))

	testArray := fortesting.ArrayForTesting()

	for _, s := range testArray {
		res := filter.Check([]byte(s))
		c.Assert(res, Equals, true)
	}

	for _, s := range testArray {
		res := filter.Check([]byte(s + "eee2"))
		c.Assert(res, Equals, false)
	}
}

func (s *filterTestSuite) TestLarge(c *C) {

	// big
	filter, err := New(10000*1000*1000, 0.001)
	c.Assert(err, IsNil)
	c.Assert(filter, NotNil)

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
		res := filter.Check([]byte(s + "eee3"))
		c.Assert(res, Equals, false)
	}
}
func (s *filterTestSuite) TestMergeEmpty(c *C) {

	filterA, err := New(1000*1000, 0.001)
	c.Assert(err, IsNil)
	c.Assert(filterA, NotNil)

	filterB, err := New(1000*1000, 0.001)
	c.Assert(err, IsNil)
	c.Assert(filterB, NotNil)

	testArray := fortesting.ArrayForTesting()

	for _, s := range testArray {
		res, err := filterA.Add([]byte(s))
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
		res := filterA.Check([]byte(s + "eee-22"))
		c.Assert(res, Equals, false)
	}
}

func (s *filterTestSuite) TestMerge(c *C) {

	filterA, err := New(1000*1000, 0.001)
	c.Assert(err, IsNil)
	c.Assert(filterA, NotNil)

	filterB, err := New(1000*1000, 0.001)
	c.Assert(err, IsNil)
	c.Assert(filterB, NotNil)

	filterC, err := New(1000*1000, 0.001)
	c.Assert(err, IsNil)
	c.Assert(filterC, NotNil)

	testArray := fortesting.ArrayForTesting()

	for i, s := range testArray {
		if i%3 == 0 {
			res, err := filterA.Add([]byte(s))
			c.Assert(err, IsNil)
			c.Assert(res, Equals, false)
		} else if i%3 == 1 {
			res, err := filterB.Add([]byte(s))
			c.Assert(err, IsNil)
			c.Assert(res, Equals, false)
		} else {
			res, err := filterC.Add([]byte(s))
			c.Assert(err, IsNil)
			c.Assert(res, Equals, false)
		}
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
		res := filterA.Check([]byte(s + "eee-21"))
		c.Assert(res, Equals, false)
	}
}
