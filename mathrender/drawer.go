package mathrender

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	parser "github.com/horriblename/latex-parser/latex"
	rw "github.com/mattn/go-runewidth"
)

// style definitions
var (
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	fg     = lipgloss.AdaptiveColor{Light: "#383838", Dark: "#AFAFAF"}

	docStyle = lipgloss.NewStyle().Foreground(fg).Align(lipgloss.Center)

	fracBorder = lipgloss.Border{
		Top:    "─",
		Bottom: "─",
	}

	fracNumerStyle = docStyle.
			BorderTop(false).
			BorderBottom(true).
			BorderStyle(fracBorder)

	fracDenomStyle = docStyle.
			BorderTop(true).
			BorderBottom(false).
			BorderStyle(fracBorder)
)

// depth-first traverse of the latex tree and dim tree in parallel
// to build the later rendered string
func (r *Renderer) DrawToBuffer(tree parser.Expr, dim *Dimensions) {
	r.Buffer = r.Prerender(tree, dim, 0, 0)
}

// x is the position of leftmost rune allowed to be written by the node
// y is the position of baseline=0 of the node
func (r *Renderer) Prerender(node parser.Expr, dim *Dimensions, x int, y int) string {
	switch n := node.(type) {
	case parser.CmdContainer:
		switch n.Command() {
		case parser.CMD_frac:
			return r.PrerenderCmdFrac(n, dim, x, y)
		}
	case parser.FlexContainer:
		return r.PrerenderFlexContainer(n, dim, x, y)
	case parser.CmdLiteral:
		content := GetVanillaString(n.Command())
		return content
		// parser.Literal interface types
	case parser.Literal:
		content := n.Content()
		switch content {
		case "+", "-", "=":
			content = " " + n.Content() + " "
		}
		return content
	}
	panic("Unhandled case in Prerender()")
}

func (r *Renderer) PrerenderFlexContainer(node parser.FlexContainer, dim *Dimensions, x int, y int) string {
	var children = []string{}
	var baseLine = []int{}
	for i, c := range node.Children() {
		children = append(children, r.Prerender(c, dim.Children[i], 0, 0)) //TODO
		baseLine = append(baseLine, dim.Children[i].BaseLine)
	}
	// println(lipgloss.JoinHorizontal(lipgloss.Center, children...))
	return JoinHorizontal(baseLine, children...)
}

func (r *Renderer) PrerenderCmdContainer(node parser.CmdContainer, dim *Dimensions, x int, y int) string {
	switch node.Command() {
	case parser.CMD_frac:
		return r.PrerenderCmdFrac(node, dim, x, y)
	}

	return ""
}

func (r *Renderer) PrerenderCmdFrac(node parser.CmdContainer, dim *Dimensions, x int, y int) string {
	arg1 := r.Prerender(node.Children()[0], dim.Children[0], x, y)
	arg2 := fracNumerStyle.Render(r.Prerender(node.Children()[1], dim.Children[1], x, y))

	return lipgloss.JoinVertical(lipgloss.Center, arg1, arg2)
}

func blockWidth(block string) int {
	println(strings.SplitN(block, "\n", 1))
	return rw.StringWidth(strings.Split(block, "\n")[0])
}
