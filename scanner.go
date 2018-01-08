package main 

type Position struct {
	line int 
	column int 
}

type SourceLocation struct {
	start Position
	end Position
	source string
}