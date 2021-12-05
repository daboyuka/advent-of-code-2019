package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	. "aoc2019/helpers"
)

type codestr []int

func toCode(s string) codestr {
	cs := make(codestr, len(s))
	for i, c := range s {
		cs[i] = int(c - '0')
	}
	return cs
}

func (cs codestr) String() string {
	sb := strings.Builder{}
	for _, b := range cs {
		sb.WriteRune(rune(b) + '0')
	}
	return sb.String()
}

func phase(in codestr) (out codestr) {
	out = make(codestr, len(in))
	sum := make([]int, len(in)+1) // cum[i] = sum of elements < i (exclusive bound), mod 10
	for i, v := range in {
		sum[i+1] = (sum[i] + v) % 10
	}

	slices := make([]int, 16)
	for i := range slices {
		slices[i] = len(out) * (i + 1) / len(slices)
	}

	var wg sync.WaitGroup
	sliceStart := 0
	for _, sliceEnd := range slices {
		wg.Add(1)
		go func(outFrom, outTo int) {
			defer wg.Done()
			for i := outFrom; i < outTo; i++ {
				digit := 0
				for from, to, add := i, 2*i+1, true; from < len(in); from, to, add = from+2*i+2, to+2*i+2, !add {
					rangeSum := sum[Min(to, len(in))] - sum[from] // sum interval is [from, to)
					if add {
						digit += rangeSum
					} else {
						digit -= rangeSum
					}
				}

				out[i] = Abs(digit) % 10
			}
		}(sliceStart, sliceEnd)
		sliceStart = sliceEnd
	}

	wg.Wait()
	return out
}

func multiphase(cs codestr, n int) codestr {
	for i := 0; i < n; i++ {
		cs = phase(cs)
		fmt.Println("completed phase", i+1, "of", n)
	}
	return cs
}

func A(lines []string) {
	cs := toCode(lines[0])
	cs = multiphase(cs, 100)
	fmt.Println(cs[:8])
}

func B(lines []string) {
	csBase := toCode(lines[0])
	var cs codestr
	for i := 0; i < 10000; i++ {
		cs = append(cs, csBase...)
	}

	fmt.Println("we've got", len(cs), "digits")

	off := Atoi(cs[:7].String())
	cs = multiphase(cs, 100)

	fmt.Println(off)
	fmt.Println(cs[off : off+8])
}

func main() {
	switch arg := os.Args[1]; arg {
	case "a":
		A(Readlines(os.Stdin))
	case "b":
		B(Readlines(os.Stdin))
	default:
		panic(arg)
	}
}
