package renderer

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	parser "github.com/horriblename/mathcha/latex"
)

const (
	CONF_RENDER_EMPTY_COMP_EXPR = true // config to enable rendering empty CompositeExpr "{}" as a space
)

// style definitions
var (
	subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	fg        = lipgloss.AdaptiveColor{Light: "#1ff2f7", Dark: "#abb2bf"}
	invert    = lipgloss.AdaptiveColor{Light: "#abb2bf", Dark: "#1ff2f7"}
	accent    = fg //lipgloss.AdaptiveColor{Light: "#264f78", Dark: "#A1BAEA"}
	accentBg  = lipgloss.Color("#555")
	highlight = lipgloss.Color("#264f78")

	docStyle       = lipgloss.NewStyle().Foreground(fg)
	focusStyle     = lipgloss.NewStyle().Foreground(accent).Background(accentBg)
	highlightStyle = focusStyle.Background(highlight).Foreground(fg)

	// underlineStyle = lipgloss.NewStyle().Underline(true)
	variableStyle = lipgloss.NewStyle().Italic(true)
)

func (r *Renderer) DrawToBuffer(tree parser.Expr) {
	r.Buffer, _ = r.Prerender(tree)
}

func (r *Renderer) Prerender(node parser.Expr) (out string, baseLevel int) {
	defer func() {
		if node == r.FocusOn && r.Focus {
			out = focusStyle.Render(out)
		}
	}()
	switch n := node.(type) {
	case *parser.TextContainer:
		return r.Prerender(n.Text)
	case *LatexCmdInput:
		str, baseLevel := r.Prerender(n.Text)
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

	case *parser.EnvExpr:
		numCols := 0
		for _, row := range n.Elts {
			if len(row) > numCols {
				numCols = len(row)
			}
		}

		colWidths := make([]int, numCols)
		rendered := make([][]string, len(n.Elts))
		baseLines := make([][]int, len(n.Elts))
		for rowIdx, row := range n.Elts {
			rendered[rowIdx] = make([]string, len(row))
			baseLines[rowIdx] = make([]int, len(row))
			for colIdx, cell := range row {
				cellStr, baseLine := r.Prerender(cell)
				rendered[rowIdx][colIdx] = cellStr
				baseLines[rowIdx][colIdx] = baseLine
				width := lipgloss.Width(cellStr)
				if width > colWidths[colIdx] {
					colWidths[colIdx] = width
				}
			}
		}

		var rows []string
		for rowIdx, row := range n.Elts {
			var cells []string
			cellBaseLines := make([]int, len(row)*2-1)
			for i := range row {
				cellStr := rendered[rowIdx][i]
				cellStr = cellStr + strings.Repeat(" ", colWidths[i]-lipgloss.Width(cellStr))
				if i != 0 {
					cells = append(cells, " ")
					cellBaseLines[i*2-1] = 0
				}
				cells = append(cells, cellStr)
				cellBaseLines[i*2] = baseLines[rowIdx][i]
			}
			rowStr := JoinHorizontal(cellBaseLines, cells...)
			rows = append(rows, rowStr)
		}

		body := lipgloss.JoinVertical(lipgloss.Top, rows...)
		height := lipgloss.Height(body)
		left := "["
		right := "]"
		if height != 1 {
			left = constructParenLike(height, "⎡", "⎢", "⎣")
			right = constructParenLike(height, "⎤", "⎥", "⎦")
		}

		return JoinHorizontal([]int{0, 0, 0}, left, body, right), 0

	case parser.CmdContainer:
		switch n.Command() {
		case parser.CMD_overline:
			return r.PrerenderCmdOverline(n)
		case parser.CMD_underline:
			return r.PrerenderCmdUnderline(n)
		case parser.CMD_frac:
			return r.PrerenderCmdFrac(n)
		case parser.CMD_superscript:
			str, _ := r.Prerender(n.Children()[0])
			return str, 1
		case parser.CMD_subscript:
			str, _ := r.Prerender(n.Children()[0])
			return str, -lipgloss.Height(str)
		case parser.CMD_sqrt:
			return r.PrerenderCmdSqrt(n)
		default:
			return "[unimplemented command container]", 0
		}

	case *parser.ParenCompExpr:
		content, baseLine := r.PrerenderFlexContainer(n)
		if n.Left == "(" && n.Right == ")" && lipgloss.Height(content) >= 2 {
			height := lipgloss.Height(content)
			left := constructParenLike(height, "⎛", "⎜", "⎝")
			right := constructParenLike(height, "⎞", "⎟", "⎠")
			return JoinHorizontal([]int{baseLine, baseLine, baseLine}, left, content, right), baseLine
		}
		return JoinHorizontal([]int{0, baseLine, 0}, n.Left, content, n.Right), baseLine
	case parser.FlexContainer:
		return r.PrerenderFlexContainer(n)
	case *parser.UnknownCmdLit: // FIXME subcase of CmdLiteral, what to do with UnknownCmdLit?
		return r.styleAndReset(underline, "?"), 0
	case parser.CmdLiteral:
		content := GetVanillaString(n.Command())
		return content, 0
		// parser.Literal interface types
	case *parser.VarLit:
		return r.styleAndReset(italic, n.Content()), 0
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

func (r *Renderer) PrerenderFlexContainer(node parser.FlexContainer) (output string, baseLine int) {
	if len(node.Children()) <= 0 {
		if CONF_RENDER_EMPTY_COMP_EXPR {
			return " ", 0
		} else {
			return "", 0
		}
	}
	var renderedChildren = make([]string, len(node.Children()))
	var baseLines = make([]int, len(node.Children()))
	var vertJoinQueue *parser.Cmd1ArgExpr // for elements that need to be rendered on top of one another superscrpit & subscript

	// init only when r.FocusOn == node?
	var selStart, selEnd = -1, -1 // [start, end] of the selection
	for index, child := range node.Children() {
		if r.HasSelection && r.FocusOn == node {
			if _, ok := child.(*Cursor); ok {
				if selStart > -1 {
					selEnd = index
				} else {
					selStart = index
				}
			}
			if selStart > -1 && selEnd == -1 {
				continue
			}
		}

		// deal with elements that render on top of eaech other
		if c, ok := child.(*parser.Cmd1ArgExpr); ok {
			switch c.Command() {
			// stack neighboring superscripts and subscripts onto each other
			case parser.CMD_subscript:
				var sup, sub string
				if vertJoinQueue != nil {
					if vertJoinQueue.Command() == parser.CMD_superscript {
						sup = renderedChildren[index-1]
						sub, baseLines[index] = r.Prerender(c)
						renderedChildren[index] = lipgloss.JoinVertical(lipgloss.Left, sup, " ", sub)
						// println(renderedChildren[index])
						renderedChildren[index-1] = ""
						continue
					}
				}

				vertJoinQueue = c
			case parser.CMD_superscript: // TODO merge above
				var sup, sub string
				if vertJoinQueue != nil {
					if vertJoinQueue.Command() == parser.CMD_subscript {
						sub = renderedChildren[index-1]
						sup, _ = r.Prerender(c)
						baseLines[index] = baseLines[index-1]
						renderedChildren[index] = lipgloss.JoinVertical(lipgloss.Left, sup, " ", sub)
						renderedChildren[index-1] = ""
						continue
					}
				}

				vertJoinQueue = c
			default:
				vertJoinQueue = nil
			}
		} else {
			vertJoinQueue = nil
		}
		renderedChildren[index], baseLines[index] = r.Prerender(child)
	}

	if 0 <= selStart && selStart < selEnd {
		str, base := r.Prerender(&parser.UnboundCompExpr{Elts: node.Children()[selStart:selEnd]})
		// FIXME workaround for highlight hiding active background
		activeBg := r.backgroundRGB(80, 80, 80)
		highlightBg := r.backgroundRGB(26, 79, 120)
		lines, _ := getLines(str)
		for i, line := range lines {
			lines[i] = highlightBg + line + activeBg
		}
		renderedChildren[selStart] = lipgloss.JoinVertical(lipgloss.Center, lines...)
		baseLines[selStart] = base
	}
	return JoinHorizontal(baseLines, renderedChildren...), min(baseLines...)
}

// TODO remove
func (r *Renderer) PrerenderCmdContainer(node parser.CmdContainer, x int, y int) (output string, baseLine int) {
	switch node.Command() {
	case parser.CMD_frac:
		return r.PrerenderCmdFrac(node)
	}

	return "[unimplemented cmd container]", 0
}

func (r *Renderer) PrerenderCmdOverline(node parser.CmdContainer) (output string, baseLevel int) {
	block, baseLevel := r.Prerender(node.Children()[0])
	lines, _ := getLines(block)
	lines[0] = "\x1b[53m" + lines[0] + "\x1b[55m"

	return lipgloss.JoinVertical(lipgloss.Center, lines...), baseLevel
}

func (r *Renderer) PrerenderCmdUnderline(node parser.CmdContainer) (output string, baseLevel int) {
	block, baseLevel := r.Prerender(node.Children()[0])
	lines, _ := getLines(block)
	// \x1b[4m sets underline, \x1b[24m unsets it
	// not using lipgloss as lipgloss ends with \x1b[0m, which resets everything
	// TODO handle double underlines?
	lines[len(lines)-1] = r.styleAndReset(underline, lines[len(lines)-1])

	return lipgloss.JoinVertical(lipgloss.Center, lines...), baseLevel
}

func (r *Renderer) PrerenderCmdFrac(node parser.CmdContainer) (output string, newBaseLevel int) {
	arg1, _ := r.Prerender(node.Children()[0])
	arg2, _ := r.Prerender(node.Children()[1])
	width := max(lipgloss.Width(arg1), lipgloss.Width(arg2))
	newBaseLevel = -lipgloss.Height(arg2)
	line := strings.Repeat("─", width)

	return lipgloss.JoinVertical(lipgloss.Center, arg1, line, arg2), newBaseLevel
}

func (r *Renderer) PrerenderCmdSqrt(node parser.CmdContainer) (output string, baseLevel int) {
	// TODO simplify adding overline escape chars
	block, baseLevel := r.Prerender(node.Children()[0])
	lines, _ := getLines(block)
	lines[0] = r.overlineAndReset(lines[0])
	block = lipgloss.JoinVertical(lipgloss.Center, lines...)
	height := lipgloss.Height(block)
	root := strings.Repeat("⎟\n", height-1) + `⎷`

	return JoinHorizontal([]int{baseLevel, baseLevel}, root, block), baseLevel
}

func constructParenLike(height int, top, mid, bot string) string {
	if height == 1 {
		return mid
	}
	return top + "\n" + strings.Repeat(mid+"\n", height-2) + bot
}
