package ui

import (
	"WMTUI/internal/config"
	"WMTUI/internal/miner"
	"WMTUI/internal/scanner"
	"WMTUI/internal/ui/logging"
	"WMTUI/internal/ui/popup"
	"WMTUI/internal/ui/siteselection"
	"WMTUI/internal/ui/table"
	"context"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MasterModel struct {
	Scanner              *scanner.Scanner
	Width                int
	Height               int
	Table                tea.Model
	Logging              tea.Model
	Popup                tea.Model
	PopupVisible         bool
	Spinner              spinner.Model
	SpinnerVisible       bool
	SiteSelection        tea.Model
	SiteSelectionVisible bool
	Context              context.Context
}

type scanStartMsg struct{}

type startSpinnerMsg struct{}

func NewModel(s *scanner.Scanner, v logging.ViewModel, t tea.Model, c config.Config) MasterModel {
	return MasterModel{
		Scanner:              s,
		Logging:              v,
		Table:                t,
		Popup:                popup.PopupModel{},
		PopupVisible:         false,
		Spinner:              NewSpinner(),
		SpinnerVisible:       false,
		SiteSelection:        siteselection.NewModel(siteselection.NewItems(c)),
		SiteSelectionVisible: false,
	}
}

func NewSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return s
}

func (m MasterModel) Init() tea.Cmd {
	return nil
}

func (m MasterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.SiteSelectionVisible {
		switch msg.(type) {
		case siteselection.SelectedSiteMsg:
			m.SiteSelectionVisible = false
			return m, func() tea.Msg { return msg }
		}

		m.SiteSelection, cmd = m.SiteSelection.Update(msg)
		return m, cmd
	}

	if m.PopupVisible {
		// NOTE: this feels kinda wrong
		switch msg.(type) {
		case popup.PopupCloseMsg:
			m.PopupVisible = false
			m.Popup = popup.PopupModel{}
			return m, nil
		}

		m.Popup, cmd = m.Popup.Update(msg)
		return m, cmd
	}

	if m.SpinnerVisible {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q":
				m.SpinnerVisible = false
				return m, nil
			}
		case table.ScanDoneMsg:
			m.SpinnerVisible = false
			return m, func() tea.Msg { return msg }
		}

		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Table, cmd = m.Table.Update(msg)
		cmds = append(cmds, cmd)
		m.Logging, cmd = m.Logging.Update(msg)
		cmds = append(cmds, cmd)
		m.SiteSelection, cmd = m.SiteSelection.Update(msg)
		cmds = append(cmds, cmd)
		m.SiteSelectionVisible = true
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			if m.Context != nil {
				m.Context.Done()
			}
			return m, tea.Quit
		case "r":
			m.Context.Done()
			return m, func() tea.Msg { return scanStartMsg{} }
		case "t":
			mnrs := []*miner.Miner{}
			for i := 0; i <= 10; i++ {
				mnrs = append(mnrs, miner.NewTestMiner(i%3, fmt.Sprintf("192.168.42.%d", i)))
			}
			m.Table, cmd = m.Table.Update(table.MinerUpdateMsg{Miners: mnrs})
		case "F":
			m.Table, cmd = m.Table.Update(table.FastbootMsg{})
			return m, cmd
		case "L":
			m.Popup = makeLimitPopup()
			m.PopupVisible = !m.PopupVisible
			return m, nil
		case "O":
			m.Table, cmd = m.Table.Update(table.SlowbootMsg{})
			return m, cmd
		case "P":
			m.Popup = makePoolPopup()
			m.PopupVisible = !m.PopupVisible
			return m, nil
		case "R":
			m.Table, cmd = m.Table.Update(table.RebootMsg{})
			return m, cmd
		case "S":
			m.Table, cmd = m.Table.Update(table.SleepMsg{})
			return m, cmd
		case "W":
			m.Popup = makeWakePopup()
			m.PopupVisible = !m.PopupVisible
			return m, nil
		default:
			m.Table, cmd = m.Table.Update(msg)
			return m, cmd
		}
	case siteselection.SelectedSiteMsg:
		m.Scanner.Conf = msg.Site
		cmds = append(cmds, func() tea.Msg { return scanStartMsg{} })
		cmds = append(cmds, func() tea.Msg { return startSpinnerMsg{} })
		return m, tea.Batch(cmds...)
	case startSpinnerMsg:
		m.SpinnerVisible = true
		return m, m.Spinner.Tick
	case scanStartMsg:
		go m.Scanner.ScanForMachines()
		ctx := context.Background()
		m.Context = ctx
		return m, nil
	case table.ScanDoneMsg:
		go m.Scanner.RefreshLoop(m.Context)
		return m, nil
	case popup.PopupCloseMsg:
		m.PopupVisible = false
		m.Popup = popup.PopupModel{}
		return m, nil
	case popup.PopupResMsg:
		return m.handlePopupResult(msg)
	case table.CommandErrMsg:
		for _, e := range msg.Errors {
			m.Logging, cmd = m.Logging.Update(logging.LoggingMsg{Err: true, Message: e.Error()})
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	case table.MinerUpdateMsg:
		m.Table, cmd = m.Table.Update(msg)
		return m, cmd
	case logging.LoggingMsg:
		m.Logging, cmd = m.Logging.Update(msg)
		return m, cmd
	}
	return m, cmd
}

