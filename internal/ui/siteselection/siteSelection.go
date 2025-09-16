package siteselection

import (
	"WMTUI/internal/config"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type item struct {
	SiteName  string
	SiteRange string
}

type SelectedSiteMsg struct {
	Site config.Site
}

func (i item) Title() string       { return i.SiteName }
func (i item) Description() string { return i.SiteRange }
func (i item) FilterValue() string { return i.SiteName }

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

type model struct {
	list         list.Model
	items        []item
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
}

func NewModel(i []item) model {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
	)

	var newI []list.Item
	for _, j := range i {
		newI = append(newI, j)
	}

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	lst := list.New(newI, delegate, 0, 0)
	lst.Title = "Sites"
	lst.Styles.Title = titleStyle
	lst.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	return model{
		list:         lst,
		items:        i,
		keys:         listKeys,
		delegateKeys: delegateKeys,
	}
}

func NewItems(sites config.Config) []item {
	var l []item

	for _, site := range sites.Sites {
		l = append(l, item{SiteName: site.Name, SiteRange: site.IPRange})
	}

	return l
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize((msg.Width-h)/3, ((msg.Height - v) / 4))

	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "return":
			s := m.items[m.list.Index()]
			return m, func() tea.Msg { return SelectedSiteMsg{Site: config.Site{Name: s.Title(), IPRange: s.Description()}} }
		case "j", "down":
			m.list.CursorDown()
			return m, nil
		case "k", "up":
			m.list.CursorUp()
			return m, nil
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return appStyle.Render(m.list.View())
}

// func main() {
// 	if _, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run(); err != nil {
// 		fmt.Println("Error running program:", err)
// 		os.Exit(1)
// 	}
// }
