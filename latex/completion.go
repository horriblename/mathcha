package latex

import (
	"github.com/derekparker/trie"
)

func NewCompletion() *trie.Trie {
	compList := trie.New()
	for k, cmd := range latexCmds {
		if k[0] != '\\' {
			continue
		}
		compList.Add(k[1:], cmd)
	}

	return compList
}
