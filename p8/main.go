package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	. "aoc2019/helpers"
)

type image [][]int // [row][col]

func ParseImages(r io.Reader, rows, cols int) (layers []image) {
	digits := Readlines(r)[0]

	nLayers := len(digits) / rows / cols
	layers = make([]image, nLayers)
	for layerIdx := range layers {
		layers[layerIdx] = make(image, rows)
		for row := 0; row < rows; row++ {
			layers[layerIdx][row] = make([]int, cols)
		}
	}

	layerIdx, row, col := 0, 0, 0
	for _, c := range digits {
		layers[layerIdx][row][col] = int(c - '0')
		col++
		if col == cols {
			col = 0
			row++
			if row == rows {
				row = 0
				layerIdx++
			}
		}
	}
	return layers
}

func (img image) Count(targetDigit int) (c int) {
	for _, rowDigs := range img {
		for _, digit := range rowDigs {
			if digit == targetDigit {
				c++
			}
		}
	}
	return c
}

func (img image) Overlay(over image) image {
	composed := make(image, len(img))
	for row, rowDigs := range over {
		composed[row] = make([]int, len(rowDigs))
		for col, dig := range rowDigs {
			if dig != 2 {
				composed[row][col] = dig
			} else {
				composed[row][col] = img[row][col]
			}
		}
	}
	return composed
}

func (img image) Render() (lines []string) {
	for _, rowDigs := range img {
		linebuf := strings.Builder{}
		for _, dig := range rowDigs {
			if dig == 1 {
				linebuf.WriteRune('█')
			} else {
				linebuf.WriteRune(' ')
			}
		}
		lines = append(lines, linebuf.String())
	}
	return lines
}

func A(layers []image) {
	minCount0 := -1
	atMinCount1, atMinCount2, atMinLayerIdx := 0, 0, 0
	for layerIdx, layer := range layers {
		if count0 := layer.Count(0); minCount0 == -1 || minCount0 > count0 {
			minCount0 = count0
			atMinCount1, atMinCount2, atMinLayerIdx = layer.Count(1), layer.Count(2), layerIdx
		}
	}

	fmt.Println(minCount0, atMinCount1, atMinCount2, atMinLayerIdx, atMinCount1*atMinCount2)
}

func B(layers []image) {
	cur := layers[0]
	for _, layer := range layers[1:] {
		cur = layer.Overlay(cur)
	}
	fmt.Println(strings.Join(cur.Render(), "\n"))
}

func main() {
	switch arg := os.Args[1]; arg {
	case "a":
		A(ParseImages(os.Stdin, 6, 25))
	case "b":
		B(ParseImages(os.Stdin, 6, 25))
	default:
		panic(arg)
	}
}
