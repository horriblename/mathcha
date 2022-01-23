// The Editor provides the interface for navigating/editing the math equation,
// that an app can then bind controls to such methods
package mathrender

import (
	tea "github.com/charmbracelet/bubbletea"
	parser "github.com/horriblename/latex-parser/latex"
	"unicode"
)

type Direction int

const (
	DIR_LEFT Direction = iota
	DIR_RIGHT
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
func (c *Cursor) Content() string       { return "" }

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
		e.renderer.Load(&parser.UnboundCompExpr{})
	}

	formatLatexTree(e.renderer.LatexTree)
	e.renderer.LatexTree.AppendChildren(e.cursor)
	e.traceStack = []parser.Container{e.renderer.LatexTree}

	e.renderer.Sync(e.getLastOnStack())
}

func (e *Editor) popStack() parser.Container {
	if len(e.traceStack) <= 1 {
		// TODO warn?
		panic("Attempting to pop Editor.traceStack when only 1 or less element is remaining")
	}
	ret := e.traceStack[len(e.traceStack)-1]
	e.traceStack[len(e.traceStack)-1] = nil
	e.traceStack = e.traceStack[:len(e.traceStack)-1]
	return ret
}

// A convenience function to move cursor to a new position and handle clean ups
// TODO remove if only stepOver*Sibling uses this
func (e *Editor) moveCursorTo(
	newParent parser.FlexContainer, // new Parent to place cursor in
	pos int, // position of the new cursor in the Parent
	ancestors []parser.Container, // A slice of Containers to insert into Editor.traceStack, the first item must be the common ancestor found in both the old and new stack
) {
	idx := e.getCursorIdxInParent()
	e.getParent().DeleteChildren(idx, idx)
	newParent.InsertChildren(pos, e.cursor)

	for i := len(e.traceStack) - 1; i >= 0; i-- {
		if e.traceStack[i] == ancestors[0] {
			e.traceStack = append(e.traceStack[:i+1], ancestors[1:]...)
		}
	}
}

// ----------------------------------------------------------------------------
// Cursor Navigation

// Navigates cursor to the left, exits a parent container if there is no left sibling
// Enters a Container if left sibling is one that allows entering
func (e *Editor) NavigateLeft() {
	idx := e.getCursorIdxInParent()
	if idx < 0 {
		panic("cursor not found in parent")
	}
	if idx == 0 {
		if len(e.traceStack) <= 1 {
			return
		}
		e.exitParent(DIR_LEFT)
	} else if prev, ok := e.getParent().Children()[idx-1].(parser.Container); ok {
		e.getParent().DeleteChildren(idx, idx)
		e.enterContainerFromRight(prev)
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
		if len(e.traceStack) <= 1 {
			return
		}
		e.exitParent(DIR_RIGHT)
	} else if next, ok := e.getParent().Children()[idx+1].(parser.Container); ok {
		e.getParent().DeleteChildren(idx, idx)
		e.enterContainerFromLeft(next)
	} else {
		e.stepOverNextSibling()
	}
}

// Convenience function for exiting a Container to the left or right
// Exits a FlexContainer plus a lingering FixedContainer, if any.
// Cursor cleanup is taken care of
// Make sure we're not exiting from the root node, no error handling for this
func (e *Editor) exitParent(direction Direction) {
	if direction != DIR_LEFT && direction != DIR_RIGHT {
		panic("exitParent got non left/right arguement")
	}
	idx := e.getCursorIdxInParent() // for loop here, might be better to pass as arguement?
	e.getParent().DeleteChildren(idx, idx)
	var exitFrom parser.Container
	exitFrom = e.popStack()
	if _, ok := e.getLastOnStack().(parser.FixedContainer); ok {
		exitFrom = e.popStack()
	}
	for i, c := range e.getParent().Children() {
		if c == exitFrom {
			idx = i
			break
		}
	}

	if direction == DIR_RIGHT {
		idx++
	}
	e.getParent().InsertChildren(idx, e.cursor)
}

