package mathrender

import (
	"strings"

	"github.com/muesli/reflow/ansi"
)

// A modified version of https://github.com/charmbracelet/lipgloss/blob/master/join.go

// JoinHorizontal is a utility function for horizontally joining two
// potentially multi-lined strings along a vertical axis. The first argument is
// the position, with 0 being all the way at the top and 1 being all the way
// at the bottom.
//
// the parameter baseLine should be passed an array of position of each block's
// vertical position of baseline
//
// Example:
//
//     blockB := "...\n...\n..."
//     blockA := "...\n...\n...\n...\n..."
//
//     // Join 20% from the top
//     str := lipgloss.JoinHorizontal(0.2, blockA, blockB)
//
//     // Join on the top edge
//     str := lipgloss.JoinHorizontal(lipgloss.Top, blockA, blockB)
//

// for the sake of clarity, level is an integer that represents how much lower/higher a block should
// be drawn compared to normal/default blocks (level 0). Blocks usually align at level 0, superscripts
// align at level 1, subscripts and fractions align at level -1 or below (depending on the block height)
//    index    level    block
//    0        1             y      1
//    1        0        3 + x  = ───────
//    2        -1                1     2
//    3        -2                ─ + xy
//    4        -3                2
type level int

func JoinHorizontal(baseline []int, strs ...string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	if len(baseline) != len(strs) {
		panic("JoinHorizontal found different numbers of baseline and input strings")
		// FIXME don't panic
	}

	var (
		// Groups of strings broken into multiple lines
		blocks = make([][]string, len(strs))

		// Max line widths for the above text blocks
		maxWidths = make([]int, len(strs))

		// Height of the combined block
		combinedHeight int

		// highest and lowest point (level) among all blocks
		highest level
		lowest  level
	)

	// Break text blocks into lines and get max widths for each text block
	for i, str := range strs {
		blocks[i], maxWidths[i] = getLines(str)
		lo := level(baseline[i])
		hi := level(len(blocks[i])) + lo - 1
		if hi > highest {
			highest = hi
		}
		if lo < lowest {
			lowest = lo
		}
	}

	// debug
	// const filename = "log"
	// f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	// if err != nil {
	// 	panic(err)
	// }

	// defer f.Close()

	// if _, err = fmt.Fprintf(f, "Horizontal Join lowest: %d, highest: %d\n", int(lowest), highest); err != nil {
	// 	panic(err)
	// }
	// debug end

	combinedHeight = int(highest-lowest) + 1
	// Add extra lines to make each side the same height
	for i := range blocks {
		if len(blocks[i]) >= combinedHeight {
			continue
		}

		extraLines := make([]string, combinedHeight-len(blocks[i]))

		// n := len(extraLines)
		// (combinedHeight + baseline[i])
		bottom := baseline[i] - int(lowest)
		// top := n - bottom

		blocks[i] = append(extraLines[bottom:], blocks[i]...)
		blocks[i] = append(blocks[i], extraLines[:bottom]...)
	}

	// Merge lines
	var b strings.Builder
	for i := range blocks[0] { // remember, all blocks have the same number of members now
		for j, block := range blocks {
			b.WriteString(block[i])

			// Also make lines the same length
			b.WriteString(strings.Repeat(" ", maxWidths[j]-ansi.PrintableRuneWidth(block[i])))
		}
		if i < len(blocks[0])-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}

// Split a string into lines, additionally returning the size of the widest
// line.
func getLines(s string) (lines []string, widest int) {
	lines = strings.Split(s, "\n")

	for _, l := range lines {
		w := ansi.PrintableRuneWidth(l)
		if widest < w {
			widest = w
		}
	}

	return lines, widest
}
