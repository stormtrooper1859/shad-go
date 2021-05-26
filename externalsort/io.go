// +build !change

package externalsort

import (
	"bufio"
	"io"
	"strings"
)

var (
	_ LineReader = &lineReader{}
	_ LineWriter = &lineWriter{}
)

type LineReader interface {
	ReadLine() (string, error)
}

type LineWriter interface {
	Write(l string) error
}

type lineReader struct {
	r *bufio.Reader
}

func (lr *lineReader) ReadLine() (string, error) {
	sb := strings.Builder{}
	isPrefix := true
	var err error
	var line []byte
	for isPrefix && err == nil {
		line, isPrefix, err = lr.r.ReadLine()
		sb.Write(line)
	}

	return sb.String(), err
}

type lineWriter struct {
	w io.Writer
}

func (lw *lineWriter) Write(l string) error {
	r := []byte(l)
	_, err := lw.w.Write(r)
	if err != nil {
		return err
	}
	_, err = lw.w.Write([]byte("\n"))
	return err
}
