package main 

// TODO implement ErrorHandler class 
type ErrorHandler interface{
	throwError(index int, line  int, col int, description string) *Error
	tolerateError(index int, line int, col int, description string) *Error
}
type Error interface{}

