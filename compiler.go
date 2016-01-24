package pi2go

import (
	"fmt"
	"io"
)

func Compile(text string, w io.Writer) {
	c := compiler{w}
	p := Parse(text)
	_, _ = p, c
	// c.Compile(p.P)
}

type compiler struct {
	w io.Writer
}

//`(a?(b).0 | b!(a).0)`

func (c *compiler) Printf(s string, args ...interface{}) {
	c.w.Write([]byte(fmt.Sprintf(s, args...)))
}

func (c *compiler) Compile(p *Process) {
	c.compileProcess(p)
}

func (c *compiler) compileProcess(p *Process) {
	if len(p.sum) > 0 {
		c.compileSum(p.sum)
	} else if len(p.par) == 1 {
		c.Compile(p.par[0])
	} else if len(p.par) > 1 {
		for _, p_ := range p.par {
			// concurrent processes
			c.Printf("go func(){\n")
			c.compileProcess(p_)
			c.Printf("}()\n")
		}
	}
}

func (c *compiler) compileSum(s []*PrefixProcess) {
	if len(s) == 1 {
		// blocking single action
		proc := s[0]
		c.compileAction(proc.action)
		c.Printf("\n")
		c.Compile(proc.proc)
	} else {
		// select block
		c.Printf("select{\n")
		for _, p := range s {
			c.Printf("case ")
			c.compileAction(p.action)
			c.Printf(":\n")
			c.Compile(p.proc)
		}
		c.Printf("}\n")
	}
}

func (c *compiler) compileAction(a *Action) {
	switch a.typ {
	case ActionTypeFire:
		c.Printf("%s <- %s", a.subject.val, a.object.val)
	case ActionTypePull:
		c.Printf("%s := <- %s", a.object.val, a.subject.val)
	}
}
