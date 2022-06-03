package tokens

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

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
