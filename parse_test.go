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
	`alph?(bbb).c!(b).d!(e).0 + g?(r).( c?(a).0 | 0 | d!(b).0)`,
	`(a?(b).0 | b!(a).0)`,
	`a!(b).(a?(b).0 + c?(d).0)`,
	`( a?(b).(c?(d).(e?(f).0 + g?(h).0) + i?(j).0) | k?(l).0 )`,
	`( a?(b).(c?(d).((e?(f).0 + g?(h).0) | (f?(e).0 + h?(g).0)) + i?(j).0) | k?(l).0 )`,
}

var tokens = []tokenType{
	tokenStringTy,
	tokenPullTy,
	tokenLeftBraceTy,
	tokenStringTy,
	tokenRightBraceTy,
	tokenDotTy,
	tokenStringTy,
	tokenFireTy,
	tokenLeftBraceTy,
	tokenStringTy,
	tokenRightBraceTy,
	tokenDotTy,
	tokenStringTy,
	tokenFireTy,
	tokenLeftBraceTy,
	tokenStringTy,
	tokenRightBraceTy,
	tokenDotTy,
	tokenZeroTy,

	tokenChoiceTy,

	tokenStringTy,
	tokenPullTy,
	tokenLeftBraceTy,
	tokenStringTy,
	tokenRightBraceTy,
	tokenDotTy,

	tokenLeftBraceTy,
	tokenStringTy,
	tokenPullTy,
	tokenLeftBraceTy,
	tokenStringTy,
	tokenRightBraceTy,
	tokenDotTy,
	tokenZeroTy,
	tokenParTy,

	tokenZeroTy,
	tokenParTy,

	tokenStringTy,
	tokenFireTy,
	tokenLeftBraceTy,
	tokenStringTy,
	tokenRightBraceTy,
	tokenDotTy,
	tokenZeroTy,
	tokenRightBraceTy,
}

func TestLexer(t *testing.T) {
	l := Lex(text[0])
	i := 0
	for tok := range l.Chan() {
		if tok.typ != tokens[i] {
			t.Fatal("Error", tok.typ, tokens[i])
		}
		i += 1
	}
}

// TODO: proper test
func TestParse(t *testing.T) {
	for _, t_ := range text {
		fmt.Println("Text:", t_)
		p := Parse(t_)
		p.run()
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
