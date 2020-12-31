package main

import (
	"fmt"
	"os"

	. "aoc2019/helpers"
	"aoc2019/intcode/v2"
)

type tile int
type pos struct{ Row, Col int }
type screen map[pos]tile

const (
	empty = tile(iota)
	wall
	block
	paddle
	ball
)

func runGame(prog intcode.Prog, screen screen) {
	io := intcode.ChanIO{
		Inputs:  make(chan int, 1),
		Outputs: make(chan int, 1),
	}
	go prog.Run(&intcode.ProgCtx{IO: &io})

	for {
		row, ok := <-io.Outputs
		if !ok {
			break
		}
		col := <-io.Outputs
		t := tile(<-io.Outputs)

		screen[pos{row, col}] = t
	}
}

func (s screen) CountTile(target tile) (c int) {
	for _, t := range s {
		if t == target {
			c++
		}
	}
	return c
}

//
//func render(out map[pos]paint) (lines []string) {
//	var minRow, minCol, maxRow, maxCol int
//	first := true
//	for p := range out {
//		if first {
//			minRow, minCol, maxRow, maxCol = p.Row, p.Col, p.Row, p.Col
//			first = false
//		} else {
//			minRow, minCol = Min(minRow, p.Row), Min(minCol, p.Col)
//			maxRow, maxCol = Max(maxRow, p.Row), Max(maxCol, p.Col)
//		}
//	}
//
//	grid := make([][]paint, maxRow-minRow+1)
//	for r := range grid {
//		grid[r] = make([]paint, maxCol-minCol+1)
//	}
//	for p, c := range out {
//		grid[p.Row-minRow][p.Col-minCol] = c
//	}
//
//	for _, lineBools := range grid {
//		linebuf := strings.Builder{}
//		for _, c := range lineBools {
//			if c {
//				linebuf.WriteRune(' ')
//			} else {
//				linebuf.WriteRune('█')
//			}
//		}
//		lines = append(lines, linebuf.String())
//	}
//	return lines
//}

func A(prog intcode.Prog) {
	screen := make(screen)
	runGame(prog, screen)
	fmt.Println(screen.CountTile(block))
}

func runGameHarder(prog intcode.Prog, screen screen) (score int) {
	ballRow, ballCol := -1, -1
	paddleRow, paddleCol := -1, -1

	joystick := 0

	var outBuf [3]int
	var outBufCur int
	prog.Run(&intcode.ProgCtx{IO: intcode.FunIO{
		Inputs: func() int {
			return joystick
		},
		Outputs: func(v int) {
			outBuf[outBufCur] = v
			outBufCur++
			if outBufCur != len(outBuf) {
				return
			}
			outBufCur = 0

			// Process command
			col, row, t := outBuf[0], outBuf[1], tile(outBuf[2])
			if row == 0 && col == -1 {
				score = int(t)
				return
			}

			if t == ball {
				ballRow, ballCol = row, col
			} else if t == paddle {
				paddleRow, paddleCol = row, col
			}

			screen[pos{row, col}] = t

			if ballCol != -1 && paddleCol != -1 {
				joystick = SignZero(ballCol - paddleCol)
				fmt.Println("paddle", paddleCol, "ball", ballCol, "joystick", joystick)
			}
		},
	}})
	return score
}

func B(prog intcode.Prog) {
	prog[0] = 2 // insert quarters

	screen := make(screen)
	score := runGameHarder(prog, screen)
	fmt.Println(score)
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
