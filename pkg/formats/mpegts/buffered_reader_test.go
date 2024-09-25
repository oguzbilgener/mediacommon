package mpegts

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBufferedReader(t *testing.T) {
	var buf bytes.Buffer

	buf.Write(bytes.Repeat([]byte{1}, 188))
	buf.Write(bytes.Repeat([]byte{2}, 188))
	buf.Write(bytes.Repeat([]byte{3}, 188))

	r := NewBufferedReader(&buf)

	byts := make([]byte, 188)
	n, err := r.Read(byts)
	require.NoError(t, err)
	require.Equal(t, 188, n)
	require.Equal(t, bytes.Repeat([]byte{1}, 188), byts[:n])

	n, err = r.Read(byts)
	require.NoError(t, err)
	require.Equal(t, 188, n)
	require.Equal(t, bytes.Repeat([]byte{2}, 188), byts[:n])

	n, err = r.Read(byts)
	require.NoError(t, err)
	require.Equal(t, 188, n)
	require.Equal(t, bytes.Repeat([]byte{3}, 188), byts[:n])

	require.Equal(t, 0, len(buf.Bytes()))

	_, err = r.Read(byts)
	require.EqualError(t, err, "received packet with size 0 not multiple of 188")
}

func TestBufferedReaderError(t *testing.T) {
	var buf bytes.Buffer

	buf.Write(bytes.Repeat([]byte{1}, 100))

	r := NewBufferedReader(&buf)
	byts := make([]byte, 188)
	_, err := r.Read(byts)
	require.EqualError(t, err, "received packet with size 100 not multiple of 188")
}
