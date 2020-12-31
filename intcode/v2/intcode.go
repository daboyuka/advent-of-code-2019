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

type (
	Instruction int
	Opcode      int
	Mode        int
)

func (in Instruction) Opcode() Opcode { return Opcode(in % 100) }
func (in Instruction) Modes(nParams int) (modes []Mode) {
	in /= 100
	for i := 0; i < nParams; i++ {
		modes = append(modes, Mode(in%10))
		in /= 10
	}
	return modes
}

func (prog Prog) Load(v int, mode Mode) int {
	switch mode {
	case 0:
		return prog[v]
	case 1:
		return v
	default:
		panic(fmt.Errorf("bad mode %d", mode))
	}
}

type ProgCtx struct {
	IO
}

func (in Instruction) Exec(prog Prog, pc int, ctx *ProgCtx) (newPC int, cont bool) {
	switch op := in.Opcode(); op {
	case 1:
		modes := in.Modes(3)
		p1, p2, p3 := prog[pc+1], prog[pc+2], prog[pc+3]
		prog[p3] = prog.Load(p1, modes[0]) + prog.Load(p2, modes[1])
		return pc + 4, true
	case 2:
		modes := in.Modes(3)
		p1, p2, p3 := prog[pc+1], prog[pc+2], prog[pc+3]
		prog[p3] = prog.Load(p1, modes[0]) * prog.Load(p2, modes[1])
		return pc + 4, true
	case 3:
		p := prog[pc+1]
		prog[p] = ctx.Input()
		return pc + 2, true
	case 4:
		modes := in.Modes(1)
		p := prog[pc+1]
		ctx.Output(prog.Load(p, modes[0]))
		return pc + 2, true
	case 5, 6:
		modes := in.Modes(2)
		p1, p2 := prog[pc+1], prog[pc+2]
		v1, v2 := prog.Load(p1, modes[0]), prog.Load(p2, modes[1])
		if (v1 != 0) == (op == 5) {
			return v2, true
		} else {
			return pc + 3, true
		}
	case 7, 8:
		modes := in.Modes(3)
		p1, p2, p3 := prog[pc+1], prog[pc+2], prog[pc+3]
		v1, v2 := prog.Load(p1, modes[0]), prog.Load(p2, modes[1])

		outV := 0
		if op == 7 && v1 < v2 {
			outV = 1
		} else if op == 8 && v1 == v2 {
			outV = 1
		}

		prog[p3] = outV
		return pc + 4, true
	case 99:
		return 0, false
	default:
		panic(in)
	}

	//Opcode 5 is jump-if-true: if the first parameter is non-zero, it sets the instruction pointer to the value from the second parameter. Otherwise, it does nothing.
	//Opcode 6 is jump-if-false: if the first parameter is zero, it sets the instruction pointer to the value from the second parameter. Otherwise, it does nothing.
	//Opcode 7 is less than: if the first parameter is less than the second parameter, it stores 1 in the position given by the third parameter. Otherwise, it stores 0.
	//Opcode 8 is equals: if the first parameter is equal to the second parameter, it stores 1 in the position given by the third parameter. Otherwise, it stores 0.
}

func (prog Prog) Run(ctx *ProgCtx) {
	pc := 0
	for {
		op := Instruction(prog[pc])
		newPC, cont := op.Exec(prog, pc, ctx)
		if !cont {
			return
		}
		pc = newPC
	}
}
