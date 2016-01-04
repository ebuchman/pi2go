package pi2go

import (
	"fmt"
	"io"
	"strings"
)

//------------------------
// process tree

type Process struct {
	names []string // new names

	// a process satisfies only one of these
	isZero bool             // zero process
	call   *ProcDefCall     // run a defined process
	sum    []*PrefixProcess // non-deterministic choice
	par    []*Process       // concurrent processes
}

type ProcDefCall struct {
	ID   string
	Args []string
}

func NewProcDefCall(id string, args ...string) *ProcDefCall {
	return &ProcDefCall{
		ID:   id,
		Args: args,
	}
}

type PrefixProcess struct {
	action *Action
	proc   *Process
}

// an action either sends or receives on a channel
type Action struct {
	typ     ActionType
	subject token
	object  token
}

func (a *Action) String() string {
	var s string
	if a.typ == ActionTypeFire {
		s = "!"
	} else {
		s = "?"
	}
	return fmt.Sprintf("%s%s(%s)", a.subject.val, s, a.object.val)
}

type ActionType int

const (
	ActionTypeFire ActionType = iota
	ActionTypePull
)

func tabs(i int) string {
	return strings.Repeat("\t", i)
}

//------------------------------------------------
// Print functions

// print the tree with lots of space
func RecursivePrint(p *Process, i int) {
	if p.isZero {
		fmt.Println(tokenZero)
		return
	}

	if len(p.names) > 0 {
		fmt.Println(tabs(i-1), p.names)
	}

	if p.call != nil {
		fmt.Printf("%s(%s)\n", p.call.ID, strings.Join(p.call.Args, ","))
	} else if len(p.sum) > 0 {
		if len(p.sum) > 1 {
			fmt.Printf("\n%s", tabs(i))
		}
		recursivePrintPrefix(p.sum[0], i+1)
		for _, p_ := range p.sum[1:] {
			fmt.Printf("\n%s", tabs(i-1))
			fmt.Println(tokenChoice)
			fmt.Printf("\n%s", tabs(i))
			recursivePrintPrefix(p_, i+1)
		}
	} else if len(p.par) > 0 {
		if len(p.par) > 1 {
			fmt.Printf("\n%s", tabs(i))
		}
		RecursivePrint(p.par[0], i+1)
		for _, p_ := range p.par[1:] {
			fmt.Printf("%s", tabs(i-1))
			fmt.Println(tokenPar)
			fmt.Printf("%s", tabs(i))
			RecursivePrint(p_, i+1)
		}

	}
}

func recursivePrintPrefix(p *PrefixProcess, i int) {
	fmt.Printf("%s.", p.action)
	RecursivePrint(p.proc, i)
}

// write the tree out on one line
type printer struct {
	io.Writer
}

func NewPrinter(w io.Writer) printer {
	return printer{w}
}

func (p printer) Printf(s string, args ...interface{}) {
	p.Write([]byte(fmt.Sprintf(s, args...)))
}

func (p printer) PrintParser(parser *parser) {
	p.printProcess(parser.P, false)
	fmt.Printf("\n")
}

func (p printer) printPrefixProcess(preProc *PrefixProcess, prefixed bool) {
	p.Printf(preProc.action.String())
	p.Printf(".")
	p.printProcess(preProc.proc, true)
}

func (p printer) printProcess(proc *Process, prefixed bool) {
	if proc.isZero {
		p.Printf("0")
		return
	}

	if len(proc.names) > 0 {
		p.Printf("new")
		p.Printf("%s", proc.names[0])
		for _, n := range proc.names[1:] {
			p.Printf(",%s", n)
		}
		p.Printf("in")
	}

	if proc.call != nil {
		p.Printf("%s(%s)", proc.call.ID, strings.Join(proc.call.Args, ","))
	} else if len(proc.sum) > 0 {
		sumL := len(proc.sum)
		if sumL > 1 {
			p.Printf("select{")
		}

		p.printPrefixProcess(proc.sum[0], prefixed)
		if sumL > 1 {
			for _, _proc := range proc.sum[1:] {
				p.Printf(" ;")
				p.printPrefixProcess(_proc, true)
			}
		}
		if sumL > 1 {
			p.Printf("}")
		}
	} else if len(proc.par) > 0 {
		p.printParallelProcess(proc.par, prefixed)
	} else {
		panic("wtf")
	}
}

func (p printer) printParallelProcess(procs []*Process, prefixed bool) {
	procL := len(procs)
	if procL > 1 {
		p.Printf("( ")
	}
	p.printProcess(procs[0], prefixed)
	if procL > 1 {
		for _, _proc := range procs[1:] {
			p.Printf(" | ")
			p.printProcess(_proc, prefixed)
		}
		p.Printf(" )")
	}
}
