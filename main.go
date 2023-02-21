// WARNING this file is mostly for testing
package main

import (
	// tea "github.com/charmbracelet/bubbletea"
	// parser "github.com/horriblename/mathcha/latex"
	// render "github.com/horriblename/mathcha/renderer"
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	ed "github.com/horriblename/mathcha/editor"
)

type model struct {
	focus   int
	editors []ed.Editor
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel() model {
	editor := ed.New()
	editor.SetFocus(true)
	return model{
		focus:   0,
		editors: []ed.Editor{*editor}, // TODO should prolly make this slice of pointers to Editors
	}
}

var keyPressed = "[waiting key]"

func (m model) CopyLatex() {
	// wayland clipboard support: https://github.com/golang-design/clipboard/issues/6
	var latex string
	if len(m.editors) == 1 {
		latex = m.editors[0].LatexSource()
	} else {
		latex = "\\begin{aligned}\n"
		for _, editor := range m.editors {
			latex += editor.LatexSource() + "\\\\\n"
		}
		latex += "\\end{aligned}"
	}

	cmd := exec.Command("wl-copy")
	cmd.Stdin = strings.NewReader(latex)
	cmd.Run()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			editor := ed.New()
			m.editors = append(m.editors, *editor)
			m.editors[m.focus].SetFocus(false)
			m.editors[m.focus], cmd = m.editors[m.focus].Update(msg)
			m.focus = len(m.editors) - 1
			m.editors[m.focus].SetFocus(true)
		case tea.KeyCtrlK:
			if m.focus > 0 {
				m.editors[m.focus].SetFocus(false)
				m.editors[m.focus], cmd = m.editors[m.focus].Update(msg)
				m.focus -= 1
				m.editors[m.focus].SetFocus(true)
			}
		case tea.KeyCtrlJ:
			if m.focus < len(m.editors)-1 {
				m.editors[m.focus].SetFocus(false)
				m.editors[m.focus], cmd = m.editors[m.focus].Update(msg)
				m.focus += 1
				m.editors[m.focus].SetFocus(true)
			}
		case tea.KeyCtrlY:
			m.CopyLatex()
			return m, nil
		case tea.KeyCtrlC, tea.KeyEsc: // chain tea command?
			m.CopyLatex()
			return m, tea.Quit
		default:
			keyPressed = msg.String()
		}

		// We handle errors just like any other message
		// case errMsg:
		// 	m.err = msg
		// 	return m, nil
	}

	m.editors[m.focus], cmd = m.editors[m.focus].Update(msg)
	return m, cmd
}

func (m model) View() string {
	editorsView := make([]string, len(m.editors))
	for _, editor := range m.editors {
		editorsView = append(editorsView, editor.View())
	}

	return fmt.Sprintf(
		"\n%s\n\n%s\n\x1b[34mKey:\x1b[33m %s\x1b[34m pressed\x1b[0m\n%s",
		strings.Join(editorsView, "\n"),
		m.editors[m.focus].LatexSource(),
		keyPressed,
		"(esc or ctrl+c to quit | ctrl+k previous line | ctrl+j next line | ctrl+y Copy Latex to clipboard (via wl-copy))",
	) + "\n"
}

func main() {
	// var latex string
	// latex = `1-\frac{1}{1-\frac{1}{1-\frac{1}{1-\frac{1}{1-\frac{1}{1-\frac{1}{1-\frac{1}{1-\frac{1}{1-\frac{1}{1-\frac{1}{2}}}}}}}}}}`
	// latex = `E = \frac\underline{12}12 mv2`
	// latex = `g(x) = \left(\frac{12}{13}\right)`
	// latex = `xyz = \text{this is a text}abc-{}+1`
	//
	// latex = `f(x) = \frac{1}{\sigma\sqrt{2\pi}}\exp\left(-\frac{1}{2}\left(\frac{x-\mu}{\sigma}\right)\right)`
	// latex = ""
	e := initialModel()
	// e.editors[e.focus].Read(latex)

	p := tea.NewProgram(e)
	if err := p.Start(); err != nil {
		// log error
	}
}
