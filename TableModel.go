package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type status int
type sortBy int
type SetFarmMsg int

const (
	MainView status = iota
	SelectedView
)

const (
	IPSort sortBy = iota
	MACSort
	THSort
	UPSort
)

type CommandKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Switch   key.Binding
	Reboot   key.Binding
	Sleep    key.Binding
	Wake     key.Binding
	Pools    key.Binding
	Limit    key.Binding
	Fast     key.Binding
	Slow     key.Binding
	Help     key.Binding
	Deselect key.Binding
	Quit     key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k CommandKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k CommandKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Help, k.Quit},           // first column
		{k.Sleep, k.Wake, k.Fast, k.Slow},        // second column
		{k.Reboot, k.Pools, k.Limit, k.Deselect}, // Third column
	}
}

var keys = CommandKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "Move Up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "Move Down"),
	),
	Switch: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "Switch to Main View"),
	),
	Reboot: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Reboot Miner(s)"),
	),
	Sleep: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Sleep Miner(s)"),
	),
	Wake: key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "Wake Miner(s)"),
	),
	Pools: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Set Pools"),
	),
	Limit: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "Set Power Limit"),
	),
	Fast: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "Set Fast Boot"),
	),
	Slow: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "Set Slow Boot"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Deselect: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Deselect Miner"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type MainKeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Switch   key.Binding
	Reboot   key.Binding
	Sleep    key.Binding
	Wake     key.Binding
	Pools    key.Binding
	Limit    key.Binding
	Fast     key.Binding
	Slow     key.Binding
	Help     key.Binding
	Deselect key.Binding
	Quit     key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k MainKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k MainKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Help, k.Quit},           // first column
		{k.Sleep, k.Wake, k.Fast, k.Slow},        // second column
		{k.Reboot, k.Pools, k.Limit, k.Deselect}, // Third column
	}
}

var keys = MainKeyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "Move Up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "Move Down"),
	),
	Switch: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "Switch to Main View"),
	),
	Reboot: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Reboot Miner(s)"),
	),
	Sleep: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Sleep Miner(s)"),
	),
	Wake: key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "Wake Miner(s)"),
	),
	Pools: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "Set Pools"),
	),
	Limit: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "Set Power Limit"),
	),
	Fast: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "Set Fast Boot"),
	),
	Slow: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "Set Slow Boot"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Deselect: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Deselect Miner"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type TableModel struct {
	focused         status
	tables          []table.Model
	mainKeys        MainKeyMap
	commandKeys     CommandKeyMap
	help            help.Model
	loaded          bool
	quiting         bool
	modelMinersList []MinerObj
	sortBy          sortBy
}

func NewTable() *TableModel {
	return &TableModel{}
}

func (t *TableModel) GenerateInitialMiner(scanwg *sync.WaitGroup, popwg *sync.WaitGroup, mainChannel chan MinerObj, swg *sync.WaitGroup) {
	mnrOChannel := make(chan MinerObj, 150)
	popChannel := make(chan MinerObj, 150)

	defer swg.Done()

	f := ChosenFarm

	ips := GenerateRangeIPs(GenIpRange(f))

	ScanRange(ips, scanwg, mnrOChannel)
	scanwg.Wait()
	close(mnrOChannel)

	PopulateRange(mnrOChannel, popwg, popChannel)
	popwg.Wait()
	close(popChannel)

	for mnr := range popChannel {
		t.modelMinersList = append(t.modelMinersList, mnr)
	}
}

func (t *TableModel) GenerateInitialMinerList() {
	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup
	var swg sync.WaitGroup

	channel := make(chan MinerObj, 150)

	swg.Add(1)

	go t.GenerateInitialMiner(&wg1, &wg2, channel, &swg)

	swg.Wait()

	close(channel)
}

func (t *TableModel) generateRows() {
	var rows []table.Row

	for _, miner := range t.modelMinersList {
		rows = append(rows, table.Row{
			miner.Miner.Ip,
			miner.Miner.Mac,
			miner.Miner.Errcode,
			fmt.Sprint(miner.Miner.UpTime),
			fmt.Sprint(miner.Miner.Hrrt),
			fmt.Sprint(miner.Miner.Wt),
			fmt.Sprint(miner.Miner.W),
			fmt.Sprint(miner.Miner.Limit),
			miner.Miner.AcvtivePool,
		})
	}

	t.tables[MainView].SetRows(rows)
}

