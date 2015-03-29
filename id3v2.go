package id3v2

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

const tagHeaderSize = 10
const frameHeaderSize = 10
const emptyHeader = "\x00\x00\x00\x00"

type Id3v2Tag struct {
	version string
	frames  map[string][]byte
	flags   byte
}

func ReadTag(r io.Reader) (Id3v2Tag, error) {
	b := make([]byte, tagHeaderSize)
	r.Read(b)
	header := string(b[:3])
	if string(header) != "ID3" {
		return Id3v2Tag{}, errors.New("No tag found")
	}

	tag := new(Id3v2Tag)
	version := int(b[3])
	revision := int(b[4])
	tag.version = fmt.Sprintf("2.%d.%d", version, revision)
	tag.flags = b[5]

	size := determineSizeOfTag(b[6:])
	tag.frames = readFrames(size, r)

	return *tag, nil
}

func determineSizeOfFrame(b []byte) int {
	return int(b[0])<<24 |
		int(b[1])<<16 |
		int(b[2])<<8 |
		int(b[3])
}

func determineSizeOfTag(b []byte) int {
	return int(b[0])<<21 |
		int(b[1])<<14 |
		int(b[2])<<7 |
		int(b[3])
}

func readFrames(size int, r io.Reader) (frames map[string][]byte) {
	frames = make(map[string][]byte)
	buf := make([]byte, size)
	readBytes := 0
	for readBytes < size {
		r.Read(buf[:frameHeaderSize])
		name := string(buf[:4])
		frameSize := determineSizeOfFrame(buf[4:8])

		if name == emptyHeader || frameSize == 0 {
			continue
		} else if frameSize > size {
			break
		}

		r.Read(buf[:frameSize])
		payload := make([]byte, frameSize)
		copy(payload, buf)
		frames[name] = payload
		readBytes += frameSize
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
		return readUtf16(b[1:])
	}
}

func readUtf16(b []byte) string {
	byteOrder := determineByteOrder(b[0])
	utf := make([]uint16, (len(b)+(2-1))/2)
	for i := 0; i+(2-1) < len(b); i += 2 {
		utf[i/2] = byteOrder.Uint16(b[i:])
	}
	if len(b)/2 < len(utf) {
		utf[len(utf)-1] = utf8.RuneError
	}
	return string(utf16.Decode(utf)[1:])
}

func determineByteOrder(b byte) binary.ByteOrder {
	if uint16(b) == 0xFEFF {
		return binary.BigEndian
	} else {
		return binary.LittleEndian
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

func (tag *Id3v2Tag) Version() string {
	return tag.version
}
