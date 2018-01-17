package error_handler

type Error struct {
	name string
	message string
	index int
	lineNumber int
	column int
	description string
}

type ErrorHandler struct{
	errors []*Error
	tolerant bool
}

func NewError(message string) *Error{
	return &Error{
		message: message,
	}
}

func NewErrorHandler() *ErrorHandler{
	return &ErrorHandler{
		errors: []*Error{},
		tolerant: false,
	}
}

func (self *ErrorHandler) recordError(error *Error) {
	self.errors = append(self.errors, error)
}


func (self *ErrorHandler) tolerate(error *Error) {
	if (self.tolerant) {
		self.recordError(error);
	} else {
		// TODO Figure out how to throw error
		//throw error;
	}
}

func (self *ErrorHandler) constructError(msg string, column int) *Error {
	error := NewError(msg);
	// TODO figure out how to throw!
	//try {
	//	throw error;
	//} catch (base) {
	//	/* istanbul ignore else */
	//	if (Object.create && Object.defineProperty) {
	//		error = Object.create(base);
	//		Object.defineProperty(error, 'column', { value: column });
	//	}
	//}
	/* istanbul ignore next */
	return error;
}

func (self *ErrorHandler) createError(index int, line int, col int, description string) *Error {
	msg := "Line " + string(line) + ": " + description;
	error := self.constructError(msg, col);
	error.index = index;
	error.lineNumber = line;
	error.description = description;
	return error;
}

func (self *ErrorHandler) throwError(index int, line int, col int, description string) {
	//throw self.createError(index, line, col, description);
}

func (self *ErrorHandler) tolerateError(index int, line int, col int, description string) {
	error := self.createError(index, line, col, description);
	if (self.tolerant) {
		self.recordError(error);
	} else {
		//throw error;
	}
}