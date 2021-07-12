package main

import (
	"fmt"
	"io"
	"reflect"

)

func main6() {
	var w io.Writer
	fmt.Println(reflect.TypeOf(w))


	var v interface{} = w
	fmt.Println(reflect.TypeOf(v))
}