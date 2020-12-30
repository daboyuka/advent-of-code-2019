package main

import (
	"fmt"
	"os"

	"aoc2019/intcode"
)

func A(prog intcode.Prog) {
	prog[1], prog[2] = 12, 2
	prog.Run()
	fmt.Println(prog[0])
}

func B(prog intcode.Prog) {
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
		A(intcode.Parse(os.Stdin))
	case "b":
		B(intcode.Parse(os.Stdin))
	default:
		panic(arg)
	}
}
