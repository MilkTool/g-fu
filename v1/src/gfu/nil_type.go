package gfu

import (
  //"log"
  "strings"
)

type NilType struct {
  BasicType
}

func (t *NilType) Bool(g *G, val Val) bool {
  return false
}

func (t *NilType) Call(g *G, val Val, args Vec, env *Env) (Val, E) {
  return g.NIL, g.E("Nil call")
}
  
func (t *NilType) Dump(val Val, out *strings.Builder) {
  out.WriteRune('_')
}

func (t *NilType) Splat(g *G, val Val, out Vec) Vec {
  return out
}

