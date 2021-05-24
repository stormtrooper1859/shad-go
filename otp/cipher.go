// +build !solution

package otp

import (
	"io"
)

const (
	readerBufferSize = 4096
	writerBufferSize = 4096
)

var (
	_ io.Reader = &cipherReader{}
	_ io.Writer = &cipherWriter{}
)

type cipherReader struct {
	reader io.Reader
	prng   io.Reader
}

type cipherWriter struct {
	writer io.Writer
	prng   io.Reader
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (c *cipherReader) Read(p []byte) (read int, err error) {
	n := len(p)

	rngBuffer := make([]byte, readerBufferSize)
	read = 0

	for read < n {
		readTo := min(read+readerBufferSize, n)
		currentRead, err := c.reader.Read(p[read:readTo])
		_, _ = c.prng.Read(rngBuffer[:currentRead])

		for j := 0; j < currentRead; j++ {
			p[read+j] ^= rngBuffer[j]
		}

		read += currentRead

		if err != nil {
			return read, err
		}
	}
	return read, nil
}

func (c *cipherWriter) Write(p []byte) (wrote int, err error) {
	n := len(p)

	rngBuffer := make([]byte, writerBufferSize)
	wrote = 0

	for wrote < n {
		writesLeft := min(writerBufferSize, n-wrote)
		_, _ = c.prng.Read(rngBuffer[:writesLeft])
		for j := 0; j < writesLeft; j++ {
			rngBuffer[j] ^= p[wrote+j]
		}
		currentWrote, err := c.writer.Write(rngBuffer[:writesLeft])
		wrote += currentWrote
		if err != nil {
			return wrote, err
		}
	}
	return wrote, nil
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return &cipherReader{r, prng}
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return &cipherWriter{w, prng}
}
