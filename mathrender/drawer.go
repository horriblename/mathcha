package mathrender

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	parser "github.com/horriblename/latex-parser/latex"
)

const (
	CONF_RENDER_EMPTY_COMP_EXPR = true // config to enable rendering empty CompositeExpr "{}" as a space
)

// style definitions
var (
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	fg     = lipgloss.AdaptiveColor{Light: "#383838", Dark: "#AFAFAF"}

	docStyle = lipgloss.NewStyle().Foreground(fg).Align(lipgloss.Center)


	fracDenomStyle = docStyle.
			BorderTop(true).
			BorderBottom(false).
			BorderStyle(fracBorder)
)

// depth-first traverse of the latex tree and dim tree in parallel
// to build the later rendered string
func (r *Renderer) DrawToBuffer(tree parser.Expr, dim *Dimensions) {
	r.Buffer = r.Prerender(tree, dim)
}

func (r *Renderer) Prerender(node parser.Expr, dim *Dimensions) string {
	switch n := node.(type) {
	case parser.CmdContainer:
		switch n.Command() {
		case parser.CMD_underline:
			return r.PrerenderCmdUnderline(n, dim)
		case parser.CMD_frac:
			return r.PrerenderCmdFrac(n, dim)
		default:
			return "[unimplemented command]"
		}
	case parser.FlexContainer:
		return r.PrerenderFlexContainer(n, dim)
	case parser.CmdLiteral:
		content := GetVanillaString(n.Command())
		return content
		// parser.Literal interface types
	case parser.Literal:
		content := n.Content()
		switch content {
		case "+", "-", "=":
			content = " " + content + " "
		}
		return content
	case *Cursor:
		return n.Content()
	case nil:
		return "[nil]"
	default:
		return "[unimplemented expression]"
	}
	// panic("Unhandled case in Prerender()")
}

func (r *Renderer) PrerenderFlexContainer(node parser.FlexContainer, dim *Dimensions) string {
	var children = []string{}
	var baseLine = []int{}
	for i, c := range node.Children() {
		children = append(children, r.Prerender(c, dim.Children[i])) //TODO
		baseLine = append(baseLine, dim.Children[i].BaseLine)
	}
	if len(children) == 0 && CONF_RENDER_EMPTY_COMP_EXPR {
		return " "
	}
	// println(lipgloss.JoinHorizontal(lipgloss.Center, children...))
	return JoinHorizontal(baseLine, children...)
}

func (r *Renderer) PrerenderCmdContainer(node parser.CmdContainer, dim *Dimensions, x int, y int) string {
	switch node.Command() {
	case parser.CMD_frac:
		return r.PrerenderCmdFrac(node, dim)
	}

	return ""
}

func (r *Renderer) PrerenderCmdFrac(node parser.CmdContainer, dim *Dimensions, x int, y int) string {
	// FIXME add proper horizontal line
	arg1 := r.Prerender(node.Children()[0], dim.Children[0])
	arg2 := r.Prerender(node.Children()[1], dim.Children[1])
	width := max(blockWidth(arg1), blockWidth(arg2))
	line := strings.Repeat("─", width)

	return lipgloss.JoinVertical(lipgloss.Center, arg1, line, arg2)
}

func blockWidth(block string) int {
	_, width := getLines(block)
	return width
}
