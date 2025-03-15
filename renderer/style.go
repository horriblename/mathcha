package renderer

import "fmt"

type color int

const ()

type style int

const (
	underline style = 4
	italic    style = 3
)
const styleResetOffset int = 20

func (self *Renderer) backgroundRGB(r int, g int, b int) string {
	if !self.Color {
		return ""
	}
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r, g, b)
}

// wraps content string with style + reset
func (self *Renderer) styleAndReset(s style, content string) string {
	if !self.Color {
		return content
	}
	return fmt.Sprintf("\x1b[%dm%s\x1b[%dm", int(s), content, int(s)+styleResetOffset)
}

func (self *Renderer) overlineAndReset(content string) string {
	if !self.Color {
		return content
	}
	return fmt.Sprintf("\x1b[53m%s\x1b[55m", content)
}
