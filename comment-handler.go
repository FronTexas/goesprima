package main 

type (
	Comment interface {
		type() string 
		value() string
	}

	Entry interface {
		comment() Comment 
		start() int
	}

	NodeInfo interface { 
		node() interface{}
		start() int
	}
)

type CommentHandler struct { 
	attach bool
	comments []Comment
	stack []NodeInfo
	leading []Entry 
	trailing Entry[]
}

func NewCommentHandler() CommentHandler {
	return CommentHandler{
		attach: false, 
		comments: []Comment{},
		stack: []NodeInfo{},
		leading: []Entry{},
		trailing: []Entry{},
	}
}


func (self CommentHandler) insertInnerComments(node, metadata interface{}){
	if node.type == Syntax.BlockStatement && len(node.body) == 0 {
		innerComments := Comment{}[]
		for i:= len(self.leading - 1); i >= 0; i-- {
			entry := self.leading[i]
			if metadata.end.offset >= entry.start() { 
				innerComments = append([]Comment{entry.comment}, innerComments...)
				self.leading = append(self.leading[:i], self.leading[i+1:]...)
				self.trailing = append(self.trailing[:i], self.trailing[i+1:]...)
			}
		}

		if len(innerComments) > 0 { 
			node.innerComments = innerComments
		}
	}
}

func (self CommentHandler) findTrailingComments(metadata interface{}) []Comment{
	trailingComments := []Comment{}

	if len(self.trailing) > 0 { 
		for i:= len(self.trailing) - 1; i >= 0; i-- {
			entry := self.trailing[i]
			if entry.start() >= metadata.end.offset {
				trailingComments = append([]Comment{entry.comment}, trailingComments)
			}
		}
		len(self.trailing) = 0 
		return trailingComments
	}

	last := self.stack[len(self.stack) - 1]
	if len(last) > 0 && len(last.node.trailingComments) > 0 { 
		firstComment := last.node.trailingComments[0]
		if len(firstComment) > 0 && firstComment.range[0] >= metadata.end.offset { 
			trailingComments = last.node.trailingComments
			// todo delete last.node.trailingComments
		}
	}
	return trailingComments	
}

func (self CommentHandler) findLeadingComments()