package main

import (
	"fmt"
	"os"

	. "aoc2019/helpers"
)

func A(lines []string) {
	vs := Ints(lines)

	fuel := 0
	for _, v := range vs {
		fuel += (v/3)-2
	}
	fmt.Println(fuel)
}

func B(lines []string) {
	vs := Ints(lines)

	// Base fuel
	fuel := 0
	for _, v := range vs {
		modFuel := (v/3)-2
		uncoveredMass := modFuel
		for {
			coverFuel := (uncoveredMass/3)-2
			if coverFuel <= 0 {
				break
			}
			modFuel += coverFuel
			uncoveredMass = coverFuel
		}

		fuel += modFuel
	}

	fmt.Println(fuel)
}

func main() {
	switch arg := os.Args[1];arg {
	case "a":
		A(Readlines(os.Stdin))
	case"b":
		B(Readlines(os.Stdin))
	default:
		panic(arg)
	}
}
