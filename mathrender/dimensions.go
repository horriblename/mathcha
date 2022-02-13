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
	AbsY     int    // the absolute position of the parent's BaseLine in the buffer, >=0
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
		dim.Lit = GetVanillaString(n.Command())
		dim.Width = rw.StringWidth(dim.Lit)
		dim.Children = nil
		return dim
		// parser.Literal interface types
	case *parser.SimpleOpLit:
		dim.Height = 1
		switch n.Content() {
		case "+", "-", "=":
			dim.Lit = " " + n.Content() + " "
		default:
			dim.Lit = n.Content()
		}
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
	case parser.CMD_underline:
		return calculateDimCmdUnderline(node)
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

func calculateDimCmdUnderline(node parser.CmdContainer) *Dimensions {
	dim := new(Dimensions)
	dim.Children = []*Dimensions{calculateDim(node.Children()[0])}
	dim.Width = dim.Children[0].Width
	dim.Height = dim.Children[0].Height
	dim.BaseLine = dim.Children[0].BaseLine
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
