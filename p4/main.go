package main

import (
	"fmt"
	"os"
	"strconv"

	. "aoc2019/helpers"
)

func searchPassword(min, max int, pred func(x int) bool) (matches int) {
	for p := min; p <= max; p++ {
		if pred(p) {
			matches++
		}
	}
	return matches
}

func A(line string) {
	aStr, bStr := Split2(line, "-")
	a, b := Atoi(aStr), Atoi(bStr)

	nMatches := searchPassword(a, b, func(pass int) bool {
		// 6 digits
		if pass < 1e5 || pass >= 1e6 {
			return false
		}

		passStr := strconv.Itoa(pass)

		// Monotonic
		prevDig := '0'
		for _, dig := range passStr {
			if dig < prevDig {
				return false
			}
			prevDig = dig
		}

		// Double-digit
		hasDoubleDig := false
		for i, dig := range passStr[1:] {
			if prevDig := rune(passStr[i]); prevDig == dig {
				hasDoubleDig = true
				break
			}
		}
		if !hasDoubleDig {
			return false
		}

		return true
	})
	fmt.Println(nMatches)
}

func B(line string) {
	aStr, bStr := Split2(line, "-")
	a, b := Atoi(aStr), Atoi(bStr)

	nMatches := searchPassword(a, b, func(pass int) bool {
		// 6 digits
		if pass < 1e5 || pass >= 1e6 {
			return false
		}

		passStr := strconv.Itoa(pass)

		// Monotonic
		prevDig := '0'
		for _, dig := range passStr {
			if dig < prevDig {
				return false
			}
			prevDig = dig
		}

		// Double-digit
		hasDoubleNotTripleDig := false
		for i, dig := range passStr {
			prevDig, nextDig, nextNextDig := 'A', 'B', 'C'
			if i > 0 {
				prevDig = rune(passStr[i-1])
			}
			if i < len(passStr)-1 {
				nextDig = rune(passStr[i+1])
			}
			if i < len(passStr)-2 {
				nextNextDig = rune(passStr[i+2])
			}

			if dig == nextDig && dig != nextNextDig && dig != prevDig {
				hasDoubleNotTripleDig = true
			}
		}
		if !hasDoubleNotTripleDig {
			fmt.Println(pass, "dig seq bad")
			return false
		}

		return true
	})
	fmt.Println(nMatches)
}

func main() {
	switch arg := os.Args[1]; arg {
	case "a":
		A(Readlines(os.Stdin)[0])
	case "b":
		B(Readlines(os.Stdin)[0])
	default:
		panic(arg)
	}
}
