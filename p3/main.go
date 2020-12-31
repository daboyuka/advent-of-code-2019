package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	. "aoc2019/helpers"
)

type move struct {
	Dir byte
	Amt int
}

type point struct{ R, U int }

func (p point) Add(m move) point {
	switch m.Dir {
	case 'U':
		return point{p.R, p.U + m.Amt}
	case 'D':
		return point{p.R, p.U - m.Amt}
	case 'R':
		return point{p.R + m.Amt, p.U}
	case 'L':
		return point{p.R - m.Amt, p.U}
	default:
		panic(m.Dir)
	}
}

type wire []point

func parseWires(r io.Reader) (ws []wire) {
	for _, line := range Readlines(r) {
		w := wire{{0, 0}}
		for i, step := range strings.Split(line, ",") {
			m := move{Dir: step[0], Amt: Atoi(step[1:])}
			w = append(w, w[i].Add(m))
		}
		ws = append(ws, w)
	}
	return ws
}

func normLine(start, end point) (point, point) {
	if vert := start.R == end.R; vert {
		if start.U < end.U {
			return start, end
		} else {
			return end, start
		}
	} else {
		if start.R < end.R {
			return start, end
		} else {
			return end, start
		}
	}
}

func IntersectLines(aStart, aEnd, bStart, bEnd point) (inters []point) {
	aStart, aEnd = normLine(aStart, aEnd)
	bStart, bEnd = normLine(bStart, bEnd)

	aVert, bVert := aStart.R == aEnd.R, bStart.R == bEnd.R
	if aVert && bVert {
		if aStart.R != bStart.R {
			return nil
		}

		// Normalize A to be before B
		if aStart.U > bStart.U {
			aStart, aEnd, bStart, bEnd = bStart, bEnd, aStart, aEnd
		}
		for u := bStart.U; u <= aEnd.U; u++ {
			inters = append(inters, point{aStart.R, u})
		}
	} else if !aVert && !bVert {
		if aStart.U != bStart.U {
			return nil
		}

		// Normalize A to be before B
		if aStart.R > bStart.R {
			aStart, aEnd, bStart, bEnd = bStart, bEnd, aStart, aEnd
		}
		for r := bStart.R; r <= aEnd.R; r++ {
			inters = append(inters, point{r, aStart.U})
		}
	} else {
		// Normalize so A is vertical
		if !aVert {
			aStart, aEnd, bStart, bEnd = bStart, bEnd, aStart, aEnd
		}
		if Between(aStart.R, bStart.R, bEnd.R) && Between(bStart.U, aStart.U, aEnd.U) {
			inters = []point{{aStart.R, bStart.U}}
		}
	}
	return inters
}

func IntersectWires(a, b wire) (allInters []point, aSegs []int, bSegs []int) {
	for ai, ap := range a[1:] {
		for bi, bp := range b[1:] {
			ap2, bp2 := a[ai], b[bi]
			if inters := IntersectLines(ap, ap2, bp, bp2); inters != nil {
				allInters = append(allInters, inters...)
				for range inters {
					aSegs, bSegs = append(aSegs, ai), append(bSegs, bi)
				}
			}
		}
	}
	return allInters, aSegs, bSegs
}

func (p point) DistTo(other point) int { return Abs(p.R-other.R) + Abs(p.U-other.U) }

func (w wire) DistTo(inter point, seg int) (dist int) {
	for i, p := range w[1 : seg+1] {
		dist += w[i].DistTo(p)
	}
	dist += w[seg].DistTo(inter)
	return dist
}

func A(wires []wire) {
	inters, _, _ := IntersectWires(wires[0], wires[1])

	minDist := -1
	pMin := point{}
	for _, inter := range inters {
		if inter == (point{}) {
			continue
		} else if dist := Abs(inter.R) + Abs(inter.U); minDist == -1 || minDist > dist {
			minDist = dist
			pMin = inter
		}
	}
	fmt.Println(minDist, pMin.R, pMin.U)
}

func B(wires []wire) {
	inters, segs1, segs2 := IntersectWires(wires[0], wires[1])

	minDist := -1
	minD1, minD2 := -1, -1
	for i, inter := range inters {
		if inter == (point{}) {
			continue
		}

		seg1, seg2 := segs1[i], segs2[i]
		d1, d2 := wires[0].DistTo(inter, seg1), wires[1].DistTo(inter, seg2)

		if dist := d1 + d2; minDist == -1 || minDist > dist {
			minDist = dist
			minD1, minD2 = d1, d2
		}
	}

	fmt.Println(minDist, minD1, minD2)
}

func main() {
	switch arg := os.Args[1]; arg {
	case "a":
		A(parseWires(os.Stdin))
	case "b":
		B(parseWires(os.Stdin))
	default:
		panic(arg)
	}
}
