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
	Buffer    string
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

	// println("w, h, b", r.Size.Width, r.Size.Height, r.Size.BaseLine)
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

func (r *Renderer) View() string {
	return r.Buffer
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
