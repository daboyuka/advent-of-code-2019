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

type (
	Instruction int
	Opcode      int
	Mode        int
	ParamKind   int
)

const (
	Load = ParamKind(iota)
	Store
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

func (in Instruction) LoadParamVals(prog Prog, pc int, pkinds []ParamKind) (vals []int) {
	modes := in.Modes(len(pkinds))

	for i, pkind := range pkinds {
		raw := prog[pc+i+1]

		var v int
		switch pkind {
		case Load:
			v = prog.Load(raw, modes[i])
		case Store:
			v = raw
		default:
			panic(pkind)
		}

		vals = append(vals, v)
	}
	return vals
}

type ProgCtx struct {
	IO
}

func (op Opcode) Params() []ParamKind {
	switch op {
	case 1, 2:
		return []ParamKind{Load, Load, Store}
	case 3:
		return []ParamKind{Store}
	case 4:
		return []ParamKind{Load}
	case 5, 6:
		return []ParamKind{Load, Load}
	case 7, 8:
		return []ParamKind{Load, Load, Store}
	case 99:
		return nil
	default:
		panic(op)
	}
}

func (in Instruction) Exec(prog Prog, pc int, ctx *ProgCtx) (newPC int, cont bool) {
	switch op := in.Opcode(); op {
	case 1, 2:
		vals := in.LoadParamVals(prog, pc, []ParamKind{Load, Load, Store})
		if op == 1 {
			prog[vals[2]] = vals[0] + vals[1]
		} else {
			prog[vals[2]] = vals[0] * vals[1]
		}
		return pc + 4, true
	case 3:
		vals := in.LoadParamVals(prog, pc, []ParamKind{Store})
		prog[vals[0]] = ctx.Input()
		return pc + 2, true
	case 4:
		vals := in.LoadParamVals(prog, pc, []ParamKind{Load})
		ctx.Output(vals[0])
		return pc + 2, true
	case 5, 6:
		vals := in.LoadParamVals(prog, pc, []ParamKind{Load, Load})
		if (vals[0] != 0) == (op == 5) {
			return vals[1], true
		} else {
			return pc + 3, true
		}
	case 7, 8:
		vals := in.LoadParamVals(prog, pc, []ParamKind{Load, Load, Store})

		cmp := false
		if op == 7 {
			cmp = vals[0] < vals[1]
		} else {
			cmp = vals[0] == vals[1]
		}

		if cmp {
			prog[vals[2]] = 1
		} else {
			prog[vals[2]] = 0
		}
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
			ctx.IO.Done()
			return
		}
		pc = newPC
	}
}
