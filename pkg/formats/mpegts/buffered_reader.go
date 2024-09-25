package mpegts

import (
	"encoding/hex"
	"fmt"
	"io"
)

// BufferedReader is a chunked reader optimized for MPEG-TS.
type BufferedReader struct {
	r io.Reader
}

// NewBufferedReader allocates a BufferedReader.
func NewBufferedReader(r io.Reader) *BufferedReader {
	return &BufferedReader{
		r: r,
	}
}

// Read implements io.Reader.
func (r *BufferedReader) Read(p []byte) (int, error) {
	n, err := io.ReadFull(r.r, p[:188])
	if n == 0 {
		return 0, io.EOF
	}
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return 0, fmt.Errorf("received packet with size %d not multiple of 188", n)
		}
	}

	fmt.Printf("Read %d bytes:\n%s\n", n, hex.Dump(p[:n]))

	return n, nil
}
