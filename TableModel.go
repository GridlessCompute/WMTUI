package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"strconv"
	"strings"
)

type status int
type sortBy int
type SetFarmMsg int

const (
	MainView status = iota
	SelectedView
)

var res map[string]interface{}

var poolmsg poolStruct

const (
	IPSort sortBy = iota
	MACSort
	THSort
	UPSort
)

type TableModel struct {
	focused         status
	tables          []table.Model
	keys            keyMap
	help            help.Model
	loaded          bool
	quiting         bool
	modelMinersList []MinerObj
	sortBy          sortBy
}

func NewTable() *TableModel {
	return &TableModel{}
}

func (m *TableModel) Init() tea.Cmd {
	return nil
}

func (m *TableModel) initTables(height int, width int) {
	columns := []table.Column{
		{Title: "IP", Width: 15},
		{Title: "Mac", Width: 20},
		{Title: "status", Width: 10},
		{Title: "Error Code", Width: 20},
		{Title: "Up time", Width: 10},
		{Title: "GHs", Width: 10},
		{Title: "WT", Width: 5},
		{Title: "W", Width: 7},
		{Title: "Limit", Width: 7},
		{Title: "Pool 1", Width: 30},
	}
	defaultTable := table.New(table.WithColumns(columns), table.WithHeight(height-25), table.WithWidth(width))

	m.tables = []table.Model{defaultTable, defaultTable}
}

func (m *TableModel) generateRows() {
	var rows []table.Row
	var statusString string
	if len(m.modelMinersList) > 0 {
		for _, miner := range m.modelMinersList {
			if miner.status == false {
				statusString = "sleeping"
			} else {
				statusString = "running"
			}

			rows = append(rows, table.Row{
				miner.Miner.Ip,
				miner.Miner.Mac,
				statusString,
				miner.Miner.Errcode,
				fmt.Sprint(miner.Miner.UpTime),
				fmt.Sprintf("%.2f", float32(miner.Miner.Hrrt)/1000),
				fmt.Sprint(miner.Miner.Wt),
				fmt.Sprint(miner.Miner.W),
				fmt.Sprint(miner.Miner.Limit),
				fmt.Sprint(miner.Miner.ActivePool),
			})
		}
	} else {
		rows = append(rows, table.Row{"", "", "", "No miners found", "", "", "", "", ""})
	}

	m.tables[MainView].SetRows(rows)
}

