package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"./gfu"
)

var prof = flag.String("prof", "", "Write CPU profile to specified file")

func main() {
	g, e := gfu.NewG()

	if e != nil {
		log.Fatal(e)
	}

	g.RootEnv.InitAbc(g)
	g.RootEnv.InitIO(g)
	g.RootEnv.InitMath(g)
	g.RootEnv.InitTime(g)

	flag.Parse()

	if *prof != "" {
		f, e := os.Create(*prof)

		if e != nil {
			log.Fatal(e)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	args := flag.Args()
	env := g.NewEnv()

	if len(args) == 0 {
		fmt.Printf("g-fu v%v\n\nPress Return twice to evaluate.\n\n  ", g.VersionStr())
		in := bufio.NewScanner(os.Stdin)
		var buf strings.Builder

		for in.Scan() {
			line := in.Text()

			if len(line) == 0 {
				if buf.Len() > 0 {
					expr := fmt.Sprintf("(try _ %v)", buf.String())
					v, e := g.EvalString(&g.MainTask, env, gfu.INIT_POS, expr)

					if e == nil {
						fmt.Printf("%v\n", g.EString(v))
					} else {
						fmt.Printf("%v\n", g.EPrintString(e))
					}
				}

				buf.Reset()
			} else {
				buf.WriteString(line)
				buf.WriteRune('\n')
			}

			fmt.Printf("  ")
		}

		if e := in.Err(); e != nil {
			log.Fatal(e)
		}
	} else {
		for _, a := range args {
			if _, e := g.Load(&g.MainTask, env, env, a); e != nil {
				log.Fatal(g.EPrintString(e))
			}
		}
	}
}
