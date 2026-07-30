// WARNING this file is mostly for testing
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/derekparker/trie"
	"github.com/horriblename/mathcha/editor"
	"github.com/horriblename/mathcha/latex"
	"github.com/horriblename/mathcha/renderer"
)

const (
	flagHelpSymbols     = "Use unicode symbols in output wherever possible"
	flagHelpSuperscript = "Use unicode superscript/subscript characters for simple scripts"
	flagHelpFile        = "Read initial formula from file; '-' for stdin"
	flagHelpHelptext    = "Help text to print below the editor"
	flagHelpPrintout    = "Internal flag for communicating with the nvim plugin"
	flagHelpLogfile     = "Print debug logs to file"
	flagHelpDebugtree   = "Print AST representation"

	flagHelpAliasSymbols     = "Alias to -symbols"
	flagHelpAliasSuperscript = "Alias to -superscript"
)

type model struct {
	cliFlags

	// current editor in focus
	focus        int
	editors      []editor.Editor
	compList     *trie.Trie
	compMatches  []string
	editorConfig *editor.EditorConfig
	showHelp     bool
}

// some CLI flags are not present here cuz they don't matter to model init
type cliFlags struct {
	helpText  *string
	printOut  *bool
	logFile   *string
	debugTree *bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func initialModel(c cliFlags, editorCfg editor.EditorConfig, initFormula string) model {
	e := editor.NewWithConfig(editorCfg, initFormula)
	e.SetFocus(true)
	return model{
		cliFlags:     c,
		focus:        0,
		editors:      []editor.Editor{*e}, // TODO should prolly make this slice of pointers to Editors
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
		latex = `\begin{aligned}` + "\n"
		for _, e := range m.editors {
			latex += e.LatexSource() + `\\` + "\n"
		}
		latex += `\end{aligned}`
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
		// case tea.KeyEnter:
		// 	editor := editor.NewWithConfig(*m.editorConfig, "")
		// 	m.editors = append(m.editors, *editor)
		// 	m.editors[m.focus].SetFocus(false)
		// 	m.editors[m.focus], cmd = m.editors[m.focus].Update(msg)
		// 	m.focus = len(m.editors) - 1
		// 	m.editors[m.focus].SetFocus(true)
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
			if m.editors[m.focus].GetState() != editor.EDIT_COMMAND {
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
	editorsView := make([]string, 0, len(m.editors))
	for _, e := range m.editors {
		editorsView = append(editorsView, e.View())
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

	tree := ""
	if *m.debugTree {
		tree = m.editors[0].Renderer().LatexTree.VisualizeTree()
	}

	return fmt.Sprintf(
		"\n%s\n\n%s\n%s\n%s",
		strings.Join(editorsView, "\n"),
		compDisplay.String(),
		m.helpSection(),
		tree,
	) + "\n"
}

func logf(s string, args ...interface{}) error {
	_, err := fmt.Fprintf(os.Stderr, s, args...)
	return err
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

const defaultHelpText = "press F1 to keybinds help"

func (m model) helpSection() string {
	if !m.showHelp {
		return *m.helpText
	} else {
		return extendedHelp
	}
}

func readFormula(fileFlag string, positionalArgs []string) string {
	switch {
	case fileFlag == "-":
		l, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic("error reading stdin: " + err.Error())
		}
		return string(l)
	case fileFlag != "":
		l, err := os.ReadFile(fileFlag)
		if err != nil {
			panic("error reading " + fileFlag + ": " + err.Error())
		}
		return string(l)
	case len(positionalArgs) > 0:
		return strings.Join(positionalArgs, " ")
	default:
		return ""
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
  mathcha [edit] [flags] [formula]   Open the TUI equation editor
  mathcha render [flags] [formula]   Render equation and print to stdout
  mathcha help                       Print this help message
`)
}

func runEdit(args []string) {
	fs := flag.NewFlagSet("edit", flag.ExitOnError)

	var useUnicode bool
	var file string
	var helptext string
	var printout bool
	var logfile string
	var debugtree bool

	fs.BoolVar(&useUnicode, "symbols", false, flagHelpSymbols)
	fs.BoolVar(&useUnicode, "s", false, flagHelpAliasSymbols)
	fs.StringVar(&file, "f", "", flagHelpFile)
	fs.StringVar(&helptext, "helptext", defaultHelpText, flagHelpHelptext)
	fs.BoolVar(&printout, "printout", false, flagHelpPrintout)
	fs.StringVar(&logfile, "logfile", "", flagHelpLogfile)
	fs.BoolVar(&debugtree, "debugtree", false, flagHelpDebugtree)

	err := fs.Parse(args)
	if err != nil {
		os.Exit(1)
	}

	editorCfg := editor.EditorConfig{
		LatexCfg: renderer.LatexSourceConfig{
			UseUnicode:         useUnicode,
			UnicodeSuperscript: false, // editor doesn't handle editing unicode scripts
		},
	}

	if logfile != "" {
		f, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			println("could not create/open log file:", err.Error())
			return
		}
		defer f.Close()
		editorCfg.Logger = log.New(f, "", log.LstdFlags)
	}

	formula := readFormula(file, fs.Args())

	e := initialModel(cliFlags{
		helpText:  &helptext,
		printOut:  &printout,
		logFile:   &logfile,
		debugTree: &debugtree,
	}, editorCfg, formula)

	p := tea.NewProgram(e,
		tea.WithInputTTY(),
		tea.WithOutput(os.Stderr),
		tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		// log error
	}

	if *e.printOut {
		// cursed magic string
		fmt.Fprint(os.Stdout, "!mAtHcHa!", e.latex())
		os.Stdout.Close()
	}
}

func runRender(args []string) {
	fs := flag.NewFlagSet("render", flag.ExitOnError)

	var useUnicode bool
	var useUnicodeSuperscript bool
	var file string

	fs.BoolVar(&useUnicode, "symbols", false, flagHelpSymbols)
	fs.BoolVar(&useUnicode, "s", false, flagHelpAliasSymbols)
	fs.BoolVar(&useUnicodeSuperscript, "superscript", true, flagHelpSuperscript)
	fs.BoolVar(&useUnicodeSuperscript, "S", true, flagHelpAliasSuperscript)
	fs.StringVar(&file, "f", "", flagHelpFile)

	err := fs.Parse(args)
	if err != nil {
		os.Exit(1)
	}

	formula := readFormula(file, fs.Args())
	if formula == "" {
		l, err := io.ReadAll(os.Stdin)
		if err != nil {
			logf("error reading stdin: %s", err.Error())
			os.Exit(1)
		}
		formula = string(l)
	}

	r := renderer.FromFormula(formula, false, useUnicodeSuperscript)
	r.Sync(nil, false)
	fmt.Print(r.Buffer)
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		runEdit(nil)
		return
	}

	switch args[0] {
	case "edit":
		runEdit(args[1:])
	case "render":
		runRender(args[1:])
	case "help", "-h", "--help":
		usage()
	default:
		runEdit(args)
	}
}
