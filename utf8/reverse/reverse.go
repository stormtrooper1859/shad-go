// +build !solution

package reverse

import "unicode/utf8"

func Reverse(input string) string {
	n := utf8.RuneCountInString(input)
	runes := make([]rune, n)
	i := 0
	for _, v := range input {
		runes[n-i-1] = v
		i++
	}
	return string(runes)
}
