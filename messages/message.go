package messages

type messages struct{
    BadImportCallArity string
    BadGetterArity string
    BadSetterArity string
    BadSetterRestParameter string
    ConstructorIsAsync string
    ConstructorSpecialMethod string
    DeclarationMissingInitializer string
    DefaultRestParameter string
    DefaultRestProperty string
    DuplicateBinding string
    DuplicateConstructor string
    DuplicateProtoProperty string
    ForInOfLoopInitializer string
    GeneratorInLegacyContext string
    IllegalBreak string
    IllegalContinue string
    IllegalExportDeclaration string
    IllegalImportDeclaration string
    IllegalLanguageModeDirective string
    IllegalReturn string
    InvalidEscapedReservedWord string
    InvalidHexEscapeSequence string
    InvalidLHSInAssignment string
    InvalidLHSInForIn string
    InvalidLHSInForLoop string
    InvalidModuleSpecifier string
    InvalidRegExp string
    LetInLexicalBinding string
    MissingFromClause string
    MultipleDefaultsInSwitch string
    NewlineAfterThrow string
    NoAsAfterImportNamespace string
    NoCatchOrFinally string
    ParameterAfterRestParameter string
    PropertyAfterRestProperty string
    Redeclaration string
    StaticPrototype string
    StrictCatchVariable string
    StrictDelete string
    StrictFunction string
    StrictFunctionName string
    StrictLHSAssignment string
	StrictLHSPostfix string
	StrictLHSPrefix string
    StrictModeWith string
    StrictOctalLiteral string
    StrictParamDupe string
    StrictParamName string
    StrictReservedWord string
    StrictVarName string
    TemplateOctalLiteral string
    UnexpectedEOS string
    UnexpectedIdentifier string
    UnexpectedNumber string
    UnexpectedReserved string
    UnexpectedString string
    UnexpectedTemplate string
    UnexpectedToken string
    UnexpectedTokenIllegal string
    UnknownLabel string
    UnterminatedRegExp string
}

var instance *messages

func GetInstance() *messages{
	if instance == nil {
		instance = &messages{
			BadImportCallArity: "Unexpected token",
			BadGetterArity: "Getter must not have any formal parameters",
			BadSetterArity: "Setter must have exactly one formal parameter",
			BadSetterRestParameter: "Setter function argument must not be a rest parameter",
			ConstructorIsAsync: "Class constructor may not be an async method",
			ConstructorSpecialMethod: "Class constructor may not be an accessor",
			DeclarationMissingInitializer: "Missing initializer in %0 declaration",
			DefaultRestParameter: "Unexpected token =",
			DefaultRestProperty: "Unexpected token =",
			DuplicateBinding: "Duplicate binding %0",
			DuplicateConstructor: "A class may only have one constructor",
			DuplicateProtoProperty: "Duplicate __proto__ fields are not allowed in object literals",
			ForInOfLoopInitializer: "%0 loop variable declaration may not have an initializer",
			GeneratorInLegacyContext: "Generator declarations are not allowed in legacy contexts",
			IllegalBreak: "Illegal break statement",
			IllegalContinue: "Illegal continue statement",
			IllegalExportDeclaration: "Unexpected token",
			IllegalImportDeclaration: "Unexpected token",
			IllegalLanguageModeDirective: "Illegal \"use strict\" directive in function with non-simple parameter list",
			IllegalReturn: "Illegal return statement",
			InvalidEscapedReservedWord: "Keyword must not contain escaped characters",
			InvalidHexEscapeSequence: "Invalid hexadecimal escape sequence",
			InvalidLHSInAssignment: "Invalid left-hand side in assignment",
			InvalidLHSInForIn: "Invalid left-hand side in for-in",
			InvalidLHSInForLoop: "Invalid left-hand side in for-loop",
			InvalidModuleSpecifier: "Unexpected token",
			InvalidRegExp: "Invalid regular expression",
			LetInLexicalBinding: "let is disallowed as a lexically bound name",
			MissingFromClause: "Unexpected token",
			MultipleDefaultsInSwitch: "More than one default clause in switch statement",
			NewlineAfterThrow: "Illegal newline after throw",
			NoAsAfterImportNamespace: "Unexpected token",
			NoCatchOrFinally: "Missing catch or finally after try",
			ParameterAfterRestParameter: "Rest parameter must be last formal parameter",
			PropertyAfterRestProperty: "Unexpected token",
			Redeclaration: "%0 \"%1\" has already been declared",
			StaticPrototype: "Classes may not have static property named prototype",
			StrictCatchVariable: "Catch variable may not be eval or arguments in strict mode",
			StrictDelete: "Delete of an unqualified identifier in strict mode.",
			StrictFunction: "In strict mode code, functions can only be declared at top level or inside a block",
			StrictFunctionName: "Function name may not be eval or arguments in strict mode",
			StrictLHSAssignment: "Assignment to eval or arguments is not allowed in strict mode",
			StrictLHSPostfix: "Postfix increment/decrement may not have eval or arguments operand in strict mode",
			StrictLHSPrefix: "Prefix increment/decrement may not have eval or arguments operand in strict mode",
			StrictModeWith: "Strict mode code may not include a with statement",
			StrictOctalLiteral: "Octal literals are not allowed in strict mode.",
			StrictParamDupe: "Strict mode function may not have duplicate parameter names",
			StrictParamName: "Parameter name eval or arguments is not allowed in strict mode",
			StrictReservedWord: "Use of future reserved word in strict mode",
			StrictVarName: "Variable name may not be eval or arguments in strict mode",
			TemplateOctalLiteral: "Octal literals are not allowed in template strings.",
			UnexpectedEOS: "Unexpected end of input",
			UnexpectedIdentifier: "Unexpected identifier",
			UnexpectedNumber: "Unexpected number",
			UnexpectedReserved: "Unexpected reserved word",
			UnexpectedString: "Unexpected string",
			UnexpectedTemplate: "Unexpected quasi %0",
			UnexpectedToken: "Unexpected token %0",
			UnexpectedTokenIllegal: "Unexpected token ILLEGAL",
			UnknownLabel: "Undefined label \"%0\"",
			UnterminatedRegExp: "Invalid regular expression: missing /",
		}
	}
	return instance
}


