package intcode

import (
	"fmt"
	"io"
	"strings"

	"aoc2019/helpers"
)

type Prog map[int]int

func Parse(r io.Reader) (prog Prog) {
	prog = make(Prog)
	for i, code := range helpers.Ints(strings.Split(helpers.Readlines(r)[0], ",")) {
		prog[i] = code
	}
	return prog
}

func (prog Prog) String() string {
	buf := strings.Builder{}
	visited := 0
	for i := 0; ; i++ {
		v, ok := prog[i]
		if !ok {
			continue
		}

		if buf.Len() > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(fmt.Sprintf("%d", v))

		visited++
		if visited == len(prog) {
			break
		}
	}
	return buf.String()
}

func (prog Prog) Copy() Prog {
	prog2 := make(Prog, len(prog))
	for i, v := range prog {
		prog2[i] = v
	}
	return prog2
}

type Opcode int

func (op Opcode) Exec(prog Prog, pc int) (newPC int, cont bool) {
	switch op {
	case 1:
		p1, p2, p3 := prog[pc+1], prog[pc+2], prog[pc+3]
		prog[p3] = prog[p1] + prog[p2]
		return pc + 4, true
	case 2:
		p1, p2, p3 := prog[pc+1], prog[pc+2], prog[pc+3]
		prog[p3] = prog[p1] * prog[p2]
		return pc + 4, true
	case 99:
		return 0, false
	default:
		panic(op)
	}
}

func (prog Prog) Run() {
	pc := 0
	for {
		op := Opcode(prog[pc])
		newPC, cont := op.Exec(prog, pc)
		if !cont {
			return
		}
		pc = newPC
	}
}
