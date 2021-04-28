// +build !solution

package speller

var (
	numbers = map[int]string{
		1:  "one",
		2:  "two",
		3:  "three",
		4:  "four",
		5:  "five",
		6:  "six",
		7:  "seven",
		8:  "eight",
		9:  "nine",
		10: "ten",
		11: "eleven",
		12: "twelve",
		13: "thirteen",
		14: "fourteen",
		15: "fifteen",
		16: "sixteen",
		17: "seventeen",
		18: "eighteen",
		19: "nineteen",
	}
	tens = map[int]string{
		2: "twenty",
		3: "thirty",
		4: "forty",
		5: "fifty",
		6: "sixty",
		7: "seventy",
		8: "eighty",
		9: "ninety",
	}
	thousands = map[int]string{
		0: "",
		1: " thousand",
		2: " million",
		3: " billion",
	}
)

func spellBelowHundred(n int) string {
	if n < 20 {
		return numbers[n]
	}

	decs, ones := n/10, n%10

	rez := tens[decs]
	if ones != 0 {
		rez += "-" + numbers[ones]
	}
	return rez
}

func spellBelowThousand(n int) string {
	rez := ""
	if n >= 100 {
		rez += spellBelowHundred(n/100) + " hundred"
	}
	n %= 100
	if n != 0 {
		if len(rez) > 0 {
			rez += " "
		}
		rez += spellBelowHundred(n)
	}
	return rez
}

func Spell(n int64) string {
	tmp := make([]int64, 0)

	var rez = ""
	if n < 0 {
		rez += "minus "
		n *= -1
	}

	if n == 0 {
		return "zero"
	}

	for n > 0 {
		tmp = append(tmp, n%1000)
		n /= 1000
	}

	for i := len(tmp) - 1; i >= 0; i-- {
		if tmp[i] != 0 {
			rez += spellBelowThousand(int(tmp[i])) + thousands[i] + " "
		}
	}

	return rez[:len(rez)-1]
}
