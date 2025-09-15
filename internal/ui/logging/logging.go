package logging

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type LoggingMsg struct {
	Err     bool
	Message string
}

type ViewModel struct {
	Content  string
	Ready    bool
	Viewport viewport.Model
}

func NewLogging() ViewModel {
	return ViewModel{
		Content:  "",
		Ready:    false,
		Viewport: viewport.New(1, 2),
	}
}

func (m ViewModel) Init() tea.Cmd {
	return nil
}

func (m ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.Ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.Viewport = viewport.New(msg.Width, (msg.Height-verticalMarginHeight)/4)
			m.Viewport.YPosition = headerHeight
			m.Viewport.SetContent(m.Content)
			m.Ready = true
		} else {
			m.Viewport.Width = msg.Width
			m.Viewport.Height = msg.Height - verticalMarginHeight
		}

	case LoggingMsg:
		m.Content += "\n"
		m.Content += msg.Message
		m.Viewport.SetContent(m.Content)
	}

	// Handle keyboard and mouse events in the viewport
	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ViewModel) View() string {
	if !m.Ready {
		return "\n  Initializing..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.Viewport.View(), m.footerView())
}

func (m ViewModel) headerView() string {
	title := titleStyle.Render("Mr. Pager")
	line := strings.Repeat("─", max(0, m.Viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m ViewModel) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.Viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.Viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

// func main() {
// 	// Load some text for our viewport
// 	content, err := os.ReadFile("artichoke.md")
// 	if err != nil {
// 		fmt.Println("could not load file:", err)
// 		os.Exit(1)
// 	}
//
// 	p := tea.NewProgram(
// 		model{content: string(content)},
// 		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
// 		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
// 	)
//
// 	if _, err := p.Run(); err != nil {
// 		fmt.Println("could not run program:", err)
// 		os.Exit(1)
// 	}
// }
