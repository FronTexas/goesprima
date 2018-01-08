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
	start Position
	end Position
	source string
}

// TODO: Find a better way to differentiate scanner.Comment and comment-handler.Comment
type Comment_scanner struct {
	multiline bool
	slice []int
	_range []int
	loc SourceLocation
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


