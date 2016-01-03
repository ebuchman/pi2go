package parse

import (
	"fmt"
)

func (t token) String() string {
	s := fmt.Sprintf("Line %-2d, Col %-2d \t %-6s \t", t.loc.line, t.loc.col, t.typ.String())
	switch t.typ {
	case tokenEOFTy:
		return s + "EOF"
	case tokenErrTy:
		return s + t.val
	}
	/*if len(t.val) > 10 {
		return fmt.Sprintf("%.10q...", t.val)
	}*/
	return s + fmt.Sprintf("%q", t.val)
}

// token types
type tokenType int

func (t tokenType) String() string {
	switch t {
	case tokenErrTy:
		return "[Error]"
	case tokenEOFTy:
		return "[EOF]"
	case tokenLeftBraceTy:
		return "[LeftBrace]"
	case tokenRightBraceTy:
		return "[RightBrace]"
	case tokenNewLineTy:
		return "[NewLine]"
	case tokenPoundTy:
		return "[Pound]"
	case tokenSpaceTy:
		return "[Space]"
	case tokenStringTy:
		return "[String]"
	case tokenFireTy:
		return "[Fire]"
	case tokenPullTy:
		return "[Pull]"
	case tokenChoiceTy:
		return "[Choice]"
	case tokenParTy:
		return "[Par]"
	case tokenDotTy:
		return "[Dot]"
	case tokenZeroTy:
		return "[Zero]"
	}
	return "[Unknown]"
}

// token types
const (
	tokenErrTy tokenType = iota // error
	tokenEOFTy                  // end of file

	tokenLeftBraceTy  // (
	tokenRightBraceTy // )
	tokenNewLineTy    // \n
	tokenPoundTy      // #
	tokenSpaceTy      // spaces
	tokenStringTy     // var names, comments

	tokenFireTy   // !
	tokenPullTy   // ?
	tokenChoiceTy //+
	tokenParTy    // |
	tokenDotTy    // .
	tokenZeroTy   // 0
)

// tokens and special chars
var (
	tokenLeftBrace  = "("
	tokenRightBrace = ")"
	tokenNewLine    = "\n"
	tokenPound      = "#"
	tokenSpace      = " "

	tokenFire   = "!"
	tokenPull   = "?"
	tokenChoice = "+"
	tokenPar    = "|"
	tokenDot    = "."
	tokenZero   = "0"

	tokenChars = "abcdefghijklmnopqrstuvwqyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_"
)
