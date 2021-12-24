package mathrender

import (
	parser "github.com/horriblename/latex-parser/latex"
)

type Dimensions struct {
	Width    int
	Height   int
	BaseLine int // the lowest point of the block, 0 by default, can go below negative
	AbsX     int // the absolute position in the buffer
	AbsY     int // the absolute position in the buffer, >=0
	Children []*Dimensions
}

func calculateDim(node parser.Expr) *Dimensions {
	dim := new(Dimensions)
	switch n := node.(type) {
	case parser.FlexContainer:
		hi := 1
		lo := 0
		var children = []*Dimensions{}
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
		return dim
	case parser.CmdContainer:
		return calculateDimCmdContainer(n)
	case parser.CmdLiteral:
		dim.Height = 1
		dim.Width = 1 //len(n.Content())
		return dim
	case parser.Literal:
		dim.Height = 1
		dim.Width = len(n.Content())
		return dim
	}

	// FIXME throw error
	dim.Width = 1
	dim.Height = 1
	dim.BaseLine = 0
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
	}
	return &dim
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
