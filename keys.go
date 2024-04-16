package main

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit key.Binding
	Up   key.Binding
	Down key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "l"),
		key.WithHelp("↑/l", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "k"),
		key.WithHelp("↓/k", "move down"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc", "quit"),
	),
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Quit, k.Up, k.Down}}
}
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Up, k.Down}
}
