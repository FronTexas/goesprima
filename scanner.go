package main

import(
	"strings"
	"regexp"
	"strconv"
	"goesprima/messages"
	"goesprima/character"
	"goesprima/token"

)

func getCharCodeAt(s string, i int) rune{
	return []rune(s)[i]
}

func hexValue(ch string) int {
	return strings.Index("0123456789abcdef", strings.ToLower(ch))
}

func octalValue (ch string) int {
	return strings.Index("01234567", ch)
}

type Position struct {
	line int
	column int
}

type SourceLocation struct {
	start *Position
	end *Position
	source string
}

// TODO: Find a better way to differentiate scanner.Comment and comment-handler.Comment
type Comment_scanner struct {
	multiline bool
	slice []int
	_range []int
	loc *SourceLocation
}

type RawToken struct {
	_type token.Token
	value_string string
	value_number float32
	pattern string
	flags string
	regex regexp.Regexp
	octal bool
	cooked string
	head bool
	tail bool
	lineNumber int
	lineStart int
	start int
	end int
}

type ScannerState struct {
	index int
	lineNumber int
	lineStart int
}

type Scanner struct {
	source string
	errorHandler ErrorHandler
	trackComment bool
	isModule bool
	index int
	lineNumber int
	lineStart int
	curlyStack []string
	length int
}

func NewScanner(code string, handler ErrorHandler) *Scanner{
	var lineNumber int
	if lineNumber = 0; len(code) > 0 {
		lineNumber = 1
	}
	return &Scanner{
		source: code,
		errorHandler: handler,
		trackComment: false,
		isModule: false,
		length: len(code),
		index: 0,
		lineNumber: lineNumber,
		lineStart : 0,
		curlyStack: []string{},
	}
}

func (self *Scanner) saveState() *ScannerState{
	return &ScannerState{
		index: self.index,
		lineNumber: self.lineNumber,
		lineStart: self.lineStart,
	}
}

func (self *Scanner) restoreState(state *ScannerState) {
	self.index = state.index
	self.lineNumber = state.lineNumber
	self.lineStart = state.lineStart
}

func (self *Scanner) eof() bool {
	return self.index >= self.length
}

// TODO: Uncomment these two functions after Error implementation is finished

 func (self *Scanner) throwUnexpectedToken(message string) *Error{
 	if message == "" {
 		message = messages.GetInstance().UnexpectedTokenIllegal
	}

 	return self.errorHandler.throwError(self.index, self.lineNumber,
 		self.index - self.lineStart + 1, message)
 }

 func (self *Scanner) tolerateUnexpectedToken(message string){
 	if message == ""{
 		message = messages.GetInstance().UnexpectedTokenIllegal
	}
 	self.errorHandler.tolerateError(self.index, self.lineNumber,
 		self.index - self.lineStart +1, message)
 }

func (self *Scanner) skipSingleLineComment(offset int) []*Comment_scanner{
	comments := []*Comment_scanner{}
	var start int
	var loc *SourceLocation
	if self.trackComment {
		start = self.index - offset
		loc = &SourceLocation{
			start: &Position{
				line: self.lineNumber,
				column: self.index - self.lineStart - offset,
			},
		}
	}

	for !self.eof() {
		ch := []rune(self.source)[self.index]
		self.index++
		// TODO implement IsLineTerminator in character.go
		if character.IsLineTerminator(ch) {
			if self.trackComment {
				loc.end = &Position{
					line: self.lineNumber,
					column: self.index - self.lineStart - 1,
				}
				entry := &Comment_scanner{
					multiline: false,
					slice: []int{start + offset, self.index - 1},
					_range: []int{start, self.index - 1},
					loc: loc,
				}
				comments = append(comments, entry)
			}

			if ch == 13 && []rune(self.source)[self.index] == 10 {
				self.index++
			}
		}
		self.lineNumber++
		self.lineStart = self.index
		return comments
	}

	if self.trackComment {
		loc.end = &Position{
			line: self.lineNumber,
			column: self.index - self.lineStart,
		}

		entry := &Comment_scanner{
			multiline: false,
			slice: []int{start + offset, self.index},
			_range: []int{start, self.index},
			loc: loc,
		}
		comments = append(comments, entry)
	}
	return comments
}

