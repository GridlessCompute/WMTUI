package popup

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type PopupModel struct {
	Caller   string
	Fields   []textinput.Model
	Selected int
}

type PopupResMsg struct {
	Caller string
	Values []string
}

func NewPopup(caller string, f []textinput.Model) PopupModel {
	return PopupModel{
		Caller:   caller,
		Fields:   f,
		Selected: 0,
	}
}

func NewInput(prompt, placeholder string) textinput.Model {
	txt := textinput.New()

	txt.Placeholder = placeholder
	txt.Prompt = prompt
	txt.CharLimit = 25
	txt.Width = 25

	return txt
}

func (m PopupModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m PopupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "escape", "esc", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.Selected = min(m.Selected+1, len(m.Fields)-1)
			m.Fields[m.Selected].Focus()
			return m, nil
		case "shift+tab":
			m.Selected = max(m.Selected-1, 0)
			m.Fields[m.Selected].Focus()
			return m, nil
		case "return", "enter":
			v := []string{}
			for _, f := range m.Fields {
				v = append(v, f.Value())
			}
			return m, func() tea.Msg { return PopupResMsg{Values: v} }
		}
	}

	m.Fields[m.Selected], cmd = m.Fields[m.Selected].Update(msg)
	return m, cmd
}

func (m PopupModel) View() string {
	s := "\n"

	for _, f := range m.Fields {
		s += f.View()
		s += "\n"
	}

	return s
}
