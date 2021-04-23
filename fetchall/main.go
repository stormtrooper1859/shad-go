// +build !solution

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func fetch(url string, ch chan<- string) {
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	_ = resp.Body.Close()
	ch <- fmt.Sprintf("url: %s", url)
}

func main() {
	urls := os.Args[1:]
	ch := make(chan string)

	for _, url := range urls {
		go fetch(url, ch)
	}

	for _ = range urls {
		fmt.Println(<-ch)
	}
	fmt.Println("Complete")
}
