package id3v2

import (
	"bytes"
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

	_, err := ReadTag(file)

	c.Check(err.Error(), Equals, "No tag found")
}

func (s *Id3v2TagSuite) TestReadValidV230Tag(c *C) {
	file := readFile("_testdata/spice.mp3", c)
	defer file.Close()

	tag, err := ReadTag(file)

	c.Assert(err, Equals, nil)
	c.Check(tag.version, Equals, "2.3.0")
	c.Check(tag.Artist(), Equals, "Xander")
	c.Check(tag.Title(), Equals, "Spice")
	c.Check(tag.Album(), Equals, "Things")
	c.Check(tag.Year(), Equals, "2015")
	c.Check(tag.Comment(), Equals, "say -v Xander")
	c.Check(tag.TrackNumber(), Equals, "1")
}

func (s *Id3v2TagSuite) TestReadSize(c *C) {
	b := []byte{0x00, 0x00, 0x00, 0x16}

	c.Assert(determineSizeOfTag(bytes.NewReader(b)), Equals, 22)
}
