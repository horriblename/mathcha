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

func (r *Renderer) View() string {
	builder := strings.Builder{}
	for _, row := range r.Buffer {
		for _, ru := range row {
			builder.WriteRune(ru)
		}
		builder.WriteByte('\n')
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
		return n.Content() // TODO return character being mapped to
	}

	return "[unknown node encountered]"
}