// Navigate cursor downwards
func (e *Editor) NavigateDown() {
	//cursorLoc := e.getParent()
	var targetContainer parser.Container
	var cursorIdx int // index, in targetContainer.Children(), of the child containing the cursor
	stackIdx := e.findEnclosingVerticallyNavigableCommand(len(e.traceStack) - 1)
LookForContainerBelow:
	for ; stackIdx > 0; stackIdx = e.findEnclosingVerticallyNavigableCommand(stackIdx - 1) {

		targetContainer = e.traceStack[stackIdx]
		cursorLoc := e.traceStack[stackIdx+1] // in current implementation all vertical navigable containers have FlexContainer as Children

		for j := 0; j < len(targetContainer.Children())-1; j++ {
			if targetContainer.Children()[j] == cursorLoc {
				cursorIdx = j
				break LookForContainerBelow
			}
		}
	}

	if stackIdx <= 0 {
		return
	}

	if n, ok := targetContainer.Children()[cursorIdx+1].(parser.Container); ok {
		idx := e.getCursorIdxInParent()
		// Clean up traceStack
		// move to dedicated function?
		e.getParent().DeleteChildren(idx, idx)
		for i := stackIdx + 1; i < len(e.traceStack); i++ {
			e.traceStack[i] = nil
		}
		e.traceStack = e.traceStack[:stackIdx+1]

		e.enterContainerFromLeft(n)
	} else {
		panic("NavigateDown: next row does not seem to be a Container type")
	}
}

// Navigate cursor upwards
func (e *Editor) NavigateUp() {
	//cursorLoc := e.getParent()
	var targetContainer parser.Container
	var cursorIdx int // index, in targetContainer.Children(), of the child containing the cursor
	stackIdx := e.findEnclosingVerticallyNavigableCommand(len(e.traceStack) - 1)
LookForContainerAbove:
	for ; stackIdx > 0; stackIdx = e.findEnclosingVerticallyNavigableCommand(stackIdx - 1) {

		targetContainer = e.traceStack[stackIdx]
		cursorLoc := e.traceStack[stackIdx+1] // in current implementation all vertical navigable containers have FlexContainer as Children

		if targetContainer.Children()[0] == cursorLoc {
			continue LookForContainerAbove
		}

		for j := 1; j < len(targetContainer.Children()); j++ {
			if targetContainer.Children()[j] == cursorLoc {
				cursorIdx = j
				break LookForContainerAbove
			}
		}
	}

	if stackIdx <= 0 {
		return
	}

	if n, ok := targetContainer.Children()[cursorIdx-1].(parser.Container); ok {
		idx := e.getCursorIdxInParent()
		// Clean up traceStack
		// move to dedicated function?
		e.getParent().DeleteChildren(idx, idx)
		for i := stackIdx + 1; i < len(e.traceStack); i++ {
			e.traceStack[i] = nil
		}
		e.traceStack = e.traceStack[:stackIdx+1]

		e.enterContainerFromLeft(n)
	} else {
		panic("NavigateDown: next row does not seem to be a Container type")
	}
}

// Moves cursor to before the previous sibling
// Throws error if there is no previous Sibling
func (e *Editor) stepOverPrevSibling() {
	idx := e.getCursorIdxInParent()
	if idx == 0 {
		panic("stepOverPrevSibling(): No siblings before cursor!")
	}

	e.moveCursorTo(e.getParent(), idx-1, []parser.Container{e.getParent()})
}

// Moves cursor to after the next sibling
// Throws error if there is no next Sibling
func (e *Editor) stepOverNextSibling() {
	idx := e.getCursorIdxInParent()
	if len(e.getParent().Children()) <= idx+1 {
		panic("stepOverNextSibling(): No siblings after cursor!")
	}

	e.moveCursorTo(e.getParent(), idx+1, []parser.Container{e.getParent()})
}

// enter a Container from the right
// The old cursor and traceStack MUST be handled before calling this
// TODO rewrite using the moveCursorTo function?
func (e *Editor) enterContainerFromRight(target parser.Container) {
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
	parent.AppendChildren(e.cursor) // TODO use e.moveCursorTo instead?
	e.traceStack = append(e.traceStack, parent)
}

// enter a Container from the left
// The old cursor and traceStack MUST be handled before calling this
// TODO rewrite using the moveCursorTo function?
func (e *Editor) enterContainerFromLeft(target parser.Container) {
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

	parent.InsertChildren(0, e.cursor)
	e.traceStack = append(e.traceStack, parent)
}

func (e *Editor) InsertFrac(detectNumerator bool) {
	arg1 := new(parser.CompositeExpr)
	arg2 := new(parser.CompositeExpr)
	frac := &parser.Cmd2ArgExpr{Type: parser.CMD_frac, Arg1: arg1, Arg2: arg2}

	idx := e.getCursorIdxInParent()

	if detectNumerator {
		i := idx - 1
		for ; i >= 0; i-- {
			sibling := e.getParent().Children()[i]
			switch sibling.(type) {
			case *parser.VarLit, *parser.NumberLit:
				arg1.InsertChildren(0, sibling)
				continue
			default: // maybe use named loop and break from here
			}
			break
		}

		if i != idx-1 {
			e.getParent().DeleteChildren(i+1, idx-1)
			idx = i + 1
		}
	}

	e.getParent().InsertChildren(idx, frac)
}

