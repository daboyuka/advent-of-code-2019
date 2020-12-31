package main

import (
	"fmt"
	"os"

	//. "aoc2019/helpers"
	"aoc2019/intcode/v2"
)

func A(prog intcode.Prog) {
	io := intcode.SliceIO{Inputs: []int{1}}
	prog.Run(&intcode.ProgCtx{IO: &io})
	fmt.Println(io.Outputs)
}

func B(prog intcode.Prog) {
	io := intcode.SliceIO{Inputs: []int{2}}
	prog.Run(&intcode.ProgCtx{IO: &io})
	fmt.Println(io.Outputs)
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
