package parser

import (
	"errors"
	"fmt"
	"strings"
)

type parser struct {
	l *lexer

	currToken token
	peekToken token

	errors []string
}

func parse(doc string) ([]annotation, error) {
	if !strings.Contains(doc, "@") {
		return nil, nil
	}

	p := &parser{
		l:      newLexer(doc),
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	s := p.parse()

	if len(p.errors) != 0 {
		errs := ""
		for _, e := range p.errors {
			errs += e + "\n"
		}

		return nil, errors.New(errs)
	}
	return s, nil
}

func (p *parser) parse() []annotation {
	var out []annotation
	for !p.currTokenIs(EOF) {
		if p.currTokenIs(AT) {
			out = append(out, p.parseAnnotation())
		}
		p.nextToken()
	}
	return out
}

func (p *parser) parseAnnotation() annotation {
	out := annotation{}
	if p.expectPeekToken(IDENT) {
		out.name = p.currToken.literal
	}
	if p.peekTokenIs(LPAREN) {
		p.nextToken()
		out.params = p.parseAnnotationParams()
	}
	return out
}

func (p *parser) parseAnnotationParams() map[string]string {
	out := map[string]string{}
	for !p.peekTokenIs(RPAREN) {
		k, v := p.parseParamPair()
		out[k] = v
		if p.peekTokenIs(EOF) {
			return nil
		}
	}
	return out
}

func (p *parser) parseParamPair() (string, string) {
	k, v := "", ""
	if p.expectPeekToken(IDENT) {
		k = p.currToken.literal
	}
	p.expectPeekToken(ASSIGN)
	if p.expectPeekToken(STRING) {
		v = p.currToken.literal
	}

	if p.peekTokenIs(COMMA) {
		p.nextToken()
	}
	return k, v
}

func (p *parser) currTokenIs(t tokenType) bool {
	return p.currToken.tokenType == t
}

func (p *parser) peekTokenIs(t tokenType) bool {
	return p.peekToken.tokenType == t
}

func (p *parser) expectPeekToken(t tokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *parser) peekError(t tokenType) {
	msg := fmt.Sprintf("expected token %s instead got %s", t, p.peekToken.tokenType)
	p.errors = append(p.errors, msg)
}

func (p *parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.nextToken()
}
