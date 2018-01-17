package tokenizer

import (
	"goesprima/scanner"
	"goesprima/token"
	"goesprima/error_handler"
	"regexp"
)

type ReaderEntry string

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

type BufferEntry struct {
    _type string
    value string
    regex struct {
        pattern string;
        flags string;
    }
    _range []int
    loc scanner.SourceLocation
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


func (self *Reader) push(tkn *RawToken) {
	if (tkn._type == token.Punctuator || tkn._type == token.Keyword) {
		if (tkn.value_string == "{") {
			self.curly = len(self.values)
		} else if (tkn.value_string == "(") {
			self.paren = len(self.values)
		}
		self.values = append(self.values, ReaderEntry(tkn.value_string))
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










