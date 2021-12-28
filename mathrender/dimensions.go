package mathrender

import (
	parser "github.com/horriblename/latex-parser/latex"
	rw "github.com/mattn/go-runewidth"
)

type Dimensions struct {
	Width    int
	Height   int
	BaseLine int    // the lowest point of the block, 0 by default, can go below negative
	AbsX     int    // the absolute position in the buffer
	AbsY     int    // the absolute position in the buffer, >=0
	Lit      string // a single line string to be drawn into the buffer if it is a Literal, else ""
	Children []*Dimensions
}

func calculateDim(node parser.Expr) *Dimensions {
	dim := new(Dimensions)
	switch n := node.(type) {
	case parser.FlexContainer:
		hi := 1
		lo := 0
		var children = []*Dimensions{} // make([]*Dimensions, len(n.Children())) FIXME y doesn't make work here??
		for _, i := range n.Children() {
			child := calculateDim(i)
			dim.Width += child.Width
			if child.BaseLine+child.Height > hi {
				hi = child.BaseLine + child.Height
			}
			if child.BaseLine < lo {
				lo = child.BaseLine
			}
			children = append(children, child)
		}
		dim.Height = hi - lo
		dim.BaseLine = lo
		dim.Children = children
		return dim
	case parser.CmdContainer:
		return calculateDimCmdContainer(n)
	case parser.CmdLiteral:
		dim.Height = 1
		dim.Width = 1 //len(n.Content())
		dim.Children = nil
		dim.Lit = string(GetVanillaRune(n.Command()))
		dim.Width = rw.StringWidth(dim.Lit)
		println("CmdLiteral and width: ", dim.Lit, dim.Width)
		println(n.Command().GetCmd())
		return dim
		// parser.Literal interface types
	case *parser.SimpleOpLit:
		dim.Height = 1
		dim.Lit = " " + n.Content() + " "
		dim.Width = rw.StringWidth(dim.Lit)
		dim.Children = nil
		return dim
	case *parser.NumberLit:
		dim.Height = 1
		dim.Lit = n.Content()
		dim.Width = rw.StringWidth(dim.Lit)
		dim.Children = nil
		return dim
	case *parser.VarLit:
		dim.Height = 1
		dim.Lit = n.Content()
		dim.Width = rw.StringWidth(dim.Lit)
		dim.Children = nil
		return dim
	}

	// FIXME throw error
	dim.Width = 1
	dim.Height = 1
	dim.BaseLine = 0
	dim.Children = nil
	return dim
}

func calculateDimCmdContainer(node parser.CmdContainer) *Dimensions {
	switch node.Command() {
	case parser.CMD_frac:
		return calculateDimCmdFrac(node)
	}

	// FIXME
	dim := new(Dimensions)
	dim.Width = 1
	dim.Height = 1
	dim.BaseLine = 0
	dim.Children = nil
	return dim
}

func calculateDimCmdFrac(node parser.CmdContainer) *Dimensions {
	// TODO assert number of children somewhere?
	dimArg1 := calculateDim(node.Children()[0])
	dimArg2 := calculateDim(node.Children()[1])

	dim := Dimensions{
		Width:    max(dimArg1.Width, dimArg2.Width),
		Height:   dimArg1.Height + dimArg2.Height + 1,
		BaseLine: -dimArg2.Height,
		Children: []*Dimensions{dimArg1, dimArg2},
	}
	return &dim
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
