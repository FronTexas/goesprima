package main

import "github.com/robertkrimen/otto/parser"

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
	}

	ArrayExpressionElement interface{}

	ArrayPatternElement interface{}
	AssignmentPattern interface{
		ArrayPatternElement
	}

	BindingIdentifier interface{
		ArrayPatternElement
		ExportableDefaultDeclaration
	}

	BindingPattern interface{
		ArrayPatternElement
		ExportableDefaultDeclaration
	}

	RestElement interface{
		ArrayPatternElement
	}

	ArrayPattern interface{
		BindingPattern
	}

	ObjectPattern interface{
		BindingPattern
	}

	Identifier interface {
		BindingIdentifier
	}

	Declaration interface{}

	AsyncFunctionDeclaration interface{
		Declaration
		ExportableNamedDeclaration
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
	}

	ImportDeclaration interface {
		Declaration
	}

	VariableDeclaration interface {
		Declaration
		ExportableNamedDeclaration
	}

	ExportableDefaultDeclaration interface {}

	//export type ExportableNamedDeclaration = AsyncFunctionDeclaration | ClassDeclaration | FunctionDeclaration | VariableDeclaration;

	ExportableNamedDeclaration interface {}


)




