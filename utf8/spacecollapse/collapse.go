// +build !solution

package spacecollapse

import "unicode"

func CollapseSpaces(input string) string {
	runes := make([]rune, 0)
	shouldSkip := false
	for _, v := range input {
		if unicode.IsSpace(v) {
			if shouldSkip {
				continue
			}
			shouldSkip = true
			runes = append(runes, ' ')
		} else {
			shouldSkip = false
			runes = append(runes, v)
		}
	}
	return string(runes)
}
