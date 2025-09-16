package popup

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()
)

type PopupModel struct {
	Caller     string
	Fields     []textinput.Model
	Selected   int
	CursorMode cursor.Mode
}

type PopupResMsg struct {
	Caller string
	Values []string
}

type PopupCloseMsg struct{}

func NewPopup(caller string, f []textinput.Model) PopupModel {
	p := PopupModel{
		Caller:   caller,
		Fields:   f,
		Selected: 0,
	}

	if len(p.Fields) > 0 {
		p.Fields[0].Focus()
		p.Fields[0].PromptStyle = focusedStyle
		p.Fields[0].TextStyle = focusedStyle
	}

	return p
}

func NewInput(prompt, initial string) textinput.Model {
	txt := textinput.New()

	txt.Prompt = prompt
	txt.Cursor.Style = cursorStyle
	txt.PromptStyle = noStyle
	txt.TextStyle = noStyle
	txt.CharLimit = 25
	txt.Width = 25
	txt.SetValue(initial)

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
			return m, func() tea.Msg { return PopupCloseMsg{} }
		case "tab":
			m.Fields[m.Selected].Blur()
			m.Fields[m.Selected].PromptStyle = noStyle
			m.Fields[m.Selected].TextStyle = noStyle

			m.Selected = min(m.Selected+1, len(m.Fields)-1)

			m.Fields[m.Selected].Focus()
			m.Fields[m.Selected].PromptStyle = focusedStyle
			m.Fields[m.Selected].TextStyle = focusedStyle

			return m, nil
		case "shift+tab":
			m.Fields[m.Selected].Blur()
			m.Fields[m.Selected].PromptStyle = noStyle
			m.Fields[m.Selected].TextStyle = noStyle

			m.Selected = max(m.Selected-1, 0)

			m.Fields[m.Selected].Focus()
			m.Fields[m.Selected].PromptStyle = focusedStyle
			m.Fields[m.Selected].TextStyle = focusedStyle

			return m, nil
		case "return", "enter":
			v := []string{}
			for _, f := range m.Fields {
				v = append(v, f.Value())
			}
			return m, func() tea.Msg { return PopupResMsg{Values: v} }
		}
	}

	cmd = m.updateInputs(msg)
	return m, cmd
}

func (m PopupModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Fields))

	for i := range m.Fields {
		m.Fields[i], cmds[i] = m.Fields[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m PopupModel) View() string {
	s := "\n"

	for _, f := range m.Fields {
		s += f.View()
		s += "\n"
	}

	return s
}
