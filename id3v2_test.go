package id3v2

import (
	. "gopkg.in/check.v1"
	"os"
	"testing"
)

func TestId3V2(t *testing.T) { TestingT(t) }

type Id3v2TagSuite struct{}

var _ = Suite(&Id3v2TagSuite{})

func readFile(fn string, c *C) *os.File {
	file, err := os.Open(fn)
	c.Assert(err, Equals, nil)
	return file
}

func (s *Id3v2TagSuite) TestReadNoTag(c *C) {
	file := readFile("_testdata/tagless-batman.mp3", c)
	defer file.Close()

	tag, err := ReadTag(file)

	c.Check(tag, Equals, Id3v2Tag{})
	c.Check(err.Error(), Equals, "No tag foound")
}