func (m *TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	var winH int
	var winW int
	switch msg := msg.(type) {
	case PowerStruct:
		newLimit, err := strconv.Atoi(msg.Power)
		if err != nil {
			fmt.Println(err)
		}
		if m.focused == MainView {
			mac := m.tables[m.focused].SelectedRow()[1]
			miner := m.FindSelectedMiner(mac)

			if newLimit >= 500 && newLimit <= 3500 {
				result, _ := SendToApi(miner.Token, "adjust_power_limit", map[string]interface{}{"power_limit": strconv.Itoa(newLimit)})
				res = result
			}
		} else if m.focused == SelectedView {
			miners := m.FindSelectedMiners()
			if newLimit >= 500 && newLimit <= 3500 {
				for _, miner := range miners {
					SendToApi(miner.Token, "adjust_power_limit", map[string]interface{}{"power_limit": strconv.Itoa(newLimit)})
				}
			}

		}

	case SetFarmMsg:
		m.ClearTables()
		//m.initTables(winH, winW)
		m.GenerateInitialMinerList()
		m.sortBy = IPSort
		m.sortByIP()
		m.generateRows()
		return m, m.Clearscreen

	case poolStruct:
		if m.focused == MainView {
			//take msg.POOLINFO and send a command with it, making custom workernames
			mac := m.tables[m.focused].SelectedRow()[1]
			miner := m.FindSelectedMiner(mac)
			poolmsg = msg
			if msg.WorkerType == "mac" || msg.WorkerType == "MAC" {
				//do thing based on mac
				minerMac := strings.Replace(miner.Miner.Mac, ":", "", -1)
				if msg.Worker1 != "" {
					msg.Worker1 = msg.Worker1 + "." + minerMac
				}
				if msg.Worker2 != "" {
					msg.Worker2 = msg.Worker2 + "." + minerMac
				}
				if msg.Worker3 != "" {
					msg.Worker3 = msg.Worker3 + "." + minerMac
				}
				// SENDTOAPI POOLS GOES BELOw
				SendToApi(miner.Token, "update_pools", map[string]interface{}{"pool1": msg.Url1, "worker1": msg.Worker1, "passwd1": msg.Psswd1, "pool2": msg.Url2, "worker2": msg.Worker2, "passwd2": msg.Psswd2, "pool3": msg.Url3, "worker3": msg.Worker3, "passwd3": msg.Psswd3})
			} else if msg.WorkerType == "ip" || msg.WorkerType == "IP" {
				//do thing based on IP
				minerIP := strings.Replace(miner.Miner.Ip, ".", "x", -1)
				if msg.Worker1 != "" {
					msg.Worker1 = msg.Worker1 + "." + minerIP
				}
				if msg.Worker2 != "" {
					msg.Worker2 = msg.Worker2 + "." + minerIP
				}
				if msg.Worker3 != "" {
					msg.Worker3 = msg.Worker3 + "." + minerIP
				}
				// SENDTOAPI POOLS GOES HERe
				res, err := SendToApi(miner.Token, "update_pools", map[string]interface{}{"pool1": msg.Url1, "worker1": msg.Worker1, "passwd1": msg.Psswd1, "pool2": msg.Url2, "worker2": msg.Worker2, "passwd2": msg.Psswd2, "pool3": msg.Url3, "worker3": msg.Worker3, "passwd3": msg.Psswd3})
				if err != nil {
					log.Println(err)
				}
				log.Println(res)
			}
		} else if m.focused == SelectedView {
			//take msg.POOLINFO and send a command with it, making custom workernames
			miners := m.FindSelectedMiners()
			if msg.WorkerType == "mac" || msg.WorkerType == "MAC" {
				for _, miner := range miners {
					minerMac := strings.Replace(miner.Miner.Mac, ":", "", -1)
					if msg.Worker1 != "" {
						msg.Worker1 = msg.Worker1 + "." + minerMac
					}
					if msg.Worker2 != "" {
						msg.Worker2 = msg.Worker2 + "." + minerMac
					}
					if msg.Worker3 != "" {
						msg.Worker3 = msg.Worker3 + "." + minerMac
					}
					// SENDTOAPI POOLS GOES BELOW
					res, _ = SendToApi(miner.Token, "update_pools", map[string]interface{}{"pool1": msg.Url1, "worker1": msg.Worker1, "passwd1": msg.Psswd1, "pool2": msg.Url2, "worker2": msg.Worker2, "passwd2": msg.Psswd2, "pool3": msg.Url3, "worker3": msg.Worker3, "passwd3": msg.Psswd3})
				}
			} else if msg.WorkerType == "ip" || msg.WorkerType == "IP" {
				for _, miner := range miners {
					minerIP := strings.Replace(miner.Miner.Ip, ".", "x", -1)
					if msg.Worker1 != "" {
						msg.Worker1 = msg.Worker1 + "." + minerIP
					}
					if msg.Worker2 != "" {
						msg.Worker2 = msg.Worker2 + "." + minerIP
					}
					if msg.Worker3 != "" {
						msg.Worker3 = msg.Worker3 + "." + minerIP
					}
					// SENDTOAPI POOLS GOES HERE
					res, _ = SendToApi(miner.Token, "update_pools", map[string]interface{}{"pool1": msg.Url1, "worker1": msg.Worker1, "passwd1": msg.Psswd1, "pool2": msg.Url2, "worker2": msg.Worker2, "passwd2": msg.Psswd2, "pool3": msg.Url3, "worker3": msg.Worker3, "passwd3": msg.Psswd3})
				}
			}
		}
	case tea.WindowSizeMsg:
		if !m.loaded {
			winH = msg.Height
			winW = msg.Width
			m.keys = keys
			m.help = help.New()
			fmt.Println("Initializing Tables")
			m.initTables(winH, winW)
			fmt.Println("Getting Miners")
			m.GenerateInitialMinerList()
			//fmt.Println("Sorting Miners")
			//m.sortBy = IPSort
			//m.sortByIP()
			fmt.Println("Generating rows")
			m.generateRows()
			m.help.Width = msg.Width
			m.quiting = false
			m.loaded = true
			return m, m.Clearscreen
		}
	case tea.KeyMsg:
		switch {
		// quit program
		case key.Matches(msg, m.keys.Quit):
			m.quiting = true
			return m, tea.Quit
			// Show full help
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			// Switch between table
		case key.Matches(msg, m.keys.Switch):
			if m.focused == MainView {
				m.focused = SelectedView
			} else {
				m.focused = MainView
			}
			// Swap selected miner between tables
		case key.Matches(msg, m.keys.Select):

			rowsV, rowsH := m.TransferRow()

			if m.focused == MainView {
				m.tables[MainView].SetRows(rowsV)
				m.tables[SelectedView].SetRows(rowsH)
			} else if m.focused == SelectedView {
				m.tables[MainView].SetRows(rowsH)
				m.tables[SelectedView].SetRows(rowsV)
			}
			// Move up in table
		case key.Matches(msg, m.keys.Up):
			m.tables[m.focused].MoveUp(1)
			// Move up in table
		case key.Matches(msg, m.keys.Down):
			m.tables[m.focused].MoveDown(1)
			// Refresh Main table
		case key.Matches(msg, m.keys.Refresh):
			if m.focused == MainView {
				tea.ClearScreen()
				m.refreshMainTable()
				//m.initTables(winH, winW)
				m.generateRows()
				return m, m.Clearscreen
			}
			// Sort Machines
		case key.Matches(msg, m.keys.Sort):
			if m.focused == MainView {
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
			}
			// Select Farm
		case key.Matches(msg, m.keys.Farm):
			Models[TableView] = m
			return Models[ListView].Update(tea.WindowSizeMsg{Width: 50, Height: 50})
			// Create New Farm
		case key.Matches(msg, m.keys.NewFarm):
			Models[TableView] = m
			return Models[FormView].Update(tea.WindowSizeMsg{})
			// Reboot Machine
		case key.Matches(msg, m.keys.Reboot):
			if m.focused == SelectedView {
				miners := m.FindSelectedMiners()
				for _, miner := range miners {
					SendToApi(miner.Token, "reboot", nil)
				}
			} else {
				mac := m.tables[m.focused].SelectedRow()[1]
				miner := m.FindSelectedMiner(mac)

				if miner.Token.IPAddress != "" {
					SendToApi(miner.Token, "reboot", nil)
				}
			}
			// Sleep Machine
		case key.Matches(msg, m.keys.Sleep):
			if m.focused == SelectedView {
				miners := m.FindSelectedMiners()
				for _, miner := range miners {
					SendToApi(miner.Token, "power_off", map[string]interface{}{"respbefore": "false"})
				}
			} else {
				mac := m.tables[m.focused].SelectedRow()[1]
				miner := m.FindSelectedMiner(mac)

				if miner.Token.IPAddress != "" {
					SendToApi(miner.Token, "power_off", map[string]interface{}{"respbefore": "false"})
				}
			}
			//  Set Pools
		case key.Matches(msg, m.keys.Pools):

			if m.focused == SelectedView {
				Models[SelectedView] = m
				return Models[PoolView].Update(tea.WindowSizeMsg{})
			} else {
				// set pools for highlighted machine
				Models[MainView] = m
				return Models[PoolView].Update(tea.WindowSizeMsg{})
			}
			// Set A Power Limit
		case key.Matches(msg, m.keys.Limit):
			if m.focused == SelectedView {
				Models[SelectedView] = m
				return Models[PowerView].Update(tea.WindowSizeMsg{})
			} else {
				// set pools for highlighted machine
				Models[MainView] = m
				return Models[PowerView].Update(tea.WindowSizeMsg{})
			}
			// Wake sleeping Machine
		case key.Matches(msg, m.keys.Wake):
			if m.focused == SelectedView {
				miners := m.FindSelectedMiners()
				for _, miner := range miners {
					SendToApi(miner.Token, "power_on", nil)
				}
			} else {
				mac := m.tables[m.focused].SelectedRow()[1]
				miner := m.FindSelectedMiner(mac)

				if miner.Token.IPAddress != "" {
					SendToApi(miner.Token, "power_on", nil)
				}
			}
			// Activate Fast Boot
		case key.Matches(msg, m.keys.Fast):
			if m.focused == SelectedView {
				miners := m.FindSelectedMiners()
				for _, miner := range miners {
					SendToApi(miner.Token, "enable_btminer_fast_boot", nil)
				}
			} else {
				mac := m.tables[m.focused].SelectedRow()[1]
				miner := m.FindSelectedMiner(mac)

				if miner.Token.IPAddress != "" {
					SendToApi(miner.Token, "enable_btminer_fast_boot", nil)
				}
			}
			// Activate Slow Boot
		case key.Matches(msg, m.keys.Slow):
			if m.focused == SelectedView {
				miners := m.FindSelectedMiners()
				for _, miner := range miners {
					SendToApi(miner.Token, "disable_btminer_fast_boot", nil)
				}
			} else {
				mac := m.tables[m.focused].SelectedRow()[1]
				miner := m.FindSelectedMiner(mac)

				if miner.Token.IPAddress != "" {
					SendToApi(miner.Token, "disable_btminer_fast_boot", nil)
				}
			}
		}
	}

	// return m, cmd
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *TableModel) View() string {
	if m.quiting {
		return ""
	}
	if m.loaded {
		switch m.focused {
		case MainView:

			helpView := m.help.View(m.keys)

			return m.tables[MainView].View() + "\n" + strings.Repeat("\n", 8) + helpView + "\n" + ChosenFarm.Name
		default:
			return m.tables[SelectedView].View() + "\n" + m.help.View(m.keys) + "\n" + ChosenFarm.Name
		}
	} else {
		return "Scanning..."
	}
}
