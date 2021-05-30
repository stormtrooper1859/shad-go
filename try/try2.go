package main

import (
	"bytes"
	"io"
)

func main2() {
	debug := false

	//var buf *bytes.Buffer
	var buf io.Writer

	if debug {
		buf = new(bytes.Buffer)
	}

	buf = (*bytes.Buffer)(nil)
	f1(buf)
}

func f1(w io.Writer) {
	if w == nil {
		println("null")
	}
	if b, ok := w.(*bytes.Buffer); ok {
		println("bb")
		if b == nil {
			println("bb2")
		}
	}
}
