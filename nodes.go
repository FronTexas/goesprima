package main

import (
	"github.com/robertkrimen/otto/parser"
	"net"
)

type (
	ArgumentListElement interface{}
	Expression interface{
		ArgumentListElement
		ArrayExpressionElement
		ExportableDefaultDeclaration
	}
	SpreadElement interface{
		ArgumentListElement
		ArrayExpressionElement
		ObjectExpressionProperty
	}

	ArrayExpressionElement interface{}

	ArrayPatternElement interface{}
	AssignmentPattern interface{
		ArrayPatternElement
		FunctionParameter
		PropertyValue
	}

	BindingIdentifier interface{
		ArrayPatternElement
		ExportableDefaultDeclaration
		FunctionParameter
		PropertyValue
	}

	BindingPattern interface{
		ArrayPatternElement
		ExportableDefaultDeclaration
		FunctionParameter
		PropertyValue
	}

	RestElement interface{
		ArrayPatternElement
		ObjectPatternProperty
	}

	ArrayPattern interface{
		BindingPattern
	}

	ObjectPattern interface{
		BindingPattern
	}

	Identifier interface {
		BindingIdentifier
		Expression
		PropertyKey
	}

	Declaration interface{
		StatementListItem
	}

	AsyncFunctionDeclaration interface{
		Declaration
		ExportableNamedDeclaration
		Statement
	}

	ClassDeclaration interface {
		Declaration
		ExportableDefaultDeclaration
		ExportableNamedDeclaration
	}

	ExportDeclaration interface {
		Declaration
	}

	FunctionDeclaration interface {
		Declaration
		ExportableDefaultDeclaration
		ExportableNamedDeclaration
		Statement
	}

	ImportDeclaration interface {
		Declaration
	}

	VariableDeclaration interface {
		Declaration
		ExportableNamedDeclaration
		Statement
	}

	ExportableDefaultDeclaration interface {
		ExportDeclaration
	}

	ExportableNamedDeclaration interface {
		ExportDeclaration
	}

	ExportAllDeclaration interface {
		ExportDeclaration
	}

	ArrayExpresion interface {
		Expression
	}

	ArrowFunctionExpression interface {
		Expression
	}

	AssignmentExpression interface {
		Expression
	}

	AsyncArrowFunctionExpression interface {
		Expression
	}

	AsyncFunctionExpression interface {
		Expression
		PropertyValue
	}

	AwaitExpression interface {
		Expression
	}

	BinaryExpression interface {
		Expression
	}

	CallExpression interface {
		Expression
	}

	ClassExpression interface {
		Expression
	}

	ComputedMemberExpression interface {
		Expression
	}

	ConditionalExpression interface {
		Expression
	}

	FunctionExpression interface {
		Expression
		PropertyValue
	}

	Literal interface {
		Expression
		PropertyKey
	}

	NewExpression interface {
		Expression
	}

	ObjectExpression interface {
		Expression
	}

	RegexLiteral interface {
		Expression
	}

	SequenceExpression interface {
		Expression
	}

	StaticMemberExpression interface {
		Expression
	}

	TaggedTemplateExpression interface {
		Expression
	}

	ThisExpression interface {
		Expression
	}


	UnaryExpression interface {
		Expression
	}

	UpdateExpression interface {
		Expression
	}


	YieldExpression interface {
		Expression
	}

	FunctionParameter interface {}

	ImportDeclarationSpecifier interface {}

	ImportDefaultSpecifier interface {
		ImportDeclarationSpecifier
	}

	ImportNamespaceSpecifier interface {
		ImportDeclarationSpecifier
	}

	ImportSpecifier interface {
		ImportDeclarationSpecifier
	}

	ObjectExpressionProperty interface {}

	Property interface {
		ObjectExpressionProperty
		ObjectPatternProperty
	}

	ObjectPatternProperty interface {}

	Statement interface {
		StatementListItem
	}

	BreakStatement interface {
		Statement
	}

	ContinueStatement interface {
		Statement
	}

	DebuggerStatement interface {
		Statement
	}

	DoWhileStatement interface {
		Statement
	}

	EmptyStatement interface {
		Statement
	}

	ExpressionStatement interface {
		Statement
	}

	Directive interface {
		Statement
	}

	ForStatement interface {
		Statement
	}

	ForInStatement interface {
		Statement
	}

	ForOfStatement interface {
		Statement
	}

	IfStatement interface {
		Statement
	}

	ReturnStatement interface {
		Statement
	}

	SwitchStaement interface {
		Statement
	}

	ThrowStatement interface {
		Statement
	}

	TryStatement interface {
		Statement
	}

	WhileStatement interface {
		Statement
	}

	WithStatement interface {
		Statement
	}

	PropertyKey interface {}

	PropertyValue interface {}

	//export type StatementListItem = Declaration | Statement;

	StatementListItem interface {}

)




