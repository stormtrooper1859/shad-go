// +build !solution

package externalsort

import (
	"bufio"
	"container/heap"
	"io"
	"os"
	"sort"
)

func NewReader(r io.Reader) LineReader {
	return &lineReader{bufio.NewReader(r)}
}

func NewWriter(w io.Writer) LineWriter {
	return &lineWriter{w}
}

func Merge(w LineWriter, readers ...LineReader) error {
	pq := make(mergePriorityQueue, 0)
	for _, reader := range readers {
		str, err := reader.ReadLine()
		if err == io.EOF {
			continue
		} else if err != nil {
			return err
		}
		pq.Push(mergePQItem{str, reader})
	}

	for pq.Len() > 0 {
		err := w.Write(pq[0].value)
		if err != nil {
			return err
		}

		next, err := pq[0].reader.ReadLine()

		if err == nil {
			pq[0].value = next
			heap.Fix(&pq, 0)
		} else if err == io.EOF {
			heap.Remove(&pq, 0)
		} else {
			return err
		}
	}

	return nil
}

func sortFile(filename string) error {
	var err error
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	data := make([]string, 0)
	lr := NewReader(f)
	var str string
	for err != io.EOF {
		str, err = lr.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		data = append(data, str)
	}

	sort.Strings(data)

	err = f.Close()
	if err != nil {
		return err
	}

	f, err = os.Create(filename)
	if err != nil {
		return err
	}

	lw := NewWriter(f)
	for _, str := range data {
		err := lw.Write(str)
		if err != nil {
			return err
		}
	}

	err = f.Close()
	return err
}

func Sort(w io.Writer, in ...string) error {
	for _, filename := range in {
		err := sortFile(filename)
		if err != nil {
			return err
		}
	}

	files := make([]*os.File, 0)
	for _, v := range in {
		f, err := os.Open(v)
		if err != nil {
			return err
		}

		files = append(files, f)
	}

	readers := make([]LineReader, len(files))
	for i := 0; i < len(files); i++ {
		readers[i] = NewReader(files[i])
	}

	err := Merge(NewWriter( w), readers...)

	for _, v := range files {
		err := v.Close()
		if err != nil {
			return err
		}
	}

	return err
}
