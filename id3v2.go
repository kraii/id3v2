package id3v2

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Id3v2Tag struct {
	version string
	frames  map[string][]byte
	flags   byte
}

func ReadTag(r io.Reader) (Id3v2Tag, error) {
	header := string(readNextNBytes(3, r))
	if string(header) != "ID3" {
		return Id3v2Tag{}, errors.New("No tag found")
	}

	tag := new(Id3v2Tag)

	version := int(readNextByte(r))	
	revision := int(readNextByte(r))
	tag.version = fmt.Sprintf("2.%d.%d", version, revision)
	tag.flags = readNextByte(r)


	size := determineSizeOfTag(r)
	tag.frames = readFrames(size, r)

	return *tag, nil
}

func readNextNBytes(n int, r io.Reader) (b []byte) {
	b = make([]byte, n)
	r.Read(b)
	return b
}

func readNextByte(r io.Reader) byte {
	b := readNextNBytes(1, r)
	return b[0]
}

func determineSizeOfFrame(r io.Reader) int {
	b := readNextNBytes(4, r)
	return int(b[0])<<24 |
		int(b[1])<<16 |
		int(b[2])<<8 |
		int(b[3])
}

func determineSizeOfTag(r io.Reader) int {
	b := readNextNBytes(4, r)
	return int(b[0])<<21 |
		int(b[1])<<14 |
		int(b[2])<<7 |
		int(b[3])
}

func readFrames(size int, r io.Reader) (frames map[string][]byte) {
	frames = make(map[string][]byte)
	for i := 0; i < size; i++ {
		name := string(readNextNBytes(4, r))
		frameSize := determineSizeOfFrame(r)
		readNextNBytes(2, r) // ignore two 1byte flags

		payload := readNextNBytes(frameSize, r)
		if frameSize != 0 {
			frames[name] = payload
		}
	}
	return
}

func (tag *Id3v2Tag) Artist() string {
	artist := tag.frames["TPE1"]
	if artist != nil {
		return convertString(artist)
	} else {
		return convertString(tag.frames["TPE2"])
	}
}

func convertString(b []byte) string {
	if len(b) <= 0 {
		return ""
	} else if b[0] == '\x00' {
		return strings.TrimLeft(string(b[1:]), "\x00")
	} else {
		panic("SHIT! UTF NICHT 8")
	}
}

func (tag *Id3v2Tag) Title() string {
	return convertString(tag.frames["TIT2"])
}

func (tag *Id3v2Tag) Album() string {
	return convertString(tag.frames["TALB"])
}

func (tag *Id3v2Tag) Year() string {
	return convertString(tag.frames["TYER"])
}

func (tag *Id3v2Tag) TrackNumber() string {
	return convertString(tag.frames["TRCK"])
}

func (tag *Id3v2Tag) Comment() string {
	return convertString(tag.frames["COMM"])
}