func (self *Scanner) skipMultiLineComment() []*Comment_scanner{
	comments := []*Comment_scanner{}
	var start int
	var loc *SourceLocation

	if self.trackComment {
		start = self.index - 2
		loc = &SourceLocation{
			start: &Position{
				line: self.lineNumber,
				column: self.index - self.lineStart - 2,
			},
		}
	}

	for !self.eof() {
		ch := []rune(self.source)[self.index]
		self.index++
		// TODO implement IsLineTerminator in character.go
		if character.IsLineTerminator(ch) {
			if ch == 0x0D && []rune(self.source)[self.index] == 0x0A{
				self.index++
			}
			self.lineNumber++
			self.index++
			self.lineStart = self.index
		}else if ch == 0x2A{
			if []rune(self.source)[self.index + 1] == 0x2F {
				self.index += 2
				if self.trackComment {
					loc.end = &Position{
						line: self.lineNumber,
						column: self.index - self.lineStart,
					}
					entry := &Comment_scanner{
						multiline: true,
						slice: []int{start + 2, self.index - 2},
						_range: []int{start, self.index},
						loc: loc,
					}
					comments = append(comments, entry)
				}
				return comments
			}
			self.index++
		}else{
			self.index++
		}
	}
    // Ran off the end of the file - the whole thing is a comment
	if self.trackComment {
		loc.end = &Position{
			line: self.lineNumber,
			column: self.index - self.lineStart,
		}

		entry := &Comment_scanner{
			multiline: true,
			slice: []int{start + 2, self.index},
			_range: []int{start, self.index},
			loc: loc,
		}
		comments = append(comments, entry)
	}
	// TODO uncomment self once the method is implemented
	// self.tolerateUnexpectedToken()
	return comments
}

func (self *Scanner) scanComments() []*Comment_scanner{
	var comments []*Comment_scanner
	if self.trackComment{
		comments = []*Comment_scanner{}
	}

	start := self.index == 0

	for !self.eof() {
		ch := getCharCodeAt(self.source, self.index)

		if character.IsWhiteSpace(ch) {
			self.index++
		}else if character.IsLineTerminator(ch){
			self.index++
			if ch == 0x0D && getCharCodeAt(self.source, self.index) == 0X0A {
				self.index++
			}
			self.lineNumber++
			self.lineStart = self.index
			start = true
		}else if ch == 0x2F {
			ch = getCharCodeAt(self.source, self.index + 1)
			if ch == 0x2F {
				self.index += 2
				comment := self.skipSingleLineComment(2)
				if self.trackComment{
					comments = append(comments, comment...)
				}
				start = true
			}else if ch == 0x2A {
				self.index += 2
				comment := self.skipSingleLineComment(3)
				if self.trackComment {
					comments = append(comments, comment...)
				}
			}else {
				break
			}
		}else if start && ch == 0x2D {
			if getCharCodeAt(self.source, self.index + 1) == 0x2D && getCharCodeAt(self.source, self.index + 2) == 0x3E {
				self.index += 3
				comment := self.skipSingleLineComment(3)
				if self.trackComment {
					comments = append(comments, comment...)
				}
			}else{
				break
			}
		}else if (ch == 0x3C && !self.isModule){
			if string([]rune(self.source)[self.index + 1 : self.index + 4]) == "!--"{
				self.index += 4
				comment := self.skipSingleLineComment(4)
				if self.trackComment {
					comments = append(comments, comment...)
				}
			}else {
				break
			}
		} else {
			break
		}
	}
	return comments
}

func (self *Scanner) isFutureReservedWord(id string) bool {
	switch id {
	case "enum":
	case "export":
	case "import":
	case "super":
		return true
	}
	return false
}

