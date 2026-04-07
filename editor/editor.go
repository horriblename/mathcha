// The Editor provides the interface for navigating/editing the math equation,
// that an app can then bind controls to such methods
package editor

import (
	"log"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	parser "github.com/horriblename/mathcha/latex"
	render "github.com/horriblename/mathcha/renderer"
)

type Direction int
type editorStates int // TODO remove?
type editorState int

const (
	DIR_LEFT Direction = iota
	DIR_RIGHT

	// TODO rename/use container type?
	EDIT_EQUATION editorState = iota // equation editing mode (normal)
	EDIT_TEXT
	EDIT_COMMAND
)

type Editor struct {
	renderer   *render.Renderer
	traceStack []parser.Container // trace our position on the tree
	cursor     *render.Cursor
	markSelect *render.Cursor
	focus      bool
	config     *EditorConfig
	banner     string // a line of text appearing below the renderer, for debugging
}

type EditorConfig struct {
	*log.Logger
	LatexCfg render.LatexSourceConfig
}

func New(formula string) *Editor {
	// TODO: detect color from tty
	renderer := render.FromFormula(formula, true)
	cursor := render.Cursor{Symbol: "\x1b[7m \x1b[27m"}
	renderer.LatexTree.AppendChildren(&cursor)
	return &Editor{
		renderer:   renderer,
		traceStack: []parser.Container{renderer.LatexTree},
		cursor:     &cursor,
		markSelect: nil,
		focus:      false,
		config: &EditorConfig{
			LatexCfg: render.LatexSourceConfig{
				UseUnicode: true,
			},
		},
	}
}

func NewWithConfig(cfg EditorConfig, formula string) *Editor {
	editor := New(formula)
	editor.config = &cfg
	return editor
}

