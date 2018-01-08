package main 

type (
	Comment interface {
		_type() string 
		value() string
		_range() [2]int
	}

	Entry interface {
		comment() Comment 
		start() int
	}

	NodeInfo struct { 
		node *Node
		start int
	}
)

type CommentHandler struct { 
	attach bool
	comments []Comment
	stack []*NodeInfo
	leading []Entry 
	trailing []Entry
}

type Node struct {
	_type string
	// I am not sure if body can be an array of Node pointers
	body []*Node 
	innerComments []Comment
	trailingComments []Comment
	leadingComments []Comment
}

type (
	Metadata struct {
		start struct {
			line int 
			column int
			offset int
		}
		end struct {
			line int 
			column int 
			offset int
		}
	}
)

func NewCommentHandler() CommentHandler {
	return CommentHandler{
		attach: false, 
		comments: []Comment{},
		stack: []*NodeInfo{},
		leading: []Entry{},
		trailing: []Entry{},
	}
}


func (self CommentHandler) insertInnerComments(node Node, metadata Metadata){
	if node._type == Syntax["BlockStatement"] && len(node.body) == 0 {
		innerComments := []Comment{}
		for i:= len(self.leading) - 1; i >= 0; i-- {
			entry := self.leading[i]
			if metadata.end.offset >= entry.start() { 
				innerComments = append([]Comment{entry.comment()}, innerComments...)
				self.leading = append(self.leading[:i], self.leading[i+1:]...)
				// splicing
				self.trailing = append(self.trailing[:i], self.trailing[i+1:]...)
			}
		}

		if len(innerComments) > 0 { 
			node.innerComments = innerComments
		}
	}
}

func (self CommentHandler) findTrailingComments(metadata Metadata) []Comment{
	trailingComments := []Comment{}

	if len(self.trailing) > 0 { 
		for i:= len(self.trailing) - 1; i >= 0; i-- {
			entry := self.trailing[i]
			if entry.start() >= metadata.end.offset {
				// unshifting
				trailingComments = append([]Comment{entry.comment()}, trailingComments...)
			}
		}
		self.trailing = []Entry{}
		return trailingComments
	}

	last := self.stack[len(self.stack) - 1]
	if last != nil && last.node.trailingComments != nil { 
		firstComment := last.node.trailingComments[0]
		if firstComment != nil && firstComment._range()[0] >= metadata.end.offset { 
			trailingComments = last.node.trailingComments
			// TODO delete last.node.trailingComments
		}
	}
	return trailingComments	
}

func (self CommentHandler) findLeadingComments(metadata Metadata) []Comment{
	leadingComments := []Comment{}
	var target *Node
	for len(self.stack) > 0 { 
		entry := self.stack[len(self.stack) - 1]
		if entry != nil && entry.start >= metadata.start.offset { 
			target = entry.node
			self.stack = self.stack[:len(self.stack)]
		}else { 
			break;
		}
	}

	if target != nil { 
		var count int 
		if count = 0; target.leadingComments != nil { 
			count = len(target.leadingComments)
		}

		for i := count - 1; i >= 0; i-- {
			comment := target.leadingComments[i]
			if comment._range()[1] <= metadata.start.offset { 
				leadingComments = append([]Comment{comment}, leadingComments...)
				target.leadingComments = append(target.leadingComments[:i], target.leadingComments[i+1:]...)
			}
		}

		if target.leadingComments != nil && len(target.leadingComments) == 0 {
			// TODO delete target.leadingComments
		}
		return leadingComments
	}

	for i := len(self.leading) - 1; i >= 0; i-- {
		entry := self.leading[i]
		if entry.start() <= metadata.start.offset {
			leadingComments = append([]Comment{entry.comment()}, leadingComments...)
			self.leading = append(self.leading[:i], self.leading[i+1:]...)
		}
	}
	return leadingComments
}

func (self CommentHandler) visitNode (node Node, metadata Metadata){
	if node._type == Syntax["Program"] && len(node.body) > 0 {
		return
	}

	self.insertInnerComments(node, metadata)

	trailingComments := self.findTrailingComments(metadata)
	leadingComments := self.findLeadingComments(metadata)

	if len(leadingComments) > 0 {
		node.leadingComments = leadingComments
	}

	if len(trailingComments) > 0 {
		node.trailingComments = trailingComments
	}

	self.stack = append(self.stack, &NodeInfo{
		&node, 
		metadata.start.offset,
	})
}
