func (self *Scanner) isRestrictedWord(id string) bool {
	return id == "eval" || id == "arguments"
}

func (self *Scanner) isKeyword(id string) bool {
	switch len(id){
	case 2:
		return id == "if" || id == "in" || id == "do"
	case 3:
		return (id == "var") || (id == "for") || (id == "new") ||
		(id == "try") || (id == "let")
	case 4:
		return (id == "self") || (id == "else") || (id == "case") ||
		(id == "void") || (id == "with") || (id == "enum")
	case 5:
		return (id == "while") || (id == "break") || (id == "catch") ||
		(id == "throw") || (id == "const") || (id == "yield") ||
		(id == "class") || (id == "super")
	case 6:
		return (id == "return") || (id == "typeof") || (id == "delete") ||
		(id == "switch") || (id == "export") || (id == "import")
	case 7:
		return (id == "default") || (id == "finally") || (id == "extends")
	case 8:
		return (id == "function") || (id == "continue") || (id == "debugger")
	case 10:
		return (id == "instanceof")
	}
	return false
}

func (self *Scanner) codePointAt(i int) rune {
	cp := getCharCodeAt(self.source, i)

	if (cp >= 0xD800 && cp <= 0xDBFF) {
		second  := getCharCodeAt(self.source, i + 1)
		if (second >= 0xDC00 && second <= 0xDFFF) {
			first := cp
			cp = (first-0xD800)*0x400 + second - 0xDC00 + 0x10000
		}
	}
	return cp
}

func (self *Scanner) scanHexEscape(prefix string) string {
	var len int
	if len = 4; prefix == "u"{
		len = 2
	}
	code := 0

	for i := 0; i < len; i++ {
		if !self.eof() && character.IsHexDigit(getCharCodeAt(self.source, self.index)) {
			code = code * 16 + hexValue(string(self.source[self.index + 1]))
			self.index += 1
		} else {
			return ""
		}
	}
	//return String.fromCharCode(code)
	toReturn := string(code)
	return toReturn
}

func (self *Scanner) scanUnicodeCodePointEscape() string {
	ch := self.source[self.index]
	code := 0

	// At least, one hex digit is required.
	if (ch == '}') {
		self.throwUnexpectedToken("")
	}

	for !self.eof() {
		ch = self.source[self.index]
		self.index += 1
		if (!character.IsHexDigit(getCharCodeAt(string(ch), 0))) {
			break
		}
		code = code * 16 + hexValue(string(ch))
	}

	if (code > 0x10FFFF || ch != '}') {
		self.throwUnexpectedToken("")
	}

	return character.FromCodePoint(rune(code))
}

func (self *Scanner) getIdentifier() string {
    start := self.index
    self.index += 1
	for !self.eof() {
		ch := getCharCodeAt(self.source, self.index)
		if (ch == 0x5C) {
			// Blackslash (U+005C) marks Unicode escape sequence.
			self.index = start
			return self.getComplexIdentifier()
		} else if (ch >= 0xD800 && ch < 0xDFFF) {
			// Need to handle surrogate pairs.
			self.index = start
			return self.getComplexIdentifier()
		}
		if (character.IsIdentifierPart(ch)) {
			self.index++
		} else {
			break
		}
	}

	return self.source[start : self.index]
}

