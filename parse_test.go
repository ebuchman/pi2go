package pi2go

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func StripSpaces(str string) string {
	return strings.Join(strings.Fields(str), "")
}

var text = []string{
	`select{alph?(bbb).c!(b).d!(e).0 ; g?(r).( c?(a).0 | MyFunc() | MyFunc(a) | MyFunc(a, b, c) | 0 | d!(b).0)}`,
	`(a?(b).0 | b!(a).0)`,
	`select{a?(b).0 ; b?(a).0}`,
	`a!(b).select{a?(b).0 ; c?(d).0}`,
	`( a?(b).select{c?(d).select{e?(f).0 ;g?(h).0} ; i?(j).0} | k?(l).0 )`,
	`( a?(b).select{c?(d).(select{e?(f).0 ;g?(h).0} | select{f?(e).0 ; h?(g).0}) ; i?(j).0} | k?(l).0 )`,
	`new a, b, c in (a?(b).0 | b?(a).0)`,
}

var tokens = []TokenType{
	tokenSelectTy,
	tokenLeftCurlBraceTy,
	tokenNameTy,
	tokenPullTy,
	tokenLeftBraceTy,
	tokenNameTy,
	tokenRightBraceTy,
	tokenDotTy,
	tokenNameTy,
	tokenFireTy,
	tokenLeftBraceTy,
	tokenNameTy,
	tokenRightBraceTy,
	tokenDotTy,
	tokenNameTy,
	tokenFireTy,
	tokenLeftBraceTy,
	tokenNameTy,
	tokenRightBraceTy,
	tokenDotTy,
	tokenZeroTy,

	tokenSemiColonTy,

	tokenNameTy,
	tokenPullTy,
	tokenLeftBraceTy,
	tokenNameTy,
	tokenRightBraceTy,
	tokenDotTy,

	tokenLeftBraceTy,
	tokenNameTy,
	tokenPullTy,
	tokenLeftBraceTy,
	tokenNameTy,
	tokenRightBraceTy,
	tokenDotTy,
	tokenZeroTy,
	tokenParTy,

	tokenCapsIDTy,
	tokenLeftBraceTy,
	tokenRightBraceTy,
	tokenParTy,

	tokenCapsIDTy,
	tokenLeftBraceTy,
	tokenNameTy,
	tokenRightBraceTy,
	tokenParTy,

	tokenCapsIDTy,
	tokenLeftBraceTy,
	tokenNameTy,
	tokenCommaTy,
	tokenNameTy,
	tokenCommaTy,
	tokenNameTy,
	tokenRightBraceTy,
	tokenParTy,

	tokenZeroTy,
	tokenParTy,

	tokenNameTy,
	tokenFireTy,
	tokenLeftBraceTy,
	tokenNameTy,
	tokenRightBraceTy,
	tokenDotTy,
	tokenZeroTy,

	tokenRightBraceTy,

	tokenRightCurlBraceTy,
}

func TestLexer(t *testing.T) {
	l := Lex(text[0])
	i := 0
	for tok := range l.Chan() {
		if tok.typ != tokens[i] {
			t.Fatalf("Got %s, expected %s. Token %d", tok.typ, tokens[i], i)
		}
		i += 1
	}
}

func TestParse(t *testing.T) {
	for _, t_ := range text {
		fmt.Println("Text:", t_)
		p := Parse(t_)

		//RecursivePrint(p.P, 0)

		b := new(bytes.Buffer)
		printer := NewPrinter(b)
		printer.PrintParser(p)
		got := StripSpaces(b.String())
		expected := StripSpaces(t_)
		if got != expected {
			t.Fatalf("Got %s \n Expected %s\n", got, expected)
		}
	}
}
