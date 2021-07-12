package main

import (
	"fmt"
	"unsafe"
)

type A struct {
	b int
	d int
}

func main() {
	c := A{
		b: 1,
	}
	fmt.Println(unsafe.Offsetof(c.d))
	fmt.Println(unsafe.Alignof(c.d))
}