func (self *Scanner) getComplexIdentifier() string {
	cp := self.codePointAt(self.index)
	id := character.FromCodePoint(cp)
	self.index += len(id)

	// '\u' (U+005C, U+0075) denotes an escaped character.
	var ch string
	if (cp == 0x5C) {
		if (getCharCodeAt(self.source, self.index) != 0x75) {
			self.throwUnexpectedToken("")
		}
		self.index++
		if (self.source[self.index] == '{') {
			self.index++
			ch = self.scanUnicodeCodePointEscape()
		} else {
			ch = self.scanHexEscape("u")
			if ch == "" || ch == "\\" || !character.IsIdentifierStart(getCharCodeAt(ch, 0)) {
				self.throwUnexpectedToken("")
			}
		}
		id = ch
	}

	for !self.eof() {
		cp = self.codePointAt(self.index)
		if (!character.IsIdentifierPart(cp)) {
			break
		}
		ch = character.FromCodePoint(cp)
		id += ch
		self.index += len(ch)

		// '\u' (U+005C, U+0075) denotes an escaped character.
		if (cp == 0x5C) {
			id = id[0:len(id) - 1]

			if (getCharCodeAt(self.source, self.index) != 0x75) {
				self.throwUnexpectedToken("")
			}
			self.index++
			if (self.source[self.index] == '{') {
				self.index++
				ch = self.scanUnicodeCodePointEscape()
			} else {
				ch = self.scanHexEscape("u")
				if ch == "" || ch == "\\" || !character.IsIdentifierPart(getCharCodeAt(ch, 0)) {
					self.throwUnexpectedToken("")
				}
			}
			id += ch
		}
	}
	return id
}

type codeOctalStruct struct {
	code int
	octal bool
}

func (self *Scanner) octalToDecimal(ch string) codeOctalStruct {
	// \0 is not octal escape sequence
	octal := (ch != "0");
	code := octalValue(ch);

	if !self.eof() && character.IsOctalDigit(getCharCodeAt(self.source, self.index)) {
		octal = true;
		code = code * 8 + octalValue(string(self.source[self.index]))
		self.index++

		// 3 digits are only allowed when string starts
		// with 0, 1, 2, 3
		if strings.Index("0123", string(ch)) >= 0 && !self.eof() && character.IsOctalDigit(getCharCodeAt(self.source, self.index)) {
			code = code * 8 + octalValue(string(self.source[self.index]));
			self.index++
		}
	}

	return codeOctalStruct{
		code,
		octal,
	};
}

// https://tc39.github.io/ecma262/#sec-names-and-keywords

func (self *Scanner) scanIdentifier() *RawToken {
	var _type token.Token
	start := self.index;

	// Backslash (U+005C) starts an escaped character.
	var id string
	if id = self.getIdentifier(); getCharCodeAt(self.source, start) == 0x5C{
	   id = self.getIdentifier()
	}

	// There is no keyword or literal with only one character.
	// Thus, it must be an identifier.
	if (len(id) == 1) {
		_type = token.Identifier;
	} else if (self.isKeyword(id)) {
		_type = token.Keyword;
	} else if (id == "null") {
		_type = token.NullLiteral;
	} else if (id == "true" || id == "false") {
		_type = token.BooleanLiteral;
	} else {
		_type = token.Identifier;
	}

	if (_type != token.Identifier && (start + len(id) != self.index)) {
		restore := self.index;
		self.index = start;
		self.tolerateUnexpectedToken(messages.GetInstance().InvalidEscapedReservedWord);
		self.index = restore;
	}

	return &RawToken{
		_type: _type,
		value_string: id,
		lineNumber: self.lineNumber,
		lineStart: self.lineStart,
		start: start,
		end: self.index,
	};
}

