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

var _tokenName []string
func GetTokenName(tkn Token) string {
	_tokenName = []string{
		BooleanLiteral: "Boolean",
		EOF: "<end>",
		Identifier: "Identifier",
		Keyword: "Keyword",
		NullLiteral: "Null",
		NumericLiteral: "Numeric",
		Punctuator: "Punctuator",
		StringLiteral: "String",
		RegularExpression : "RegularExpression",
		Template: "Template",
	}
	return _tokenName[tkn]

}


