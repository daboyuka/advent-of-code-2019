package helpers

import (
	"strconv"
	"strings"
)

func Ints(lines []string) (out []int) {
	for _, line := range lines {
		out = append(out, Atoi(line))
	}
	return out
}

func Atoi(s string) int {
	x, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return x
}

func AtoiSafe(s string) int {
	x, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return x
}

func Between(x, low, high int) bool {
	return x >= low && x <= high
}

func Split2(s, sep string) (a, b string) {
	idx := strings.Index(s, sep)
	if idx == -1 {
		return s, ""
	} else {
		return s[:idx], s[idx+len(sep):]
	}
}

func Split3(s, sep1, sep2 string) (a, b, c string) {
	if idx1 := strings.Index(s, sep1); idx1 != -1 {
		a = s[:idx1]
		b, c = Split2(s[idx1+len(sep1):], sep2)
	} else {
		b = ""
		a, c = Split2(s, sep2)
	}
	return a, b, c
}

func Max(xs ...int) int {
	max := xs[0]
	for _, x := range xs[1:] {
		if max < x {
			max = x
		}
	}
	return max
}

func Min(xs ...int) int {
	min := xs[0]
	for _, x := range xs[1:] {
		if min > x {
			min = x
		}
	}
	return min
}

func Sum(xs ...int) (sum int) {
	for _, x := range xs {
		sum += x
	}
	return sum
}

// x in [min, max) or min or max
func Clamp(x, min, max int) int {
	if x < min {
		return min
	} else if x >= max {
		return max - 1
	} else {
		return x
	}
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func SignZero(x int) int {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	} else {
		return 0
	}
}

func GCD(a, b int) int {
	a, b = Abs(a), Abs(b)
	if a < b {
		a, b = b, a
	}

	if b == 0 {
		return a
	} else {
		return GCD(b, a%b)
	}
}

func LCM(a, b int) int {
	return a / GCD(a, b) * b
}
