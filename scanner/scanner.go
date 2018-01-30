package scanner

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
	Line   int
	Column int
}

type SourceLocation struct {
	Start  *Position
	End    *Position
	Source string
}

// TODO: Find a better way to differentiate scanner.Comment and comment-handler.Comment
type Comment_scanner struct {
	Multiline bool
	Slice []int
	Range []int
	Loc *SourceLocation
}

type RawToken struct {
	Type         token.Token
	Value_string string
	Value_number float32
	Pattern      string
	Flags        string
	Regex        regexp.Regexp
	Octal        bool
	Cooked       string
	Head         bool
	Tail         bool
	LineNumber   int
	LineStart    int
	Start        int
	End          int
}

type ScannerState struct {
	Index int
	LineNumber int
	LineStart int
}

type Scanner struct {
	Source string
	ErrorHandler ErrorHandler
	TrackComment bool
	IsModule bool
	Index int
	LineNumber int
	LineStart int
	CurlyStack []string
	Length int
}

func NewScanner(code string, handler ErrorHandler) *Scanner{
	var LineNumber int
	if LineNumber = 0; len(code) > 0 {
		LineNumber = 1
	}
	return &Scanner{
		Source: code,
		ErrorHandler: handler,
		TrackComment: false,
		IsModule: false,
		Length: len(code),
		Index: 0,
		LineNumber: LineNumber,
		LineStart : 0,
		CurlyStack: []string{},
	}
}

func (self *Scanner) SaveState() *ScannerState{
	return &ScannerState{
		Index: self.Index,
		LineNumber: self.LineNumber,
		LineStart: self.LineStart,
	}
}

func (self *Scanner) RestoreState(state *ScannerState) {
	self.Index = state.Index
	self.LineNumber = state.LineNumber
	self.LineStart = state.LineStart
}

func (self *Scanner) Eof() bool {
	return self.Index >= self.Length
}

// TODO: Uncomment these two functions after Error implementation is finished

func (self *Scanner) throwUnexpectedToken(message string) *Error{
	if message == "" {
		message = messages.GetInstance().UnexpectedTokenIllegal
	}

	return self.ErrorHandler.throwError(self.Index, self.LineNumber,
		self.Index - self.LineStart + 1, message)
}

func (self *Scanner) tolerateUnexpectedToken(message string){
	if message == ""{
		message = messages.GetInstance().UnexpectedTokenIllegal
	}
	self.ErrorHandler.tolerateError(self.Index, self.LineNumber,
		self.Index - self.LineStart +1, message)
}

