// The Editor provides the interface for navigating/editing the math equation,
// that an app can then bind controls to such methods
package mathrender

import (
	tea "github.com/charmbracelet/bubbletea"
	parser "github.com/horriblename/latex-parser/latex"
	"unicode"
)

// the Cursor object implements a zero-width parser.Literal TODO rn its still a normal character
type Cursor struct {
	offsetX int // offset the position of the cursor
}

//
type Editor struct {
	renderer   *Renderer
	traceStack []parser.Container // trace our position on the tree
	cursor     *Cursor
}

func (c *Cursor) Pos() parser.Pos { return parser.Pos(0) } // FIXME remove
func (c *Cursor) End() parser.Pos { return parser.Pos(0) }

func (c *Cursor) VisualizeTree() string { return "Cursor        " }
func (c *Cursor) Content() string       { return "\x1b[47mI\x1b[40m" } // TODO implement zero-width cursor

// TODO remove?
func (e Editor) Init() tea.Cmd {
	e.cursor = new(Cursor)
	return nil
}

func (e *Editor) Read(latex string) {
	// load latex input
	if latex != "" {
		p := parser.Parser{}
		p.Init(latex)
		// e.renderer.Load(p.GetTree()) // FIXME why doesn't this work
		e.renderer = &Renderer{LatexTree: p.GetTree()}
		// p (Parser object) can be discarded now
	} else {
		e.renderer.Load(&parser.TopLevelExpr{})
	}

	formatLatexTree(e.renderer.LatexTree)
	e.renderer.LatexTree.AppendChild(e.cursor)
	e.traceStack = []parser.Container{e.renderer.LatexTree}

	e.renderer.Sync()
}

// ----------------------------------------------------------------------------
// Navigation utils

// Navigates cursor to the left, exits a parent container if there is no left sibling
// Enters a Container if left sibling is one that allows entering
func (e *Editor) NavigateLeft() {
	idx := e.getCursorIdxInParent()
	if idx < 0 {
		panic("cursor not found in parent")
	}
	if idx == 0 {
		// exit parent
	} else if prev, ok := e.getParent().Children()[idx-1].(parser.Container); ok {
		e.getParent().DeleteChildren(idx, idx)
		e.EnterContainerFromRight(prev)
	} else {
		e.stepOverPrevSibling()
	}
}

// Navigates cursor to the right, exits a parent container if there is no right sibling
// Enters a Container if right sibling is one that allows entering
func (e *Editor) NavigateRight() {
	idx := e.getCursorIdxInParent()
	if idx < 0 {
		panic("cursor not found in parent")
	}
	if idx+1 >= len(e.getParent().Children()) {
		// exit parent
	} else if next, ok := e.getParent().Children()[idx+1].(parser.Container); ok {
		e.getParent().DeleteChildren(idx, idx)
		e.EnterContainerFromLeft(next)
	} else {
		e.stepOverNextSibling()
	}
}

// Moves cursor to before the previous sibling
// Throws error if there is no previous Sibling
func (e *Editor) stepOverPrevSibling() {
	idx := e.getCursorIdxInParent()
	if idx == 0 {
		panic("stepOverPrevSibling(): No siblings before cursor!")
	}

	if n, ok := e.getParent().(parser.FlexContainer); ok {
		n.DeleteChildren(idx, idx)
		n.InsertChild(idx-1, e.cursor)
	} else {
		panic("stepOverPrevSibling(): cursor not in a FlexContainer")
	}
}

// Moves cursor to after the next sibling
// Throws error if there is no next Sibling
func (e *Editor) stepOverNextSibling() {
	idx := e.getCursorIdxInParent()
	if len(e.getParent().Children()) <= idx+1 {
		panic("stepOverNextSibling(): No siblings after cursor!")
	}

	if n, ok := e.getParent().(parser.FlexContainer); ok {
		n.DeleteChildren(idx, idx)
		n.InsertChild(idx+1, e.cursor)
	} else {
		panic("stepOverNextSibling(): cursor not in a FlexContainer")
	}
}

func (e *Editor) EnterContainerFromRight(target parser.Container) {
	var parent parser.FlexContainer
	switch t := target.(type) {
	case parser.FixedContainer:
		e.traceStack = append(e.traceStack, t)
		if m, ok := t.Children()[0].(parser.FlexContainer); ok { // TODO pick different 'children Container' based on command type?
			parent = m
		} else {
			panic("Editor attempted to enter a FixedContainer type with a non FlexContainer as first Child")
		}
	case parser.FlexContainer:
		parent = t
	default:
		panic("Editor attempted to enter a non Fixed- or FlexContainer")
	}
	parent.AppendChild(e.cursor)
	e.traceStack = append(e.traceStack, parent)
}

