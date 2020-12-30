package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	. "aoc2019/helpers"
)

type intcodeProg map[int]int

func ParseIntcodeProg(r io.Reader) (prog intcodeProg) {
	prog = make(intcodeProg)
	for i, code := range Ints(strings.Split(Readlines(r)[0], ",")) {
		prog[i] = code
	}
	return prog
}

func (prog intcodeProg) String() string {
	buf := strings.Builder{}
	visited := 0
	for i := 0; ; i++ {
		v, ok := prog[i]
		if !ok {
			continue
		}

		if buf.Len() > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(fmt.Sprintf("%d", v))

		visited++
		if visited == len(prog) {
			break
		}
	}
	return buf.String()
}

func (prog intcodeProg) Copy() intcodeProg {
	prog2 := make(intcodeProg, len(prog))
	for i, v := range prog {
		prog2[i] = v
	}
	return prog2
}

func (prog intcodeProg) Run() (ok bool) {
	pos := 0
	for {
		switch op := prog[pos]; op {
		case 1:
			p1, p2, p3 := prog[pos+1], prog[pos+2], prog[pos+3]
			prog[p3] = prog[p1] + prog[p2]
			pos += 4
		case 2:
			p1, p2, p3 := prog[pos+1], prog[pos+2], prog[pos+3]
			prog[p3] = prog[p1] * prog[p2]
			pos += 4
		case 99:
			return true
		default:
			return false
		}
	}
}

func A(prog intcodeProg) {
	prog[1], prog[2] = 12, 2
	prog.Run()
	fmt.Println(prog[0])
}

func B(prog intcodeProg) {
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			prog := prog.Copy()
			prog[1], prog[2] = noun, verb
			prog.Run()
			if prog[0] == 19690720 {
				fmt.Println(noun, verb, 100*noun+verb)
				return
			}
		}
	}
}

func main() {
	switch arg := os.Args[1]; arg {
	case "a":
		A(ParseIntcodeProg(os.Stdin))
	case "b":
		B(ParseIntcodeProg(os.Stdin))
	default:
		panic(arg)
	}
}
