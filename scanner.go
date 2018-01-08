package main 

import(
	"strings"
	"regexp"
)

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

// func (self *Scanner) throwUnexpectedToken(message string) *Error{
// 	return self.errorHandler.throwError(self.index, self.lineNumber,
// 		self.index - self.lineStart + 1, message)
// }

// func (self *Scanner) tolerateUnexpectedToken(message){
// 	self.errorHandler.tolerateError(self.index, self.lineNumber,
// 		self.index - self.lineStart +1, message)
// }

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
		if IsLineTerminator(ch) {
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
		if IsLineTerminator(ch) {
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
	// TODO uncomment this once the method is implemented
	// self.tolerateUnexpectedToken()
	return comments
}



































