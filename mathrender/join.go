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
		alignAt        int

		// The vertical level of which the blocks should be aligned at
		highest int
		lowest  int
	)

	// Break text blocks into lines and get max widths for each text block
	for i, str := range strs {
		blocks[i], maxWidths[i] = getLines(str)
		lo := baseline[i]
		hi := len(blocks[i]) + lo
		if hi > highest {
			highest = hi
		}
		if lo < lowest {
			lowest = lo
		}
	}

	combinedHeight = highest - lowest
	alignAt = highest - 1
	// Add extra lines to make each side the same height
	for i := range blocks {
		if len(blocks[i]) >= combinedHeight {
			continue
		}

		extraLines := make([]string, combinedHeight-len(blocks[i]))

		n := len(extraLines)
		// (combinedHeight + baseline[i])
		split := alignAt //- len(blocks[i]) + 2
		top := n - split
		bottom := n - top

		blocks[i] = append(extraLines[top:], blocks[i]...)
		blocks[i] = append(blocks[i], extraLines[bottom:]...)
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
