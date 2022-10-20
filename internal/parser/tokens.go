package parser

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	STRING = "STRING"

	AT     = "@"
	ASSIGN = "="
	COMMA  = ","

	LPAREN = "("
	RPAREN = ")"
)

type tokenType string

type token struct {
	tokenType tokenType
	literal   string
}
