package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	. "aoc2019/helpers"
	"aoc2019/intcode/v2"
)

type tile int
type pos struct{ Row, Col int }
type screen map[pos]tile

const (
	unvisited = tile(iota)
	wall
	droid
	oxygen
	visitedBase  // never appears; visitedX = visited + any dir
	visitedNorth // visitedX = visited, entered by moving X
	visitedSouth
	visitedWest
	visitedEast
)

func tileVisited(d dir) tile { return visitedBase + tile(d) }

func tileVisitedDir(t tile) dir { return dir(t - visitedBase) }

type dir int

const (
	north = dir(1 + iota)
	south
	west
	east
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

func (d dir) Right() dir { return d.Left().Left().Left() }

func (d dir) String() string { return [...]string{"north", "south", "west", "east"}[d-1] }

func (p pos) Add(d dir) pos {
	switch d {
	case north:
		return pos{p.Row - 1, p.Col}
	case west:
		return pos{p.Row, p.Col - 1}
	case south:
		return pos{p.Row + 1, p.Col}
	case east:
		return pos{p.Row, p.Col + 1}
	}
	panic(fmt.Errorf("bad dir %d", d))
}

func runDroid(prog intcode.Prog, screen screen) (oxyPos pos) {
	io := intcode.ChanIO{
		Inputs:  make(chan int, 1),
		Outputs: make(chan int, 1),
	}
	go prog.Run(&intcode.ProgCtx{IO: &io})

	at := pos{0, 0}
	facing := north
	screen[at] = tileVisited(facing)
	hit := 0
	for {
		nextpos := at.Add(facing)
		next := screen[nextpos]
		if next == wall || (next >= visitedBase && tileVisitedDir(next) == facing) { // if we hit a wall or loop, turn left
			hit++
			if hit == 4 {
				return oxyPos
			}
			facing = facing.Right()
			fmt.Println("turned", facing)
			continue
		}

		io.Inputs <- int(facing)   // move
		status, ok := <-io.Outputs // read status
		if !ok {
			panic(fmt.Errorf("unexpected end of program"))
		}

		// If we hit a wall, record it and do next step
		if status == 0 {
			screen[nextpos] = wall
			fmt.Println("hit wall to", facing)
			//for _, s := range screen.render(at) {
			//	fmt.Println(s)
			//}
			continue
		}

		// Otherwise, move
		at = nextpos
		screen[at] = tileVisited(facing)
		facing = facing.Left()
		hit = 0

		fmt.Println("moved", facing)
		//for _, s := range screen.render(at) {
		//	fmt.Println(s)
		//}

		if status == 2 {
			oxyPos = at
		} else if at == (pos{0, 0}) && facing == north {
			fmt.Println("back to start")
			return oxyPos
		}
	}
}

func (s screen) render(at pos) (lines []string) {
	var minRow, minCol, maxRow, maxCol int
	first := true
	for p := range s {
		if first {
			minRow, minCol, maxRow, maxCol = p.Row, p.Col, p.Row, p.Col
			first = false
		} else {
			minRow, minCol = Min(minRow, p.Row), Min(minCol, p.Col)
			maxRow, maxCol = Max(maxRow, p.Row), Max(maxCol, p.Col)
		}
	}

	grid := make([][]tile, maxRow-minRow+1)
	for r := range grid {
		grid[r] = make([]tile, maxCol-minCol+1)
	}
	for p, t := range s {
		if p == at {
			t = droid
		}
		grid[p.Row-minRow][p.Col-minCol] = t
	}

	for _, lineTiles := range grid {
		linebuf := strings.Builder{}
		for _, t := range lineTiles {
			switch t {
			case unvisited:
				linebuf.WriteRune(' ')
			case wall:
				linebuf.WriteRune('█')
			case oxygen:
				linebuf.WriteRune('O')
			case droid:
				linebuf.WriteRune('D')
			case visitedNorth:
				linebuf.WriteRune('⇧')
			case visitedSouth:
				linebuf.WriteRune('⇩')
			case visitedWest:
				linebuf.WriteRune('⇦')
			case visitedEast:
				linebuf.WriteRune('⇨')
			}
		}
		lines = append(lines, linebuf.String())
	}
	return lines
}

func shortestPath(from, to pos, s screen) int {
	type toVisit struct {
		p    pos
		dist int
	}

	visited := make(map[pos]int, len(s))

	sortedQ := []toVisit{{from, 0}}
	for len(sortedQ) > 0 {
		at := sortedQ[0]
		sortedQ = sortedQ[1:]

		if at.p == to {
			return at.dist
		}

		_, ok := visited[at.p]
		if ok {
			continue
		}
		visited[at.p] = at.dist

		fmt.Println("visited", at.p, "dist", at.dist)

		for _, d := range allDirs {
			next := toVisit{at.p.Add(d), at.dist + 1}
			if _, ok := visited[next.p]; !ok && s[next.p] != wall {

				insertAt := sort.Search(len(sortedQ), func(idx int) bool { return next.dist < sortedQ[idx].dist })
				sortedQ = append(sortedQ[:insertAt], toVisit{})
				copy(sortedQ[insertAt+1:], sortedQ[insertAt:])
				sortedQ[insertAt] = next

			}
		}
	}

	panic(fmt.Errorf("didn't find target somehow"))
}

func A(prog intcode.Prog) {
	scr := make(screen)
	oxyPos := runDroid(prog, scr)
	fmt.Println(oxyPos)
	for _, s := range scr.render(pos{}) {
		fmt.Println(s)
	}
	dist := shortestPath(pos{}, oxyPos, scr)
	fmt.Println(dist)
}

func floodFill(from pos, s screen) int {
	toFill := []pos{from}

	for round := 0; ; round++ {
		for _, p := range toFill {
			s[p] = oxygen
		}

		var nextFill []pos
		for _, p := range toFill {
			for _, d := range allDirs {
				p2 := p.Add(d)
				if s[p2] != wall && s[p2] != oxygen {
					nextFill = append(nextFill, p2)
				}
			}
		}

		toFill = nextFill

		fmt.Println("after round", round)
		for _, s := range s.render(pos{9999, 9999}) {
			fmt.Println(s)
		}

		if len(toFill) == 0 {
			return round
		}
	}
}

func B(prog intcode.Prog) {
	scr := make(screen)
	oxyPos := runDroid(prog, scr)
	fmt.Println(oxyPos)
	for _, s := range scr.render(pos{}) {
		fmt.Println(s)
	}

	finalRound := floodFill(oxyPos, scr)
	fmt.Println(finalRound)
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
