package latex

import "strings"

/* ----------------------------------------------------------------------------
   Interface

   There are 3 main classes of nodes: Expression nodes, ...?
   nodes, and declaration nodes. The node names usually match the corresponding
   latex command. The node fields correspond to the individual parts of the
   respective commmands

   All nodes contain position information marking the beginning of the
   corresponding source text segment; it is accessible via the Pos accessor
   method.

   All node types implement the Node interface
*/
// 3 Basic node type/interfaces: Expr, Container and Literal

// TODO after completing error handling component, see if we still need Pos()
// and End(), they're not implemented yet
type Node interface {
	Pos() Pos
	End() Pos
	VisualizeTree() string
}

type Expr interface {
	Node
	//exprNode()
}

// A branch in the syntax tree
type Container interface {
	Node
	Children() []Expr
}

// A leaf in the syntax tree
type Literal interface {
	Node
	Content() string
}

// ----------------------------------------------------------------------------
// more interfaces...
// these interfaces take after Container or Literal and includes some other
// functions of their own. Each expression node struct defined later should
// implement one of the following interfaces, or Literal TODO poor desc.

// Referring to containers that have indefinite amount of children
type FlexContainer interface {
	Container
	AppendChild(Expr)
	DeleteChildren(int, int)
	InsertChild(int, Expr)
	// TODO
	Identifier() string // temporary solution to identify the concrete type
}

// Containers with a fixed number of children.
// Used for commands that take n arguements e.g. \frac takes 2 arguements
// and same goes for superscript ^ and subscript _, which take 1 arguement
type FixedContainer interface {
	Container
	Parameters() int // number of children
	SetArg(int, Expr)
}

type CmdLiteral interface {
	Literal // the only reason this is here is to identify UnknownCmdLit via Content()
	Command() LatexCmd
}

type CmdContainer interface {
	FixedContainer
	Command() LatexCmd
}

// ---
type Comment struct {
	Percent Pos    // position of "%" starting in a Comment
	Text    string // commented text (excluding \n)
}

func (c *Comment) Pos() Pos { return c.Percent }
func (c *Comment) End() Pos { return Pos(int(c.Percent) + len(c.Text)) }

type TextLit struct {
	CmdText  Pos    // position of "\text"
	From, To Pos    // position of "{" / "}" or position of single-character
	Content  string // contained
}