func (e *Editor) EnterContainerFromLeft(target parser.Container) {
	var parent parser.FlexContainer
	switch t := target.(type) {
	case parser.FixedContainer:
		e.traceStack = append(e.traceStack, t)
		if m, ok := t.Children()[0].(parser.FlexContainer); ok { // TODO pick different 'children Container' based on command type?
			parent = m
		} else {
			panic("Editor attempted to enter a FixedContainer type with a non FlexContainer as first Child")
		}
	case parser.FlexContainer:
		parent = t
	default:
		panic("Editor attempted to enter a non Fixed- or FlexContainer")
	}

	parent.InsertChild(0, e.cursor)
	e.traceStack = append(e.traceStack, parent)
}

func (e *Editor) DeleteBack() {
	idx := e.getCursorIdxInParent()
	if idx == 0 {
		// TODO exit container
		return
	}
	e.getParent().DeleteChildren(idx-1, idx-1)
}

// returns the cursor position relative to the parent
func (e *Editor) getCursorIdxInParent() int {
	for i, n := range e.getParent().Children() {
		if n == e.cursor {
			return i
		}
	}
	return -1
}

// ---
// Key input handlers
func (e *Editor) handleLetter(letter rune) {
	idx := e.getCursorIdxInParent()
	switch parent := e.getParent().(type) {
	case parser.FlexContainer:
		parent.InsertChild(idx, &parser.VarLit{Source: string(letter)})
	}
}

func (e *Editor) handleDigit(digit rune) {
	idx := e.getCursorIdxInParent()
	switch parent := e.getParent().(type) {
	case parser.FlexContainer:
		parent.InsertChild(idx, &parser.NumberLit{Source: string(digit)})
	}
}

func (e *Editor) handleRest(char rune) {
	// TODO handle special characters _, ^ etc
	idx := e.getCursorIdxInParent()
	switch parent := e.getParent().(type) {
	case parser.FlexContainer:
		parent.InsertChild(idx, &parser.SimpleOpLit{Source: string(char)})
	}
}

// ---
// utilities

// getParent() returns the last Container on the traceStack, if it's not a FlexContainer,
// the function will panic
func (e *Editor) getParent() parser.FlexContainer {
	// TODO error handler for len(e.traceStack) == 0?
	if n, ok := e.traceStack[len(e.traceStack)-1].(parser.FlexContainer); ok {
		// TODO
		return n
	}
	panic("Parent is not a FlexContainer")
}

func (e Editor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			e.NavigateLeft()
			e.renderer.Sync()
		case tea.KeyRight:
			e.NavigateRight()
			e.renderer.Sync()
		case tea.KeyBackspace:
			e.DeleteBack()
			e.renderer.Sync()
		case tea.KeyCtrlC:
			return e, tea.Quit
		case tea.KeyRunes:
			if len(msg.Runes) == 1 {
				// TODO add if in text block add to text
				switch {
				case unicode.IsLetter(msg.Runes[0]):
					e.handleLetter(msg.Runes[0])
				case unicode.IsDigit(msg.Runes[0]):
					e.handleDigit(msg.Runes[0])
				default:
					e.handleRest(msg.Runes[0])
				}
				e.renderer.Sync()
			}
		}
	}
	return e, nil
}

func (e Editor) View() string {
	return e.renderer.View() + "\n" + ProduceLatex(e.renderer.LatexTree)
}

// Search the tree for any FixedContainer type that has children that is
// not surrounded in "{}" (CompositeExpr), then wraps them in a CompositeExpr.
// This is to make the editing part easier
func formatLatexTree(tree parser.Expr) {
	// TODO also convert (), [], \{\} to \left...\right

	if n, ok := tree.(parser.FixedContainer); ok {
		for i, child := range n.Children() {
			if _, ok := child.(parser.FlexContainer); !ok {
				n.SetArg(i, &parser.CompositeExpr{Elts: []parser.Expr{child}})
			}
		}
	}

	if n, ok := tree.(parser.Container); ok {
		for _, child := range n.Children() {
			formatLatexTree(child)
		}
	}
}
