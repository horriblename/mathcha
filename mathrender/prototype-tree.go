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
	Buffer  [][]rune
	Size    *Dimensions
	DrawerX int // position of drawer cursor
	DrawerY int
}

// Rendering is a 3 step process:
// 1. build a separate tree with all the dimensions, then create a [][]rune buffer of appropriate size
// 2. walk through the ast and dimensions tree in parallel and write the contents to the buffer
// 3. combine the [][]rune buffer into a string
func (r *Renderer) Load(tree parser.Expr) {
	// step 1: create [][]rune buffer of appropriate size
	r.Size = calculateDim(tree)
	r.Buffer = make([][]rune, r.Size.Height)
	for i := range r.Buffer {
		r.Buffer[i] = make([]rune, r.Size.Width)
	}
	// step 2: write characters
	println("w, h, b", r.Size.Width, r.Size.Height, r.Size.BaseLine)
}

// TODO So right now both calculateDim*() and Prerender*() are mutual
// recursive functions, maybe merge them afterwards?

func ProduceLatex(node parser.Expr) string {
	latex := ""
	suffix := ""
	if n, ok := node.(parser.FlexContainer); ok {
		if n.Identifier() == "{" { // CompositeExpr
			latex = "{"
			suffix = "}"
		}
		for _, c := range n.Children() {
			latex += ProduceLatex(c)
		}
		return latex + suffix
	}

	if n, ok := node.(parser.CmdContainer); ok {
		latex = n.Command().GetCmd() + " "
		for _, c := range n.Children() {
			latex += ProduceLatex(c)
		}
		return latex
	}

	if n, ok := node.(parser.Literal); ok {
		return n.Content()
	}
	if n, ok := node.(parser.CmdLiteral); ok {
		return n.Content() // TODO return character being mapped to
	}

	return "[unknown node encountered]"
}
