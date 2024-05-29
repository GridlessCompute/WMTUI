package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"sort"
	"sync"
)

func (m *TableModel) GenerateInitialMiner(scanwg *sync.WaitGroup, popwg *sync.WaitGroup, mainChannel chan MinerObj, swg *sync.WaitGroup) {
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
		m.modelMinersList = append(m.modelMinersList, mnr)
	}

	fmt.Printf("NUMBER OF MINERS FOUND %d\n", len(m.modelMinersList))
}

func (m *TableModel) GenerateInitialMinerList() {
	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup
	var swg sync.WaitGroup

	channel := make(chan MinerObj, 150)

	swg.Add(1)

	go m.GenerateInitialMiner(&wg1, &wg2, channel, &swg)

	swg.Wait()

	close(channel)
}

func (m *TableModel) Clearscreen() tea.Msg {
	return tea.ClearScreen()
}

func (m *TableModel) refreshMainTable() {
	var wgSeven sync.WaitGroup
	var oldMinerSlice []MinerObj
	var newMinerSlice []MinerObj

	newMinerChannel := make(chan MinerObj, 150)

	oldMinerSlice = append(oldMinerSlice, m.modelMinersList...)

	PopulateRangeSlice(oldMinerSlice, &wgSeven, newMinerChannel)

	wgSeven.Wait()

	close(newMinerChannel)

	for miner := range newMinerChannel {
		newMinerSlice = append(newMinerSlice, miner)
	}

	m.modelMinersList = m.modelMinersList[:0]
	m.modelMinersList = append(m.modelMinersList, newMinerSlice...)
}

func (m *TableModel) ClearTables() {
	m.tables[MainView].SetRows(nil)
	m.tables[SelectedView].SetRows(nil)
	m.modelMinersList = m.modelMinersList[:0]
}

func (m *TableModel) sortByIP() {
	sort.Slice(m.modelMinersList, func(i, j int) bool {
		// TODO: Do I add a field that's just the miners last octet to help filer properly?
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

func (m *TableModel) FindSelectedMiner(mac string) MinerObj {
	allMiners := m.modelMinersList

	for _, miner := range allMiners {
		if miner.Miner.Mac == mac {
			return miner
		}
	}
	return *new(MinerObj)
}
