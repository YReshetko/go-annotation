package annotation

import (
	"fmt"
	"strings"

	"github.com/YReshetko/go-annotation/internal/annotation/lexer"
	"github.com/YReshetko/go-annotation/internal/annotation/tokens"
)

func Parse(doc string) ([]Annotation, bool) {
	if !strings.Contains(doc, "@") {
		return nil, false
	}

	p := newParser(lexer.NewLexer(doc))

	s := p.parse()

	if len(p.errors) != 0 {
		errs := ""
		for _, e := range p.errors {
			errs += e + "\n"
		}
		panic(errs)

	}
	return s, true
}

func newParser(l *lexer.Lexer) *parser {
	p := &parser{
		l:      l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()
	return p
}

type parser struct {
	l *lexer.Lexer

	currToken tokens.Token
	peekToken tokens.Token

	errors []string
}

func (p *parser) parse() []Annotation {
	out := []Annotation{}
	for !p.currTokenIs(tokens.EOF) {
		if p.currTokenIs(tokens.AT) {
			out = append(out, p.parseAnnotation())
		}
		p.nextToken()
	}
	return out
}

func (p *parser) parseAnnotation() Annotation {
	out := Annotation{}
	if p.expectPeekToken(tokens.IDENT) {
		out.name = p.currToken.Literal
	}
	if p.peekTokenIs(tokens.LPAREN) {
		p.nextToken()
		out.params = p.parseAnnotationParams()
	}
	return out
}

func (p *parser) parseAnnotationParams() map[string]string {
	out := map[string]string{}
	for !p.peekTokenIs(tokens.RPAREN) {
		k, v := p.parseParamPair()
		out[k] = v
	}
	return out
}

func (p *parser) parseParamPair() (string, string) {
	k, v := "", ""
	if p.expectPeekToken(tokens.IDENT) {
		k = p.currToken.Literal
	}
	p.expectPeekToken(tokens.ASSIGN)
	if p.expectPeekToken(tokens.STRING) {
		v = p.currToken.Literal
	}

	if p.peekTokenIs(tokens.COMMA) {
		p.nextToken()
	}
	return k, v
}

func (p *parser) currTokenIs(t tokens.TokenType) bool {
	return p.currToken.Type == t
}

func (p *parser) peekTokenIs(t tokens.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *parser) expectPeekToken(t tokens.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *parser) peekError(t tokens.TokenType) {
	msg := fmt.Sprintf("expected token %s instead got %s", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}
