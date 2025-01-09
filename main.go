// WARNING this file is mostly for testing
package main

import (
	// tea "github.com/charmbracelet/bubbletea"
	// parser "github.com/horriblename/mathcha/latex"
	// render "github.com/horriblename/mathcha/renderer"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/derekparker/trie"
	"github.com/horriblename/mathcha/editor"
	ed "github.com/horriblename/mathcha/editor"
	"github.com/horriblename/mathcha/latex"
	"github.com/horriblename/mathcha/renderer"
)

type model struct {
	// current editor in focus
	focus        int
	editors      []ed.Editor
	compList     *trie.Trie
	compMatches  []string
	editorConfig *ed.EditorConfig
	showHelp     bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel(editorCfg ed.EditorConfig, initFormula string) model {
	editor := ed.NewWithConfig(editorCfg, initFormula)
	editor.SetFocus(true)
	return model{
		focus:        0,
		editors:      []ed.Editor{*editor}, // TODO should prolly make this slice of pointers to Editors
		compList:     latex.NewCompletion(),
		editorConfig: &editorCfg,
	}
}

var keyPressed = "[waiting key]"

func (m model) latex() string {
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

	return latex
}

func (m model) CopyLatex() {
	// wayland clipboard support: https://github.com/golang-design/clipboard/issues/6
	latex := m.latex()

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
			editor := ed.NewWithConfig(*m.editorConfig, "")
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
		case tea.KeyCtrlC: // chain tea command?
			m.CopyLatex()
			return m, tea.Quit
		case tea.KeyTab, tea.KeyShiftTab:
			if m.editors[m.focus].GetState() != ed.EDIT_COMMAND {
				break
			}
			lead := m.editors[m.focus].FocusedTextField().BuildString()
			m.compMatches = m.compList.FuzzySearch(lead)
			return m, nil

		case tea.KeyF1:
			m.showHelp = !m.showHelp
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
		"\n%s\n\n%s\n%s",
		strings.Join(editorsView, "\n"),
		compDisplay.String(),
		m.helpSection(),
	) + "\n"
}

var extendedHelp = `
Editor
------
` + editor.KeybindsHelp + `

General
-------
	F1 toggles keybinds help
	ctrl+c to quit
	ctrl+k previous line
	ctrl+j next line
	ctrl+y Copy Latex to clipboard (via wl-copy)
`

func (m model) helpSection() string {
	if !m.showHelp {
		return "press F1 to keybinds help"
	} else {
		return extendedHelp
	}
}

func main() {
	var useUnicode bool
	flag.BoolVar(&useUnicode, "symbols", false, `Use unicode symbols in latex output wherever possible. e.g. output "α" in place of "\alpha"`)
	flag.BoolVar(&useUnicode, "s", false, `Use unicode symbols in latex output wherever possible. e.g. output "α" in place of "\alpha"`)
	file := flag.String("f", "", "Read initial formula from file; use '-' to read from stdin")
	flag.Parse()

	editorCfg := ed.EditorConfig{
		LatexCfg: renderer.LatexSourceConfig{
			UseUnicode: useUnicode,
		},
	}

	var latex string
	switch *file {
	case "":
	case "-":
		l, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic("error reading stdin: " + err.Error())
		}
		latex = string(l)
	default:
		l, err := os.ReadFile(*file)
		if err != nil {
			panic("error reading " + *file + ": " + err.Error())
		}
		latex = string(l)
	}

	e := initialModel(editorCfg, latex)

	p := tea.NewProgram(e)
	if _, err := p.Run(); err != nil {
		// log error
	}
}