// WIP
func (t *TableModel) refreshMainTable() {
	var wgSeven sync.WaitGroup
	var oldMinerSlice []MinerObj
	var newMinerSlice []MinerObj

	newMinerChannel := make(chan MinerObj, 150)

	oldMinerSlice = append(oldMinerSlice, t.modelMinersList...)

	PopulateRangeSlice(oldMinerSlice, &wgSeven, newMinerChannel)

	wgSeven.Wait()

	close(newMinerChannel)

	for miner := range newMinerChannel {
		newMinerSlice = append(newMinerSlice, miner)
	}

	t.modelMinersList = t.modelMinersList[:0]
	t.modelMinersList = append(t.modelMinersList, newMinerSlice...)
}

func (m *TableModel) initTables(height int) {
	columns := []table.Column{
		{Title: "IP", Width: 15},
		{Title: "Mac", Width: 20},
		{Title: "Error Code", Width: 25},
		{Title: "Up time", Width: 10},
		{Title: "Hrrt", Width: 10},
		{Title: "WT", Width: 10},
		{Title: "W", Width: 10},
		{Title: "Limit", Width: 10},
		{Title: "Pool 1", Width: 50},
	}
	defaultTable := table.New(table.WithColumns(columns), table.WithHeight(height-50))

	m.tables = []table.Model{defaultTable, defaultTable}
}

func (m *TableModel) ClearTables() {
	m.tables[MainView].SetRows(nil)
	m.tables[SelectedView].SetRows(nil)
	m.modelMinersList = m.modelMinersList[:0]
}

func (m *TableModel) sortByIP() {
	sort.Slice(m.modelMinersList, func(i, j int) bool {
		return m.modelMinersList[i].Miner.Ip < m.modelMinersList[j].Miner.Ip
	})
}

func (m *TableModel) sortByMAC() {
	sort.Slice(m.modelMinersList, func(i, j int) bool {
		return m.modelMinersList[i].Miner.Mac < m.modelMinersList[j].Miner.Mac
	})
}

func (m *TableModel) sortByTH() {
	sort.Slice(m.modelMinersList, func(i, j int) bool {
		return m.modelMinersList[i].Miner.Hrrt < m.modelMinersList[j].Miner.Hrrt
	})
}

func (m *TableModel) sortByUP() {
	sort.Slice(m.modelMinersList, func(i, j int) bool {
		return m.modelMinersList[i].Miner.UpTime < m.modelMinersList[j].Miner.UpTime
	})
}

func (m *TableModel) TransferRow() ([]table.Row, []table.Row) {
	// Seems to transfer all rows
	var tableH table.Model
	var newRowsV []table.Row
	var newRowsH []table.Row

	tableV := m.tables[m.focused]

	if m.focused == MainView {
		tableH = m.tables[SelectedView]
	} else {
		tableH = m.tables[MainView]
	}

	rowsV := tableV.Rows()
	newRowsH = append(newRowsH, tableH.Rows()...)

	mac := m.tables[m.focused].SelectedRow()[1]

	for _, r := range rowsV {
		if r[1] != mac {
			newRowsV = append(newRowsV, r)

		} else {
			newRowsH = append(newRowsH, r)
		}
	}

	return newRowsV, newRowsH
}

func (m *TableModel) FindSelectedMiners() []MinerObj {
	var miners []MinerObj

	selectedMiners := m.tables[SelectedView].Rows()

	allMiners := m.modelMinersList

	for _, row := range selectedMiners {
		for _, miner := range allMiners {
			if miner.Miner.Mac == row[1] {
				miners = append(miners, miner)
			}
		}
	}
	return miners
}

func (m TableModel) Init() tea.Cmd {
	return nil
}

