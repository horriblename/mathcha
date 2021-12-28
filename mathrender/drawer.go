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

// func blockConcat(blocks [][][]string, levels []int) [][]string {
// 	if len(blocks) != len(levels) {
// 		//FIXME proper error handling
// 		println("Error in blockConcat(): number of blocks and levels don't match")
// 	}
// 	if len(blocks) == 1 {
// 		// TODO perhaps it would be better to panic when len(blocks) == 0
// 		return blocks[0]
// 	}

// 	ret := blocks[0]
// 	baseLv := levels[0] // levels count upwards, and can go below 0
// 	height := len(ret)

// 	// first loop checks the vertical size of the blocks to allocate adequate space
// 	// FIXME we're not using arrays here, so this might not be useful
// 	for i := 1; i < len(blocks); i++ {
// 		low := levels[i]           // lower limit of current block
// 		up := low + len(blocks[i]) // upper limit of current block
// 		if low < baseLv {
// 			baseLv = low
// 		}
// 		if up > height+baseLv {
// 			height = up - baseLv
// 		}
// 	}

// 	for i := 1; i < len(blocks); i++ {
// 		block := blocks[i]
// 		bh := len(block) // height of block
// 		lv := levels[i]
// 		startFrom := height - bh - (lv - baseLv) // row index to start from

// 		if height > bh {
// 			// grow upwards by bh - len(ret)
// 			w := lineWidth(block)
// 			padding := strings.Repeat("-", w)

// 			for j := 0; j < startFrom; j++ {
// 				ret[j] = append(ret[j], padding)
// 			}
// 			for j := 0; j < height+lv-1; j++ {
// 				ret[j] = append(ret[j], padding)
// 			}
// 		}

// 		for j := 0; j < bh; j++ {
// 			ret[j+startFrom] = append(ret[j+startFrom], block[j]...)
// 		}

// 	}
// 	return ret
// }

func lineWidth(line []rune) int {
	w := 0
	// for i := range line {
	// 	w += len(line)
	// }
	return w
}