/* ----------------------------------------------------------------------------
   Expresions

   An expression is represented by a tree consisting of one or more of the
   folowing concrete expression nodes.
*/
type (
	// A BadExpr node is a placeholder for an expression containing syntax errors for which a correct expression node cannot be created.
	BadExpr struct {
		From, To Pos
		source   string
	}

	// A EmptyExpr node is a placeholder to mark the termination of a previous expression
	EmptyExpr struct {
		From, To Pos
		Type     Token
	}

	// A NumberLit node represents a literal consisting of digits
	NumberLit struct {
		From, To Pos
		Source   string // literal string; e.g. 23x
	}

	// A VarLit node represents a literal consisting of alphabets
	VarLit struct {
		From, To Pos
		Source   string
	}

	// A Composite node represents a composite bracea surrounded { expression }
	CompositeExpr struct {
		Type       Expr   // literal type; or nil ?
		Lbrace     Pos    // Position of "{"
		Elts       []Expr // list of composite elements; or nil
		Rbrace     Pos    // position of "}"
		Incomplete bool   // true if (source) expressions are missing in Elts
	}

	// A UnboundCompExpr is basically the same as CompositeLit but without brackets "{}"
	UnboundCompExpr struct {
		From, To Pos
		Elts     []Expr // list of composite elements; or nil
	}

	// A SimpleOpLit node represents a simple operator literal
	SimpleOpLit struct {
		From, To Pos
		Source   string // e.g. + - =
	}

	// IncompleteCmdLit node is a placeholder for an incomplete command
	// It is treated as a SimpleCmdLit, without any special grammar
	// TODO it should be treated as some kind of TextLit
	IncompleteCmdLit struct {
		Backslash Pos    // Position of "\"
		Source    string // he command string including backslash
		To        Pos    // position of the last character
	}

	// UnknownCmdLit node is a placeholder for an unrecognized command
	// It is treated as a SimpleCmdLit, without any special grammar
	UnknownCmdLit struct {
		Backslash Pos    // Position of "\"
		Source    string // he command string including backslash
		To        Pos    // position of the last character
	}

	// A SimpleCmdLit node is a simple command that behaves like any other simple literal e.g. \times
	SimpleCmdLit struct {
		Backslash Pos    // Position of "\"
		Source    string // the command string including backslash
		Type      LatexCmd
		To        Pos // position of last character
	}

	// A SuperExpr node represents a superscript expression
	// TODO generalize
	SuperExpr struct {
		Symbol Pos  // position of "^"
		X      Expr // superscripted expression
		Close  Pos  // position of "}" if its a composite expression, otherwise the character
	}

	// A SubExpr node represents a subscript expression
	SubExpr struct {
		Symbol Pos  // position of "_"
		X      Expr // superscripted expression
		Close  Pos  // position of "}" if its a composite expression, otherwise the character
	}

	// TODO maybe remove variable source if not used
	// A Cmd1ArgExpr node represents a command that takes 1 arguement e.g. \underline
	Cmd1ArgExpr struct {
		source    string // remove ?
		Type      LatexCmd
		Backslash Pos // position of "\"
		Arg1      Expr
		To        Pos
	}

	Cmd2ArgExpr struct {
		source    string // remove?
		Type      LatexCmd
		Backslash Pos // position of "\"
		Arg1      Expr
		Arg2      Expr
		To        Pos
	}
)

func (x *BadExpr) Pos() Pos         { return x.From }
func (x *EmptyExpr) Pos() Pos       { return x.From }
func (x *NumberLit) Pos() Pos       { return x.From }
func (x *VarLit) Pos() Pos          { return x.From }
func (x *CompositeExpr) Pos() Pos   { return x.Lbrace }
func (x *UnboundCompExpr) Pos() Pos { return x.From }
func (x *SimpleOpLit) Pos() Pos     { return x.From }
func (x *UnknownCmdLit) Pos() Pos   { return x.Backslash }
func (x *SimpleCmdLit) Pos() Pos    { return x.Backslash }
func (x *SuperExpr) Pos() Pos       { return x.Symbol }
func (x *SubExpr) Pos() Pos         { return x.Symbol }
func (x *Cmd1ArgExpr) Pos() Pos     { return x.Backslash }
func (x *Cmd2ArgExpr) Pos() Pos     { return x.Backslash }

func (x *BadExpr) End() Pos         { return x.To }
func (x *EmptyExpr) End() Pos       { return x.To }
func (x *NumberLit) End() Pos       { return x.To }
func (x *VarLit) End() Pos          { return x.To }
func (x *CompositeExpr) End() Pos   { return x.Lbrace }
func (x *UnboundCompExpr) End() Pos { return x.To }
func (x *SimpleOpLit) End() Pos     { return x.From }
func (x *UnknownCmdLit) End() Pos   { return x.To }
func (x *SimpleCmdLit) End() Pos    { return x.To }
func (x *SuperExpr) End() Pos       { return x.Close }
func (x *SubExpr) End() Pos         { return x.Close }
func (x *Cmd1ArgExpr) End() Pos     { return x.To }
func (x *Cmd2ArgExpr) End() Pos     { return x.To }

// Container method definitions
func (x *CompositeExpr) Children() []Expr   { return x.Elts }
func (x *UnboundCompExpr) Children() []Expr { return x.Elts }

