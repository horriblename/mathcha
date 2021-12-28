package mathrender

import (
	"fmt"
	parser "github.com/horriblename/latex-parser/latex"
)

// breadth-first traverse of the latex tree and dim tree in parallel
// to build the later rendered 2D rune buffer
func (r *Renderer) DrawToBuffer(tree parser.Expr, dim *Dimensions) {
	treeQueue := []parser.Expr{tree}
	dimQueue := []*Dimensions{dim}
	cumulX := dim.AbsX                 // cumulative sum of width of previous siblings
	y := dim.Height + dim.BaseLine - 1 // absolute position of the "baseline relative to parent node"
	dim.AbsY = y

	for len(treeQueue) > 0 {
		//fmt.Println("loop ", treeQueue[0].VisualizeTree())
		if dimQueue[0] == nil { // use nil to mark a new set of siblings ahead
			fmt.Println("width ", dimQueue[1].Width)
			dimQueue = dimQueue[1:]
			y = dimQueue[0].AbsY // FIXME not sure..
			cumulX = dimQueue[0].AbsX
		}
		t := treeQueue[0]
		d := dimQueue[0]
		r.Prerender(t, d, cumulX, y)

		switch n := t.(type) {
		case parser.FixedContainer:
			treeQueue = append(treeQueue, n.Children()...)
			for _, c := range d.Children {
				//c.AbsX = cumulX
				//c.AbsY = y // FIXME not sure..
				dimQueue = append(dimQueue, nil) // mark new set of siblings
				dimQueue = append(dimQueue, c)
			}

		case parser.FlexContainer:
			treeQueue = append(treeQueue, n.Children()...)
			dimQueue = append(dimQueue, nil) // mark new set of siblings
			dimQueue = append(dimQueue, d.Children...)

			// btw this is the only way i can think of to pass the absolute
			// positions to the children without backtracking
			for _, c := range d.Children {
				c.AbsX = cumulX
				c.AbsY = y // FIXME not sure..
			}
		}

		cumulX += d.Width
		treeQueue = treeQueue[1:]
		dimQueue = dimQueue[1:]
	}
}

// x is the position of leftmost rune allowed to be written by the node
// y is the position of baseline=0 of the node
func (r *Renderer) Prerender(node parser.Expr, dim *Dimensions, x int, y int) {
	switch n := node.(type) {
	case parser.CmdContainer:
		switch n.Command() {
		case parser.CMD_frac:
			r.PrerenderCmdFrac(n, dim, x, y)
		}
	case parser.Container:
		r.PrerenderContainer(n, dim, x, y)
	// case parser.CmdLiteral:
	// 	content := n.Content() // TODO map to unicode
	// 	println("rendered CmdLiteral: ", content)
	// 	// parser.Literal interface types
	case parser.Literal:
		// for i := 0; i < dim.Width && i < len(content); i++ {
		i := 0
		for _, char := range dim.Lit {
			// FIXME ^runes cannot be looped normally
			r.Buffer[y][x+i] = rune(char)
			i++
		}
		println("rendered Literal: ", dim.Lit)
		// case *parser.SimpleOpLit:
		// 	content := n.Content()
		// 	for i := 0; i < dim.Width && i < len(content)-2; i++ {
		// 		r.Buffer[y][x+i+1] = rune(content[i])
		// 	}
		// 	println("rendered Literal: ", content)
		// case *parser.NumberLit:
		// 	content := n.Content()
		// 	for i := 0; i < dim.Width && i < len(content); i++ {
		// 		r.Buffer[y][x+i] = rune(content[i])
		// 	}
		// 	println("rendered Literal: ", content)
		// case *parser.VarLit:
		// 	content := n.Content()
		// 	for i := 0; i < dim.Width && i < len(content); i++ {
		// 		r.Buffer[y][x+i] = rune(content[i])
		// 	}
		// 	println("rendered Literal: ", content)
	}
}

func (r *Renderer) PrerenderContainer(node parser.Container, dim *Dimensions, x int, y int) {

	// TODO
	return
}

func (r *Renderer) PrerenderCmdContainer(node parser.CmdContainer, dim *Dimensions, x int, y int) {
	switch node.Command() {
	case parser.CMD_frac:
		r.PrerenderCmdFrac(node, dim, x, y) // TODO x y
	}

	return
}

func (r *Renderer) PrerenderCmdFrac(node parser.CmdContainer, dim *Dimensions, x int, y int) {
	l := rune('─') // \u2500
	for i := 0; i < dim.Width; i++ {
		r.Buffer[y][x+i] = l
	}

	arg1 := dim.Children[0]
	arg2 := dim.Children[1]
	arg1.AbsX = x
	arg2.AbsX = x
	arg1.AbsY = y + arg1.BaseLine - 1
	arg2.AbsY = y + arg2.Height + arg2.BaseLine //FIXME minus arg1.baseline
}

func lineWidth(line []rune) int {
	w := 0
	// for i := range line {
	// 	w += len(line)
	// }
	return w
}
