package editor

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func typeKeys(e *Editor, keys ...interface{}) {
	for _, k := range keys {
		var msg tea.Msg
		switch k := k.(type) {
		case rune:
			msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{k}}
		case string:
			for _, r := range k {
				updated, _ := e.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
				*e = updated
			}
			continue
		case tea.KeyType:
			msg = tea.KeyMsg{Type: k}
		}
		updated, _ := e.Update(msg)
		*e = updated
	}
}

func TestInsertLeft(t *testing.T) {
	tests := []struct {
		desc     string
		keys     string // typed after "\left" + space
		expected string
	}{
		{
			desc:     "explicit pair",
			keys:     "{.",
			expected: "\\left\\{\\right.",
		},
		{
			desc:     "single paren auto-matches",
			keys:     "(",
			expected: "\\left(\\right)",
		},
		{
			desc:     "single bracket auto-matches",
			keys:     "[",
			expected: "\\left[\\right]",
		},
		{
			desc:     "single brace auto-matches",
			keys:     "{",
			expected: "\\left\\{\\right\\}",
		},
		{
			desc:     "bar repeats",
			keys:     "|",
			expected: "\\left|\\right|",
		},
		{
			desc:     "escaped bar repeats",
			keys:     "\\|",
			expected: "\\left\\|\\right\\|",
		},
		{
			desc:     "mismatched pair",
			keys:     "(|",
			expected: "\\left(\\right|",
		},
		{
			desc:     "empty input defaults to parens",
			keys:     "",
			expected: "\\left(\\right)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			e := New("")
			typeKeys(e, '\\', 'l', 'e', 'f', 't', tea.KeySpace)
			typeKeys(e, tc.keys)
			typeKeys(e, tea.KeySpace)

			if got := e.LatexSource(); got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}

			// cursor should be inside the parens
			state := e.GetState()
			if state != EDIT_EQUATION {
				t.Errorf("expected EDIT_EQUATION inside parens, got %v", state)
			}
		})
	}
}

func TestInsertLeftWithContent(t *testing.T) {
	e := New("")
	typeKeys(e, '\\', 'l', 'e', 'f', 't', tea.KeySpace)
	typeKeys(e, "{.")
	typeKeys(e, tea.KeySpace)
	typeKeys(e, 'x')
	typeKeys(e, tea.KeyEscape)

	expected := "\\left\\{x\\right."
	if got := e.LatexSource(); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestMatShortcut(t *testing.T) {
	e := New("")
	typeKeys(e, '\\', 'm', 'a', 't', tea.KeySpace)

	expected := "\\begin{matrix}\n\n\\end{matrix}"
	if got := e.LatexSource(); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}

	// cursor should be inside the matrix cell
	if state := e.GetState(); state != EDIT_EQUATION {
		t.Errorf("expected EDIT_EQUATION inside matrix, got %v", state)
	}
}

func selectAll(e *Editor) {
	e.NavigateToBeginning()
	for e.selectRight() {
	}
}

func TestWrapSelectionInCommand(t *testing.T) {
	tests := []struct {
		desc     string
		keys     []interface{}
		expected string
	}{
		{
			desc:     "one-argument command",
			keys:     []interface{}{'\\', 's', 'q', 'r', 't', tea.KeySpace},
			expected: `\sqrt {abc}`,
		},
		{
			desc:     "two-argument command",
			keys:     []interface{}{'\\', 'f', 'r', 'a', 'c', tea.KeySpace},
			expected: `\frac {abc}{}`,
		},
		{
			desc:     "matrix shortcut",
			keys:     []interface{}{'\\', 'm', 'a', 't', tea.KeySpace},
			expected: "\\begin{matrix}\nabc\n\\end{matrix}",
		},
		{
			desc:     "begin with environment name",
			keys:     []interface{}{'\\', 'b', 'e', 'g', 'i', 'n', tea.KeySpace, 'm', 'a', 't', 'r', 'i', 'x', tea.KeySpace},
			expected: "\\begin{matrix}\nabc\n\\end{matrix}",
		},
		{
			desc:     "left/right delimiters",
			keys:     []interface{}{'\\', 'l', 'e', 'f', 't', tea.KeySpace, '(', tea.KeySpace},
			expected: `\left(abc\right)`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			e := New("")
			typeKeys(e, "abc")
			selectAll(e)
			typeKeys(e, tc.keys...)

			if got := e.LatexSource(); got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}