func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	var winH int
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			winH = msg.Height
			fmt.Println("Initializing Tables")
			m.initTables(winH)
			fmt.Println("Getting Miners")
			m.GenerateInitialMinerList()
			fmt.Println("Sorting Miners")
			m.sortBy = IPSort
			m.sortByIP()
			fmt.Println("Generating rows")
			m.generateRows()
			m.help.Width = msg.Width
			m.quiting = false
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quiting = true
			return m, tea.Quit
		}
	case SetFarmMsg:
		m.ClearTables()
		m.initTables(winH)
		m.GenerateInitialMinerList()
		m.sortBy = IPSort
		m.sortByIP()
		m.generateRows()
	}

	switch m.focused {
	case MainView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "tab":
				m.focused = SelectedView
			case "enter":
				rowsV, rowsH := m.TransferRow()

				if m.focused == MainView {
					m.tables[MainView].SetRows(rowsV)
					m.tables[SelectedView].SetRows(rowsH)
				} else if m.focused == SelectedView {
					m.tables[MainView].SetRows(rowsH)
					m.tables[SelectedView].SetRows(rowsV)
				}
			case "up":
				m.tables[m.focused].MoveUp(1)
			case "down":
				m.tables[m.focused].MoveDown(1)
			case "r":
				m.refreshMainTable()
				m.generateRows()
			case "`":
				// TODO: Implement Sort -> 90% done IPSort is a bit janky due to the nature of strings
				switch m.sortBy {
				case IPSort:
					m.sortBy = MACSort
					m.sortByMAC()
					m.generateRows()
				case MACSort:
					m.sortBy = THSort
					m.sortByTH()
					m.generateRows()
				case THSort:
					m.sortBy = UPSort
					m.sortByUP()
					m.generateRows()
				case UPSort:
					m.sortBy = IPSort
					m.sortByIP()
					m.generateRows()
				}
			case "f":
				Models[TableView] = m
				return Models[ListView].Update(tea.WindowSizeMsg{Width: 50, Height: 50})
			case "n":
				Models[TableView] = m

				return Models[FormView].Update(tea.WindowSizeMsg{})
			}
		}
	case SelectedView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.commandKeys.Deselect):
				rowsV, rowsH := m.TransferRow()

				if m.focused == MainView {
					m.tables[MainView].SetRows(rowsV)
					m.tables[SelectedView].SetRows(rowsH)
				} else if m.focused == SelectedView {
					m.tables[MainView].SetRows(rowsH)
					m.tables[SelectedView].SetRows(rowsV)
				}
			case key.Matches(msg, m.commandKeys.Switch):
				m.focused = MainView
			case key.Matches(msg, m.commandKeys.Up):
				m.tables[m.focused].MoveUp(1)
			case key.Matches(msg, m.commandKeys.Down):
				m.tables[m.focused].MoveDown(1)
			case key.Matches(msg, m.commandKeys.Reboot):
				// TODO: Implement Reboot
				miners := m.FindSelectedMiners()
				for _, miner := range miners {
					SendToApi(miner.Token, "reboot", nil)
				}
			case key.Matches(msg, m.commandKeys.Sleep):
				// TODO: Implement Sleep
				miners := m.FindSelectedMiners()
				for _, miner := range miners {
					SendToApi(miner.Token, "power_off", map[string]interface{}{"respbefore": "false"})
				}
			case key.Matches(msg, m.commandKeys.Pools):
				// TODO: Implement Pools
			case key.Matches(msg, m.commandKeys.Limit):
				// TODO: Implement Power Limit
			case key.Matches(msg, m.commandKeys.Wake):
				// TODO: Implement Wake
				miners := m.FindSelectedMiners()
				for _, miner := range miners {
					SendToApi(miner.Token, "power_on", nil)
				}
			case key.Matches(msg, m.commandKeys.Fast):
				// TODO: Implement Fastboot Toggle
				miners := m.FindSelectedMiners()
				for _, miner := range miners {
					SendToApi(miner.Token, "enable_btminer_fast_boot", nil)
				}
			case key.Matches(msg, m.commandKeys.Slow):
				// TODO: Implement Slow Boot
				miners := m.FindSelectedMiners()
				for _, miner := range miners {
					SendToApi(miner.Token, "disable_btminer_fast_boot", nil)
				}
				// case "b":
				// 	// TODO: Implemet Blink Toggle
				// 	// Need function to fins selected miners
			}
		}
	}

	// return m, cmd
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (t TableModel) View() string {
	if t.quiting {
		return ""
	}
	if t.loaded {
		switch t.focused {
		case MainView:

			helpView := t.help.View(t.commandKeys)
			t.help.ShowAll = true

			// return "\n" + status + strings.Repeat("\n", height) + helpView

			return fmt.Sprintln(ChosenFarm) + "\n" + t.tables[MainView].View() + "\n" + strings.Repeat("\n", 8) + helpView
		default:
			return t.tables[SelectedView].View() + "\n" + t.help.View(t.commandKeys)
		}
	} else {
		return "Scanning..."
	}
}
