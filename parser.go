package pi2go

import (
	"fmt"
)

type parser struct {
	l    *lexer
	last token

	peekCount int // 1 if we've peeked

	P *Process // top level process
}

func Parse(input string) *parser {
	l := Lex(input)
	p := &parser{
		l: l,
		P: new(Process),
	}
	p.parseProcess(p.P, true, true)
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

func (p *parser) expect(typ TokenType) token {
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

//------------------------------------------------------------------------

// parse possibly concurrent processes
func (p *parser) parseProcess(proc *Process, acceptChoice, acceptPar bool) {
	proc1 := new(Process)

	t := p.next()
	switch t.typ {
	case tokenZeroTy:
		proc1.isZero = true
	case tokenNewTy:
		name := p.expect(tokenNameTy)
		proc1.names = append(proc1.names, name.val)
		for p.next().typ == tokenCommaTy {
			name := p.expect(tokenNameTy)
			proc1.names = append(proc1.names, name.val)
		}
		p.backup()
		p.expect(tokenInTy)
		p.parseProcess(proc1, true, true)
	case tokenSelectTy:
		p.expect(tokenLeftCurlBraceTy)
		p.parseSum(proc1, true, acceptPar)
		p.expect(tokenRightCurlBraceTy)
	case tokenLeftBraceTy:
		p.parseProcess(proc1, true, true)
		p.expect(tokenRightBraceTy)
	case tokenNameTy:
		p.backup()
		preProc := new(PrefixProcess)
		preProc.proc = new(Process)
		p.parsePrefixProc(preProc)
		proc1.sum = append(proc1.sum, preProc)
	case tokenCapsIDTy:
		p.expect(tokenLeftBraceTy)
		var names []string
		t2 := p.next()
		if t2.typ != tokenRightBraceTy {
			// grab the list of names (possibly just one)
			names = append(names, t2.val)
			for p.peek().typ != tokenRightBraceTy {
				p.expect(tokenCommaTy)
				t2 = p.next()
				names = append(names, t2.val)
			}
			p.expect(tokenRightBraceTy)
		}
		proc1.call = NewProcDefCall(t.val, names...)
	}

	proc.par = append(proc.par, proc1)

	// if there's a "|", parse the next concurrent process
	if acceptPar && p.peek().typ == tokenParTy {
		p.next()
		p.parseProcess(proc, acceptChoice, acceptPar)
	}
}

func (p *parser) parsePrefixProc(preProc *PrefixProcess) {
	subject := p.expect(tokenNameTy)

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
	object := p.expect(tokenNameTy)
	p.expect(tokenRightBraceTy)
	p.expect(tokenDotTy)

	preProc.action = &Action{typ, subject, object}

	p.parseProcess(preProc.proc, false, false)

}

func (p *parser) parseSum(sp *Process, acceptChoice, acceptPar bool) {
	preProc := new(PrefixProcess)
	preProc.proc = new(Process)

	p.parsePrefixProc(preProc)

	sp.sum = append(sp.sum, preProc)

	if acceptChoice && p.peek().typ == tokenSemiColonTy {
		p.next()
		p.parseSum(sp, acceptChoice, acceptPar)
	}
}
