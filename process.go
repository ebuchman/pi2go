package parse

import (
	"fmt"
	"io"
)

type Process struct {
	processes []*process // concurrent processes
}

type process struct {
	isZero bool     // if the process is 0
	sum    *sum     // if its a sum
	proc   *Process // if its a process in brackets
}

type sum struct {
	processes []*prefixedProcess
}

type prefixedProcess struct {
	action  *action
	process *Process
}

// an action either sends or receives on a channel
type action struct {
	typ     ActionType
	subject token
	object  token
}

func (a action) String() string {
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

//------------------------------------------------
// Print functions

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
	p.printParallelProcess(parser.P, false)
	fmt.Printf("\n")
}

func (p printer) printPrefixProcess(proc *prefixedProcess, prefixed bool) {
	p.Printf(proc.action.String())
	p.Printf(".")
	p.printParallelProcess(proc.process, true)
}

func (p printer) printProcess(proc *process, prefixed bool) {
	if proc.isZero {
		p.Printf("0")
	} else if proc.sum != nil && len(proc.sum.processes) > 0 {
		if prefixed && len(proc.sum.processes) > 1 {
			p.Printf("( ")
		}
		p.printPrefixProcess(proc.sum.processes[0], prefixed)
		if len(proc.sum.processes) > 1 {
			for _, _proc := range proc.sum.processes[1:] {
				p.Printf(" + ")
				p.printPrefixProcess(_proc, true)
			}
		}
		if prefixed && len(proc.sum.processes) > 1 {
			p.Printf(" )")
		}
	} else if proc.proc != nil && len(proc.proc.processes) > 0 {
		p.printParallelProcess(proc.proc, prefixed)
	} else {
		panic("wtf")
	}
}

func (p printer) printParallelProcess(proc *Process, prefixed bool) {
	if len(proc.processes) > 1 {
		p.Printf("( ")
	}
	p.printProcess(proc.processes[0], prefixed)
	if len(proc.processes) > 1 {
		for _, _proc := range proc.processes[1:] {
			p.Printf(" | ")
			p.printProcess(_proc, prefixed)
		}
		p.Printf(" )")
	}
}
