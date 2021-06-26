package main

import "fmt"

func main4() {
	s := "string, строка"
	for k, v := range s {
		fmt.Printf("%v %c\n", k, v)
	}
	fmt.Println(s[12:])
}
