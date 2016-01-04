package pi2go

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
type TokenType int

//go:generate stringer -type=TokenType

// token types
const (
	tokenErrTy TokenType = iota // error
	tokenEOFTy                  // end of file

	tokenLeftBraceTy        // (
	tokenRightBraceTy       // )
	tokenLeftSquareBraceTy  // [
	tokenRightSquareBraceTy // ]
	tokenLeftCurlBraceTy    // {
	tokenRightCurlBraceTy   // }
	tokenNewLineTy          // \n
	tokenPoundTy            // #
	tokenSpaceTy            // spaces
	tokenStringTy           // var names, comments

	tokenFireTy      // !
	tokenPullTy      // ?
	tokenChoiceTy    //+
	tokenParTy       // |
	tokenDotTy       // .
	tokenCommaTy     // ,
	tokenSemiColonTy // ;
	tokenEqualsTy    // =

	tokenZeroTy    // 0
	tokenTauTy     // tau
	tokenEpsilonTy // epsilon
	tokenNewTy     // new
	tokenInTy      // in
	tokenSelectTy  // select
	tokenNameTy    // channel_Name
	tokenCapsIDTy  // ProcDef
)

// tokens and special chars
var (
	tokenLeftBrace        = "("
	tokenRightBrace       = ")"
	tokenLeftSquareBrace  = "["
	tokenRightSquareBrace = "]"
	tokenLeftCurlBrace    = "{"
	tokenRightCurlBrace   = "}"
	tokenNewLine          = "\n"
	tokenPound            = "#"
	tokenSpace            = " "

	tokenFire      = "!"
	tokenPull      = "?"
	tokenChoice    = "+"
	tokenPar       = "|"
	tokenDot       = "."
	tokenComma     = ","
	tokenEquals    = "="
	tokenSemiColon = ";"

	tokenZero    = "0"
	tokenTau     = "tau"
	tokenEpsilon = "epsilon"
	tokenNew     = "new"
	tokenIn      = "in"
	tokenSelect  = "select"

	tokenLower  = "abcdefghijklmnopqrstuvwqyz"
	tokenUpper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	tokenLetter = tokenLower + tokenUpper
	tokenDigit  = "1234567890"
	tokenChar   = tokenLetter + tokenDigit + "_"
)