func (x *Cmd1ArgExpr) Children() []Expr { return []Expr{x.Arg1} }
func (x *Cmd2ArgExpr) Children() []Expr { return []Expr{x.Arg1, x.Arg2} }

// FlexContainer methods
func (x *CompositeExpr) AppendChild(child Expr)   { x.Elts = append(x.Elts, child) }
func (x *UnboundCompExpr) AppendChild(child Expr) { x.Elts = append(x.Elts, child) }

// Deletes Children from index, to index, inclusive
func (x *CompositeExpr) DeleteChildren(from int, to int) {
	if from < 0 || to >= len(x.Children()) {
		panic("DeleteChildren(): index out of range!")
	}
	if from > to {
		panic("DeleteChildren(): 'from' cannot be larger than 'to'")
	}
	l := to - from + 1
	copy(x.Elts[from:], x.Elts[to+1:])
	for i := range x.Elts[len(x.Elts)-l:] {
		x.Elts[len(x.Elts)-l+i] = nil // garbage collection
	}
	x.Elts = x.Elts[:len(x.Elts)-l]
}
func (x *UnboundCompExpr) DeleteChildren(from int, to int) {
	if from < 0 || to >= len(x.Children()) {
		panic("DeleteChildren(): index out of range!")
	}
	if from > to {
		panic("DeleteChildren(): 'from' cannot be larger than 'to'")
	}
	l := to - from + 1
	copy(x.Elts[from:], x.Elts[to+1:])
	for i := range x.Elts[len(x.Elts)-l:] {
		x.Elts[len(x.Elts)-l+i] = nil // garbage collection
	}
	x.Elts = x.Elts[:len(x.Elts)-l]
	// copy(x.Elts[i:], x.Elts[i+1:])
	// x.Elts[len(x.Elts)-1] = nil // or the zero value of T
	// x.Elts = x.Elts[:len(x.Elts)-1]
}

// FIXME do I really need to repeat these exact same functions for each struct??
// Insert child at index; the new child has the index 'at'
func (x *CompositeExpr) InsertChild(at int, child Expr) {
	if at < 0 || at > len(x.Children()) {
		panic("InsertChild(): invalid index for 'at'")
	}
	if at == len(x.Children()) {
		x.AppendChild(child)
		return
	}
	x.Elts = append(x.Elts[:at+1], x.Elts[at:]...)
	x.Elts[at] = child
}
func (x *UnboundCompExpr) InsertChild(at int, child Expr) {
	if at < 0 || at > len(x.Children()) {
		panic("InsertChild(): invalid index for 'at'")
	}
	if at == len(x.Children()) {
		x.AppendChild(child)
		return
	}
	x.Elts = append(x.Elts[:at+1], x.Elts[at:]...)
	x.Elts[at] = child
}

func (x *CompositeExpr) Identifier() string   { return "{" }
func (x *UnboundCompExpr) Identifier() string { return "" }

// FixedContainer methods
func (x *Cmd1ArgExpr) Parameters() int { return 1 }
func (x *Cmd2ArgExpr) Parameters() int { return 2 }

func (x *Cmd1ArgExpr) SetArg(index int, expr Expr) {
	if index > 0 {
		panic("SetArg(): index out of range")
	}
	x.Arg1 = expr
}
func (x *Cmd2ArgExpr) SetArg(index int, expr Expr) {
	if index > 1 {
		panic("SetArg(): index out of range")
	}
	switch index {
	case 0:
		x.Arg1 = expr
	case 1:
		x.Arg2 = expr
	}
}

// Literal method definitions
func (x *BadExpr) Content() string       { return x.source }
func (x *EmptyExpr) Content() string     { return "" }
func (x *NumberLit) Content() string     { return x.Source }
func (x *VarLit) Content() string        { return x.Source }
func (x *SimpleOpLit) Content() string   { return x.Source }
func (x *UnknownCmdLit) Content() string { return x.Source }
func (x *SimpleCmdLit) Content() string  { return x.Source }

