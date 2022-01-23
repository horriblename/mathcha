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
	FocusOn   parser.Container
	Size      *Dimensions
}

// Rendering is a 2 step process: TODO merge the process?
// 1. build a separate tree with all the dimensions
// 2.
func (r *Renderer) Load(tree parser.FlexContainer) {
	r.LatexTree = tree
	r.Sync(nil)
}

// rerender the latex tree
func (r *Renderer) Sync(focus parser.Container) {
	r.Size = calculateDim(r.LatexTree)
	r.FocusOn = focus
	r.DrawToBuffer(r.LatexTree, r.Size)
}

func (r *Renderer) View() string {
	return r.Buffer
}

func ProduceLatex(node parser.Expr) string {
	latex := ""
	suffix := ""
	switch n := node.(type) {
	case *parser.TextContainer: // TODO CmdContainer subtype
		latex = "\\text {" + n.Text.BuildString() + "}"
	case parser.FlexContainer:
		// FIXME type switch for cases CompositeExpr and others
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
	case parser.CmdLiteral:
		return n.Content() + " "
	case parser.Literal:
		return n.Content()
	}

	return "[unknown node encountered]"
}
