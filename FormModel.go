package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type FormModel struct {
	Name  textinput.Model
	Start textinput.Model
	End   textinput.Model
}

func NewForm() *FormModel {
	form := &FormModel{Name: textinput.New(), Start: textinput.New(), End: textinput.New()}

	form.Name.Placeholder = "Farm Name"
	form.Start.Placeholder = "IP Start"
	form.End.Placeholder = "IP End"

	form.Name.Focus()
	return form
}

func (m FormModel) Init() tea.Cmd {
	return textinput.Blink
}

func NewFarm(name, root, length string) FarmStruct {
	return FarmStruct{Name: name, Start: root, End: length}
}

func (m FormModel) CreateFarm() tea.Msg {
	farm := NewFarm(m.Name.Value(), m.Start.Value(), m.End.Value())
	return farm
}

func (m FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, textinput.Blink
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return Models[TableView].Update(nil)
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			if m.Name.Focused() {
				m.Name.Blur()
				m.Start.Focus()
			} else if m.Start.Focused() {
				m.Start.Blur()
				m.End.Focus()
			} else {
				m.End.Blur()
				m.Name.Focus()
			}
			return m, textinput.Blink
		case "enter":
			return Models[ListView], m.CreateFarm
		}
	}
	if m.Name.Focused() {
		m.Name, cmd = m.Name.Update(msg)
		return m, cmd
	} else if m.Start.Focused() {
		m.Start, cmd = m.Start.Update(msg)
		return m, cmd
	} else if m.End.Focused() {
		m.End, cmd = m.End.Update(msg)
		return m, cmd
	}

	return m, cmd
}

func (m FormModel) View() string {
	return m.Name.View() + "\n" + m.Start.View() + "\n" + m.End.View()
}
