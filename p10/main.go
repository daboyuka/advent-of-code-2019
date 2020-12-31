package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"sort"

	. "aoc2019/helpers"
)

type asteroids [][]bool // [row][col]

func ParseAsteroids(r io.Reader) (asters asteroids) {
	lines := Readlines(r)
	asters = make(asteroids, len(lines))
	for i, line := range lines {
		asters[i] = make([]bool, 0, len(line))
		for _, c := range line {
			asters[i] = append(asters[i], c == '#')
		}
	}
	return asters
}

type coord struct{ Row, Col int }

func (asters asteroids) CountLOS(from coord) (los int) {
	losAngles := map[[2]int]bool{}
	for row, asterRow := range asters {
		for col, isAster := range asterRow {
			if !isAster || (row == from.Row && col == from.Col) {
				continue // don't count non-asteroids or self
			}

			deltaRow, deltaCol := row-from.Row, col-from.Col
			deltaGCD := GCD(deltaRow, deltaCol)
			losAngles[[2]int{deltaRow / deltaGCD, deltaCol / deltaGCD}] = true
		}
	}
	return len(losAngles)
}

func (asters asteroids) CollectLOS(from coord) (astersByLOS map[coord][]int) { // [slope] -> depths
	astersByLOS = map[coord][]int{}
	for row, asterRow := range asters {
		for col, isAster := range asterRow {
			if !isAster || (row == from.Row && col == from.Col) {
				continue // don't count non-asteroids or self
			}

			deltaRow, deltaCol := row-from.Row, col-from.Col
			deltaGCD := GCD(deltaRow, deltaCol)

			slope := coord{deltaRow / deltaGCD, deltaCol / deltaGCD}
			astersByLOS[slope] = append(astersByLOS[slope], deltaGCD)
		}
	}
	return astersByLOS
}

func A(asters asteroids) {
	maxC := 0
	maxRow, maxCol := 0, 0
	for row, asterRow := range asters {
		for col, isAster := range asterRow {
			if isAster {
				if c := asters.CountLOS(coord{row, col}); maxC < c {
					maxC = c
					maxRow, maxCol = row, col
				}
			}
		}
	}
	fmt.Println(maxRow, maxCol, maxC)
}

// angle from a to b
func vecAngle(a, b coord) float64 {
	abcos := a.Row*b.Row + a.Col*b.Col
	absin := a.Row*b.Col - a.Col*b.Row
	angle := math.Atan2(float64(absin), float64(abcos))
	if angle < 0 {
		angle += 2 * math.Pi
	}
	return angle
}

func (c coord) Add(other coord) coord { return coord{c.Row + other.Row, c.Col + other.Col} }
func (c coord) Mul(x int) coord       { return coord{c.Row * x, c.Col * x} }

func B(asters asteroids) {
	maxC := 0
	maxRow, maxCol := 0, 0
	for row, asterRow := range asters {
		for col, isAster := range asterRow {
			if isAster {
				if c := asters.CountLOS(coord{row, col}); maxC < c {
					maxC = c
					maxRow, maxCol = row, col
				}
			}
		}
	}

	base := coord{maxRow, maxCol}
	los := asters.CollectLOS(base)

	type losCoord struct {
		SweepAngle float64 // angle in radians in range [0, 2*pi) measured from vertical (-1, 0) in the sweep direction; i.e., pi/2 == 90 deg means (0, 1)
		Slope      coord
		Depth      int
		AsterCoord coord
	}

	var losCoords []losCoord
	for slope, losAsters := range los {
		sort.Ints(losAsters)
		for asterDepth, slopeMul := range losAsters {
			lc := losCoord{
				SweepAngle: vecAngle(slope, coord{-1, 0}),
				Slope:      slope,
				Depth:      asterDepth,
				AsterCoord: slope.Mul(slopeMul).Add(base),
			}
			losCoords = append(losCoords, lc)
		}
	}
	sort.Slice(losCoords, func(i, j int) bool {
		a, b := losCoords[i], losCoords[j]
		if a.Depth < b.Depth {
			return true
		} else if a.Depth > b.Depth {
			return false
		} else {
			return a.SweepAngle < b.SweepAngle
		}
	})

	fmt.Println(base, losCoords[200-1], losCoords[200-1].AsterCoord.Col*100+losCoords[200-1].AsterCoord.Row)
}

func main() {
	switch arg := os.Args[1]; arg {
	case "a":
		A(ParseAsteroids(os.Stdin))
	case "b":
		B(ParseAsteroids(os.Stdin))
	default:
		panic(arg)
	}
}
