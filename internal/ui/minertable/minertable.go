package minertable

import (
	"WMTUI/internal/miner"
	"fmt"

	"github.com/GridlessCompute/wmapi/client"
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

type Pool struct {
	URL      string
	Worker   string
	Password string
}

type CommandErrMsg struct{ Errors []error }

type RebootMsg struct{}

type SleepMsg struct{}

type PoolMsg struct{ Pools []Pool }

type LimitMsg struct {
	Vlt   float64
	Freq  float64
	Limit int
}

type WakeMsg struct{ Limit int }

type FastbootMsg struct{}

type SlowbootMsg struct{}

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
	case FastbootMsg:
		errs := m.fastboot()
		if len(errs) == 0 {
			return m, nil
		}
		return m, func() tea.Msg { return CommandErrMsg{Errors: errs} }
	case LimitMsg:
		errs := m.limit(msg.Freq, msg.Vlt, msg.Limit)
		if len(errs) == 0 {
			return m, nil
		}
		return m, func() tea.Msg { return CommandErrMsg{Errors: errs} }
	case MinerUpdateMsg:
		m.MinerList = msg.Miners
		m.Table.SetRows(makeNewRows(m.MinerList))
		return m, nil
	case PoolMsg:
		errs := m.pools(msg.Pools)
		if len(errs) == 0 {
			return m, nil
		}
		return m, func() tea.Msg { return CommandErrMsg{Errors: errs} }
	case RebootMsg:
		errs := m.reboot()
		if len(errs) == 0 {
			return m, nil
		}
		return m, func() tea.Msg { return CommandErrMsg{Errors: errs} }
	case SleepMsg:
		errs := m.sleep()
		if len(errs) == 0 {
			return m, nil
		}
		return m, func() tea.Msg { return CommandErrMsg{Errors: errs} }
	case SlowbootMsg:
		errs := m.slowboot()
		if len(errs) == 0 {
			return m, nil
		}
		return m, func() tea.Msg { return CommandErrMsg{Errors: errs} }
	case WakeMsg:
		errs := m.wake(msg.Limit)
		if len(errs) == 0 {
			return m, nil
		}
		return m, func() tea.Msg { return CommandErrMsg{Errors: errs} }
	}

	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m MinerTableModel) View() string {
	return m.Table.View()
}

func (m MinerTableModel) fastboot() []error {
	errors := []error{}
	for _, mnr := range m.MinerList {
		if mnr.Selected && mnr.Type == 1 {
			_, err := mnr.API.WM.Write.EnableFastboot()
			if err != nil {
				e := fmt.Errorf("failed to enable fastboot for %s: %w", mnr.IP, err)
				errors = append(errors, e)
			}
		}
	}
	return errors
}

func (m MinerTableModel) limit(freq, vlt float64, limit int) []error {
	errors := []error{}
	for _, mnr := range m.MinerList {
		if mnr.Selected {
			switch mnr.Type {
			case 1:
				_, err := mnr.API.WM.Write.AdjPowerLimit(limit)
				if err != nil {
					e := fmt.Errorf("failed to sleep %s: %w", mnr.IP, err)
					errors = append(errors, e)
				}
			case 2:
				if err := mnr.API.Epic.Post.Tune(freq, vlt); err != nil {
					e := fmt.Errorf("failed to sleep %s: %w", mnr.IP, err)
					errors = append(errors, e)
				}
			}
		}
	}
	return errors
}

func (m MinerTableModel) pools(p []Pool) []error {
	var pool1 Pool
	var pool2 Pool
	var pool3 Pool

	errors := []error{}

	switch len(p) {
	case 1:
		pool1.URL = p[0].URL
		pool1.Worker = p[0].Worker
		pool1.Password = p[0].Password
	case 2:
		pool1.URL = p[0].URL
		pool1.Worker = p[0].Worker
		pool1.Password = p[0].Password
		pool2.URL = p[1].URL
		pool2.Worker = p[1].Worker
		pool2.Password = p[1].Password
	case 3:
		pool1.URL = p[0].URL
		pool1.Worker = p[0].Worker
		pool1.Password = p[0].Password
		pool2.URL = p[1].URL
		pool2.Worker = p[1].Worker
		pool2.Password = p[1].Password
		pool3.URL = p[2].URL
		pool3.Worker = p[2].Worker
		pool3.Password = p[2].Password
	}

	for _, mnr := range m.MinerList {
		if mnr.Selected {
			switch mnr.Type {
			case 1:
				_, err := mnr.API.WM.Write.Pools(
					client.Pool{URL: pool1.URL, Worker: pool1.Worker, Password: pool1.Password},
					client.Pool{URL: pool2.URL, Worker: pool2.Worker, Password: pool2.Password},
					client.Pool{URL: pool3.URL, Worker: pool3.Worker, Password: pool3.Password},
				)
				if err != nil {
					e := fmt.Errorf("failed to set pools for %s: %w", mnr.IP, err)
					errors = append(errors, e)
				}
			case 2:
				// TODO: findout if epic has a pools endpoint, Dont know why it wouldnt
			}
		}
	}
	return errors
}

func (m MinerTableModel) reboot() []error {
	errors := []error{}
	for _, mnr := range m.MinerList {
		if mnr.Selected {
			switch mnr.Type {
			case 1:
				_, err := mnr.API.WM.Write.RebootSystem()
				if err != nil {
					e := fmt.Errorf("failed to reboot %s: %w", mnr.IP, err)
					errors = append(errors, e)
					continue
				}
			case 2:
				if err := mnr.API.Epic.Post.Reboot(0); err != nil {
					e := fmt.Errorf("failed to reboot %s: %w", mnr.IP, err)
					errors = append(errors, e)
					continue
				}
			}
		}
	}
	return errors
}

func (m MinerTableModel) sleep() []error {
	errors := []error{}
	for _, mnr := range m.MinerList {
		if mnr.Selected {
			switch mnr.Type {
			case 1:
				_, err := mnr.API.WM.Write.AdjPowerLimit(0)
				if err != nil {
					e := fmt.Errorf("failed to sleep %s: %w", mnr.IP, err)
					errors = append(errors, e)
				}
			case 2:
				if err := mnr.API.Epic.Post.Miner(false); err != nil {
					e := fmt.Errorf("failed to sleep %s: %w", mnr.IP, err)
					errors = append(errors, e)
				}
			}
		}
	}
	return errors
}

func (m MinerTableModel) slowboot() []error {
	errors := []error{}
	for _, mnr := range m.MinerList {
		if mnr.Selected && mnr.Type == 1 {
			_, err := mnr.API.WM.Write.Disablefastboot()
			if err != nil {
				e := fmt.Errorf("failed to disable fastboot for %s: %w", mnr.IP, err)
				errors = append(errors, e)
			}
		}
	}
	return errors
}

func (m MinerTableModel) wake(l int) []error {
	errors := []error{}
	for _, mnr := range m.MinerList {
		if mnr.Selected {
			switch mnr.Type {
			case 1:
				_, err := mnr.API.WM.Write.AdjPowerLimit(l)
				if err != nil {
					e := fmt.Errorf("failed to sleep %s: %w", mnr.IP, err)
					errors = append(errors, e)
				}
			case 2:
				if err := mnr.API.Epic.Post.Miner(true); err != nil {
					e := fmt.Errorf("failed to sleep %s: %w", mnr.IP, err)
					errors = append(errors, e)
				}
			}
		}
	}
	return errors
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