// https://tc39.github.io/ecma262/#sec-punctuators
func(self *Scanner) scanPunctuator() *RawToken {
	start := self.index;

	// Check for most common single-character punctuators.
	str := string(self.source[self.index]);
	switch (str) {

		case "(":
		case "{":
			if (str == "{") {
				self.curlyStack = append(self.curlyStack, "{")
			}
			self.index++;
			break;
		case ".":
			self.index++;
			if (self.source[self.index] == '.' && self.source[self.index + 1] == '.') {
				// Spread operator: ...
				self.index += 2;
				str = "...";
			}
			break;
		case "}":
			self.index++;
			self.curlyStack = self.curlyStack[:len(self.curlyStack) - 1]
			break;
		case ")":
		case ";":
		case ",":
		case "[":
		case "]":
		case ":":
		case "?":
		case "~":
			self.index++;
			break;
		default:
			// 4-character punctuator.
			str = self.source[self.index : 4]
			if (str == ">>>=") {
				self.index += 4
			} else {
				// 3-character punctuators.
				str = str[0:3]
				if str == "===" || str == "!==" || str == ">>>" ||
					str == "<<=" || str == ">>=" || str == "**=" {
					self.index += 3
				} else {
					// 2-character punctuators.
					str = str[0:2]
					if str == "&&" || str == "||" || str == "==" || str == "!=" ||
						str == "+=" || str == "-=" || str == "*=" || str == "/=" ||
						str == "++" || str == "--" || str == "<<"|| str == ">>" ||
						str == "&="|| str == "|="|| str == "^=" || str == "%=" ||
						str == "<=" || str == ">="|| str == "=>" || str == "**" {
						self.index += 2
					} else {
						// 1-character punctuators.
						str = string(self.source[self.index])
						if strings.Index("<>=!+-*%&|^/", str) >= 0 {
							self.index++
						}
					}
				}
			}
	}

	if (self.index == start) {
		self.throwUnexpectedToken("");
	}

	return &RawToken{
		_type: token.Punctuator,
		value_string: str,
		lineNumber: self.lineNumber,
		lineStart: self.lineStart,
		start: start,
		end: self.index,
	};
}

// https://tc39.github.io/ecma262/#sec-literals-numeric-literals
func (self *Scanner) scanHexLiteral(start int) *RawToken {
	num := "";

	for !self.eof() {
		if !character.IsHexDigit(getCharCodeAt(self.source, self.index)) {
			break;
		}
		num += string(self.source[self.index])
		self.index++
	}

	if len(num) == 0 {
		self.throwUnexpectedToken("");
	}

	if (character.IsIdentifierStart(getCharCodeAt(self.source, self.index))) {
		self.throwUnexpectedToken("");
	}
	value_number, _ := strconv.ParseInt("0x" + num, 16, 0)
	return &RawToken{
		_type: token.NumericLiteral,
		value_number: float32(value_number),
		lineNumber: self.lineNumber,
		lineStart: self.lineStart,
		start: start,
		end: self.index,
	};
}

func (self *Scanner) scanBinaryLiteral(start int) *RawToken {
	num := "";
	var ch rune

	for !self.eof() {
		ch = rune(self.source[self.index]);
		if ch != '0' && ch != '1' {
			break;
		}
		num += string(self.source[self.index]);
		self.index++
	}

	if len(num) == 0 {
		// only 0b or 0B
		self.throwUnexpectedToken("");
	}

	if (!self.eof()) {
		ch = getCharCodeAt(self.source, self.index);
		/* istanbul ignore else */
		if (character.IsIdentifierStart(ch) || character.IsDecimalDigit(ch)) {
			self.throwUnexpectedToken("");
		}
	}

	value_number, _ := strconv.ParseInt(num, 2, 0)
	return &RawToken{
		_type: token.NumericLiteral,
		value_number: float32(value_number),
		lineNumber: self.lineNumber,
		lineStart: self.lineStart,
		start: start,
		end: self.index,
	};
}

func (self *Scanner) scanOctalLiteral(prefix string, start int) *RawToken {
	num := "";
	octal := false;

	if (character.IsOctalDigit(getCharCodeAt(prefix, 0))) {
		octal = true;
		num = "0" + string(self.source[self.index]);
		self.index++
	} else {
		self.index++
	}

	for !self.eof() {
		if !character.IsOctalDigit(getCharCodeAt(self.source, self.index)) {
			break;
		}
		num += string(self.source[self.index]);
		self.index++
	}

	if (!octal && len(num) == 0) {
		// only 0o or 0O
		self.throwUnexpectedToken("");
	}

	if character.IsIdentifierStart(getCharCodeAt(self.source, self.index)) || character.IsDecimalDigit(getCharCodeAt(self.source, self.index)) {
		self.throwUnexpectedToken("");
	}
	value_number, _ := strconv.ParseInt(num, 2, 0)
	return &RawToken{
		_type: token.NumericLiteral,
		value_number: float32(value_number),
		octal: octal,
		lineNumber: self.lineNumber,
		lineStart: self.lineStart,
		start: start,
		end: self.index,
	};
}

