package id3v2

import (
	"io"
)

type Id3v2Tag struct {
	toast string
}

func ReadTag(r io.ReadSeeker) (Id3v2Tag, error) {
	return Id3v2Tag{}, nil
}