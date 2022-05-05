// Some extended parser.Expr that are useful to the editor
package renderer

import (
	parser "github.com/horriblename/mathcha/latex"
)

// the Cursor object implements a zero-width parser.Literal TODO rn its still a normal character
type Cursor struct {
	Symbol string // appearance of the cursor
}

type LatexCmdInput struct {
	Text *parser.TextStringWrapper
}

func (c *Cursor) Pos() parser.Pos { return parser.Pos(0) } // FIXME remove
func (c *Cursor) End() parser.Pos { return parser.Pos(0) }

func (c *Cursor) VisualizeTree() string { return "Cursor" + c.Symbol }
func (c *Cursor) Content() string       { return c.Symbol }

func (x *LatexCmdInput) Pos() parser.Pos         { return 0 }
func (x *LatexCmdInput) End() parser.Pos         { return 0 }
func (x *LatexCmdInput) Children() []parser.Expr { return []parser.Expr{x.Text} }
func (x *LatexCmdInput) Parameters() int         { return 1 }
func (x *LatexCmdInput) SetArg(index int, expr parser.Expr) {
	if index > 0 {
		panic("SetArg(): index out of range")
	}
	// FIXME this is awful
	if n, ok := expr.(*parser.TextStringWrapper); ok {
		x.Text = n
	} else {
		panic("TextContainer.SetArg: expected TextStringWrapper")
	}
}
func (x *LatexCmdInput) VisualizeTree() string { return "TextContainer " + x.Text.VisualizeTree() }