func (self *Scanner) isImplicitOctalLiteral() bool {
	// Implicit octal, unless there is a non-octal digit.
	// (Annex B.1.1 on Numeric Literals)
	for i := self.index + 1; i < self.length; i++ {
		ch := self.source[i];
		if ch == '8' || ch == '9' {
			return false;
		}
		if !character.IsOctalDigit(getCharCodeAt(string(ch), 0)) {
			return true;
		}
	}

	return true;
}

func (self *Scanner) scanNumericLiteral() *RawToken {
	start := self.index;
	ch := self.source[start];

	// TODO figure out how to use assert instead of if statement
	if !character.IsDecimalDigit(getCharCodeAt(string(ch), 0)) && !(ch == '.'){
		panic("Numeric literal must start with a decimal digit or a decimal point")
	}

	var num string;
	if (ch != '.') {
		num = string(self.source[self.index]);
		self.index++
		ch = self.source[self.index];

		// Hex number starts with '0x'.
		// Octal number starts with '0'.
		// Octal number in ES6 starts with '0o'.
		// Binary number in ES6 starts with '0b'.
		if num == "0" {
			if ch == 'x' || ch == 'X' {
				self.index++;
				return self.scanHexLiteral(start);
			}
			if ch == 'b' || ch == 'B' {
				self.index++
				return self.scanBinaryLiteral(start);
			}
			if ch == 'o' || ch == 'O' {
				return self.scanOctalLiteral(string(ch), start);
			}

			if &ch != nil && character.IsOctalDigit(getCharCodeAt(string(ch), 0)) {
				if (self.isImplicitOctalLiteral()) {
					return self.scanOctalLiteral(string(ch), start);
				}
			}
		}

		for (character.IsDecimalDigit(getCharCodeAt(self.source, self.index))) {
			num += string(self.source[self.index])
			self.index++
		}
		ch = self.source[self.index];
	}

	if (ch == '.') {
		num += string(self.source[self.index]);
		self.index++
		for character.IsDecimalDigit(getCharCodeAt(self.source, self.index)) {
			num += string(self.source[self.index]);
			self.index++
		}
		ch = self.source[self.index];
	}

	if ch == 'e' || ch == 'E' {
		num += string(self.source[self.index]);
		self.index++

		ch = self.source[self.index];
		if (ch == '+' || ch == '-') {
			num += string(self.source[self.index])
			self.index++
		}
		if (character.IsDecimalDigit(getCharCodeAt(self.source, self.index))) {
			for (character.IsDecimalDigit(getCharCodeAt(self.source, self.index))) {
				num += string(self.source[self.index]);
				self.index++
			}
		} else {
			self.throwUnexpectedToken("");
		}
	}

	if (character.IsIdentifierStart(getCharCodeAt(self.source, self.index))) {
		self.throwUnexpectedToken("");
	}
	value_number, _ := strconv.ParseFloat(num,32)
	return &RawToken{
		_type: token.NumericLiteral,
		value_number: float32(value_number),
		lineNumber: self.lineNumber,
		lineStart: self.lineStart,
		start: start,
		end: self.index,
	};
}

