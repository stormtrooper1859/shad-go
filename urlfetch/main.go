// +build !solution

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	urls := os.Args[1:]

	for _, url := range urls {
		resp, err := http.Get(url)
		check(err)

		s, err := ioutil.ReadAll(resp.Body)
		check(err)

		fmt.Println(string(s))
	}
}
