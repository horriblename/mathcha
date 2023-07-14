# mathcha
A [mathquill](http://mathquill.com/)-like interactive math input in a terminal, built with the charm TUI library.

[![asciicast](https://asciinema.org/a/rph5qHyOpcjpQ4pvQJUFZqUDI.svg)](https://asciinema.org/a/rph5qHyOpcjpQ4pvQJUFZqUDI)

## Usage

Keybinds:

- Arrow keys/ `Ctrl+b/f/n/p` for basic cursor navigation
- `Alt` + left/right to start or extend selection
  - in selection mode, parenthesis `(`/`)` and divide `/` keys will wrap the selected block in the corresponding command
- `Tab` to go out a block
- Press `\` to enter a `\command` (e.g. `\alpha` or `\frac`), when you're done, hit `Space`. While entering a command, you can hit `Tab` to see a list of available commands (there's no autocomplete, you still have to type it out yourself)
- `Enter` for a new equation in a new line
- `Ctrl+k` to go to previous line, `Ctrl+j` to go to next line

## Supported Symbols and Commands
There is no standard support table or even a goal, if I ever feel like turning this into a serious project, I would start from KaTeX, but that probably won't happen :P.

Take a look at [command_def.go](latex/command_def.go)for a list of all recognized latex commands.

---

I started this project when I was learning parsing, and this remains as a toy project in which I add stuff when I feel like doing so. Feature requests and whatnot is welcome, but I probably won't do anything.

Other interesting stuff I might try in the future:
- reactive-programming styled rendering
- autocomplete
- go-specific memory optimizations
