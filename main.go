// WARNING this file is mostly for testing
package main

import (
	// tea "github.com/charmbracelet/bubbletea"
	// parser "github.com/horriblename/mathcha/latex"
	// render "github.com/horriblename/mathcha/renderer"
	"flag"
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/derekparker/trie"
	ed "github.com/horriblename/mathcha/editor"
	"github.com/horriblename/mathcha/latex"
	"github.com/horriblename/mathcha/renderer"
)

type model struct {
	focus        int
	editors      []ed.Editor
	compList     *trie.Trie
	compMatches  []string
	editorConfig *ed.EditorConfig
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel(editorCfg ed.EditorConfig) model {
	editor := ed.NewWithConfig(editorCfg)
	editor.SetFocus(true)
	return model{
		focus:        0,
		editors:      []ed.Editor{*editor}, // TODO should prolly make this slice of pointers to Editors
		compList:     latex.NewCompletion(),
		editorConfig: &editorCfg,
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
			editor := ed.NewWithConfig(*m.editorConfig)
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
		case tea.KeyTab, tea.KeyShiftTab:
			if m.editors[m.focus].GetState() != ed.EDIT_COMMAND {
				break
			}
			lead := m.editors[m.focus].FocusedTextField().BuildString()
			m.compMatches = m.compList.FuzzySearch(lead)
			return m, nil

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

	var compDisplay strings.Builder
	var displayLen int
	for _, match := range m.compMatches {
		cmd := latex.MatchLatexCmd("\\" + match)
		compDisplay.WriteString("\x1b[34m")
		if cmd.IsVanillaSym() {
			compDisplay.WriteString(renderer.GetVanillaString(cmd))
		} else {
			compDisplay.WriteRune(' ')
		}
		compDisplay.WriteString(" \x1b[33m")
		compDisplay.WriteString(match)
		compDisplay.WriteString("   ")

		displayLen += len(match) + 5
		// hard line width limit
		if displayLen > 500 {
			break
		}
	}

	return fmt.Sprintf(
		"\n%s\n\n%s\n%s\n%s",
		strings.Join(editorsView, "\n"),
		m.editors[m.focus].LatexSource(),
		compDisplay.String(),
		"(esc or ctrl+c to quit | ctrl+k previous line | ctrl+j next line | ctrl+y Copy Latex to clipboard (via wl-copy))",
	) + "\n"
}

func main() {
	var useUnicode bool
	flag.BoolVar(&useUnicode, "symbols", false, `Use unicode symbols in latex output wherever possible. e.g. output "α" in place of "\alpha"`)
	flag.BoolVar(&useUnicode, "s", false, `Use unicode symbols in latex output wherever possible. e.g. output "α" in place of "\alpha"`)
	flag.Parse()

	editorCfg := ed.EditorConfig{
		LatexCfg: renderer.LatexSourceConfig{
			UseUnicode: useUnicode,
		},
	}
	e := initialModel(editorCfg)

	p := tea.NewProgram(e)
	if err := p.Start(); err != nil {
		// log error
	}
}
