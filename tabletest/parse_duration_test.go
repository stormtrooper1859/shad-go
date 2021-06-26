package tabletest_test

import (
	"fmt"
	"gitlab.com/slon/shad-go/tabletest"
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	var tests = []struct {
		input     string
		result    time.Duration
		withError bool
	}{
		{"12m 32s", 0 * time.Second, true},
		{"9999999999999999999999ns", 0 * time.Second, true},
		{"9223372036854775807ns", 9223372036854775807 * time.Nanosecond, false},
		{"9223372036854775809ns", 0 * time.Second, true},
		{".9999999999999999999999h", 1 * time.Hour, false},
		{".9223372036854775809h", 3320413933267 * time.Nanosecond, false},
		{"12m32s", 752 * time.Second, false},
		{"1m1000s", 1060 * time.Second, false},
		{"-2m3s", -123 * time.Second, false},
		{"0h", 0 * time.Second, false},
		{"0", 0 * time.Second, false},
		{"", 0 * time.Second, true},
		{"+1.2h", 72 * time.Minute, false},
		{"-1.2h", -72 * time.Minute, false},
		{"1.2h.m", 0 * time.Second, true},
		{"1.2.2.2h", 0 * time.Second, true},
		{"1.2us", 1200 * time.Nanosecond, false},
		{"-1.2us", -1200 * time.Nanosecond, false},
		{"kek", 0 * time.Nanosecond, true},
	}

	for _, test := range tests {
		got, err := tabletest.ParseDuration(test.input)
		if got != test.result || ((err != nil) != test.withError) {
			t.Errorf("ParseDuration(%q) = %v", test.input, got)
			fmt.Println("Error: ", err)
		}
	}
}
