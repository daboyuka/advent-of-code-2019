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

func (prog Prog) Load(ctx *ProgCtx, v int, mode Mode) int {
	switch mode {
	case Position:
		return prog[v]
	case Value:
		return v
	case Relative:
		return prog[ctx.RB+v]
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
	Position = Mode(0)
	Value    = Mode(1)
	Relative = Mode(2)
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

func (in Instruction) LoadParamVals(prog Prog, ctx *ProgCtx, pkinds []ParamKind) (vals []int) {
	modes := in.Modes(len(pkinds))

	for i, pkind := range pkinds {
		raw := prog[ctx.PC+i+1]

		var v int
		switch pkind {
		case Load:
			v = prog.Load(ctx, raw, modes[i])
		case Store:
			v = raw
			if modes[i] == Relative {
				v += ctx.RB
			}
		default:
			panic(pkind)
		}

		vals = append(vals, v)
	}
	return vals
}

type ProgCtx struct {
	IO

	PC int
	RB int
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

func (in Instruction) Exec(prog Prog, ctx *ProgCtx) (cont bool) {
	switch op := in.Opcode(); op {
	case 1, 2:
		vals := in.LoadParamVals(prog, ctx, []ParamKind{Load, Load, Store})
		if op == 1 {
			prog[vals[2]] = vals[0] + vals[1]
		} else {
			prog[vals[2]] = vals[0] * vals[1]
		}
		ctx.PC += 4
	case 3:
		vals := in.LoadParamVals(prog, ctx, []ParamKind{Store})
		prog[vals[0]] = ctx.Input()
		ctx.PC += 2
	case 4:
		vals := in.LoadParamVals(prog, ctx, []ParamKind{Load})
		ctx.Output(vals[0])
		ctx.PC += 2
	case 5, 6:
		vals := in.LoadParamVals(prog, ctx, []ParamKind{Load, Load})
		if (vals[0] != 0) == (op == 5) {
			ctx.PC = vals[1]
		} else {
			ctx.PC += 3
		}
	case 7, 8:
		vals := in.LoadParamVals(prog, ctx, []ParamKind{Load, Load, Store})

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
		ctx.PC += 4
	case 9:
		vals := in.LoadParamVals(prog, ctx, []ParamKind{Load})
		ctx.RB += vals[0]
		ctx.PC += 2
	case 99:
		return false
	default:
		panic(in)
	}

	return true
}

func (prog Prog) Run(ctx *ProgCtx) {
	ctx.PC = 0
	for {
		op := Instruction(prog[ctx.PC])
		if cont := op.Exec(prog, ctx); !cont {
			ctx.IO.Done()
			return
		}
	}
}
