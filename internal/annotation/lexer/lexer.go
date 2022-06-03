package lexer

import "github.com/YReshetko/go-annotation/internal/annotation/tokens"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	// The sign of reading completion is ch=0, in ASCII it's 'NUL'.
	// So, we need initially read character to move on if input is not empty
	l.readChar()
	return l
}

func (l *Lexer) NextToken() tokens.Token {
	tok := tokens.Token{}
	l.skipWhitespace()
	switch l.ch {
	case '@':
		tok = l.newToken(tokens.AT, "@")
	case '=':
		tok = l.newToken(tokens.ASSIGN, "=")
	case ',':
		tok = l.newToken(tokens.COMMA, ",")
	case ')':
		tok = l.newToken(tokens.RPAREN, ")")
	case '(':
		tok = l.newToken(tokens.LPAREN, "(")
	case '"':
		tok = l.newToken(tokens.STRING, l.readString())
	case 0:
		return l.newToken(tokens.EOF, "")
	default:
		if l.isLetter(l.ch) {
			lit := l.readIdentifier()
			return l.newToken(tokens.IDENT, lit)
		} else {
			tok = l.newToken(tokens.ILLEGAL, string(l.ch))
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func (l *Lexer) isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for l.isLetter(l.ch) || l.isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	out := ""
	isContinue := true
	for isContinue {
		l.readChar()
		switch l.ch {
		case '\\':
			switch l.peekChar() {
			case '"':
				out = out + string('"')
			case 'n':
				out = out + string('\n')
			case 't':
				out = out + string('\t')
			}
			l.readChar()
		case '"':
			isContinue = false
		case 0:
			// TODO error lexer
			isContinue = false
		default:
			out = out + string(l.ch)
		}
	}
	return out
}

func (l *Lexer) newToken(tokenType tokens.TokenType, literal string) tokens.Token {
	return tokens.Token{
		Type:    tokenType,
		Literal: literal,
	}
}
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
