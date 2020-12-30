package helpers

import (
	"bufio"
	"io"
)

func Readlines(r io.Reader) (lines []string) {
	rbuf := bufio.NewReader(r)
	for {
		switch line, _, err := rbuf.ReadLine(); err {
		case nil:
			lines = append(lines, string(line))
		case io.EOF:
			return
		default:
			panic(err)
		}
	}
}

func ReadLinegroups(r io.Reader) (linegroups [][]string) {
	lines := Readlines(r)

	var curGroup []string
	for _, line := range lines {
		if line == "" {
			if len(curGroup) > 0 {
				linegroups = append(linegroups, curGroup)
				curGroup = nil
			}
		} else {
			curGroup = append(curGroup, line)
		}
	}
	if len(curGroup) > 0 {
		linegroups = append(linegroups, curGroup)
	}
	return linegroups
}
