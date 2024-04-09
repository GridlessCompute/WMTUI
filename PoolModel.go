package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type PoolModel struct {
	url1       textinput.Model
	worker1    textinput.Model
	psswd1     textinput.Model
	url2       textinput.Model
	worker2    textinput.Model
	psswd2     textinput.Model
	url3       textinput.Model
	worker3    textinput.Model
	psswd3     textinput.Model
	workerType textinput.Model
}

type poolStruct struct {
	Url1       string
	Worker1    string
	Psswd1     string
	Url2       string
	Worker2    string
	Psswd2     string
	Url3       string
	Worker3    string
	Psswd3     string
	WorkerType string
}

func NewPoolForm() *PoolModel {
	form := &PoolModel{
		url1: textinput.New(), worker1: textinput.New(), psswd1: textinput.New(),
		url2: textinput.New(), worker2: textinput.New(), psswd2: textinput.New(),
		url3: textinput.New(), worker3: textinput.New(), psswd3: textinput.New(),
		workerType: textinput.New(),
	}

	form.url1.Placeholder = "Pool 1 Url"
	form.worker1.Placeholder = "Raw Worker Name"
	form.psswd1.Placeholder = "Pool Password"
	form.url2.Placeholder = "Pool 2 Url"
	form.worker2.Placeholder = "Raw Worker Name"
	form.psswd2.Placeholder = "Pool Password"
	form.url3.Placeholder = "Pool 3 Url"
	form.worker3.Placeholder = "Raw Worker Name"
	form.psswd3.Placeholder = "Pool Password"
	form.workerType.Placeholder = "(MAC) Worker / (IP) Worker"

	form.url1.Focus()
	return form
}

func (m PoolModel) Init() tea.Cmd {
	return textinput.Blink
}

// How do I get the Ip And Mac
// Potentially send with a message
func NewPoolInfo(url1, url2, url3, worker1, worker2, worker3, psswd1, psswd2, psswd3, workerFormatting string) poolStruct {
	var ps poolStruct
	// ps.machineMac = miner.Miner.Mac
	// ps.machineIP = miner.Miner.Ip

	ps.Worker1 = worker1
	ps.Worker2 = worker2
	ps.Worker3 = worker3

	ps.Url1 = url1
	ps.Url2 = url2
	ps.Url3 = url3

	ps.Psswd1 = psswd1
	ps.Psswd2 = psswd2
	ps.Psswd3 = psswd3

	return ps

}

// func NewFarm(name, root, length string) FarmStruct {
// 	return FarmStruct{Name: name, Start: root, End: length}
// }

func (m PoolModel) CreatePool() tea.Msg {
	pool := NewPoolInfo(
		m.url1.Value(), m.worker1.Value(), m.psswd1.Value(),
		m.url2.Value(), m.worker2.Value(), m.psswd2.Value(),
		m.url3.Value(), m.worker3.Value(), m.psswd3.Value(),
		m.workerType.Value(),
	)

	return pool
}

