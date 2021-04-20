// +build !solution

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	fileNames := os.Args[1:]

	linesCounter := map[string]int{}

	for _, fileName := range fileNames {
		fileData, _ := ioutil.ReadFile(fileName)
		fileLines := strings.Split(string(fileData), "\n")

		for _, line := range fileLines {
			linesCounter[line] += 1
		}
	}

	for key, value := range linesCounter {
		if value > 1 {
			fmt.Printf("%d\t%s\n", value, key)
		}
	}
}
