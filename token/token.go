package token

type Token int

const (
	_ Token = iota
	BooleanLiteral
	EOF
	Identifier
	Keyword
	NullLiteral
	NumericLiteral
	Punctuator
	StringLiteral
	RegularExpression
	Template
)

