package ui

import (
	"WMTUI/internal/miner"
	"WMTUI/internal/minerScanner"
	"WMTUI/internal/ui/logging"
	"WMTUI/internal/ui/minertable"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MasterModel struct {
	Scanner *minerScanner.Scanner
	Table   tea.Model
	Logging tea.Model
}

func NewModel(s *minerScanner.Scanner, v logging.ViewModel, t tea.Model) MasterModel {
	return MasterModel{
		Scanner: s,
		Logging: v,
		Table:   t,
	}
}

func (m MasterModel) Init() tea.Cmd {
	return nil
}

func (m MasterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
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
			m.Table, cmd = m.Table.Update(minertable.MinerUpdateMsg{Miners: mnrs})
		case "F":
			m.Table, cmd = m.Table.Update(minertable.FastbootMsg{})
			return m, cmd
		case "L":
			// TODO: make this take an input on press, then send msg
			m.Table, cmd = m.Table.Update(minertable.LimitMsg{Limit: 3000, Vlt: 11.5, Freq: 300})
			return m, cmd
		case "O":
			m.Table, cmd = m.Table.Update(minertable.SlowbootMsg{})
			return m, cmd
		case "P":
			// TODO: make this take an input on press, then send msg
			m.Table, cmd = m.Table.Update(minertable.PoolMsg{Pools: []minertable.Pool{}})
			return m, cmd
		case "R":
			m.Table, cmd = m.Table.Update(minertable.RebootMsg{})
			return m, cmd
		case "S":
			m.Table, cmd = m.Table.Update(minertable.SleepMsg{})
			return m, cmd
		case "W":
			// TODO: make this take an input on press, then send msg
			m.Table, cmd = m.Table.Update(minertable.WakeMsg{Limit: 3000})
			return m, cmd
		default:
			m.Table, cmd = m.Table.Update(msg)
			return m, cmd

		}
	case minertable.CommandErrMsg:
		for _, e := range msg.Errors {
			m.Logging, cmd = m.Logging.Update(logging.LoggingMsg{Err: true, Message: e.Error()})
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	case minertable.MinerUpdateMsg:
		m.Table, cmd = m.Table.Update(msg)
		return m, cmd
	case logging.LoggingMsg:
		fmt.Println("Got it?")
		m.Logging, cmd = m.Logging.Update(msg)
		return m, cmd
	}

	// m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m MasterModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Top, fmt.Sprintf("%4s", m.Table.View()), m.Logging.View())
}
