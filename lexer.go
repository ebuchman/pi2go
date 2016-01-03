package pi2go

// This lexer is heavily inspired by Rob Pike's "Lexical Scanning in Go"

import (
	"fmt"
	"log"
	"strings"
)

type lexStateFunc func(*lexer) lexStateFunc

// the lexer object
type lexer struct {
	input  string // input string to lex
	length int    // length of the input string
	pos    int    // current pos
	start  int    // start of current token

	line        int // current line number
	lastNewLine int // pos of last new line

	tokens chan token // channel to emit tokens over

	temp string // a place to hold eg. commands
}

// a token
type token struct {
	typ tokenType
	val string

	loc location
}

// location for error reporting
type location struct {
	line int
	col  int
}

// Lex the input, returning the lexer
// Tokens can be fetched off the channel
func Lex(input string) *lexer {
	l := &lexer{
		input:  input,
		length: len(input),
		pos:    0,
		tokens: make(chan token, 2),
	}
	go l.run()
	return l
}

func (l *lexer) Error(s string) lexStateFunc {
	return func(l *lexer) lexStateFunc {
		// TODO: print location data too
		log.Println(s)
		return nil
	}
}

// Return the tokens channel
func (l *lexer) Chan() chan token {
	return l.tokens
}

// Run the lexer
// This is the most beautiful function in the world
func (l *lexer) run() {
	for state := lexStateStart; state != nil; state = state(l) {
		// :D
	}
	close(l.tokens)
}

// Return next character in the string
// To hell with utf8 :p
func (l *lexer) next() string {
	if l.pos >= l.length {
		return ""
	}
	b := l.input[l.pos : l.pos+1]
	l.pos += 1
	return b
}

// backup a step
func (l *lexer) backup() {
	l.pos -= 1
}

// peek ahead a character without consuming
func (l *lexer) peek() string {
	s := l.next()
	l.backup()
	return s
}

// consume a token and push out on the channel
func (l *lexer) emit(ty tokenType) {
	l.tokens <- token{
		typ: ty,
		val: l.input[l.start:l.pos],
		loc: location{
			line: l.line,
			col:  l.pos - l.lastNewLine,
		},
	}
	l.start = l.pos
}

func (l *lexer) accept(options string) bool {
	if strings.Contains(options, l.next()) {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(options string) bool {
	i := 0
	for s := l.next(); strings.Contains(options, s); s = l.next() {
		i += 1
	}
	l.backup()
	return i > 0
}

// Starting state
func lexStateStart(l *lexer) lexStateFunc {
	// check the one character tokens
	t := l.next()
	switch t {
	case "":
		return nil
	case tokenLeftBrace:
		l.emit(tokenLeftBraceTy)
		return lexStateStart
	case tokenRightBrace:
		l.emit(tokenRightBraceTy)
		return lexStateStart
	case tokenNewLine:
		return lexStateNewLine
	case tokenPound:
		l.emit(tokenPoundTy)
		return lexStateComment
	case tokenFire:
		l.emit(tokenFireTy)
		return lexStateStart
	case tokenPull:
		l.emit(tokenPullTy)
		return lexStateStart
	case tokenChoice:
		l.emit(tokenChoiceTy)
		return lexStateStart
	case tokenPar:
		l.emit(tokenParTy)
		return lexStateStart
	case tokenDot:
		l.emit(tokenDotTy)
		return lexStateStart
	case tokenZero:
		l.emit(tokenZeroTy)
		return lexStateStart
	}
	l.backup()

	// skip spaces
	if isSpace(l.peek()) {
		return lexStateSpace
	}

	return lexStateExpressions

	return nil
}

func isSpace(s string) bool {
	return s == " " || s == "\t"
}

func lexStateExpressions(l *lexer) lexStateFunc {
	s := l.next()

	// check for chars
	if strings.Contains(tokenChars, s) {
		l.backup()
		return lexStateString
	}

	return l.Error(fmt.Sprintf("Invalid char: %s", s))
}

func lexStateNewLine(l *lexer) lexStateFunc {
	//for s := tokenNewLine; s == tokenNewLine; s = l.next() {
	//}
	//l.backup()
	l.emit(tokenNewLineTy)
	l.line += 1
	l.lastNewLine = l.pos
	return lexStateStart
}

// Scan past spaces
func lexStateSpace(l *lexer) lexStateFunc {
	for s := l.next(); isSpace(s); s = l.next() {
	}
	l.backup()
	l.start = l.pos
	return lexStateStart
}

// In a comment. Scan to new line
func lexStateComment(l *lexer) lexStateFunc {
	for r := ""; r != tokenNewLine; r = l.next() {
	}
	l.backup()
	l.emit(tokenStringTy)
	return lexStateStart
}

// a string
func lexStateString(l *lexer) lexStateFunc {
	if !l.acceptRun(tokenChars) {
		return l.Error("Expected a string")
	}
	l.emit(tokenStringTy)
	return lexStateStart
}

// error!
func lexStateErr(l *lexer) lexStateFunc {
	l.emit(tokenErrTy)
	return nil
}
