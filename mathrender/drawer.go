package mathrender

import (
	"fmt"
	parser "github.com/horriblename/latex-parser/latex"
)

func (r *Renderer) Visit(node parser.Expr, dim *Dimensions) Visitor {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case parser.Container:
		r.PrerenderContainer(n, dim, r.DrawerX, r.DrawerY)
		return r
	}

	r.Prerender(node, dim, r.DrawerX, r.DrawerY)
	return nil
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
	case parser.CmdLiteral:
		content := n.Content() // TODO map to unicode
		println("rendered CmdLiteral: ", content)
	case parser.Literal:
		content := n.Content()
		println("rendered Literal: ", content)
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
	fmt.Println("PrerenderCmdFrac...")

	y = y + dim.Children[1].Height
	a := rune('-')
	for i := range r.Buffer[y][x:dim.Width] {
		r.Buffer[y][i] = a
	}
}

func View(block [][]string) {
	for _, line := range block {
		l := ""
		for _, col := range line {
			l += col
		}
		fmt.Println(l)
	}
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
