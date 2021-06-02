// +build !solution

package varfmt

import (
	"fmt"
	"strconv"
	"strings"
)

func Sprintf(format string, args ...interface{}) string {
	sb := strings.Builder{}
	stringArgs := make([]string, len(args))

	prev := -1
	globalPatternIndex := 0

	for i, r := range format {
		if r == '{' {
			prev = i
		}
		if prev != -1 {
			if r == '}' {
				currentPatterIndex := globalPatternIndex
				if prev+1 != i {
					currentPatterIndex, _ = strconv.Atoi(format[prev+1 : i])
				}

				if stringArgs[currentPatterIndex] == "" {
					stringArgs[currentPatterIndex] = fmt.Sprintf("%v", args[currentPatterIndex])
				}
				sb.WriteString(stringArgs[currentPatterIndex])

				prev = -1
				globalPatternIndex += 1
			}
		} else {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}