func (m PoolModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	// case MinerMessage:
	// Take the message and call a function to do some shit
	// need to look at other FormModel to figure this out
	case tea.WindowSizeMsg:
		return m, textinput.Blink
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return Models[TableView].Update(nil)
		case "ctrl+c":
			return m, tea.Quit
		// THis is super ugly
		case "tab":
			if m.url1.Focused() {
				m.url1.Blur()
				m.worker1.Focus()
				m.psswd1.Blur()
				m.url2.Blur()
				m.worker2.Blur()
				m.psswd2.Blur()
				m.url3.Blur()
				m.worker3.Blur()
				m.psswd3.Blur()
				m.workerType.Blur()
			} else if m.worker1.Focused() {
				m.url1.Blur()
				m.worker1.Blur()
				m.psswd1.Focus()
				m.url2.Blur()
				m.worker2.Blur()
				m.psswd2.Blur()
				m.url3.Blur()
				m.worker3.Blur()
				m.psswd3.Blur()
				m.workerType.Blur()
			} else if m.psswd1.Focused() {
				m.url1.Blur()
				m.worker1.Blur()
				m.psswd1.Blur()
				m.url2.Focus()
				m.worker2.Blur()
				m.psswd2.Blur()
				m.url3.Blur()
				m.worker3.Blur()
				m.psswd3.Blur()
				m.workerType.Blur()
			} else if m.url2.Focused() {
				m.url1.Blur()
				m.worker1.Blur()
				m.psswd1.Blur()
				m.url2.Blur()
				m.worker2.Focus()
				m.psswd2.Blur()
				m.url3.Blur()
				m.worker3.Blur()
				m.psswd3.Blur()
				m.workerType.Blur()
			} else if m.worker2.Focused() {
				m.url1.Blur()
				m.worker1.Blur()
				m.psswd1.Blur()
				m.url2.Blur()
				m.worker2.Blur()
				m.psswd2.Focus()
				m.url3.Blur()
				m.worker3.Blur()
				m.psswd3.Blur()
				m.workerType.Blur()
			} else if m.psswd2.Focused() {
				m.url1.Blur()
				m.worker1.Blur()
				m.psswd1.Blur()
				m.url2.Blur()
				m.worker2.Blur()
				m.psswd2.Blur()
				m.url3.Focus()
				m.worker3.Blur()
				m.psswd3.Blur()
				m.workerType.Blur()
			} else if m.url3.Focused() {
				m.url1.Blur()
				m.worker1.Blur()
				m.psswd1.Blur()
				m.url2.Blur()
				m.worker2.Blur()
				m.psswd2.Blur()
				m.url3.Blur()
				m.worker3.Focus()
				m.psswd3.Blur()
				m.workerType.Blur()
			} else if m.worker3.Focused() {
				m.url1.Blur()
				m.worker1.Blur()
				m.psswd1.Blur()
				m.url2.Blur()
				m.worker2.Blur()
				m.psswd2.Blur()
				m.url3.Blur()
				m.worker3.Blur()
				m.psswd3.Focus()
				m.workerType.Blur()
			} else if m.psswd3.Focused() {
				m.url1.Blur()
				m.worker1.Blur()
				m.psswd1.Blur()
				m.url2.Blur()
				m.worker2.Blur()
				m.psswd2.Blur()
				m.url3.Blur()
				m.worker3.Blur()
				m.psswd3.Blur()
				m.workerType.Focus()
			} else if m.workerType.Focused() {
				m.url1.Focus()
				m.worker1.Blur()
				m.psswd1.Blur()
				m.url2.Blur()
				m.worker2.Blur()
				m.psswd2.Blur()
				m.url3.Blur()
				m.worker3.Blur()
				m.psswd3.Blur()
				m.workerType.Blur()
			}
			return m, textinput.Blink
		case "enter":
			return Models[ListView], m.CreatePool
		}
	}
	if m.url1.Focused() {
		m.url1, cmd = m.url1.Update(msg)
		return m, cmd
	} else if m.worker1.Focused() {
		m.worker1, cmd = m.worker1.Update(msg)
		return m, cmd
	} else if m.psswd1.Focused() {
		m.psswd1, cmd = m.psswd1.Update(msg)
		return m, cmd
	} else if m.url2.Focused() {
		m.url2, cmd = m.url2.Update(msg)
		return m, cmd
	} else if m.worker2.Focused() {
		m.worker2, cmd = m.worker2.Update(msg)
		return m, cmd
	} else if m.psswd2.Focused() {
		m.psswd2, cmd = m.psswd2.Update(msg)
		return m, cmd
	} else if m.url3.Focused() {
		m.url3, cmd = m.url3.Update(msg)
		return m, cmd
	} else if m.worker3.Focused() {
		m.worker3, cmd = m.worker3.Update(msg)
		return m, cmd
	} else if m.psswd3.Focused() {
		m.psswd3, cmd = m.psswd3.Update(msg)
		return m, cmd
	} else if m.workerType.Focused() {
		m.workerType, cmd = m.workerType.Update(msg)
		return m, cmd
	}

	return m, cmd
}

func (m PoolModel) View() string {
	return m.url1.View() + "\n" + m.worker1.View() + "\n" + m.psswd1.View() + "\n" + m.url2.View() + "\n" + m.worker2.View() + "\n" + m.psswd2.View() + "\n" + m.url3.View() + "\n" + m.worker3.View() + "\n" + m.psswd3.View() + "\n" + m.workerType.View()
}
