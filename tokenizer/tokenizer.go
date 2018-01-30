package tokenizer

import (
	"goesprima/scanner"
	"goesprima/token"
	"goesprima/error_handler"
)

type ReaderEntry string

type _regexBufferEntry struct{
	pattern string
	flags string
}

type BufferEntry struct {
    _type string
    value string
    regex _regexBufferEntry
    _range []int
    loc *scanner.SourceLocation
}

type Reader struct {
	values []ReaderEntry
	curly int
	paren int
}

func NewReader() *Reader{
	return &Reader{
		values: []ReaderEntry{},
		curly: -1,
		paren: -1,
	}
}

func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// A function following one of those tokens is an expression.
func (self *Reader) beforeFunctionExpression(t string) bool{
	return Index([]string{"(", "{", "[", "in", "typeof", "instanceof", "new",
		"return", "case", "delete", "throw", "void",
		// assignment operators
		"=", "+=", "-=", "*=", "**=", "/=", "%=", "<<=", ">>=", ">>>=",
		"&=", "|=", "^=", ",",
		// binary/unary operators
		"+", "-", "*", "**", "/", "%", "++", "--", "<<", ">>", ">>>", "&",
		"|", "^", "!", "~", "&&", "||", "?", ":", "===", "==", ">=",
		"<=", "<", ">", "!=", "!=="}, t) >= 0;
}
func (self *Reader) isRegexStart() bool {
	previous := self.values[len(self.values) - 1]
	regex := len(previous) != 0

	switch previous {
		case "self":
		case "]":
			regex = false;
			break;

		case ")":
			keyword := self.values[self.paren - 1];
			regex = (keyword == "if" || keyword == "while" || keyword == "for" || keyword == "with");
			break;
		case "}":
			// Dividing a function by anything makes little sense,
			// but we have to check for that.
			regex = true;
			if (self.values[self.curly - 3] == "function") {
				// Anonymous function, e.g. function(){} /42
				check := self.values[self.curly - 4];
				if &check != nil {
					regex = !self.beforeFunctionExpression(string(check))
				}else{
					regex =  false;
				}
			} else if (self.values[self.curly - 4] == "function") {
				// Named function, e.g. function f(){} /42/
				check := self.values[self.curly - 5];
				if &check != nil {
					regex = !self.beforeFunctionExpression(string(check))
				}else{
					regex = true;
				}
			}
			break;
		default:
			break;
	}

	return regex;
}


func (self *Reader) push(tkn *scanner.RawToken) {
	if (tkn.Type == token.Punctuator || tkn.Type == token.Keyword) {
		if (tkn.Value_string == "{") {
			self.curly = len(self.values)
		} else if (tkn.Value_string == "(") {
			self.paren = len(self.values)
		}
		self.values = append(self.values, ReaderEntry(tkn.Value_string))
	} else {
		self.values = append(self.values, "")
	}
}

type Config struct {
	tolerant bool
	comment bool
	_range bool
	loc bool
}

type Tokenizer struct {
	ErrorHandler *error_handler.ErrorHandler
	Scanner *scanner.Scanner
	TrackRange bool
	TrackLoc bool
	Buffer []*BufferEntry
	Reader *Reader
}

func NewTokenizer(code string, config Config) *Tokenizer {
	var tolerant bool
	var trackComment bool
	var trackRange bool
	var trackLoc bool

	if &config != nil {
		tolerant = config.tolerant
		trackComment = config.comment
		trackRange = config._range
		trackLoc = config.loc
	}else{
		tolerant = false
		trackComment = false
		trackRange = false
	}

	errorHandler := &error_handler.ErrorHandler{
		[]*error_handler.Error{},
		tolerant,
	}
	_scanner := scanner.NewScanner(code, errorHandler)
	_scanner.TrackComment = trackComment


	return &Tokenizer{
		ErrorHandler: errorHandler,
		Scanner: _scanner,
		TrackRange: trackRange,
		TrackLoc: trackLoc,
		Buffer: []*BufferEntry{},
		Reader: NewReader(),
	}
}

func (self *Tokenizer) errors() []*error_handler.Error {
	return self.ErrorHandler.Errors;
}

func (self *Tokenizer) getNextToken() *BufferEntry{
	if (len(self.Buffer) == 0) {

		comments := self.Scanner.ScanComments()
		if (self.Scanner.TrackComment) {
			for i := 0; i < len(comments); i++ {
				e := comments[i];
				value := self.Scanner.Source[e.Slice[0] : e.Slice[1]];
				var _type string
				if e.Multiline {
					_type = "BlockComment"
				}else {
					_type = "LineComment"
				}
				comment:= &BufferEntry {
					_type: _type,
					value: value,
				};
				if (self.TrackRange) {
					comment._range = e.Range;
				}
				if (self.TrackLoc) {
					comment.loc = e.Loc;
				}
				self.Buffer = append(self.Buffer, comment)
			}
		}

		if (!self.Scanner.Eof()) {

			var loc scanner.SourceLocation
			if (self.TrackLoc) {
				loc = scanner.SourceLocation{
					Start: &scanner.Position{
						Line: self.Scanner.LineNumber,
						Column: self.Scanner.Index,
					},
				}
			}
			var _token *scanner.RawToken
			maybeRegex := (self.Scanner.Source[self.Scanner.Index] == '/') && self.Reader.isRegexStart();
			if maybeRegex {
				state := self.Scanner.SaveState();
				var err error
				_token, err = self.Scanner.ScanRegExp()
				if err != nil {
					self.Scanner.RestoreState(state)
					_token = self.Scanner.Lex()
				}
			} else {
				_token = self.Scanner.Lex();
			}
			self.Reader.push(_token)
			entry:= &BufferEntry{
				_type: token.GetTokenName(_token.Type),
				value: self.Scanner.Source[_token.Start : _token.End],
			};
			if (self.TrackRange) {
				entry._range = []int{_token.Start, _token.End}
			}
			if (self.TrackLoc) {
				loc.End = &scanner.Position{
					Line: self.Scanner.LineNumber,
					Column: self.Scanner.Index- self.Scanner.LineStart,
				};
				entry.loc = &loc;
			}
			if (_token.Type == token.RegularExpression) {
				pattern := _token.Pattern;
				flags := _token.Flags;
				entry.regex = _regexBufferEntry{ pattern, flags };
			}
			self.Buffer = append(self.Buffer, entry)
		}
	}

	toReturn := self.Buffer[0]
	if len(self.Buffer) > 1 {
		self.Buffer = self.Buffer[1:]
	}else{
		self.Buffer = []*BufferEntry{}
	}
	return toReturn
}












