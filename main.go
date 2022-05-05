package main

import (
	// tea "github.com/charmbracelet/bubbletea"
	// parser "github.com/horriblename/mathcha/latex"
	// render "github.com/horriblename/mathcha/renderer"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	ed "github.com/horriblename/mathcha/editor"
)

type model struct {
	editor ed.Editor
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel() model {
	return model{*ed.New()}
}

var keyPressed = "[waiting key]"

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		default:
			keyPressed = msg.String()
		}

		// We handle errors just like any other message
		// case errMsg:
		// 	m.err = msg
		// 	return m, nil
	}

	m.editor, cmd = m.editor.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"\n%s\n\n%s\n\x1b[34mKey:\x1b[33m %s\x1b[34m pressed\x1b[0m\n%s",
		m.editor.View(),
		m.editor.LatexSource(),
		keyPressed,
		"(esc or ctrl+c to quit)",
	) + "\n"
}

func main() {
	var latex string
	latex = `1-\frac{1}{1-\frac{1}{1-\frac{1}{1-\frac{1}{1-\frac{1}{1-\frac{1}{1-\frac{1}{1-\frac{1}{1-\frac{1}{1-\frac{1}{2}}}}}}}}}}`
	latex = `E = \frac\underline{12}12 mv2`
	latex = `g(x) = \left(\frac{12}{13}\right)`
	latex = `xyz = \text{this is a text}abc-{}+1`

	latex = `f(x) = \frac{1}{\sigma\sqrt{2\pi}}\exp\left(-\frac{1}{2}\left(\frac{x-\mu}{\sigma}\right)\right)`
	latex = "f(x)"
	e := initialModel()
	e.editor.Read(latex)

	p := tea.NewProgram(e)
	if err := p.Start(); err != nil {
		// log error
	}
}
