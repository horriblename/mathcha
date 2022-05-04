package renderer

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
	Buffer       string
	LatexTree    parser.FlexContainer
	FocusOn      parser.Container // the container in which the cursor is, a better implementation would be letting Render functions return a 'focused' flag when cursor is found
	HasSelection bool             // whether there is a selection in FocusOn
	Size         *Dimensions
}

// Rendering is a 2 step process: TODO merge the process?
// 1. build a separate tree with all the dimensions
// 2.
func (r *Renderer) Load(tree parser.FlexContainer) {
	r.LatexTree = tree
	r.Sync(nil, false)
}

// rerender the latex tree
func (r *Renderer) Sync(focus parser.Container, selected bool /*whether there is a selection*/) {
	r.Size = calculateDim(r.LatexTree)
	r.FocusOn = focus
	r.HasSelection = selected
	r.DrawToBuffer(r.LatexTree, r.Size)
}

func (r *Renderer) View() string {
	return r.Buffer
}

// possible optimisation: pass the strings.Builder object by reference into the recursive
// function to avoid multiple Builder instances. Then returning string is no longer needed,
// we just use the original Builder to get the string instead
func ProduceLatex(node parser.Expr) string {
	latex := ""
	suffix := ""
	switch n := node.(type) {
	case *parser.TextContainer: // TODO CmdContainer subtype
		return "\\text {" + n.Text.BuildString() + "}"
	case *parser.ParenCompExpr: // TODO FlexContainer subtype
		builder := strings.Builder{}
		builder.WriteString("\\left" + n.Left)
		for _, c := range n.Children() {
			builder.WriteString(ProduceLatex(c))
		}
		builder.WriteString("\\right" + n.Right)
		return builder.String()

	case parser.FlexContainer:
		// FIXME type switch for cases CompositeExpr and others
		// FIXME use strings.Builder instead
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
	case *Cursor:
		return ""
	case parser.Literal:
		return n.Content()
	default:
		return "[unknown node encountered]"
	}
}
