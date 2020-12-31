package main

import (
	"fmt"
	"io"
	"os"

	. "aoc2019/helpers"
)

type pos3 [3]int

func ParsePositions(r io.Reader) (ps []pos3) {
	for _, line := range Readlines(r) {
		xs, ys, zs := Split3(line[1:len(line)-1], ", ", ", ")
		x, y, z := Atoi(xs[2:]), Atoi(ys[2:]), Atoi(zs[2:])
		ps = append(ps, pos3{x, y, z})
	}
	return ps
}

func gravityDim(ps []pos3, vs []pos3, dim int) {
	for i := range ps {
		for j := range ps[i+1:] {
			p1, p2, v1, v2 := &ps[i][dim], &ps[j+i+1][dim], &vs[i][dim], &vs[j+i+1][dim]
			s := SignZero(*p1 - *p2)
			*v1 -= s
			*v2 += s
		}
	}
}

func velocityDim(ps []pos3, vs []pos3, dim int) {
	for i := range ps {
		p, v := &ps[i][dim], &vs[i][dim]
		*p += *v
	}
}

func gravity(ps []pos3, vs []pos3) {
	for dim := range (pos3{}) {
		gravityDim(ps, vs, dim)
	}
}

func velocity(ps []pos3, vs []pos3) {
	for dim := range (pos3{}) {
		velocityDim(ps, vs, dim)
	}
}

func simulate(ps []pos3, vs []pos3, steps int) {
	fmt.Println(0, ps, vs)
	for i := 0; i < steps; i++ {
		gravity(ps, vs)
		velocity(ps, vs)
		fmt.Println(i+1, ps, vs)
	}
}

func energy(ps []pos3, vs []pos3) (energy int) {
	for i, p := range ps {
		v := vs[i]
		pot := Abs(p[0]) + Abs(p[1]) + Abs(p[2])
		kin := Abs(v[0]) + Abs(v[1]) + Abs(v[2])
		energy += pot * kin
	}
	return energy
}

func A(ps []pos3) {
	vs := make([]pos3, len(ps))
	simulate(ps, vs, 1000)

	fmt.Println(energy(ps, vs))
}

func cycleTimeDim(ps [4]pos3, vs [4]pos3, dim int) (steps int) {
	type dimCfg struct {
		P [4]int
		V [4]int
	}

	seens := map[dimCfg]int{}
	addSeen := func(step int) (lastSeenAgo int) {
		var cfg dimCfg
		for i, p := range ps {
			v := vs[i]
			cfg.P[i], cfg.V[i] = p[dim], v[dim]
		}
		lastStep, prevSeen := seens[cfg]
		seens[cfg] = step
		if prevSeen {
			return step - lastStep
		} else {
			return 0
		}
	}

	for step := 1; ; step++ {
		gravity(ps[:], vs[:])
		velocity(ps[:], vs[:])
		if cycle := addSeen(step); cycle > 0 {
			return cycle
		}
	}
}

func B(ps []pos3) {
	var psFixed, vsFixed [4]pos3
	for i, p := range ps {
		psFixed[i] = p
	}

	var cycleTimes [3]int
	for i := range (pos3{}) {
		cycleTimes[i] = cycleTimeDim(psFixed, vsFixed, i)
	}

	fmt.Println(cycleTimes)

	grandCycleTime := 1
	for _, cycleTime := range cycleTimes {
		gcd := GCD(grandCycleTime, cycleTime)
		fmt.Println(cycleTime, gcd)
		grandCycleTime = grandCycleTime / gcd * cycleTime
	}
	fmt.Println(grandCycleTime)
}

func main() {
	switch arg := os.Args[1]; arg {
	case "a":
		A(ParsePositions(os.Stdin))
	case "b":
		B(ParsePositions(os.Stdin))
	default:
		panic(arg)
	}
}
