package bechars

import (
	"src.userspace.com.au/lexer"
)

const (
	_ lexer.TokenType = iota
	tBREStart
	tBREEnd
	tRangeStart
	tRangeDash
	tRangeEnd
	tCharacter
	tClass
	tNot
)

func startState(l *lexer.Lexer) lexer.StateFunc {
	l.SkipWhitespace()
	r := l.Next()
	if r != '[' {
		return l.Error("expecting [")
	}
	l.Emit(tBREStart)
	return beFirstState
}

// Handle the first characters of the BRE.
func beFirstState(l *lexer.Lexer) lexer.StateFunc {
	switch l.Next() {
	case '^':
		l.Emit(tNot)
		// - or ] After ^ is literal
		if l.Accept("-]") {
			l.Emit(tCharacter)
		}
		return beState
	case ']':
		// Check for empty BRE
		if l.Peek() == lexer.EOFRune {
			l.Emit(tBREEnd)
			return nil
		}
		l.Emit(tCharacter)
		return beState
	case '-':
		l.Emit(tCharacter)
		return beState
	default:
		l.Backup()
		return beState
	}
}

func beState(l *lexer.Lexer) lexer.StateFunc {
	switch r := l.Next(); {
	case r == ']':
		l.Emit(tBREEnd)
		return nil
	case r == ':':
		return classState
	case r == '-':
		return l.Error("parse error, unexpected '-'")
	case r == '\\':
		if l.Accept("ux") {
			return unicodeState
		}
		l.Emit(tCharacter)
		return beState
	case r == lexer.EOFRune:
		return l.Error("parse error, unexpected EOF")
	default:
		if l.Peek() == '-' {
			l.Emit(tRangeStart)
			l.Accept("-")
			l.Emit(tRangeDash)
			if l.Accept("-][^") {
				return l.Error("parse error, invalid range end")
			}
			l.Next()
			l.Emit(tRangeEnd)
		} else {
			l.Emit(tCharacter)
		}
		return beState
	}
}

func classState(l *lexer.Lexer) lexer.StateFunc {
	// TODO
	l.AcceptRun("abcdefghijklmnopqrstuvwxyz")
	if !l.Accept(":") {
		return l.Error("parse error, expecting ':'")
	}
	l.Emit(tClass)
	return beState
}

func unicodeState(l *lexer.Lexer) lexer.StateFunc {
	// TODO valid code point
	if n := l.AcceptRun("0123456789abcdef"); n > 0 {
		l.Emit(tCharacter)
	}
	return beState
}
