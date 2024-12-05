package mpegts

import (
	"fmt"
	"io"
	"os"
	"time"
)

// BufferedReader is a buffered reader optimized for MPEG-TS.
type BufferedReader struct {
	r            io.Reader
	midbuf       []byte
	midbufpos    int
	target       bool
	loggingStart *time.Time
	tsFile       *os.File
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
		if r.loggingStart == nil {
			now := time.Now()
			r.loggingStart = &now
		}
		if r.tsFile == nil {
			r.tsFile, err = os.OpenFile("/solink/logs/mediamtx/tsfile.ts", os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Failed to open ts file", err)
				r.tsFile = nil
			}
		}
		now := time.Now()
		if now.Sub(*r.loggingStart) < 30*time.Second {
			if r.tsFile != nil {
				r.tsFile.Write(r.midbuf[:mn])
			}
			fmt.Printf("Read %d bytes\n", mn)
		}
	}

	r.midbuf = r.midbuf[:mn]
	n := copy(p, r.midbuf)
	r.midbufpos = n
	return n, nil
}
