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
	subtle   = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	fg       = lipgloss.AdaptiveColor{Light: "#383838", Dark: "#AFAFAF"}
	invert   = lipgloss.AdaptiveColor{Light: "#AFAFAF", Dark: "#383838"}
	accent   = lipgloss.AdaptiveColor{Light: "#579AD1", Dark: "#A1BAEA"}
	accentBg = lipgloss.Color("#505570")

	docStyle   = lipgloss.NewStyle().Foreground(fg)
	focusStyle = lipgloss.NewStyle().Foreground(accent).Background(accentBg)

	underlineStyle = lipgloss.NewStyle().Underline(true)
	cursorStyle    = lipgloss.NewStyle().
			Foreground(invert).
			Background(lipgloss.Color("#FFFFFF"))
	variableStyle = lipgloss.NewStyle().Italic(true)
)

// depth-first traverse of the latex tree and dim tree in parallel
// to build the later rendered string
func (r *Renderer) DrawToBuffer(tree parser.Expr, dim *Dimensions) {
	r.Buffer, _ = r.Prerender(tree, dim)
}

func (r *Renderer) Prerender(node parser.Expr, dim *Dimensions) (out string, baseLevel int) {
	defer func() {
		if node == r.FocusOn {
			out = focusStyle.Render(out)
		}
	}()
	switch n := node.(type) {
	case *parser.TextContainer:
		return r.Prerender(n.Text, dim)
	case *LatexCmdInput:
		str, baseLevel := r.Prerender(n.Text, dim)
		return "\\" + str, baseLevel
	case *parser.TextStringWrapper:
		if CONF_RENDER_EMPTY_COMP_EXPR {
			var builder strings.Builder
			for _, i := range n.Runes {
				switch r := i.(type) {
				case parser.RawRuneLit:
					builder.WriteRune(rune(r))
				case *Cursor:
					builder.WriteString(r.Content())
				default: // panic?
				}
			}
			return builder.String(), 0
		}
		return n.BuildString(), 0

	case parser.CmdContainer:
		switch n.Command() {
		case parser.CMD_underline:
			return r.PrerenderCmdUnderline(n, dim)
		case parser.CMD_frac:
			return r.PrerenderCmdFrac(n, dim)
		case parser.CMD_superscript:
			str, _ := r.Prerender(n.Children()[0], dim)
			return str, 1
		case parser.CMD_subscript:
			str, _ := r.Prerender(n.Children()[0], dim)
			return str, -lipgloss.Height(str)
		case parser.CMD_sqrt:
			return r.PrerenderCmdSqrt(n, dim)
		default:
			return "[unimplemented command container]", 0
		}

	case *parser.ParenCompExpr:
		content, baseLine := r.PrerenderFlexContainer(n, dim)
		if n.Left == "(" && n.Right == ")" && lipgloss.Height(content) >= 2 {
			height := lipgloss.Height(content)
			left := "╭\n" + strings.Repeat("│\n", height-2) + "╰"
			right := "╮\n" + strings.Repeat("│\n", height-2) + "╯"
			return JoinHorizontal([]int{baseLine, baseLine, baseLine}, left, content, right), baseLine
		}
		return JoinHorizontal([]int{0, baseLine, 0}, n.Left, content, n.Right), baseLine
	case parser.FlexContainer:
		return r.PrerenderFlexContainer(n, dim)
	case *parser.UnknownCmdLit: // FIXME subcase of CmdLiteral, what to do with UnknownCmdLit?
		return "\x1b[4m?\x1b[24m", 0
	case parser.CmdLiteral:
		content := GetVanillaString(n.Command())
		return content, 0
		// parser.Literal interface types
	case *parser.VarLit:
		return "\x1b[3m" + n.Content() + "\x1b[23m", 0 // apply italic(3) then unset italic(23)
	// case *Cursor:
	// 	return "\x1b[7m \x1b[27m", 0 // set bg color as white(47) then set bg color to default(49)
	case parser.Literal:
		content := n.Content()
		switch content {
		case "+", "-", "=":
			content = " " + content + " "
		}
		return content, 0
	case nil:
		// TODO handle error?
		return "[nil]", 0
	default:
		return "[unimplemented expression]", 0
	}
	// panic("Unhandled case in Prerender()")
}

func (r *Renderer) PrerenderFlexContainer(node parser.FlexContainer, dim *Dimensions) (output string, baseLine int) {
	if len(node.Children()) <= 0 {
		if CONF_RENDER_EMPTY_COMP_EXPR {
			return " ", 0
		} else {
			return "", 0
		}
	}
	var renderedChildren = make([]string, len(node.Children()))
	var baseLines = make([]int, len(node.Children()))
	for i, c := range node.Children() {
		renderedChildren[i], baseLines[i] = r.Prerender(c, dim)
	}
	return JoinHorizontal(baseLines, renderedChildren...), min(baseLines...)
}

// TODO remove
func (r *Renderer) PrerenderCmdContainer(node parser.CmdContainer, dim *Dimensions, x int, y int) (output string, baseLine int) {
	switch node.Command() {
	case parser.CMD_frac:
		return r.PrerenderCmdFrac(node, dim)
	}

	return "[unimplemented cmd container]", 0
}

func (r *Renderer) PrerenderCmdUnderline(node parser.CmdContainer, dim *Dimensions) (output string, baseLevel int) {
	block, baseLevel := r.Prerender(node.Children()[0], dim.Children[0])
	lines, _ := getLines(block)
	// lines[len(lines)-1] = underlineStyle.Copy().Render(lines[len(lines)-1])
	// lines[len(lines)-1] = termenv.String(lines[len(lines)-1]).Underline().String()
	// \x1b[4m sets underline, \x1b[24m unsets it
	lines[len(lines)-1] = "\x1b[4m" + lines[len(lines)-1] + "\x1b[24m"

	return lipgloss.JoinVertical(lipgloss.Center, lines...), baseLevel
}

func (r *Renderer) PrerenderCmdFrac(node parser.CmdContainer, dim *Dimensions) (output string, newBaseLevel int) {
	arg1, _ := r.Prerender(node.Children()[0], dim)
	arg2, _ := r.Prerender(node.Children()[1], dim)
	width := max(lipgloss.Width(arg1), lipgloss.Width(arg2))
	newBaseLevel = -lipgloss.Height(arg2)
	line := strings.Repeat("─", width)

	return lipgloss.JoinVertical(lipgloss.Center, arg1, line, arg2), newBaseLevel
}

func (r *Renderer) PrerenderCmdSqrt(node parser.CmdContainer, dim *Dimensions) (output string, baseLevel int) {
	// TODO simplify adding overline escape chars
	block, baseLevel := r.Prerender(node.Children()[0], dim)
	lines, _ := getLines(block)
	lines[0] = "\x1b[53m" + lines[0] + "\x1b[55m"
	block = lipgloss.JoinVertical(lipgloss.Center, lines...)
	height := lipgloss.Height(block)
	root := strings.Repeat("│\n", height-1) + `√`

	return JoinHorizontal([]int{baseLevel, baseLevel}, root, block), baseLevel
}
