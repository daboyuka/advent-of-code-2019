package main

import (
	"fmt"
	"os"

	//. "aoc2019/helpers"
	"aoc2019/intcode/v2"
)

func runAmps(prog intcode.Prog, phases [5]int, feedback bool) int {
	var progs [5]intcode.Prog
	var ios [5]intcode.ChanIO
	var inChans [5]chan int
	for i := range progs {
		progs[i] = prog.Copy()
		inChans[i] = make(chan int, 2)
		inChans[i] <- phases[i]
	}

	inChans[0] <- 0 // initial input

	mainOutput := make(chan int, 2)
	for i := range progs {
		if i == len(progs)-1 {
			ios[i] = intcode.ChanIO{Inputs: inChans[i], Outputs: mainOutput}
		} else {
			ios[i] = intcode.ChanIO{Inputs: inChans[i], Outputs: inChans[i+1]}
		}
	}
	for i, prog := range progs {
		go prog.Run(&intcode.ProgCtx{IO: &ios[i]})
	}

	if !feedback {
		return <-mainOutput
	}

	lastOutputV := 0
	for v := range mainOutput {
		lastOutputV = v
		inChans[0] <- v
	}
	return lastOutputV
}

// [min, max)
func allCombosFrom(cur [5]int, at int, min, max int, f func(phases [5]int)) {
	if at == len(cur) {
		f(cur)
		return
	}

loop:
	for try := min; try < max; try++ {
		for _, v := range cur[:at] {
			if v == try {
				continue loop
			}
		}
		cur[at] = try
		allCombosFrom(cur, at+1, min, max, f)
	}
}

func A(prog intcode.Prog) {
	maxPower := 0
	var maxPowerPhases [5]int
	allCombosFrom([5]int{}, 0, 0, 5, func(phases [5]int) {
		power := runAmps(prog, phases, false)
		if maxPower < power {
			maxPower = power
			maxPowerPhases = phases
		}
	})
	fmt.Println(maxPowerPhases, maxPower)
}

func B(prog intcode.Prog) {
	maxPower := 0
	var maxPowerPhases [5]int
	allCombosFrom([5]int{}, 0, 5, 10, func(phases [5]int) {
		power := runAmps(prog, phases, true)
		if maxPower < power {
			maxPower = power
			maxPowerPhases = phases
		}
	})
	fmt.Println(maxPowerPhases, maxPower)
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