// https://tc39.github.io/ecma262/#sec-literals-string-literals
func (self *Scanner) scanStringLiteral() *RawToken {
	start := self.index;
	quote := self.source[start];
	if !(quote == '\'') && !(quote == '"') {
		panic("String literal must starts with a quote")
	}

	self.index++;
	octal := false;
	var str string

	for !self.eof() {
		ch := self.source[self.index]
		self.index++

		if ch == quote {
			// TODO quote supposed to be empty, not a space
			quote = ' ';
			break;
		} else if ch == '\\' {
			ch = self.source[self.index];
			self.index++
			if (&ch != nil || !character.IsLineTerminator(getCharCodeAt(string(ch), 0))) {
				switch (ch) {
					case 'u':
						if (self.source[self.index] == '{') {
							self.index++;
							str += self.scanUnicodeCodePointEscape();
						} else {
							unescapedChar := self.scanHexEscape(string(ch));
							if unescapedChar == "" {
								self.throwUnexpectedToken("");
							}
							str += unescapedChar;
						}
						break;
					case 'x':
						unescaped := self.scanHexEscape(string(ch));
						if unescaped == "" {
							self.throwUnexpectedToken(messages.GetInstance().InvalidHexEscapeSequence);
						}
						str += unescaped;
						break;
					case 'n':
						str += string('\n');
						break;
					case 'r':
						str += string('\r');
						break;
					case 't':
						str += string('\t');
						break;
					case 'b':
						str += string('\b');
						break;
					case 'f':
						str += string('\f');
						break;
					case 'v':
						str += string('\x0B');
						break;
					case '8':
					case '9':
						str += string(ch);
						self.tolerateUnexpectedToken("");
						break;

					default:
						if len(string(ch)) > 0 && character.IsOctalDigit(getCharCodeAt(string(ch), 0)) {
							octToDec := self.octalToDecimal(string(ch));
							octal = octToDec.octal || octal;
							str += string(octToDec.code);
						} else {
							str += string(ch);
						}
						break;
				}
			} else {
				self.lineNumber++;
				if (ch == '\r' && self.source[self.index] == '\n') {
					self.index++;
				}
				self.lineStart = self.index;
			}
		} else if (character.IsLineTerminator(getCharCodeAt(string(ch), 0))) {
			break;
		} else {
			str += string(ch);
		}
	}

	if (len(string(quote)) != 0) {
		self.index = start;
		self.throwUnexpectedToken("");
	}

	return &RawToken{
		_type: token.StringLiteral,
		value_string: str,
		octal: octal,
		lineNumber: self.lineNumber,
		lineStart: self.lineStart,
		start: start,
		end: self.index,
	};
}

// https://tc39.github.io/ecma262/#sec-template-literal-lexical-components

func (self *Scanner) scanTemplate() *RawToken {
	cooked := "";
	terminated := false;
	start := self.index;

	head := (self.source[start] == '`');
	tail := false;
	rawOffset := 2;

	self.index++;

	for !self.eof() {
		ch := self.source[self.index];
		self.index++
		if (ch == '`') {
			rawOffset = 1;
			tail = true;
			terminated = true;
			break;
		} else if (ch == '$') {
			if (self.source[self.index] == '{') {
				self.curlyStack = append(self.curlyStack, "${")
				self.index++;
				terminated = true;
				break;
			}
			cooked += string(ch);
		} else if (ch == '\\') {
			ch = self.source[self.index];
			self.index++
			if (!character.IsLineTerminator(getCharCodeAt(string(ch), 0))) {
				switch (ch) {
					case 'n':
						cooked += "\n";
						break;
					case 'r':
						cooked += "\r";
						break;
					case 't':
						cooked += "\t";
						break;
					case 'u':
						if (self.source[self.index] == '{') {
							self.index++;
							cooked += self.scanUnicodeCodePointEscape();
						} else {
							restore := self.index;
							unescapedChar := self.scanHexEscape(string(ch));
							if (&unescapedChar != nil) {
								cooked += unescapedChar;
							} else {
								self.index = restore;
								cooked += string(ch);
							}
						}
						break;
					case 'x':
						unescaped := self.scanHexEscape(string(ch));
						if (&unescaped == nil) {
							self.throwUnexpectedToken(messages.GetInstance().InvalidHexEscapeSequence);
						}
						cooked += unescaped;
						break;
					case 'b':
						cooked += "\b";
						break;
					case 'f':
						cooked += "\f";
						break;
					case 'v':
						cooked += "\v";
						break;

					default:
						if (ch == '0') {
							if character.IsDecimalDigit(getCharCodeAt(self.source, self.index)) {
								// Illegal: \01 \02 and so on
								self.throwUnexpectedToken(messages.GetInstance().TemplateOctalLiteral);
							}
							cooked += "\0";
						} else if (character.IsOctalDigit(getCharCodeAt(string(ch), 0))) {
							// Illegal: \1 \2
							self.throwUnexpectedToken(messages.GetInstance().TemplateOctalLiteral);
						} else {
							cooked += string(ch);
						}
						break;
				}
			} else {
				self.lineNumber++
				if (ch == '\r' && self.source[self.index] == '\n') {
					self.index++
				}
				self.lineStart = self.index;
			}
		} else if character.IsLineTerminator(getCharCodeAt(string(ch), 0)) {
			self.lineNumber++
			if (ch == '\r' && self.source[self.index] == '\n') {
				self.index++
			}
			self.lineStart = self.index;
			cooked += "\n";
		} else {
			cooked += string(ch);
		}
	}



	if (!terminated) {
		self.throwUnexpectedToken("");
	}

	if (!head) {
		self.curlyStack = self.curlyStack[:len(self.curlyStack) - 1]
	}

	return &RawToken{
		_type: token.Template,
		value_string: self.source[start + 1 : self.index - rawOffset],
		cooked: cooked,
		head: head,
		tail: tail,
		lineNumber: self.lineNumber,
		lineStart: self.lineStart,
		start: start,
		end: self.index,
	};
}



