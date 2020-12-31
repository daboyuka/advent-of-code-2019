package main

import (
	"fmt"
	"os"
	"strings"

	. "aoc2019/helpers"
	"aoc2019/intcode/v2"
)

type turn bool
type paint bool
type dir int

type pos struct{ Row, Col int }

func (d dir) Turn(t turn) dir {
	if t {
		return (d + 1) % 4
	} else {
		return (d + 3) % 4
	}
}

func (p pos) Move(d dir) pos {
	switch d {
	case north:
		return pos{p.Row - 1, p.Col}
	case south:
		return pos{p.Row + 1, p.Col}
	case east:
		return pos{p.Row, p.Col + 1}
	case west:
		return pos{p.Row, p.Col - 1}
	default:
		panic(fmt.Errorf("%d", d))
	}
}

const (
	black = paint(false)
	white = paint(true)
)

const (
	left  = turn(false)
	right = turn(true)
)

const (
	north = dir(iota)
	east
	south
	west
)

type move struct {
	paint
	dir
}

func runPainter(prog intcode.Prog, initial map[pos]paint) (out map[pos]paint) {
	var p pos
	var d dir
	out = initial

	io := intcode.ChanIO{
		Inputs:  make(chan int, 1),
		Outputs: make(chan int, 1),
	}
	go prog.Run(&intcode.ProgCtx{IO: &io})

	for {
		inV := 0
		if out[p] {
			inV = 1
		}
		io.Inputs <- inV

		outV1, ok := <-io.Outputs
		if !ok {
			break
		}
		outV2 := <-io.Outputs

		c := black
		if outV1 == 1 {
			c = white
		}

		t := left
		if outV2 == 1 {
			t = right
		}

		out[p] = c
		d = d.Turn(t)
		p = p.Move(d)
	}

	return out
}

func render(out map[pos]paint) (lines []string) {
	var minRow, minCol, maxRow, maxCol int
	first := true
	for p := range out {
		if first {
			minRow, minCol, maxRow, maxCol = p.Row, p.Col, p.Row, p.Col
			first = false
		} else {
			minRow, minCol = Min(minRow, p.Row), Min(minCol, p.Col)
			maxRow, maxCol = Max(maxRow, p.Row), Max(maxCol, p.Col)
		}
	}

	grid := make([][]paint, maxRow-minRow+1)
	for r := range grid {
		grid[r] = make([]paint, maxCol-minCol+1)
	}
	for p, c := range out {
		grid[p.Row-minRow][p.Col-minCol] = c
	}

	for _, lineBools := range grid {
		linebuf := strings.Builder{}
		for _, c := range lineBools {
			if c {
				linebuf.WriteRune(' ')
			} else {
				linebuf.WriteRune('█')
			}
		}
		lines = append(lines, linebuf.String())
	}
	return lines
}

func A(prog intcode.Prog) {
	out := runPainter(prog, map[pos]paint{})
	fmt.Println(len(out))
}

func B(prog intcode.Prog) {
	out := runPainter(prog, map[pos]paint{{0, 0}: white})

	lines := render(out)
	for _, line := range lines {
		fmt.Println(line)
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
