package mpegts

import (
	"encoding/hex"
	"fmt"
	"io"
)

// BufferedReader is a buffered reader optimized for MPEG-TS.
type BufferedReader struct {
	r         io.Reader
	midbuf    []byte
	midbufpos int
	target    bool
}

// NewBufferedReader allocates a BufferedReader.
func NewBufferedReader(r io.Reader, target bool) *BufferedReader {
	return &BufferedReader{
		r:      r,
		midbuf: make([]byte, 0, 1500),
		target: target,
	}
}

// Read implements io.Reader.
func (r *BufferedReader) Read(p []byte) (int, error) {
	if r.midbufpos < len(r.midbuf) {
		n := copy(p, r.midbuf[r.midbufpos:])
		r.midbufpos += n
		return n, nil
	}

	mn, err := r.r.Read(r.midbuf[:cap(r.midbuf)])
	if err != nil {
		return 0, err
	}

	if (mn % 188) != 0 {
		return 0, fmt.Errorf("received packet with size %d not multiple of 188", mn)
	}

	if r.target {
		fmt.Printf("Read %d bytes:\n%s\n", mn, hex.Dump(r.midbuf[:mn]))
	}

	r.midbuf = r.midbuf[:mn]
	n := copy(p, r.midbuf)
	r.midbufpos = n
	return n, nil
}