func (self *Scanner) testRegExp(pattern string, flags string){
	// TODO implement self after I understnd regex in Golang
}


func (self *Scanner) scanRegExpBody() string {
	ch := self.source[self.index];
	if ch != '/'{
		panic("Regular expression literal must start with a slash")
	}

	str := string(self.source[self.index]);
	self.index++
	classMarker := false;
	terminated := false;

	for !self.eof() {
		ch = self.source[self.index]
		self.index++
		str += string(ch)
		if (ch == '\\') {
			ch = self.source[self.index];
			self.index++
			// https://tc39.github.io/ecma262/#sec-literals-regular-expression-literals
			if character.IsLineTerminator(getCharCodeAt(string(ch), 0)) {
				self.throwUnexpectedToken(messages.GetInstance().UnterminatedRegExp);
			}
			str += string(ch);
		} else if character.IsLineTerminator(getCharCodeAt(string(ch), 0)) {
			self.throwUnexpectedToken(messages.GetInstance().UnterminatedRegExp);
		} else if classMarker {
			if (ch == ']') {
				classMarker = false;
			}
		} else {
			if (ch == '/') {
				terminated = true;
				break;
			} else if (ch == '[') {
				classMarker = true;
			}
		}
	}

	if (!terminated) {
		self.throwUnexpectedToken(messages.GetInstance().UnterminatedRegExp);
	}

	// Exclude leading and trailing slash.
	return string(str[1 : len(str) - 2])
}

func (self *Scanner) scanRegExpFlags() string {
	str := ""
	flags := ""
	for (!self.eof()) {
		ch := self.source[self.index];
		if !character.IsIdentifierPart(getCharCodeAt(string(ch), 0)) {
			break;
		}

		self.index++
		if (ch == '\\' && !self.eof()) {
			ch = self.source[self.index];
			if (ch == 'u') {
				self.index++
				restore := self.index;
				char := self.scanHexEscape("u");
				if (&char != nil) {
					flags += char
					for restore = self.index; restore < self.index; restore++ {
						str += string(self.source[restore]);
						str += "\\u";
					}
				} else {
					self.index = restore;
					flags += "u";
					str += "\\u"
				}
				self.tolerateUnexpectedToken("");
			} else {
				str += "\\"
				self.tolerateUnexpectedToken("");
			}
		} else {
			flags += string(ch);
			str += string(ch);
		}
	}

	return flags;
}




















