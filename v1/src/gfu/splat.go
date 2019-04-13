package gfu

import (
  //"log"
  "strings"
)

type Splat struct {
  Wrap
}

func NewSplat(val Val) (s Splat) {
  s.Wrap.val = val
  return s
}

func (s Splat) Call(g *G, args Vec, env *Env) (Val, E) {
  return s, nil 
}

func (s Splat) Dump(out *strings.Builder) {
  s.val.Dump(out)
  out.WriteString("..")
}

func (s Splat) Eq(g *G, rhs Val) bool {
  return s.val.Is(g, rhs.(Splat).val)
}

func (s Splat) Eval(g *G, env *Env) (v Val, e E) {
  v, e = s.val.Eval(g, env)

  if e != nil {
    return nil, e
  }
  
  return NewSplat(v), e
}

func (s Splat) Is(g *G, rhs Val) bool {
  return s == rhs
}

func (s Splat) Quote(g *G, env *Env) (v Val, e E) {
  if v, e = s.val.Quote(g, env); e != nil {
    return nil, e
  }

  return NewSplat(v), nil
}

func (s Splat) Splat(g *G, out Vec) Vec {
  v := s.val

  if _, ok := v.(Vec); !ok {
    return append(out, s)
  }

  return v.Splat(g, out)
}

func (s Splat) Type(g *G) *Type {
  return &g.SplatType
}

