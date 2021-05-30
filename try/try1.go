package main

import "fmt"

func main1() {
	a := map[string]int{}
	for i := 0; i < 1024; i++ {
		a[fmt.Sprintf("ok%d", i)] = i
	}
	rez := 0
	for _, v := range a {
		rez += v
	}
	fmt.Println(rez)
}