// CmdLiteral, CmdContainer method definitions
func (x *UnknownCmdLit) Command() LatexCmd { return CMD_UNKNOWN }
func (x *SimpleCmdLit) Command() LatexCmd  { return x.Type }
func (x *Cmd1ArgExpr) Command() LatexCmd   { return x.Type }
func (x *Cmd2ArgExpr) Command() LatexCmd   { return x.Type }

//

// ----------------------------------------------------------------------------
// VisualizeTree, naive approach, only for debugging purposes
func (x *UnboundCompExpr) VisualizeTree() string {
	tree := "$\n"
	for _, el := range x.Children() {
		branch := (el).VisualizeTree()
		splits := strings.Split(branch, "\n")
		tree += "├───" + splits[0] + "\n"
		if len(splits) == 1 {
			continue
		}
		for _, line := range splits[1:] {
			tree += "|   " + line + "\n"
		}
	}
	return tree
}

// VisualizeTree, naive approach, only for debugging purposes
func (x *CompositeExpr) VisualizeTree() string {
	tree := "{\n"
	for _, el := range x.Children() {
		branch := (el).VisualizeTree()
		splits := strings.Split(branch, "\n")
		tree += "├───" + splits[0] + "\n"
		if len(splits) == 1 {
			continue
		}
		for _, line := range splits[1:] {
			tree += "|   " + line + "\n"
		}
	}
	tree += "└───}\n"
	return tree
}

func (x *SuperExpr) VisualizeTree() string {
	tree := "^\n"
	branch := x.X.VisualizeTree()
	splits := strings.Split(branch, "\n")
	tree += "├───" + splits[0] + "\n"
	if len(splits) >= 1 {
		for _, line := range splits[1:] {
			tree += "|   " + line + "\n"
		}
	}
	return tree
}

func (x *SubExpr) VisualizeTree() string {
	tree := "_\n"
	branch := x.X.VisualizeTree()
	splits := strings.Split(branch, "\n")
	tree += "├───" + splits[0] + "\n"
	if len(splits) >= 1 {
		for _, line := range splits[1:] {
			tree += "|   " + line + "\n"
		}
	}
	return tree
}

func (x *Cmd1ArgExpr) VisualizeTree() string {
	tree := x.source + "\n"
	branch := x.Arg1.VisualizeTree()
	splits := strings.Split(branch, "\n")
	tree += "├───" + splits[0] + "\n"
	if len(splits) >= 1 {
		for _, line := range splits[1:] {
			tree += "|   " + line + "\n"
		}
	}
	return tree
}

func (x *Cmd2ArgExpr) VisualizeTree() string {
	tree := "\\frac\n"
	branch := x.Arg1.VisualizeTree()
	splits := strings.Split(branch, "\n")
	tree += "├───" + splits[0] + "\n"
	if len(splits) >= 1 {
		for _, line := range splits[1:] {
			tree += "|   " + line + "\n"
		}
	}

	branch = x.Arg2.VisualizeTree()
	splits = strings.Split(branch, "\n")
	tree += "├───" + splits[0] + "\n"
	if len(splits) >= 1 {
		for _, line := range splits[1:] {
			tree += "|   " + line + "\n"
		}
	}

	return tree
}

func (x *BadExpr) VisualizeTree() string       { return "BadExpr" }
func (x *EmptyExpr) VisualizeTree() string     { return "EmptyExpr" }
func (x *NumberLit) VisualizeTree() string     { return "NumberLit     " + x.Source }
func (x *VarLit) VisualizeTree() string        { return "VarLit        " + x.Source }
func (x *SimpleOpLit) VisualizeTree() string   { return "SimpleOpLit   " + x.Source }
func (x *SimpleCmdLit) VisualizeTree() string  { return "SimpleCmdLit  " + x.Source }
func (x *UnknownCmdLit) VisualizeTree() string { return "UnknownCmdLit " + x.Source }
