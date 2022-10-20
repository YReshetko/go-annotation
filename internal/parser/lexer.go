package parser

type lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func newLexer(input string) *lexer {
	l := &lexer{
		input: input,
	}
	// The sign of reading completion is ch=0, in ASCII it's 'NUL'.
	// So, we need initially read character to move on if input is not empty
	l.readChar()
	return l
}

func (l *lexer) nextToken() token {
	tok := token{}
	l.skipWhitespace()
	switch l.ch {
	case '@':
		tok = l.newToken(AT, "@")
	case '=':
		tok = l.newToken(ASSIGN, "=")
	case ',':
		tok = l.newToken(COMMA, ",")
	case ')':
		tok = l.newToken(RPAREN, ")")
	case '(':
		tok = l.newToken(LPAREN, "(")
	case '"':
		tok = l.newToken(STRING, l.readString())
	case 0:
		return l.newToken(EOF, "")
	default:
		if l.isLetter(l.ch) {
			lit := l.readIdentifier()
			return l.newToken(IDENT, lit)
		} else {
			tok = l.newToken(ILLEGAL, string(l.ch))
		}
	}
	l.readChar()
	return tok
}

func (l *lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *lexer) isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func (l *lexer) isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *lexer) readIdentifier() string {
	position := l.position
	for l.isLetter(l.ch) || l.isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *lexer) readString() string {
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
		case 0: // End of string lexer was interrupted unexpectedly
			isContinue = false
		default:
			out = out + string(l.ch)
		}
	}
	return out
}

func (l *lexer) newToken(tokenType tokenType, literal string) token {
	return token{
		tokenType: tokenType,
		literal:   literal,
	}
}
func (l *lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
