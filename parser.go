package parse

import (
	"fmt"
)

type parser struct {
	l    *lexer
	last token

	peekCount int // 1 if we've peeked

	P *Process // top level concurrent processes
}

func Parse(input string) *parser {
	l := Lex(input)
	p := &parser{
		l: l,
		P: new(Process),
	}
	return p
}

func (p *parser) next() token {
	if p.peekCount == 1 {
		p.peekCount = 0
		return p.last

	}
	p.last = <-p.l.Chan()
	return p.last
}

func (p *parser) peek() token {
	if p.peekCount == 1 {
		return p.last
	}
	p.next()
	p.peekCount = 1
	return p.last
}

func (p *parser) backup() {
	p.peekCount = 1
}

func (p *parser) run() {
	p.parseProcess(p.P, true, true)
}

func (p *parser) expect(typ tokenType) token {
	t := p.next()
	if t.typ != typ {
		start, fin := 0, len(p.l.input) // assumes only one line!
		col := t.loc.col
		if col-5 > 0 {
			start = col - 5
		}
		if col+5 < fin {
			fin = col + 5
		}

		context := p.l.input[start:fin]
		panic(fmt.Sprintf("Got %s, expected %s. Location %v. Context %s", t.typ, typ, t.loc, context))
	}
	return t
}

// parse possibly concurrent processes
func (p *parser) parseProcess(proc *Process, acceptChoice, acceptPar bool) {
	proc1 := new(process)

	// a concurrent process is either 0, a sum, or in brackets
	t := p.next()
	switch t.typ {
	case tokenZeroTy:
		proc1.isZero = true
	case tokenLeftBraceTy:
		proc1.proc = new(Process)
		p.parseProcess(proc1.proc, true, true)
		p.expect(tokenRightBraceTy)
	default:
		p.backup()
		proc1.sum = new(sum)
		p.parseSum(proc1.sum, acceptChoice, acceptPar)
	}

	proc.Append(proc1)

	// if there's a "|", parse the next concurrent process
	if acceptPar && p.peek().typ == tokenParTy {
		p.next()
		p.parseProcess(proc, acceptChoice, acceptPar)
	}
}

func (p *parser) parseSum(s *sum, acceptChoice, acceptPar bool) {
	proc := new(prefixedProcess)
	proc.proc = new(Process)

	subject := p.expect(tokenStringTy)

	var typ ActionType
	t := p.next()
	switch t.typ {
	case tokenFireTy:
		typ = ActionTypeFire
	case tokenPullTy:
		typ = ActionTypePull
	default:
		// XXX: error!
	}

	p.expect(tokenLeftBraceTy)
	object := p.expect(tokenStringTy)
	p.expect(tokenRightBraceTy)
	p.expect(tokenDotTy)

	p.parseProcess(proc.proc, false, false)

	proc.action = &action{typ, subject, object}

	s.Append(proc)

	if acceptChoice && p.peek().typ == tokenChoiceTy {
		p.next()
		p.parseSum(s, acceptChoice, acceptPar)
	}
}

/*
func parseStateStart(p *parser) parseStateFunc {
	t := p.next()
	// scan past spaces, new lines, and comments
	switch t.typ {
	case tokenErrTy:
		return nil
	case tokenNewLineTy, tokenSpaceTy:
		return parseStateStart
	case tokenPoundTy:
		return parseStateComment
	}

	return parseStateProcess
}*/
