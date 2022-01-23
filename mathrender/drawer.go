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
	invert = lipgloss.AdaptiveColor{Light: "#AFAFAF", Dark: "#383838"}
	accent = lipgloss.AdaptiveColor{Light: "#579AD1", Dark: "#579AD1"}

	docStyle   = lipgloss.NewStyle().Foreground(fg)
	focusStyle = lipgloss.NewStyle().Foreground(accent)

	underlineStyle = lipgloss.NewStyle().Underline(true)
	cursorStyle    = lipgloss.NewStyle().
			Foreground(invert).
			Background(lipgloss.Color("#FFFFFF"))
	variableStyle = lipgloss.NewStyle().Italic(true)
)

// depth-first traverse of the latex tree and dim tree in parallel
// to build the later rendered string
func (r *Renderer) DrawToBuffer(tree parser.Expr, dim *Dimensions) {
	r.Buffer = r.Prerender(tree, dim)
}

func (r *Renderer) Prerender(node parser.Expr, dim *Dimensions) (out string) {
	defer func() {
		if node == r.FocusOn {
			out = focusStyle.Render(out)
		}
	}()
	switch n := node.(type) {
	case *parser.TextContainer:
		return n.Text.Content()
	case parser.CmdContainer:
		switch n.Command() {
		case parser.CMD_underline:
			return r.PrerenderCmdUnderline(n, dim)
		case parser.CMD_frac:
			return r.PrerenderCmdFrac(n, dim)
		default:
			return "[unimplemented command]"
		}
	case *parser.ParenCompExpr:
		content := r.PrerenderFlexContainer(n, dim)
		return JoinHorizontal([]int{dim.BaseLine, dim.Children[0].BaseLine, dim.BaseLine}, n.Left, content, n.Right)
	case parser.FlexContainer:
		return r.PrerenderFlexContainer(n, dim)
	case parser.CmdLiteral:
		content := GetVanillaString(n.Command())
		return content
		// parser.Literal interface types
	case *parser.VarLit:
		return "\x1b[3m" + n.Content() + "\x1b[23m" // apply italic(3) then unset italic(23)
	case *Cursor:
		return "\x1b[47m \x1b[49m" // set bg color as white(47) then set bg color to default(49)
	case parser.Literal:
		content := n.Content()
		switch content {
		case "+", "-", "=":
			content = " " + content + " "
		}
		return content
	case nil:
		// TODO handle error?
		return "[nil]"
	default:
		return "[unimplemented expression]"
	}
	// panic("Unhandled case in Prerender()")
}

func (r *Renderer) PrerenderFlexContainer(node parser.FlexContainer, dim *Dimensions) string {
	var children = make([]string, len(node.Children()))
	var baseLine = make([]int, len(node.Children()))
	for i, c := range node.Children() {
		children[i] = r.Prerender(c, dim.Children[i])
		baseLine[i] = dim.Children[i].BaseLine
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

func (r *Renderer) PrerenderCmdUnderline(node parser.CmdContainer, dim *Dimensions) string {
	block := r.Prerender(node.Children()[0], dim.Children[0])
	lines, _ := getLines(block)
	// lines[len(lines)-1] = underlineStyle.Copy().Render(lines[len(lines)-1])
	// lines[len(lines)-1] = termenv.String(lines[len(lines)-1]).Underline().String()
	// \x1b[4m sets underline, \x1b[24m unsets it
	lines[len(lines)-1] = "\x1b[4m" + lines[len(lines)-1] + "\x1b[24m"

	return lipgloss.JoinVertical(lipgloss.Center, lines...)
}

func (r *Renderer) PrerenderCmdFrac(node parser.CmdContainer, dim *Dimensions) string {
	arg1 := r.Prerender(node.Children()[0], dim.Children[0])
	arg2 := r.Prerender(node.Children()[1], dim.Children[1])
	width := max(lipgloss.Width(arg1), lipgloss.Width(arg2))
	line := strings.Repeat("─", width)

	return lipgloss.JoinVertical(lipgloss.Center, arg1, line, arg2)
}
