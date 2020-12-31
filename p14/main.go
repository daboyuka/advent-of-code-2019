package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	. "aoc2019/helpers"
)

type reagent string

type reagentQuant struct {
	Kind  reagent
	Quant int
}

type reaction struct {
	Ins []reagentQuant
	Out reagentQuant
}

func ParseReagentQuant(s string) reagentQuant {
	a, b := Split2(s, " ")
	return reagentQuant{Kind: reagent(b), Quant: Atoi(a)}
}

func ParseReactions(r io.Reader) (reactions map[reagent]reaction) {
	reactions = make(map[reagent]reaction)
	for _, line := range Readlines(r) {
		insS, outS := Split2(line, " => ")

		r := reaction{Out: ParseReagentQuant(outS)}
		for _, inS := range strings.Split(insS, ", ") {
			r.Ins = append(r.Ins, ParseReagentQuant(inS))
		}

		reactions[r.Out.Kind] = r
	}
	return reactions
}

// Sorts reactions such that all inputs for a reaction are later in the slice
func TopoSort(reactions map[reagent]reaction, final reagent) (sorted []reaction) {
	// produce reverse sort
	visited := make(map[reagent]bool, len(reactions))
	var dft func(reagent)
	dft = func(product reagent) {
		if visited[product] {
			return
		}
		visited[product] = true
		r := reactions[product]
		for _, in := range r.Ins {
			dft(in.Kind)
		}
		sorted = append(sorted, r)
	}
	dft(final)

	// reverse slice
	for i := range sorted[:len(sorted)/2] {
		sorted[i], sorted[len(sorted)-i-1] = sorted[len(sorted)-i-1], sorted[i]
	}
	return sorted
}

func Produce(rawInput, output reagent, outputQ int, reactions map[reagent]reaction) (inputQ int) {
	need := map[reagent]int{output: outputQ}
	ordered := TopoSort(reactions, output)
	for _, r := range ordered {
		toMake := need[r.Out.Kind]
		if toMake > 0 {
			mult := (toMake + r.Out.Quant - 1) / r.Out.Quant
			need[r.Out.Kind] -= r.Out.Quant * mult
			for _, in := range r.Ins {
				need[in.Kind] += in.Quant * mult
			}
		}
	}
	return need[rawInput]
}

func A(reactions map[reagent]reaction) {
	need := map[reagent]int{"FUEL": 1}
	ordered := TopoSort(reactions, "FUEL")
	for _, r := range ordered {
		toMake := need[r.Out.Kind]
		if toMake > 0 {
			mult := (toMake + r.Out.Quant - 1) / r.Out.Quant
			need[r.Out.Kind] -= r.Out.Quant * mult
			for _, in := range r.Ins {
				need[in.Kind] += in.Quant * mult
			}
		}
	}

	fmt.Println(Produce("ORE", "FUEL", 1, reactions))
}

func B(reactions map[reagent]reaction) {
	const ore = 1e12

	approxOrePerFuel := Produce("ORE", "FUEL", 1, reactions)

	maxFuel := 0
	maxNeed := 0
	for fuel := ore / approxOrePerFuel; ; fuel++ {
		need := Produce("ORE", "FUEL", fuel, reactions)
		fmt.Println(need, fuel)
		if need > ore {
			break
		}
		maxFuel, maxNeed = fuel, need
	}
	fmt.Println(maxNeed, maxFuel)
}

func main() {
	switch arg := os.Args[1]; arg {
	case "a":
		A(ParseReactions(os.Stdin))
	case "b":
		B(ParseReactions(os.Stdin))
	default:
		panic(arg)
	}
}