func (m MasterModel) View() string {
	normal := lipgloss.JoinVertical(lipgloss.Top, fmt.Sprintf("%4s", m.Table.View()), m.Logging.View())

	if m.SiteSelectionVisible {
		s := m.SiteSelection.View()
		bg := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Render(s)

		cntr := lipgloss.Place(
			(m.Width*6)/10,
			m.Height,
			lipgloss.Center,
			lipgloss.Center,
			bg,
		)

		return lipgloss.JoinVertical(lipgloss.Left, cntr, m.Logging.View())
	}

	if m.PopupVisible {
		p := m.Popup.View()
		popupBg := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Render(p)

		centeredPopup := lipgloss.Place(
			(m.Width*6)/10,
			m.Height,
			lipgloss.Center,
			lipgloss.Center,
			popupBg,
		)

		return lipgloss.JoinVertical(lipgloss.Left, centeredPopup, m.Logging.View())
	}

	if m.SpinnerVisible {
		s := m.Spinner.View()
		bg := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Render("Scanning for\n", "machines\n", s)

		cntr := lipgloss.Place(
			m.Width,
			(m.Height*6)/10,
			lipgloss.Center,
			lipgloss.Center,
			bg,
		)

		return lipgloss.JoinVertical(lipgloss.Left, cntr, m.Logging.View())
	}

	return normal
}

func makePoolPopup() popup.PopupModel {
	p1u := popup.NewInput("Pool 1 URL: ", "")
	p1w := popup.NewInput("Pool 1 Worker: ", "")
	p1p := popup.NewInput("Pool 1 Pass: ", "")
	p2u := popup.NewInput("Pool 2 URL: ", "")
	p2w := popup.NewInput("Pool 2 Worker: ", "")
	p2p := popup.NewInput("Pool 2 Pass: ", "")
	p3u := popup.NewInput("Pool 3 URL: ", "")
	p3w := popup.NewInput("Pool 3 Worker: ", "")
	p3p := popup.NewInput("Pool 3 Pass: ", "")
	return popup.NewPopup("L", []textinput.Model{
		p1u, p1w, p1p, p2u, p2w, p2p, p3u, p3w, p3p,
	})
}

func makeLimitPopup() popup.PopupModel {
	t := popup.NewInput("Power Limit: ", "3000")
	t2 := popup.NewInput("Voltage: ", "11.5")
	t3 := popup.NewInput("Frequency: ", "300")
	return popup.NewPopup("P", []textinput.Model{t, t2, t3})
}

func makeWakePopup() popup.PopupModel {
	t := popup.NewInput("Power Limit: ", "3000")
	return popup.NewPopup("W", []textinput.Model{t})
}

func (m MasterModel) handlePopupResult(msg popup.PopupResMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.Caller {
	case "L":
		l, _ := strconv.Atoi(msg.Values[0])
		l = max(500, l)
		l = min(3600, l)

		v, _ := strconv.ParseFloat(msg.Values[1], 64)
		v = max(10, v)
		v = min(13, v)

		f, _ := strconv.ParseFloat(msg.Values[2], 64)
		f = max(250, f)
		f = min(700, f)

		m.Table, cmd = m.Table.Update(table.LimitMsg{Limit: l, Vlt: v, Freq: f})
		return m, cmd
	case "P":
		var pool1 table.Pool
		var pool2 table.Pool
		var pool3 table.Pool

		pool1.URL = msg.Values[0]
		pool1.Worker = msg.Values[1]
		pool1.Password = msg.Values[2]

		pool2.URL = msg.Values[3]
		pool2.Worker = msg.Values[4]
		pool2.Password = msg.Values[5]

		pool3.URL = msg.Values[6]
		pool3.Worker = msg.Values[7]
		pool3.Password = msg.Values[8]

		m.Table, cmd = m.Table.Update(table.PoolMsg{Pools: []table.Pool{pool1, pool2, pool3}})
		return m, cmd
	case "W":
		l, _ := strconv.Atoi(msg.Values[0])
		l = max(500, l)
		l = min(3600, l)

		m.Table, cmd = m.Table.Update(table.WakeMsg{Limit: l})
		return m, cmd
	}

	return m, nil
}
