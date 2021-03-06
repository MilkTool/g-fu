package gfu

import (
	"bufio"
	//"log"
	"strings"
)

type ArgType int

const (
	ARG_PLAIN ArgType = 0
	ARG_OPT   ArgType = 1
	ARG_SPLAT ArgType = 2
)

type Arg struct {
	arg_type ArgType
	str_id   string
	id       *Sym
	opt_val  Val
}

func (a *Arg) Init(id *Sym) *Arg {
	a.id = id
	return a
}

func A(id string) (a Arg) {
	a.str_id = id
	return a
}

func AOpt(id string, val Val) (a Arg) {
	a.str_id = id
	a.opt_val = val
	a.arg_type = ARG_OPT
	return a
}

func ASplat(id string) (a Arg) {
	a.str_id = id
	a.arg_type = ARG_SPLAT
	return a
}

func (a Arg) DumpId(out *bufio.Writer) {
	if a.id == nil {
		out.WriteString("n/a")
	} else {
		out.WriteString(a.id.name)
	}
}

func (a Arg) Dump(g *G, out *bufio.Writer) E {
	switch a.arg_type {
	case ARG_OPT:
		out.WriteRune('(')
		a.DumpId(out)

		if a.opt_val != nil {
			out.WriteRune(' ')

			if e := g.Dump(a.opt_val, out); e != nil {
				return e
			}
		}

		out.WriteRune(')')
	case ARG_SPLAT:
		a.DumpId(out)
		out.WriteString("..")
	default:
		a.DumpId(out)
	}

	return nil
}

type Args []Arg

func (as Args) Dump(g *G, out *bufio.Writer) E {
	out.WriteRune('(')

	for i, a := range as {
		if i > 0 {
			out.WriteRune(' ')
		}

		a.Dump(g, out)
	}

	out.WriteRune(')')
	return nil
}

func (as Args) EString(g *G) string {
	var out strings.Builder
	w := bufio.NewWriter(&out)

	if e := as.Dump(g, w); e != nil {
		s, _ := g.String(e)
		return s

	}

	w.Flush()
	return out.String()
}

type ArgList struct {
	items    Args
	min, max int
}

func (l *ArgList) Init(g *G, args Args) *ArgList {
	nargs := len(args)

	if nargs == 0 {
		return l
	}

	l.items = args
	l.min, l.max = nargs, nargs

	for i, a := range l.items {
		if a.arg_type == ARG_OPT {
			l.min--
		}

		if a.id == nil && len(a.str_id) > 0 {
			l.items[i].id = g.Sym(a.str_id)
		}
	}

	a := l.items[nargs-1]

	if a.arg_type == ARG_SPLAT {
		l.min--
		l.max = -1
	}

	return l
}

func (l *ArgList) Check(g *G, args Vec) E {
	nargs := len(args)

	if (l.min != -1 && nargs < l.min) || (l.max != -1 && nargs > l.max) {
		return g.E("Arg mismatch: %v %v", l.items.EString(g), args)
	}

	return nil
}

func (l *ArgList) Fill(g *G, args Vec) Vec {
	for i := len(args); i < len(l.items); i++ {
		a := l.items[i]

		if a.arg_type != ARG_OPT {
			break
		}

		args = append(args, a.OptVal(g))
	}

	return args
}

func (a Arg) OptVal(g *G) Val {
	v := a.opt_val

	if v == nil {
		v = &g.NIL
	}

	return v
}

func (l *ArgList) LetVars(g *G, env *Env, args Vec) E {
	nargs := len(args)

	for i, a := range l.items {
		if a.id == nil {
			continue
		}

		if a.arg_type == ARG_SPLAT {
			var v Vec

			if i < nargs {
				v = make(Vec, nargs-i)
				copy(v, args[i:])
			}

			if e := env.Let(g, a.id, v); e != nil {
				return e
			}

			break
		}

		var v Val

		if i < nargs {
			v = args[i]
		} else {
			v = a.OptVal(g)
		}

		if e := env.Let(g, a.id, v); e != nil {
			return e
		}
	}

	return nil
}

func ParseArgs(g *G, task *Task, env *Env, in Vec, args_env *Env) (Args, E) {
	var e E
	var out Args

	for _, v := range in {
		var a Arg

		if v == &g.NIL {
			// Skip
		} else if id, ok := v.(*Sym); ok {
			a.id = id
		} else if vv, ok := v.(Vec); ok {
			if len(vv) < 2 {
				return nil, g.E("Invalid arg: %v", vv)
			}

			a.arg_type = ARG_OPT
			a.id = vv[0].(*Sym)

			if a.opt_val, e = g.Eval(task, env, vv[1], args_env); e != nil {
				return nil, e
			}
		} else if sv, ok := v.(Splat); ok {
			a.arg_type = ARG_SPLAT
			a.id = sv.val.(*Sym)
		} else {
			return nil, g.E("Invalid arg: %v", v)
		}

		out = append(out, a)
	}

	return out, nil
}
