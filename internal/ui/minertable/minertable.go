package minertable

import (
	"WMTUI/internal/miner"
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

var Columns = []table.Column{
	{Title: "X", Width: 1},
	{Title: "Type", Width: 4},
	{Title: "IP", Width: 15},
	{Title: "Mac", Width: 17},
	{Title: "Status", Width: 10},
	{Title: "Errors", Width: 20},
	{Title: "Up Time", Width: 10},
	{Title: "GHs", Width: 10},
	{Title: "WTH", Width: 10},
	{Title: "Power", Width: 10},
	{Title: "Limit", Width: 10},
}

// type Miner struct {
// 	Selected   bool
// 	Type       int
// 	IP         string
// 	Mac        string
// 	Status     string
// 	Errors     string
// 	UpTime     time.Duration
// 	Hashrate   float64
// 	WTH        float64
// 	Power      float64
// 	PowerLimit float64
// 	Pool       string
// }

type MinerTableModel struct {
	Table     table.Model
	loaded    bool
	MinerList []*miner.Miner
}

type MinerUpdateMsg struct {
	Miners []*miner.Miner
}

func NewMinerTableModel() tea.Model {
	tbl := table.New(table.WithColumns(Columns), table.WithHeight(15))

	return MinerTableModel{
		Table:     tbl,
		loaded:    false,
		MinerList: []*miner.Miner{},
	}
}

func (m MinerTableModel) Init() tea.Cmd {
	return nil
}

func (m MinerTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.Table.MoveUp(1)
			return m, nil
		case "K", "shift + up":
			m.Table.MoveUp(5)
			return m, nil
		case "down", "j":
			m.Table.MoveDown(1)
			return m, nil
		case "J", "shift + down":
			m.Table.MoveDown(5)
		case "s":
			// Sort
		case "return", "enter":
			mnrI := m.Table.Cursor()
			m.MinerList[mnrI].Selected = !m.MinerList[mnrI].Selected
			m.Table.SetRows(makeNewRows(m.MinerList))
			return m, nil
		case "delete", "backspace":
			for i, mnr := range m.MinerList {
				mnr.Selected = false
				m.MinerList[i] = mnr
			}
			m.Table.SetRows(makeNewRows(m.MinerList))
			return m, nil
		}

	case MinerUpdateMsg:
		m.MinerList = msg.Miners
		m.Table.SetRows(makeNewRows(m.MinerList))
		return m, nil
	}

	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m MinerTableModel) View() string {
	return m.Table.View()
}

func makeNewRows(mnrs []*miner.Miner) []table.Row {
	var rows []table.Row

	for _, mnr := range mnrs {
		sel := ""

		if mnr.Selected {
			sel = "x"
		}

		typ := ""
		switch mnr.Type {
		case 0:
			typ = "NAM"
		case 1:
			typ = "WM"
		case 2:
			typ = "EPIC"
		}
		rows = append(rows, table.Row{
			sel,
			typ,
			mnr.IP,
			mnr.Mac,
			mnr.Status,
			mnr.Errors,
			mnr.UpTime.String(),
			fmt.Sprintf("%.2f", mnr.Hashrate),
			fmt.Sprintf("%.2f", mnr.Efficiency),
			fmt.Sprintf("%.2f", mnr.Power),
			fmt.Sprintf("%.2f", mnr.PowerLimit),
		})
	}

	return rows
}