func (e *Editor) DeleteBack() {
	idx := e.getCursorIdxInParent()
	if idx == 0 {
		// TODO exit container
		return
	}
	if n, ok := e.getParent().Children()[idx-1].(parser.Container); ok {
		e.getParent().DeleteChildren(idx, idx)
		e.enterContainerFromRight(n)
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
// Keyboard input handlers
func (e *Editor) handleLetter(letter rune) {
	idx := e.getCursorIdxInParent()
	switch parent := e.getParent().(type) {
	case parser.FlexContainer:
		parent.InsertChildren(idx, &parser.VarLit{Source: string(letter)})
	}
}

func (e *Editor) handleDigit(digit rune) {
	idx := e.getCursorIdxInParent()
	switch parent := e.getParent().(type) {
	case parser.FlexContainer:
		parent.InsertChildren(idx, &parser.NumberLit{Source: string(digit)})
	}
}

func (e *Editor) handleRest(char rune) {
	// TODO handle special characters _, ^ etc
	idx := e.getCursorIdxInParent()
	switch char {
	case '/':
		e.InsertFrac(true)
		e.NavigateLeft()
		e.NavigateDown()
		return
	case ' ':
		e.getParent().InsertChildren(idx, &parser.SimpleCmdLit{Type: parser.CMD_SPACE, Source: `\ `})
		// case '\\':
		//    e.getParent().InsertChild(idx, &parser.IncompleteCmdLit{})
	}
	switch parent := e.getParent().(type) {
	case parser.FlexContainer:
		parent.InsertChildren(idx, &parser.SimpleOpLit{Source: string(char)})
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
	panic("Editor.getParent(): Parent is not a FlexContainer")
}

// Same as getParent() but returns a general Container type instead
func (e *Editor) getLastOnStack() parser.Container {
	return e.traceStack[len(e.traceStack)-1]
}

// Find the closest enclosing FixedContainer that has vertical navigation controls and
// returns its index within the traceStack
func (e *Editor) findEnclosingVerticallyNavigableCommand(searchFrom int) (index int) {
	if searchFrom >= len(e.traceStack) {
		panic("findEnclosingVerticallyNavigableCommand(): searchFrom index >= len(traceStack)")
	}
	if searchFrom < 0 {
		searchFrom = len(e.traceStack)
	}
	for ; searchFrom > 0; searchFrom-- {
		if n, ok := e.traceStack[searchFrom].(*parser.Cmd2ArgExpr); ok {
			switch n.Command() { // TODO
			case parser.CMD_frac, parser.CMD_binom:
				return searchFrom
			}
		}
	}
	return -1
}

func (e Editor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if n, ok := e.getParent().(*parser.TextStringWrapper); ok {
		// handle key events while in a TextStringWrapper, move somewhere else
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyLeft:
				e.NavigateLeft()
			case tea.KeyRight:
				e.NavigateRight()
			case tea.KeyDown:
				e.NavigateDown()
			case tea.KeyUp:
				e.NavigateUp()
			case tea.KeyBackspace:
				e.DeleteBack()
			case tea.KeyTab:
				e.exitParent(DIR_RIGHT)

			case tea.KeyCtrlC:
				return e, tea.Quit
			case tea.KeyRunes:
				switch {
				case unicode.IsLetter(msg.Runes[0]), unicode.IsDigit(msg.Runes[0]):
					idx := e.getCursorIdxInParent()
					n.InsertChildren(idx, parser.RawRuneLit(msg.Runes[0]))
				}
			}
		}
		e.renderer.Sync(e.getLastOnStack())
		return e, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			e.NavigateLeft()
		case tea.KeyRight:
			e.NavigateRight()
		case tea.KeyDown:
			e.NavigateDown()
		case tea.KeyUp:
			e.NavigateUp()
		case tea.KeyBackspace:
			e.DeleteBack()
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
			}
		default:
			return e, nil
		}
	}
	e.renderer.Sync(e.getLastOnStack())
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

	switch n := tree.(type) {
	case *parser.TextContainer:
	case parser.FixedContainer:
		for i, child := range n.Children() {
			if _, ok := child.(parser.FlexContainer); !ok {
				n.SetArg(i, &parser.CompositeExpr{Elts: []parser.Expr{child}})
			}
		}

	case parser.FlexContainer:
		for _, child := range n.Children() {
			formatLatexTree(child)
		}
	}
}