func (e *Editor) Read(latex string) {
	// load latex input
	if latex != "" {
		ast := parser.Parse(latex)
		// e.renderer.Load(p.GetTree()) // FIXME why doesn't this work
		e.renderer = &render.Renderer{LatexTree: ast}
		// p (Parser object) can be discarded now
	} else {
		e.renderer.Load(&parser.UnboundCompExpr{})
	}

	formatLatexTree(e.renderer.LatexTree)
	e.renderer.LatexTree.AppendChildren(e.cursor)
	e.traceStack = []parser.Container{e.renderer.LatexTree}

	e.renderer.Sync(e.getLastOnStack(), false)
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

func (e *Editor) SetFocus(f bool) {
	if e.focus == f {
		return
	}
	e.focus = f
	e.renderer.Focus = f
	e.renderer.Sync(e.getLastOnStack(), false)
}

// gets the 'state' of the editor, e.g. inserting a command/text node, or in normal equation node
func (e *Editor) GetState() editorState {
	var state editorState
	switch e.getParent().(type) {
	case *parser.TextStringWrapper:
		if _, ok := e.traceStack[len(e.traceStack)-2].(*render.LatexCmdInput); ok {
			state = EDIT_COMMAND
		} else {
			state = EDIT_TEXT
		}
	case parser.FlexContainer:
		state = EDIT_EQUATION
	default:
		// TODO
		panic("getState could not match parent to an expected type")
	}

	return state
}

func (e *Editor) hasSelection() bool {
	return e.markSelect != nil
}

func (e *Editor) cancelSelection() {
	markIdx := e.getSelectionIdxInParent()
	e.getParent().DeleteChildren(markIdx, markIdx)
	e.markSelect = nil
}

func (e *Editor) deleteSelection() {
	idx := e.getCursorIdxInParent()
	mark := e.getSelectionIdxInParent()
	if mark < idx {
		e.getParent().DeleteChildren(mark, idx-1)
	} else {
		e.getParent().DeleteChildren(idx+1, mark)
	}
	e.markSelect = nil
}

func (e *Editor) deleteToStart() {
	if e.markSelect != nil {
		e.deleteSelection()
	}

	idx := e.getCursorIdxInParent()
	if idx == 0 {
		return
	}

	e.getParent().DeleteChildren(0, idx-1)
}

func (e *Editor) deleteToEnd() {
	if e.markSelect != nil {
		e.deleteSelection()
	}

	idx := e.getCursorIdxInParent()
	lastIdx := len(e.getParent().Children()) - 1
	if idx == lastIdx {
		return
	}

	e.getParent().DeleteChildren(idx+1, lastIdx)
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
// Does nothing and returns false if we're the first child of the root node
func (e *Editor) NavigateLeft() bool {
	idx := e.getCursorIdxInParent()
	if idx < 0 {
		return false
	}
	if idx == 0 {
		if len(e.traceStack) <= 1 {
			return true
		} else if 2 <= len(e.traceStack) {
			if env, ok := e.traceStack[len(e.traceStack)-2].(*parser.EnvExpr); ok {
				row, col := env.FindCell(e.getParent().(*parser.UnboundCompExpr))
				var targetCell *parser.UnboundCompExpr
				if col > 0 {
					targetCell = env.Elts[row][col-1]
				} else if row > 0 {
					targetCell = env.Elts[row-1][len(env.Elts[row-1])-1]
				} else {
					e.exitParent(DIR_LEFT)
					return true
				}
				idx := e.getCursorIdxInParent()
				e.getParent().DeleteChildren(idx, idx)
				targetCell.AppendChildren(e.cursor)
				e.traceStack[len(e.traceStack)-1] = targetCell
				return true
			}
		}
		e.exitParent(DIR_LEFT)
	} else if prev, ok := e.getParent().Children()[idx-1].(parser.Container); ok {
		e.getParent().DeleteChildren(idx, idx)
		e.enterContainerFromRight(prev)
	} else {
		e.stepOverPrevSibling()
	}

	return true
}

// Navigates cursor to the right, exits a parent container if there is no right sibling
// Enters a Container if right sibling is one that allows entering
// Does nothing and returns false if we're the last child of the root node
func (e *Editor) NavigateRight() bool {
	idx := e.getCursorIdxInParent()
	if idx < 0 {
		return false
	}
	if idx+1 >= len(e.getParent().Children()) {
		if len(e.traceStack) <= 1 {
			return true
		} else if 2 <= len(e.traceStack) {
			if env, ok := e.traceStack[len(e.traceStack)-2].(*parser.EnvExpr); ok {
				row, col := env.FindCell(e.getParent().(*parser.UnboundCompExpr))
				var targetCell *parser.UnboundCompExpr
				if col < len(env.Elts[row])-1 {
					targetCell = env.Elts[row][col+1]
				} else if row < len(env.Elts)-1 {
					targetCell = env.Elts[row+1][0]
				} else {
					e.exitParent(DIR_RIGHT)
					return true
				}
				idx := e.getCursorIdxInParent()
				e.getParent().DeleteChildren(idx, idx)
				targetCell.InsertChildren(0, e.cursor)
				e.traceStack[len(e.traceStack)-1] = targetCell
				return true
			}
		}
		e.exitParent(DIR_RIGHT)
	} else if next, ok := e.getParent().Children()[idx+1].(parser.Container); ok {
		e.getParent().DeleteChildren(idx, idx)
		e.enterContainerFromLeft(next)
	} else {
		e.stepOverNextSibling()
	}
	return true
}

// Convenience function for exiting a Container to the left or right
// Exits a FlexContainer plus a lingering FixedContainer/EnvExpr, if any.
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
	switch e.getLastOnStack().(type) {
	case parser.FixedContainer, *parser.EnvExpr:
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

// May return nil if out of range (for last or first row).
//
// cursorLoc is the child of container that contains the cursor as a descendant.
//
// Panics if container is a [parser.FlexContainer] or [parser.Cmd1ArgExpr].
func (e *Editor) containerGetSiblingRow(container parser.Container, cursorLoc parser.Container, up bool) parser.FlexContainer {
	switch n := container.(type) {
	case *parser.Cmd2ArgExpr:
		if up && n.Arg2 == cursorLoc {
			return n.Arg1.(parser.FlexContainer)
		} else if !up && n.Arg1 == cursorLoc {
			return n.Arg2.(parser.FlexContainer)
		}
		return nil
	case *parser.EnvExpr:
		row, col := n.FindCell(cursorLoc.(*parser.UnboundCompExpr))
		if up && row > 0 {
			last := len(n.Elts[row-1]) - 1
			return n.Elts[row-1][min(col, last)]
		} else if !up && row < len(n.Elts)-1 {
			last := len(n.Elts[row+1]) - 1
			return n.Elts[row+1][min(col, last)]
		}
	}
	return nil
}

// Navigate cursor vertically (up or down)
func (e *Editor) navigateVertical(up bool) {
	var targetContainer parser.Container
	var targetRow parser.FlexContainer

	stackIdx := e.findEnclosingVerticallyNavigableCommand(len(e.traceStack) - 1)
	for ; stackIdx > 0; stackIdx = e.findEnclosingVerticallyNavigableCommand(stackIdx - 1) {
		targetContainer = e.traceStack[stackIdx]
		cursorLoc := e.traceStack[stackIdx+1]
		targetRow = e.containerGetSiblingRow(targetContainer, cursorLoc, up)
		if targetRow != nil {
			break
		}
	}

	if targetRow == nil {
		return
	}

	idx := e.getCursorIdxInParent()
	e.getParent().DeleteChildren(idx, idx)
	for i := stackIdx + 1; i < len(e.traceStack); i++ {
		e.traceStack[i] = nil
	}
	e.traceStack = e.traceStack[:stackIdx+1]

	e.enterContainerFromLeft(targetRow)
}

func (e *Editor) NavigateDown() { e.navigateVertical(false) }
func (e *Editor) NavigateUp()   { e.navigateVertical(true) }

// Moves cursor to before the previous sibling
// Returns false if no previous node
func (e *Editor) stepOverPrevSibling() bool {
	idx := e.getCursorIdxInParent()
	if idx == 0 {
		return false
	}

	e.moveCursorTo(e.getParent(), idx-1, []parser.Container{e.getParent()})
	return true
}

// Moves cursor to after the next sibling
// Returns false if no next node
func (e *Editor) stepOverNextSibling() bool {
	idx := e.getCursorIdxInParent()
	if len(e.getParent().Children()) <= idx+1 {
		return false
	}

	e.moveCursorTo(e.getParent(), idx+1, []parser.Container{e.getParent()})
	return true
}

// Start or extend selection to the left
// If at first child, do nothing and return false
func (e *Editor) selectLeft() bool {
	idx := e.getCursorIdxInParent()
	if idx == 0 {
		return false
	}

	if e.markSelect == nil {
		e.markSelect = new(render.Cursor)
		e.getParent().InsertChildren(idx+1, e.markSelect)
	}
	_ = e.stepOverPrevSibling() || e.NavigateLeft()
	return true
}

// Start or extend selection to the right
// If at last child, do nothing and return false
func (e *Editor) selectRight() bool {
	idx := e.getCursorIdxInParent()
	if idx >= len(e.getLastOnStack().Children())-1 {
		return false
	}

	if e.markSelect == nil {
		e.markSelect = new(render.Cursor)
		e.getParent().InsertChildren(idx, e.markSelect)
		idx++
	}

	_ = e.stepOverNextSibling() || e.NavigateRight()
	return true
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
	case *parser.EnvExpr:
		lastRow := len(t.Elts) - 1
		parent = t.Elts[lastRow][len(t.Elts[lastRow])-1]
		e.traceStack = append(e.traceStack, t)
	default:
		panic("Editor attempted to enter a non EnvExpr or Fixed- or FlexContainer")
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
	case *parser.EnvExpr:
		parent = t.Elts[0][0]
		e.traceStack = append(e.traceStack, t)
	default:
		panic("Editor attempted to enter a non EnvExpr or Fixed- or FlexContainer")
	}

	parent.InsertChildren(0, e.cursor)
	e.traceStack = append(e.traceStack, parent)
}

func (e *Editor) NavigateToBeginning() {
	idx := e.getCursorIdxInParent()
	parent := e.getParent()
	parent.DeleteChildren(idx, idx)
	parent.InsertChildren(0, e.cursor)
}

func (e *Editor) NavigateToEnd() {
	idx := e.getCursorIdxInParent()
	parent := e.getParent()
	parent.DeleteChildren(idx, idx)
	parent.AppendChildren(e.cursor)
}

func (e *Editor) InsertCmd(cmd string) {
	kind := parser.MatchLatexCmd(cmd)
	idx := e.getCursorIdxInParent()
	switch {
	case cmd == "\\begin":
		envNameStr := "matrix" // TODO
		envName := parser.GetEnvName(envNameStr)
		cell := &parser.UnboundCompExpr{}
		node := &parser.EnvExpr{
			Name: envName,
			Elts: [][]*parser.UnboundCompExpr{
				{cell},
			},
		}
		e.getParent().DeleteChildren(idx, idx)
		e.getParent().InsertChildren(idx, node)
		e.traceStack = append(e.traceStack, node)
		cell.AppendChildren(e.cursor)
		e.traceStack = append(e.traceStack, cell)
	case kind.IsTextCmd():
		node := &parser.TextContainer{Text: &parser.TextStringWrapper{}}
		e.getParent().DeleteChildren(idx, idx)
		e.getParent().InsertChildren(idx, node)
		e.enterContainerFromRight(node)
	case kind.TakesOneArg():
		node := &parser.Cmd1ArgExpr{Type: kind, Arg1: new(parser.CompositeExpr)}
		e.getParent().DeleteChildren(idx, idx)
		e.getParent().InsertChildren(idx, node)
		e.enterContainerFromRight(node)
	case kind.TakesTwoArg():
		node := &parser.Cmd2ArgExpr{Type: kind, Arg1: new(parser.CompositeExpr), Arg2: new(parser.CompositeExpr)}
		e.getParent().DeleteChildren(idx, idx)
		e.getParent().InsertChildren(idx, node)
		e.enterContainerFromRight(node)
	case kind.TakesRawStrArg():
	case kind.IsVanillaSym():
		node := &parser.SimpleCmdLit{Type: kind, Source: cmd}
		e.getParent().InsertChildren(idx, node)

	case kind.IsEnclosing():
	default:
		node := &parser.UnknownCmdLit{Source: cmd}
		e.getParent().InsertChildren(idx, node)
	}
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

// delete the parent container and possibly the parent's FixedContainer
// and insert the parent's children to the node one level above:
//
//	x + \frac{2y}{3} --[flatten '\frac']-> x + 2y3
//
// In an EnvExpr, newlines and column alignment `&` can be erased:
//
//	\begin{matrix} a & b & c \end{matrix}
//	--[flatten at b]->
//	\begin{matrix} ab & c \end{matrix}
func (e *Editor) flattenDeleteParent() {
	idx := e.getCursorIdxInParent()
	if forest, ok := e.traceStack[len(e.traceStack)-2].(parser.FixedContainer); ok {
		e.exitParent(DIR_LEFT)
		idx = e.getCursorIdxInParent()
		e.getParent().DeleteChildren(idx+1, idx+1)

		for i := len(forest.Children()) - 1; i >= 0; i-- {
			if c, ok := forest.Children()[i].(parser.FlexContainer); ok {
				e.getParent().InsertChildren(idx+1, c.Children()...)
			} else {
				e.getParent().InsertChildren(idx+1, c)
			}
		}
	} else if forest, ok := e.traceStack[len(e.traceStack)-2].(*parser.EnvExpr); ok {
		oldCell := e.getParent()
		for r, row := range forest.Elts {
			for c, cell := range row {
				if cell == oldCell {
					if c == 0 {
						forest.Elts[r-1] = append(forest.Elts[r-1], row...)
						for r1 := r; r < len(forest.Elts); r1++ {
							forest.Elts[r1] = forest.Elts[r1+1]
						}
						forest.Elts = forest.Elts[:len(forest.Elts)-1]
					} else {
						idx := e.getCursorIdxInParent()
						oldCell.DeleteChildren(idx, idx)
						row[c-1].AppendChildren(e.cursor)
						row[c-1].AppendChildren(row[c].Elts...)
						for col := c; col < len(row); col++ {
							row[col] = row[col+1]
						}
						forest.Elts[r] = row[:len(row)-1]
					}
					return
				}
			}
		}
	} else {
		deleting := e.getParent()
		e.exitParent(DIR_LEFT)
		idx = e.getCursorIdxInParent()
		e.getParent().DeleteChildren(idx+1, idx+1)

		e.getParent().InsertChildren(idx+1, deleting.Children()...)
	}
}

func (e *Editor) DeleteBack() {
	if e.markSelect != nil {
		e.deleteSelection()
		return
	}

	idx := e.getCursorIdxInParent()
	if idx == 0 {
		if len(e.traceStack) <= 1 {
			return
		}
		e.flattenDeleteParent()
		return
	} else {
		switch n := e.getParent().Children()[idx-1].(type) {
		case *parser.ParenCompExpr:
			e.getParent().DeleteChildren(idx, idx)
			e.enterContainerFromRight(n)
		case parser.Container:
			e.getParent().DeleteChildren(idx-1, idx-1)
		default:
			e.getParent().DeleteChildren(idx-1, idx-1)
		}
	}
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

func (e *Editor) getSelectionIdxInParent() int {
	for i, n := range e.getParent().Children() {
		if n == e.markSelect {
			return i
		}
	}
	return -1
}

// ---
// Keyboard input handlers
func (e *Editor) handleLetter(letter rune) {
	kind := e.GetState()
	sel := e.hasSelection()

	if sel {
		e.deleteSelection()
	}

	idx := e.getCursorIdxInParent()
	if kind == EDIT_EQUATION {
		e.getParent().InsertChildren(idx, &parser.VarLit{Source: string(letter)})
	} else {
		e.getParent().InsertChildren(idx, parser.RawRuneLit(letter))
	}
}

func (e *Editor) handleDigit(digit rune) {
	kind := e.GetState()
	sel := e.hasSelection()
	if sel {
		e.deleteSelection()
	}

	idx := e.getCursorIdxInParent()
	if kind == EDIT_EQUATION {
		e.getParent().InsertChildren(idx, &parser.NumberLit{Source: string(digit)})
	} else {
		// TODO handle case where LatexCmdInput is second on stack
		e.getParent().InsertChildren(idx, parser.RawRuneLit(digit))
	}
}

func (e *Editor) handleRest(char rune) {
	// TODO handle special characters _, ^ etc
	kind := e.GetState()
	sel := e.hasSelection()
	idx := e.getCursorIdxInParent()
	if sel {
		switch char {
		case '/': // similar to case '(', ')', merge?
			mark := e.getSelectionIdxInParent()
			numerator := new(parser.CompositeExpr)
			block := parser.Cmd2ArgExpr{Type: parser.CMD_frac, Arg1: numerator, Arg2: new(parser.CompositeExpr)}
			if mark < idx {
				mark, idx = idx, mark
			}
			numerator.AppendChildren(e.getParent().Children()[idx+1 : mark]...)
			e.getParent().DeleteChildren(idx, mark)
			e.getParent().InsertChildren(idx, &block)

			e.enterContainerFromRight(&block)
			e.NavigateDown()

			e.markSelect = nil
			return
		case '&':
			if len(e.traceStack) < 2 {
				break // TODO
			}
			// TODO
			// if env, ok := e.traceStack[len(e.traceStack)-2].(*parser.EnvExpr); ok {
			// 	r, c := env.FindCell(e.getParent().(*parser.UnboundCompExpr))
			// 	idx := e.getCursorIdxInParent()
			// 	mark := e.getSelectionIdxInParent()
			// 	if mark < idx {
			// 		mark, idx = idx, mark
			// 	}
			// 	oldCell := e.getParent()
			// 	cell := env.InsertCell(r, c+1)
			//
			// 	el := oldCell.Children()[idx+1 : mark]
			// 	cell.AppendChildren(e.cursor)
			// 	cell.AppendChildren(el...)
			//
			// 	e.getParent().DeleteChildren(idx, idx)
			// 	e.deleteSelection()
			// 	e.popStack()
			//
			// 	e.traceStack = append(e.traceStack, cell)
			// 	cell.AppendChildren(e.cursor)
			// 	return
			// }
		case '(', ')':
			mark := e.getSelectionIdxInParent()
			block := parser.ParenCompExpr{Left: "(", Right: ")"}
			idxAt := DIR_LEFT
			if mark < idx {
				mark, idx = idx, mark
				idxAt = DIR_RIGHT
			}
			block.AppendChildren(e.getParent().Children()[idx+1 : mark]...)
			e.getParent().DeleteChildren(idx, mark)
			e.getParent().InsertChildren(idx, &block)

			if idxAt == DIR_LEFT {
				e.enterContainerFromLeft(&block)
			} else {
				e.enterContainerFromRight(&block)
			}

			e.markSelect = nil
			return

		default:
			mark := e.getSelectionIdxInParent()
			if mark < idx {
				mark, idx = idx, mark
			}
			e.getParent().DeleteChildren(idx, mark-1)
			e.markSelect = nil
		}
	}

	// is equation and no selection
	if kind == EDIT_EQUATION {
		switch char {
		case '^', '_':
			brace := new(parser.CompositeExpr)
			block := &parser.Cmd1ArgExpr{Type: parser.MatchLatexCmd(string(char)), Arg1: brace}
			e.getParent().DeleteChildren(idx, idx)
			e.getParent().InsertChildren(idx, block)
			e.enterContainerFromLeft(block)
		case '*':
			dot := &parser.SimpleCmdLit{Source: `\cdot`, Type: parser.CMD_cdot}
			e.getParent().InsertChildren(idx, dot)
		case '(':
			block := &parser.ParenCompExpr{Left: "(", Right: ")"}
			e.getParent().DeleteChildren(idx, idx)
			e.getParent().InsertChildren(idx, block)
			e.enterContainerFromRight(block)

		case '\\':
			field := &render.LatexCmdInput{Text: new(parser.TextStringWrapper)}
			e.getParent().DeleteChildren(idx, idx)
			e.getParent().InsertChildren(idx, field)

			e.traceStack = append(e.traceStack, field, field.Text)
			e.getParent().AppendChildren(e.cursor)

		case '&':
			if len(e.traceStack) < 2 {
				break // TODO
			}
			if env, ok := e.traceStack[len(e.traceStack)-2].(*parser.EnvExpr); ok {
				r, c := env.FindCell(e.getParent().(*parser.UnboundCompExpr))
				idx := e.getCursorIdxInParent()
				cell := env.InsertCell(r, c+1)

				e.getParent().DeleteChildren(idx, idx)
				cell.AppendChildren(e.cursor)

				e.popStack()
				e.traceStack = append(e.traceStack, cell)
				return
			}

		case '/':
			e.InsertFrac(true)
			e.NavigateLeft()
			e.NavigateDown()
		case ' ':
			e.getParent().InsertChildren(idx, &parser.SimpleCmdLit{Type: parser.CMD_SPACE, Source: `\ `})
		default:
			e.getParent().InsertChildren(idx, &parser.SimpleOpLit{Source: string(char)})
		}
		return
	}

	// guranteed to pass check
	if n, ok := e.getLastOnStack().(*parser.TextStringWrapper); ok {
		n.InsertChildren(idx, parser.RawRuneLit(char))
	}
}

// not in use yet
func (e *Editor) handlePaste(v string) {
	idx := e.getCursorIdxInParent()

	// TODO error handling, when the given string is not valid latex
	ast := parser.Parse(v)
	e.getParent().InsertChildren(idx, ast.Children()...)
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
		switch n := e.traceStack[searchFrom].(type) {
		case *parser.Cmd2ArgExpr:
			switch n.Command() { // TODO
			case parser.CMD_frac, parser.CMD_binom:
				return searchFrom
			}
		case *parser.EnvExpr:
			return searchFrom
		}
	}
	return -1
}

const KeybindsHelp = `
	Arrow keys - move around
	ctrl+p / ctrl+n / ctrl+f / ctrl+b - Up / Down / Left / Right
	alt+p / alt+n / alt+f / alt+b - Move around without entering a node
	alt + Left/Right - Select Text
	alt + w/W - Select Text

	ctrl + u - delete to start of node
	ctrl + k - delete to end of node
	`

func (e Editor) Update(msg tea.Msg) (Editor, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft, tea.KeyCtrlB:
			if msg.Alt {
				e.selectLeft()
				e.renderer.Sync(e.getLastOnStack(), true)
				return e, nil
			} else if e.markSelect != nil {
				e.cancelSelection()
			}
			e.NavigateLeft()
		case tea.KeyRight, tea.KeyCtrlF:
			if msg.Alt {
				e.selectRight()
				e.renderer.Sync(e.getLastOnStack(), true)
				return e, nil
			} else if e.markSelect != nil {
				e.cancelSelection()
			}
			e.NavigateRight()
		case tea.KeyDown, tea.KeyCtrlN:
			if e.markSelect != nil {
				e.cancelSelection()
			}
			e.NavigateDown()
		case tea.KeyUp, tea.KeyCtrlP:
			if e.markSelect != nil {
				e.cancelSelection()
			}
			e.NavigateUp()
		case tea.KeyHome, tea.KeyCtrlA:
			e.NavigateToBeginning()
		case tea.KeyEnd, tea.KeyCtrlE:
			e.NavigateToEnd()
		case tea.KeyBackspace:
			e.DeleteBack()
		case tea.KeyTab: // TODO complete command when in a LatexCmdInput
			if len(e.traceStack) <= 1 {
				return e, nil
			}
			e.exitParent(DIR_RIGHT)
		case tea.KeyCtrlC:
			return e, tea.Quit
		case tea.KeyEnter:
			switch e.GetState() {
			case EDIT_COMMAND:
				e.realizeCommand()
			case EDIT_TEXT:
				e.exitParent(DIR_RIGHT)
			case EDIT_EQUATION:
				if len(e.traceStack) < 2 {
					break
				}
				if env, ok := e.traceStack[len(e.traceStack)-2].(*parser.EnvExpr); ok {
					r, _ := env.FindCell(e.getParent().(*parser.UnboundCompExpr))
					cell := env.InsertRow(r + 1)
					idx := e.getCursorIdxInParent()
					e.getParent().DeleteChildren(idx, idx)
					e.popStack()
					e.traceStack = append(e.traceStack, cell)
					cell.AppendChildren(e.cursor)
				}
			}

		case tea.KeyCtrlU:
			e.deleteToStart()

		case tea.KeyCtrlK:
			e.deleteToEnd()

		case tea.KeySpace:
			idx := e.getCursorIdxInParent()
			switch e.GetState() {
			case EDIT_EQUATION:
				e.getParent().InsertChildren(idx, &parser.SimpleCmdLit{
					Type:   parser.CMD_SPACE,
					Source: `\ `,
				})

			case EDIT_COMMAND:
				e.realizeCommand()

			case EDIT_TEXT:
				if n, ok := e.getLastOnStack().(*parser.TextStringWrapper); ok {
					n.InsertChildren(idx, parser.RawRuneLit(' '))
				}
			}
		case tea.KeyRunes:
			if len(msg.Runes) == 1 {
				if msg.Alt {
					switch msg.Runes[0] {
					case 'b':
						_ = e.stepOverPrevSibling() || e.NavigateLeft()
					case 'f':
						_ = e.stepOverNextSibling() || e.NavigateRight()
					case 'w':
						e.selectRight()
					case 'W':
						e.selectLeft()
					}
					break
				}
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
	e.renderer.Sync(e.getLastOnStack(), e.markSelect != nil)
	return e, nil
}

// panics if the cursor is not in EDIT_COMMAND mode
func (e *Editor) realizeCommand() {
	if e.GetState() != EDIT_COMMAND {
		panic("realizeCommand called outside of command editing mode")
	}

	if n, ok := e.getLastOnStack().(*parser.TextStringWrapper); ok {
		cmd := "\\" + n.BuildString()
		e.exitParent(DIR_RIGHT)
		idx := e.getCursorIdxInParent() // TODO error handling for idx == 0?
		e.getParent().DeleteChildren(idx-1, idx-1)
		if len(cmd) <= 1 {
			cmd += " "
		}
		e.InsertCmd(cmd)
	}
}

func (e Editor) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, e.renderer.View(), e.banner)
}

// Search the tree for any FixedContainer type that has children that is
// not surrounded in "{}" (CompositeExpr), then wraps them in a CompositeExpr.
// This is to make the editing part easier
func formatLatexTree(tree parser.Expr) {
	// TODO also convert (), [], \{\} to \left...\right

	switch n := tree.(type) {
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

func (e Editor) LatexSource() string {
	return e.config.LatexCfg.ProduceLatex(e.renderer.LatexTree)
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
