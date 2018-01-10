package main

import(
	"strings"
	"regexp"
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
	_type Token
	value_string string
	value_number int
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
			if ch == 0x0D && []rune(self.source)[self.index + 1] == 0x0A{
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
		code = code * 8 + octalValue(string(self.source[self.index+1]))
		self.index++

		// 3 digits are only allowed when string starts
		// with 0, 1, 2, 3
		if strings.Index("0123", string(ch)) >= 0 && !self.eof() && character.IsOctalDigit(getCharCodeAt(self.source, self.index)) {
			code = code * 8 + octalValue(string(self.source[self.index + 1]));
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































