package renderer

import (
	"strings"

	parser "github.com/horriblename/mathcha/latex"
)

type Renderer struct {
	Buffer       string
	LatexTree    parser.FlexContainer
	FocusOn      parser.Container // the container in which the cursor is, a better implementation would be letting Render functions return a 'focused' flag when cursor is found
	HasSelection bool             // whether there is a selection in FocusOn
	Focus        bool             // whether the widget itself is focused
}

func New() Renderer {
	root := &parser.UnboundCompExpr{}
	return Renderer{
		Buffer:       "",
		LatexTree:    root,
		FocusOn:      root,
		HasSelection: false,
		Focus:        false,
	}
}

func (r *Renderer) Load(tree parser.FlexContainer) {
	r.LatexTree = tree
	r.Sync(nil, false)
}

// rerender the latex tree
func (r *Renderer) Sync(focus parser.Container, selected bool /*whether there is a selection*/) {
	r.FocusOn = focus
	r.HasSelection = selected
	r.DrawToBuffer(r.LatexTree)
}

func (r *Renderer) View() string {
	return r.Buffer
}

// possible optimisation: pass the strings.Builder object by reference into the recursive
// function to avoid multiple Builder instances. Then returning string is no longer needed,
// we just use the original Builder to get the string instead
func ProduceLatex(node parser.Expr, useUnicode bool) string {
	latex := ""
	suffix := ""
	switch n := node.(type) {
	case *parser.TextContainer: // TODO CmdContainer subtype
		return "\\text {" + n.Text.BuildString() + "}"
	case *parser.ParenCompExpr: // TODO FlexContainer subtype
		builder := strings.Builder{}
		builder.WriteString("\\left" + n.Left)
		for _, c := range n.Children() {
			builder.WriteString(ProduceLatex(c, useUnicode))
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
			latex += ProduceLatex(c, useUnicode)
		}
		return latex + suffix
	case parser.CmdContainer:
		latex = n.Command().GetCmd() + " "
		for _, c := range n.Children() {
			latex += ProduceLatex(c, useUnicode)
		}
		return latex
	case parser.CmdLiteral:
		if useUnicode {
			unicode := GetVanillaString(n.Command())
			if len(unicode) == 1 {
				return unicode
			}
			return n.Content() + " "
		}
		return n.Content() + " "
	case *Cursor:
		return ""
	case parser.Literal:
		return n.Content()
	default:
		return "[unknown node encountered]"
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(numbers ...int) int {
	if len(numbers) == 0 {
		panic("min was passed 0 parameters")
	} else if len(numbers) == 1 {
		return numbers[0]
	}

	m := numbers[0]
	for _, n := range numbers {
		if n < m {
			m = n
		}
	}
	return m
}
