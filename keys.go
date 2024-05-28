package main

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Switch  key.Binding
	Reboot  key.Binding
	Sleep   key.Binding
	Wake    key.Binding
	Pools   key.Binding
	Limit   key.Binding
	Fast    key.Binding
	Slow    key.Binding
	Help    key.Binding
	Refresh key.Binding
	Select  key.Binding
	Sort    key.Binding
	Farm    key.Binding
	NewFarm key.Binding
	Quit    key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Help, k.Refresh, k.Quit}, // first column
		{k.Sleep, k.Wake, k.Fast, k.Slow, k.Sort}, // second column
		{k.Reboot, k.Pools, k.Limit, k.Select},    // Third column
		{k.Farm, k.NewFarm},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "Move Up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "Move Down"),
	),
	Switch: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "Switch to Views"),
	),
	Reboot: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Reboot Miner(s)"),
	),
	Sleep: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Sleep Miner(s)"),
	),
	Wake: key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "Wake Miner(s)"),
	),
	Pools: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Set Pools"),
	),
	Limit: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "Set Power Limit"),
	),
	Fast: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "Set Fast Boot"),
	),
	Slow: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "Set Slow Boot"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Select/Deselect Miner"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "Refresh Miners"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Sort: key.NewBinding(
		key.WithKeys("`"),
		key.WithHelp("`", "Sort"),
	),
	Farm: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "select farm"),
	),
	NewFarm: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new farm"),
	),
}
