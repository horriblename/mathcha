package mathrender

import (
	"fmt"
	"strings"

	parser "github.com/horriblename/latex-parser/latex"
)

// FIXME
var x = fmt.Println
var _ = strings.Title

// TODO maybe add Dimensions.Init

type Renderer struct {
	Buffer    [][]rune
	LatexTree parser.FlexContainer
	Size      *Dimensions
}

// Rendering is a 3 step process:
// 1. build a separate tree with all the dimensions, then create a [][]rune buffer of appropriate size
//    - in this step the rendered string for Literals are also computed and stored in Dimensions.Lit
// 2. walk through the ast and dimensions tree in parallel and write the contents to the buffer
// 3. combine the [][]rune buffer into a string
func (r *Renderer) Load(tree parser.FlexContainer) {
	// step 1: create [][]rune buffer of appropriate size
	r.LatexTree = tree
	r.Size = calculateDim(r.LatexTree)
	r.Buffer = make([][]rune, r.Size.Height)
	for i := range r.Buffer {
		r.Buffer[i] = make([]rune, r.Size.Width)
		for j := range r.Buffer[i] { // TODO there should be better ways to do this
			r.Buffer[i][j] = ' '
		}
	}
	// step 2: write characters
	println("w, h, b", r.Size.Width, r.Size.Height, r.Size.BaseLine)
	r.DrawToBuffer(r.LatexTree, r.Size)
}

// rerender the latex tree
func (r *Renderer) Sync() {
	r.Size = calculateDim(r.LatexTree)
	// TODO
	// if r.Size.Width != len(r.Buffer[0]){
	//      r.Buffer[0]
	//   }

}

// fill a [][]rune buffer with whitespace. slices are passed by reference so
// I shouldn't need to return it?
func blankBuffer(buffer [][]rune, width int, height int) {
	for i := range buffer {
		buffer[i] = make([]rune, width)
		for j := range buffer[i] { // TODO there should be better ways to do this
			buffer[i][j] = ' '
		}
	}
}

func (r *Renderer) View() string {
	builder := strings.Builder{}
	for _, row := range r.Buffer {
		builder.WriteString("│")
		for _, ru := range row {
			builder.WriteRune(ru)
		}
		builder.WriteString("│\n")
	}

	return builder.String()
}

func ProduceLatex(node parser.Expr) string {
	latex := ""
	suffix := ""
	switch n := node.(type) {
	case parser.FlexContainer:
		if n.Identifier() == "{" { // CompositeExpr
			latex = "{"
			suffix = "}"
		}
		for _, c := range n.Children() {
			latex += ProduceLatex(c)
		}
		return latex + suffix
	case parser.CmdContainer:
		latex = n.Command().GetCmd() + " "
		for _, c := range n.Children() {
			latex += ProduceLatex(c)
		}
		return latex
	case parser.Literal:
		return n.Content()
	case parser.CmdLiteral:
		return n.Content()
	}

	return "[unknown node encountered]"
}
