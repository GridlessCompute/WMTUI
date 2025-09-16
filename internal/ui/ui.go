package ui

import (
	"WMTUI/internal/miner"
	"WMTUI/internal/scanner"
	"WMTUI/internal/ui/logging"
	"WMTUI/internal/ui/popup"
	"WMTUI/internal/ui/table"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MasterModel struct {
	Scanner      *scanner.Scanner
	Table        tea.Model
	Logging      tea.Model
	Popup        tea.Model
	Width        int
	Height       int
	PopupVisible bool
}

func NewModel(s *scanner.Scanner, v logging.ViewModel, t tea.Model) MasterModel {
	return MasterModel{
		Scanner:      s,
		Logging:      v,
		Table:        t,
		Popup:        popup.PopupModel{},
		PopupVisible: false,
	}
}

func (m MasterModel) Init() tea.Cmd {
	return nil
}

func (m MasterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

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

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Table, cmd = m.Table.Update(msg)
		cmds = append(cmds, cmd)
		m.Logging, cmd = m.Logging.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "r":
			go m.Scanner.StartScanning()
			return m, nil
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
		fmt.Println("Got it?")
		m.Logging, cmd = m.Logging.Update(msg)
		return m, cmd
	}
	return m, cmd
}

func (m MasterModel) View() string {
	normal := lipgloss.JoinVertical(lipgloss.Top, fmt.Sprintf("%4s", m.Table.View()), m.Logging.View())

	if !m.PopupVisible {
		return normal
	}

	p := m.Popup.View()
	popupBg := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Render(p)

	centeredPopup := lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		popupBg,
	)

	return lipgloss.JoinVertical(lipgloss.Left, normal, centeredPopup)
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
