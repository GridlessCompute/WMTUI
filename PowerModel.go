package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type PowerModel struct {
	Power textinput.Model
}

type PowerStruct struct {
	Power string
}

func NewPower() *PowerModel {
	form := &PowerModel{Power: textinput.New()}

	form.Power.Placeholder = "Power Limit: 500 - 3500"

	form.Power.Focus()
	return form
}

func (m PowerModel) Init() tea.Cmd {
	return textinput.Blink
}

func NewPowerLimit(powerLimit string) PowerStruct {
	return PowerStruct{Power: powerLimit}
}

func (m PowerModel) CreatePowerLimit() tea.Msg {
	pwr := NewPowerLimit(m.Power.Value())
	return pwr
}

func (m PowerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "enter":
			return Models[MainView], m.CreatePowerLimit
		}
	}
	if m.Power.Focused() {
		m.Power, cmd = m.Power.Update(msg)
		return m, cmd
	} else {
		return m, cmd
	}

}

func (m PowerModel) View() string {
	return m.Power.View()
}
