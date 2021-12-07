package main

import (
	"fmt"
	"os"
	"strings"

	. "aoc2019/helpers"
	"aoc2019/intcode/v2"
)

type tile byte
type pos struct{ Row, Col int }
type screen [][]tile

type posdir struct {
	pos
	dir
}

const (
	scaffold = tile('#')
	space    = tile('.')
	botNorth = tile('^')
	botSouth = tile('V')
	botEast  = tile('>')
	botWest  = tile('<')
	botDead  = tile('X')
)

type dir int

const (
	north = dir(1 + iota)
	south
	west
	east
	numDirs
)

var allDirs = [...]dir{north, south, west, east}

func (d dir) Left() dir {
	switch d {
	case north:
		return west
	case west:
		return south
	case south:
		return east
	case east:
		return north
	}
	panic(fmt.Errorf("bad dir %d", d))
}

func (d dir) BotTile() tile {
	switch d {
	case north:
		return botNorth
	case south:
		return botSouth
	case west:
		return botWest
	case east:
		return botEast
	default:
		panic(fmt.Errorf("bad dir %d", d))
	}
}

func (d dir) Right() dir { return d.Left().Left().Left() }

func (d dir) String() string { return [...]string{"north", "south", "west", "east"}[d-1] }

func (p pos) Add(d dir, amt int) pos {
	switch d {
	case north:
		return pos{p.Row - amt, p.Col}
	case west:
		return pos{p.Row, p.Col - amt}
	case south:
		return pos{p.Row + amt, p.Col}
	case east:
		return pos{p.Row, p.Col + amt}
	}
	panic(fmt.Errorf("bad dir %d", d))
}

func runCamera(prog intcode.Prog) (scr screen, bot posdir) {
	row, col := 0, 0
	prog.Run(&intcode.ProgCtx{IO: intcode.FunIO{
		Outputs: func(tRaw int) {
			t := tile(tRaw)
			if t == '\n' {
				row++
				col = 0
				return
			} else if t == botNorth {
				bot = posdir{pos{row, col}, north}
			} else if t == botSouth {
				bot = posdir{pos{row, col}, south}
			} else if t == botEast {
				bot = posdir{pos{row, col}, east}
			} else if t == botWest {
				bot = posdir{pos{row, col}, west}
			}
			if len(scr) <= row {
				scr = append(scr, []tile(nil))
			}
			scr[row] = append(scr[row], t)
			col++
		},
	}})
	return scr, bot
}

func (s screen) String() string {
	sb := strings.Builder{}
	for _, line := range s {
		for _, t := range line {
			sb.WriteByte(byte(t))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func (s screen) At(p pos) tile {
	if p.Row < 0 || p.Row >= len(s) {
		return space
	}
	row := s[p.Row]
	if p.Col < 0 || p.Col >= len(row) {
		return space
	}
	return row[p.Col]
}
func (s screen) Set(p pos, t tile) {
	if p.Row < 0 || p.Row >= len(s) {
		panic(fmt.Errorf("pos (%d, %d) out of bounds", p.Row, p.Col))
	}
	row := s[p.Row]
	if p.Col < 0 || p.Col >= len(row) {
		panic(fmt.Errorf("pos (%d, %d) out of bounds", p.Row, p.Col))
	}
	row[p.Col] = t
}

func findIntersections(scr screen) (inters []pos) {
	for i, line := range scr[1 : len(scr)-1] {
		for j, t := range line[1 : len(line)-1] {
			if t == '#' && scr[i+1][j] == '#' && scr[i][j+1] == '#' && scr[i+1][j+2] == '#' && scr[i+2][j+1] == '#' {
				inters = append(inters, pos{i + 1, j + 1})
			}
		}
	}
	return inters
}

func A(prog intcode.Prog) {
	scr, _ := runCamera(prog)
	fmt.Println(scr)
	inters := findIntersections(scr)
	fmt.Println(inters)
	for _, inter := range inters {
		scr[inter.Row][inter.Col] = 'O'
	}
	fmt.Println(scr)

	sum := 0
	for _, inter := range inters {
		sum += inter.Row * inter.Col
	}
	fmt.Println(sum)
}

func runBot(prog intcode.Prog, subM, subA, subB, subC string) {
	data := subM + "\n" + subA + "\n" + subB + "\n" + subC + "\n" + "n" + "\n"

	sb := strings.Builder{}
	lastVal := 0
	prog.Run(&intcode.ProgCtx{IO: intcode.FunIO{
		Inputs: func() int {
			v := int(data[0])
			data = data[1:]
			return v
		},
		Outputs: func(tRaw int) {
			lastVal = tRaw
			if tRaw == '\n' {
				sb.Reset()
			} else {
				sb.WriteByte(byte(tRaw))
			}
		},
	}})
	fmt.Println(lastVal)
}

func simBotSubroutine(scr screen, at posdir, prog string, paint tile) (posdir, bool) {
	for len(prog) > 0 {
		step, rest := Split2(prog, ",")
		prog = rest

		switch step {
		case "L":
			at.dir = at.dir.Left()
		case "R":
			at.dir = at.dir.Right()
		default:
			for i := 0; i < Atoi(step); i++ {
				if t := scr.At(at.pos); t == scaffold || t == botNorth || t == botSouth || t == botWest || t == botEast {
					scr.Set(at.pos, paint)
				} else {
					scr.Set(at.pos, '2')
				}

				at.pos = at.pos.Add(at.dir, 1)
				if scr.At(at.pos) == space {
					scr.Set(at.pos, 'X')
					return posdir{}, false
				}
			}
		}
	}
	scr.Set(at.pos, at.dir.BotTile())
	return at, true
}

func simBot(scr screen, at posdir, subM, subA, subB, subC string) {
	for ok := true; ok && len(subM) > 0; {
		step, rest := Split2(subM, ",")
		subM = rest
		switch step {
		case "A":
			at, ok = simBotSubroutine(scr, at, subA, 'A')
			fmt.Println(scr)
		case "B":
			at, ok = simBotSubroutine(scr, at, subB, 'B')
			fmt.Println(scr)
		case "C":
			at, ok = simBotSubroutine(scr, at, subC, 'C')
			fmt.Println(scr)
		}
	}

	fmt.Println(scr)
}

func B(prog intcode.Prog) {
	progBot := prog.Copy()
	progBot[0] = 2 // bot mode

	scr, bot := runCamera(prog)
	fmt.Println(scr)
	fmt.Println(bot)

	subM := "A,A,B,C,B,C,B,C,C,A"
	subA := "L,10,R,8,R,8"
	subB := "L,10,L,12,R,8,R,10"
	subC := "R,10,L,12,R,10"

	if len(subM) > 20 || len(subA) > 20 || len(subB) > 20 || len(subC) > 20 {
		panic(fmt.Errorf("bad programs %t %t %t %t", len(subM) > 20, len(subA) > 20, len(subB) > 20, len(subC) > 20))
	}

	//simBot(scr, bot, subM, subA, subB, subC)
	runBot(progBot, subM, subA, subB, subC)
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
