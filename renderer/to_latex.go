package renderer

import (
	"strings"
	"unicode"
	"unicode/utf8"

	parser "github.com/horriblename/mathcha/latex"
)

type LatexSourceConfig struct {
	UseUnicode bool
}

// possible optimisation: pass the strings.Builder object by reference into the recursive
// function to avoid multiple Builder instances. Then returning string is no longer needed,
// we just use the original Builder to get the string instead
func (cfg *LatexSourceConfig) ProduceLatex(node parser.Expr) string {
	latex := ""
	suffix := ""
	switch n := node.(type) {
	case *parser.TextContainer: // TODO CmdContainer subtype
		return "\\text {" + n.Text.BuildString() + "}"
	case *parser.ParenCompExpr: // TODO FlexContainer subtype
		builder := strings.Builder{}
		builder.WriteString("\\left" + n.Left)
		for _, c := range n.Children() {
			builder.WriteString(cfg.ProduceLatex(c))
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
			latex += cfg.ProduceLatex(c)
		}
		return latex + suffix
	case parser.CmdContainer:
		latex = n.Command().GetCmd()
		if unicode.IsLetter(rune(latex[len(latex)-1])) {
			latex += " "
		}

		for _, c := range n.Children() {
			latex += cfg.ProduceLatex(c)
		}
		return latex
	case parser.CmdLiteral:
		if cfg.UseUnicode {
			renderedString := GetVanillaString(n.Command())
			if utf8.RuneCountInString(renderedString) == 1 {
				return renderedString
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