func (self *Scanner) skipSingleLineComment(offset int) []*Comment_scanner{
	comments := []*Comment_scanner{}
	var start int
	var loc *SourceLocation
	if self.TrackComment {
		start = self.Index - offset
		loc = &SourceLocation{
			Start: &Position{
				Line:   self.LineNumber,
				Column: self.Index - self.LineStart - offset,
			},
		}
	}

	for !self.Eof() {
		ch := []rune(self.Source)[self.Index]
		self.Index++
		// TODO implement IsLineTerminator in character.go
		if character.IsLineTerminator(ch) {
			if self.TrackComment {
				loc.End = &Position{
					Line:   self.LineNumber,
					Column: self.Index - self.LineStart - 1,
				}
				entry := &Comment_scanner{
					multiline: false,
					slice: []int{start + offset, self.Index - 1},
					_range: []int{start, self.Index - 1},
					loc: loc,
				}
				comments = append(comments, entry)
			}

			if ch == 13 && []rune(self.Source)[self.Index] == 10 {
				self.Index++
			}
		}
		self.LineNumber++
		self.LineStart = self.Index
		return comments
	}

	if self.TrackComment {
		loc.End = &Position{
			Line:   self.LineNumber,
			Column: self.Index - self.LineStart,
		}

		entry := &Comment_scanner{
			multiline: false,
			slice: []int{start + offset, self.Index},
			_range: []int{start, self.Index},
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

	if self.TrackComment {
		start = self.Index - 2
		loc = &SourceLocation{
			Start: &Position{
				Line:   self.LineNumber,
				Column: self.Index - self.LineStart - 2,
			},
		}
	}

	for !self.Eof() {
		ch := []rune(self.Source)[self.Index]
		self.Index++
		// TODO implement IsLineTerminator in character.go
		if character.IsLineTerminator(ch) {
			if ch == 0x0D && []rune(self.Source)[self.Index] == 0x0A{
				self.Index++
			}
			self.LineNumber++
			self.Index++
			self.LineStart = self.Index
		}else if ch == 0x2A{
			if []rune(self.Source)[self.Index + 1] == 0x2F {
				self.Index += 2
				if self.TrackComment {
					loc.End = &Position{
						Line:   self.LineNumber,
						Column: self.Index - self.LineStart,
					}
					entry := &Comment_scanner{
						multiline: true,
						slice: []int{start + 2, self.Index - 2},
						_range: []int{start, self.Index},
						loc: loc,
					}
					comments = append(comments, entry)
				}
				return comments
			}
			self.Index++
		}else{
			self.Index++
		}
	}
	// Ran off the End of the file - the whole thing is a comment
	if self.TrackComment {
		loc.End = &Position{
			Line:   self.LineNumber,
			Column: self.Index - self.LineStart,
		}

		entry := &Comment_scanner{
			multiline: true,
			slice: []int{start + 2, self.Index},
			_range: []int{start, self.Index},
			loc: loc,
		}
		comments = append(comments, entry)
	}
	// TODO uncomment self once the method is implemented
	// self.tolerateUnexpectedToken()
	return comments
}

func (self *Scanner) ScanComments() []*Comment_scanner{
	var comments []*Comment_scanner
	if self.TrackComment{
		comments = []*Comment_scanner{}
	}

	start := self.Index == 0

	for !self.Eof() {
		ch := getCharCodeAt(self.Source, self.Index)

		if character.IsWhiteSpace(ch) {
			self.Index++
		}else if character.IsLineTerminator(ch){
			self.Index++
			if ch == 0x0D && getCharCodeAt(self.Source, self.Index) == 0X0A {
				self.Index++
			}
			self.LineNumber++
			self.LineStart = self.Index
			start = true
		}else if ch == 0x2F {
			ch = getCharCodeAt(self.Source, self.Index + 1)
			if ch == 0x2F {
				self.Index += 2
				comment := self.skipSingleLineComment(2)
				if self.TrackComment{
					comments = append(comments, comment...)
				}
				start = true
			}else if ch == 0x2A {
				self.Index += 2
				comment := self.skipSingleLineComment(3)
				if self.TrackComment {
					comments = append(comments, comment...)
				}
			}else {
				break
			}
		}else if start && ch == 0x2D {
			if getCharCodeAt(self.Source, self.Index + 1) == 0x2D && getCharCodeAt(self.Source, self.Index + 2) == 0x3E {
				self.Index += 3
				comment := self.skipSingleLineComment(3)
				if self.TrackComment {
					comments = append(comments, comment...)
				}
			}else{
				break
			}
		}else if (ch == 0x3C && !self.IsModule){
			if string([]rune(self.Source)[self.Index + 1 : self.Index + 4]) == "!--"{
				self.Index += 4
				comment := self.skipSingleLineComment(4)
				if self.TrackComment {
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
	cp := getCharCodeAt(self.Source, i)

	if (cp >= 0xD800 && cp <= 0xDBFF) {
		second  := getCharCodeAt(self.Source, i + 1)
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
		if !self.Eof() && character.IsHexDigit(getCharCodeAt(self.Source, self.Index)) {
			code = code * 16 + hexValue(string(self.Source[self.Index + 1]))
			self.Index += 1
		} else {
			return ""
		}
	}
	//return String.fromCharCode(code)
	toReturn := string(code)
	return toReturn
}

func (self *Scanner) scanUnicodeCodePointEscape() string {
	ch := self.Source[self.Index]
	code := 0

	// At least, one hex digit is required.
	if (ch == '}') {
		self.throwUnexpectedToken("")
	}

	for !self.Eof() {
		ch = self.Source[self.Index]
		self.Index += 1
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
	start := self.Index
	self.Index += 1
	for !self.Eof() {
		ch := getCharCodeAt(self.Source, self.Index)
		if (ch == 0x5C) {
			// Blackslash (U+005C) marks Unicode escape sequence.
			self.Index = start
			return self.getComplexIdentifier()
		} else if (ch >= 0xD800 && ch < 0xDFFF) {
			// Need to handle surrogate pairs.
			self.Index = start
			return self.getComplexIdentifier()
		}
		if (character.IsIdentifierPart(ch)) {
			self.Index++
		} else {
			break
		}
	}

	return self.Source[start : self.Index]
}

func (self *Scanner) getComplexIdentifier() string {
	cp := self.codePointAt(self.Index)
	id := character.FromCodePoint(cp)
	self.Index += len(id)

	// '\u' (U+005C, U+0075) denotes an escaped character.
	var ch string
	if (cp == 0x5C) {
		if (getCharCodeAt(self.Source, self.Index) != 0x75) {
			self.throwUnexpectedToken("")
		}
		self.Index++
		if (self.Source[self.Index] == '{') {
			self.Index++
			ch = self.scanUnicodeCodePointEscape()
		} else {
			ch = self.scanHexEscape("u")
			if ch == "" || ch == "\\" || !character.IsIdentifierStart(getCharCodeAt(ch, 0)) {
				self.throwUnexpectedToken("")
			}
		}
		id = ch
	}

	for !self.Eof() {
		cp = self.codePointAt(self.Index)
		if (!character.IsIdentifierPart(cp)) {
			break
		}
		ch = character.FromCodePoint(cp)
		id += ch
		self.Index += len(ch)

		// '\u' (U+005C, U+0075) denotes an escaped character.
		if (cp == 0x5C) {
			id = id[0:len(id) - 1]

			if (getCharCodeAt(self.Source, self.Index) != 0x75) {
				self.throwUnexpectedToken("")
			}
			self.Index++
			if (self.Source[self.Index] == '{') {
				self.Index++
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
	// \0 is not Octal escape sequence
	octal := (ch != "0");
	code := octalValue(ch);

	if !self.Eof() && character.IsOctalDigit(getCharCodeAt(self.Source, self.Index)) {
		octal = true;
		code = code * 8 + octalValue(string(self.Source[self.Index]))
		self.Index++

		// 3 digits are only allowed when string starts
		// with 0, 1, 2, 3
		if strings.Index("0123", string(ch)) >= 0 && !self.Eof() && character.IsOctalDigit(getCharCodeAt(self.Source, self.Index)) {
			code = code * 8 + octalValue(string(self.Source[self.Index]));
			self.Index++
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
	start := self.Index;

	// Backslash (U+005C) starts an escaped character.
	var id string
	if id = self.getIdentifier(); getCharCodeAt(self.Source, start) == 0x5C{
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

	if (_type != token.Identifier && (start + len(id) != self.Index)) {
		restore := self.Index;
		self.Index = start;
		self.tolerateUnexpectedToken(messages.GetInstance().InvalidEscapedReservedWord);
		self.Index = restore;
	}

	return &RawToken{
		Type:         _type,
		Value_string: id,
		LineNumber:   self.LineNumber,
		LineStart:    self.LineStart,
		Start:        start,
		End:          self.Index,
	};
}

// https://tc39.github.io/ecma262/#sec-punctuators
func(self *Scanner) scanPunctuator() *RawToken {
	start := self.Index;

	// Check for most common single-character punctuators.
	str := string(self.Source[self.Index]);
	switch (str) {

	case "(":
	case "{":
		if (str == "{") {
			self.CurlyStack = append(self.CurlyStack, "{")
		}
		self.Index++;
		break;
	case ".":
		self.Index++;
		if (self.Source[self.Index] == '.' && self.Source[self.Index + 1] == '.') {
			// Spread operator: ...
			self.Index += 2;
			str = "...";
		}
		break;
	case "}":
		self.Index++;
		self.CurlyStack = self.CurlyStack[:len(self.CurlyStack) - 1]
		break;
	case ")":
	case ";":
	case ",":
	case "[":
	case "]":
	case ":":
	case "?":
	case "~":
		self.Index++;
		break;
	default:
		// 4-character punctuator.
		str = self.Source[self.Index : 4]
		if (str == ">>>=") {
			self.Index += 4
		} else {
			// 3-character punctuators.
			str = str[0:3]
			if str == "===" || str == "!==" || str == ">>>" ||
				str == "<<=" || str == ">>=" || str == "**=" {
				self.Index += 3
			} else {
				// 2-character punctuators.
				str = str[0:2]
				if str == "&&" || str == "||" || str == "==" || str == "!=" ||
					str == "+=" || str == "-=" || str == "*=" || str == "/=" ||
					str == "++" || str == "--" || str == "<<"|| str == ">>" ||
					str == "&="|| str == "|="|| str == "^=" || str == "%=" ||
					str == "<=" || str == ">="|| str == "=>" || str == "**" {
					self.Index += 2
				} else {
					// 1-character punctuators.
					str = string(self.Source[self.Index])
					if strings.Index("<>=!+-*%&|^/", str) >= 0 {
						self.Index++
					}
				}
			}
		}
	}

	if (self.Index == start) {
		self.throwUnexpectedToken("");
	}

	return &RawToken{
		Type:         token.Punctuator,
		Value_string: str,
		LineNumber:   self.LineNumber,
		LineStart:    self.LineStart,
		Start:        start,
		End:          self.Index,
	};
}

// https://tc39.github.io/ecma262/#sec-literals-numeric-literals
func (self *Scanner) scanHexLiteral(start int) *RawToken {
	num := "";

	for !self.Eof() {
		if !character.IsHexDigit(getCharCodeAt(self.Source, self.Index)) {
			break;
		}
		num += string(self.Source[self.Index])
		self.Index++
	}

	if len(num) == 0 {
		self.throwUnexpectedToken("");
	}

	if (character.IsIdentifierStart(getCharCodeAt(self.Source, self.Index))) {
		self.throwUnexpectedToken("");
	}
	value_number, _ := strconv.ParseInt("0x" + num, 16, 0)
	return &RawToken{
		Type:         token.NumericLiteral,
		Value_number: float32(value_number),
		LineNumber:   self.LineNumber,
		LineStart:    self.LineStart,
		Start:        start,
		End:          self.Index,
	};
}

func (self *Scanner) scanBinaryLiteral(start int) *RawToken {
	num := "";
	var ch rune

	for !self.Eof() {
		ch = rune(self.Source[self.Index]);
		if ch != '0' && ch != '1' {
			break;
		}
		num += string(self.Source[self.Index]);
		self.Index++
	}

	if len(num) == 0 {
		// only 0b or 0B
		self.throwUnexpectedToken("");
	}

	if (!self.Eof()) {
		ch = getCharCodeAt(self.Source, self.Index);
		/* istanbul ignore else */
		if (character.IsIdentifierStart(ch) || character.IsDecimalDigit(ch)) {
			self.throwUnexpectedToken("");
		}
	}

	value_number, _ := strconv.ParseInt(num, 2, 0)
	return &RawToken{
		Type:         token.NumericLiteral,
		Value_number: float32(value_number),
		LineNumber:   self.LineNumber,
		LineStart:    self.LineStart,
		Start:        start,
		End:          self.Index,
	};
}

func (self *Scanner) scanOctalLiteral(prefix string, start int) *RawToken {
	num := "";
	octal := false;

	if (character.IsOctalDigit(getCharCodeAt(prefix, 0))) {
		octal = true;
		num = "0" + string(self.Source[self.Index]);
		self.Index++
	} else {
		self.Index++
	}

	for !self.Eof() {
		if !character.IsOctalDigit(getCharCodeAt(self.Source, self.Index)) {
			break;
		}
		num += string(self.Source[self.Index]);
		self.Index++
	}

	if (!octal && len(num) == 0) {
		// only 0o or 0O
		self.throwUnexpectedToken("");
	}

	if character.IsIdentifierStart(getCharCodeAt(self.Source, self.Index)) || character.IsDecimalDigit(getCharCodeAt(self.Source, self.Index)) {
		self.throwUnexpectedToken("");
	}
	value_number, _ := strconv.ParseInt(num, 2, 0)
	return &RawToken{
		Type:         token.NumericLiteral,
		Value_number: float32(value_number),
		Octal:        octal,
		LineNumber:   self.LineNumber,
		LineStart:    self.LineStart,
		Start:        start,
		End:          self.Index,
	};
}

func (self *Scanner) isImplicitOctalLiteral() bool {
	// Implicit Octal, unless there is a non-Octal digit.
	// (Annex B.1.1 on Numeric Literals)
	for i := self.Index + 1; i < self.Length; i++ {
		ch := self.Source[i];
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
	start := self.Index;
	ch := self.Source[start];

	// TODO figure out how to use assert instead of if statement
	if !character.IsDecimalDigit(getCharCodeAt(string(ch), 0)) && !(ch == '.'){
		panic("Numeric literal must Start with a decimal digit or a decimal point")
	}

	var num string;
	if (ch != '.') {
		num = string(self.Source[self.Index]);
		self.Index++
		ch = self.Source[self.Index];

		// Hex number starts with '0x'.
		// Octal number starts with '0'.
		// Octal number in ES6 starts with '0o'.
		// Binary number in ES6 starts with '0b'.
		if num == "0" {
			if ch == 'x' || ch == 'X' {
				self.Index++;
				return self.scanHexLiteral(start);
			}
			if ch == 'b' || ch == 'B' {
				self.Index++
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

		for (character.IsDecimalDigit(getCharCodeAt(self.Source, self.Index))) {
			num += string(self.Source[self.Index])
			self.Index++
		}
		ch = self.Source[self.Index];
	}

	if (ch == '.') {
		num += string(self.Source[self.Index]);
		self.Index++
		for character.IsDecimalDigit(getCharCodeAt(self.Source, self.Index)) {
			num += string(self.Source[self.Index]);
			self.Index++
		}
		ch = self.Source[self.Index];
	}

	if ch == 'e' || ch == 'E' {
		num += string(self.Source[self.Index]);
		self.Index++

		ch = self.Source[self.Index];
		if (ch == '+' || ch == '-') {
			num += string(self.Source[self.Index])
			self.Index++
		}
		if (character.IsDecimalDigit(getCharCodeAt(self.Source, self.Index))) {
			for (character.IsDecimalDigit(getCharCodeAt(self.Source, self.Index))) {
				num += string(self.Source[self.Index]);
				self.Index++
			}
		} else {
			self.throwUnexpectedToken("");
		}
	}

	if (character.IsIdentifierStart(getCharCodeAt(self.Source, self.Index))) {
		self.throwUnexpectedToken("");
	}
	value_number, _ := strconv.ParseFloat(num,32)
	return &RawToken{
		Type:         token.NumericLiteral,
		Value_number: float32(value_number),
		LineNumber:   self.LineNumber,
		LineStart:    self.LineStart,
		Start:        start,
		End:          self.Index,
	};
}

// https://tc39.github.io/ecma262/#sec-literals-string-literals
func (self *Scanner) scanStringLiteral() *RawToken {
	start := self.Index;
	quote := self.Source[start];
	if !(quote == '\'') && !(quote == '"') {
		panic("String literal must starts with a quote")
	}

	self.Index++;
	octal := false;
	var str string

	for !self.Eof() {
		ch := self.Source[self.Index]
		self.Index++

		if ch == quote {
			// TODO quote supposed to be empty, not a space
			quote = ' ';
			break;
		} else if ch == '\\' {
			ch = self.Source[self.Index];
			self.Index++
			if (&ch != nil || !character.IsLineTerminator(getCharCodeAt(string(ch), 0))) {
				switch (ch) {
				case 'u':
					if (self.Source[self.Index] == '{') {
						self.Index++;
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
				self.LineNumber++;
				if (ch == '\r' && self.Source[self.Index] == '\n') {
					self.Index++;
				}
				self.LineStart = self.Index;
			}
		} else if (character.IsLineTerminator(getCharCodeAt(string(ch), 0))) {
			break;
		} else {
			str += string(ch);
		}
	}

	if (len(string(quote)) != 0) {
		self.Index = start;
		self.throwUnexpectedToken("");
	}

	return &RawToken{
		Type:         token.StringLiteral,
		Value_string: str,
		Octal:        octal,
		LineNumber:   self.LineNumber,
		LineStart:    self.LineStart,
		Start:        start,
		End:          self.Index,
	};
}

// https://tc39.github.io/ecma262/#sec-template-literal-lexical-components

func (self *Scanner) scanTemplate() *RawToken {
	cooked := "";
	terminated := false;
	start := self.Index;

	head := (self.Source[start] == '`');
	tail := false;
	rawOffset := 2;

	self.Index++;

	for !self.Eof() {
		ch := self.Source[self.Index];
		self.Index++
		if (ch == '`') {
			rawOffset = 1;
			tail = true;
			terminated = true;
			break;
		} else if (ch == '$') {
			if (self.Source[self.Index] == '{') {
				self.CurlyStack = append(self.CurlyStack, "${")
				self.Index++;
				terminated = true;
				break;
			}
			cooked += string(ch);
		} else if (ch == '\\') {
			ch = self.Source[self.Index];
			self.Index++
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
					if (self.Source[self.Index] == '{') {
						self.Index++;
						cooked += self.scanUnicodeCodePointEscape();
					} else {
						restore := self.Index;
						unescapedChar := self.scanHexEscape(string(ch));
						if (&unescapedChar != nil) {
							cooked += unescapedChar;
						} else {
							self.Index = restore;
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
						if character.IsDecimalDigit(getCharCodeAt(self.Source, self.Index)) {
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
				self.LineNumber++
				if (ch == '\r' && self.Source[self.Index] == '\n') {
					self.Index++
				}
				self.LineStart = self.Index;
			}
		} else if character.IsLineTerminator(getCharCodeAt(string(ch), 0)) {
			self.LineNumber++
			if (ch == '\r' && self.Source[self.Index] == '\n') {
				self.Index++
			}
			self.LineStart = self.Index;
			cooked += "\n";
		} else {
			cooked += string(ch);
		}
	}



	if (!terminated) {
		self.throwUnexpectedToken("");
	}

	if (!head) {
		self.CurlyStack = self.CurlyStack[:len(self.CurlyStack) - 1]
	}

	return &RawToken{
		Type:         token.Template,
		Value_string: self.Source[start + 1 : self.Index - rawOffset],
		Cooked:       cooked,
		Head:         head,
		Tail:         tail,
		LineNumber:   self.LineNumber,
		LineStart:    self.LineStart,
		Start:        start,
		End:          self.Index,
	};
}



func (self *Scanner) testRegExp(pattern string, flags string){
	// TODO implement self after I understnd Regex in Golang
}


func (self *Scanner) scanRegExpBody() string {
	ch := self.Source[self.Index];
	if ch != '/'{
		panic("Regular expression literal must Start with a slash")
	}

	str := string(self.Source[self.Index]);
	self.Index++
	classMarker := false;
	terminated := false;

	for !self.Eof() {
		ch = self.Source[self.Index]
		self.Index++
		str += string(ch)
		if (ch == '\\') {
			ch = self.Source[self.Index];
			self.Index++
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
	for (!self.Eof()) {
		ch := self.Source[self.Index];
		if !character.IsIdentifierPart(getCharCodeAt(string(ch), 0)) {
			break;
		}

		self.Index++
		if (ch == '\\' && !self.Eof()) {
			ch = self.Source[self.Index];
			if (ch == 'u') {
				self.Index++
				restore := self.Index;
				char := self.scanHexEscape("u");
				if (&char != nil) {
					flags += char
					for restore = self.Index; restore < self.Index; restore++ {
						str += string(self.Source[restore]);
						str += "\\u";
					}
				} else {
					self.Index = restore;
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

func (self *Scanner) ScanRegExp() (*RawToken, error)  {
	// TODO implement self after implementing testRegExp
}


func (self *Scanner) Lex() *RawToken {
	if self.Eof() {
		return &RawToken{
			Type:         token.EOF,
			Value_string: "",
			LineNumber:   self.LineNumber,
			LineStart:    self.LineStart,
			Start:        self.Index,
			End:          self.Index,
		};
	}

	cp := getCharCodeAt(self.Source, self.Index);

	if (character.IsIdentifierStart(cp)) {
		return self.scanIdentifier();
	}

	// Very common: ( and ) and ;
	if (cp == 0x28 || cp == 0x29 || cp == 0x3B) {
		return self.scanPunctuator();
	}

	// String literal starts with single quote (U+0027) or double quote (U+0022).
	if (cp == 0x27 || cp == 0x22) {
		return self.scanStringLiteral();
	}

	// Dot (.) U+002E can also Start a floating-point number, hence the need
	// to check the next character.
	if (cp == 0x2E) {
		if character.IsDecimalDigit(getCharCodeAt(self.Source, self.Index + 1)) {
			return self.scanNumericLiteral();
		}
		return self.scanPunctuator();
	}

	if (character.IsDecimalDigit(cp)) {
		return self.scanNumericLiteral();
	}

	// Template literals Start with ` (U+0060) for template head
	// or } (U+007D) for template middle or template Tail.
	if (cp == 0x60 || (cp == 0x7D && self.CurlyStack[len(self.CurlyStack) - 1] == "${")) {
		return self.scanTemplate();
	}

	// Possible identifier Start in a surrogate pair.
	if (cp >= 0xD800 && cp < 0xDFFF) {
		if (character.IsIdentifierStart(self.codePointAt(self.Index))) {
			return self.scanIdentifier();
		}
	}

	return self.scanPunctuator();
}