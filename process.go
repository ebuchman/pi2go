package parse

import (
	"fmt"
	"io"
)

//------------------------
// process tree

type Process []*process // concurrent processes

func (p *Process) Append(proc ...*process) {
	*p = append(*p, proc...)
}

func (p *Process) Len() int {
	if p == nil {
		return 0
	}
	return len(*p)
}

type process struct {
	isZero bool     // if the process is 0
	sum    *sum     // if its a sum
	proc   *Process // if its a process in brackets
}

type sum []*prefixedProcess

func (s *sum) Append(proc ...*prefixedProcess) {
	*s = append(*s, proc...)
}

func (s *sum) Len() int {
	if s == nil {
		return 0
	}
	return len(*s)
}

type prefixedProcess struct {
	action *action
	proc   *Process
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
	p.printParallelProcess(proc.proc, true)
}

func (p printer) printProcess(proc *process, prefixed bool) {
	if proc.isZero {
		p.Printf("0")
	} else if proc.sum.Len() > 0 {
		if prefixed && proc.sum.Len() > 1 {
			p.Printf("( ")
		}
		sum := *proc.sum
		p.printPrefixProcess(sum[0], prefixed)
		if proc.sum.Len() > 1 {
			for _, _proc := range sum[1:] {
				p.Printf(" + ")
				p.printPrefixProcess(_proc, true)
			}
		}
		if prefixed && proc.sum.Len() > 1 {
			p.Printf(" )")
		}
	} else if proc.proc.Len() > 0 {
		p.printParallelProcess(proc.proc, prefixed)
	} else {
		panic("wtf")
	}
}

func (p printer) printParallelProcess(proc *Process, prefixed bool) {
	if proc.Len() > 1 {
		p.Printf("( ")
	}
	p.printProcess((*proc)[0], prefixed)
	if proc.Len() > 1 {
		for _, _proc := range (*proc)[1:] {
			p.Printf(" | ")
			p.printProcess(_proc, prefixed)
		}
		p.Printf(" )")
	}
}
