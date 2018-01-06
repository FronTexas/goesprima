package main 

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// describe: should enumerate all syntax node types
func TestSyntax(t *testing.T){
	assert := assert.New(t)
	expected := map[string]string{
        	"ArrayExpression": "ArrayExpression",
                "ArrayPattern": "ArrayPattern",
                "ArrowFunctionExpression": "ArrowFunctionExpression",
                "AssignmentExpression": "AssignmentExpression",
                "AssignmentPattern": "AssignmentPattern",
                "AwaitExpression": "AwaitExpression",
                "BinaryExpression": "BinaryExpression",
                "BlockStatement": "BlockStatement",
                "BreakStatement": "BreakStatement",
                "CallExpression": "CallExpression",
                "CatchClause": "CatchClause",
                "ClassBody": "ClassBody",
                "ClassDeclaration": "ClassDeclaration",
                "ClassExpression": "ClassExpression",
                "ConditionalExpression": "ConditionalExpression",
                "ContinueStatement": "ContinueStatement",
                "DebuggerStatement": "DebuggerStatement",
                "DoWhileStatement": "DoWhileStatement",
                "EmptyStatement": "EmptyStatement",
                "ExportAllDeclaration": "ExportAllDeclaration",
                "ExportDefaultDeclaration": "ExportDefaultDeclaration",
                "ExportNamedDeclaration": "ExportNamedDeclaration",
                "ExportSpecifier": "ExportSpecifier",
                "ExpressionStatement": "ExpressionStatement",
                "ForInStatement": "ForInStatement",
                "ForOfStatement": "ForOfStatement",
                "ForStatement": "ForStatement",
                "FunctionDeclaration": "FunctionDeclaration",
                "FunctionExpression": "FunctionExpression",
                "Identifier": "Identifier",
                "IfStatement": "IfStatement",
                "Import": "Import",
                "ImportDeclaration": "ImportDeclaration",
                "ImportDefaultSpecifier": "ImportDefaultSpecifier",
                "ImportNamespaceSpecifier": "ImportNamespaceSpecifier",
                "ImportSpecifier": "ImportSpecifier",
                "LabeledStatement": "LabeledStatement",
                "Literal": "Literal",
                "LogicalExpression": "LogicalExpression",
                "MemberExpression": "MemberExpression",
                "MetaProperty": "MetaProperty",
                "MethodDefinition": "MethodDefinition",
                "NewExpression": "NewExpression",
                "ObjectExpression": "ObjectExpression",
                "ObjectPattern": "ObjectPattern",
                "Program": "Program",
                "Property": "Property",
                "RestElement": "RestElement",
                "ReturnStatement": "ReturnStatement",
                "SequenceExpression": "SequenceExpression",
                "SpreadElement": "SpreadElement",
                "Super": "Super",
                "SwitchCase": "SwitchCase",
                "SwitchStatement": "SwitchStatement",
                "TaggedTemplateExpression": "TaggedTemplateExpression",
                "TemplateElement": "TemplateElement",
                "TemplateLiteral": "TemplateLiteral",
                "ThisExpression": "ThisExpression",
                "ThrowStatement": "ThrowStatement",
                "TryStatement": "TryStatement",
                "UnaryExpression": "UnaryExpression",
                "UpdateExpression": "UpdateExpression",
                "VariableDeclaration": "VariableDeclaration",
                "VariableDeclarator": "VariableDeclarator",
                "WhileStatement": "WhileStatement",
                "WithStatement": "WithStatement",
                "YieldExpression": "YieldExpression",
	}
	assert.Equal(Syntax, expected)
}

func TestParse(t *testing.T){
        // should not accept zero parameter
        assert.Panics(t, func() {goesprima.parse()})
}

