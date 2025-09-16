package scanner

import (
	"WMTUI/internal/config"
	"WMTUI/internal/miner"
	"WMTUI/internal/ui/logging"
	"WMTUI/internal/ui/table"
	"context"
	"fmt"
	"net"
	"net/netip"
	"sync"
	"time"

	"github.com/GridlessCompute/epicapi"
	"github.com/GridlessCompute/wmapi"
	tea "github.com/charmbracelet/bubbletea"
)

type Scanner struct {
	Conf        config.Site
	Machines    []*miner.Miner
	RefreshTime time.Duration
	Program     *tea.Program
}

func NewScanner(c config.Site, r int) Scanner {
	return Scanner{
		Conf:        c,
		Machines:    []*miner.Miner{},
		RefreshTime: time.Duration(r * int(time.Second)),
		Program:     nil,
	}
}

func (s *Scanner) SendMsg(msg tea.Msg) {
	s.Program.Send(msg)
}

func (s *Scanner) SendLog(e bool, str string) {
	s.SendMsg(logging.LoggingMsg{
		Err:     e,
		Message: str,
	})
}

func (s *Scanner) SetProgram(p *tea.Program) {
	s.Program = p
}

func (s *Scanner) ScanForMachines() {
	prefix, err := netip.ParsePrefix(s.Conf.IPRange)
	if err != nil {
		// fmt.Println(err)
		s.SendLog(true, err.Error())
		return
	}

	machinesChan := make(chan *miner.Miner)
	var wg sync.WaitGroup

	addr := prefix.Addr()
	for prefix.Contains(addr) {
		wg.Add(1)
		go func(ip netip.Addr) {
			defer wg.Done()
			m := miner.Miner{
				IP: ip.String(),
			}
			DetermineMinerType(&m, ip.String())
			if m.Type != 0 {
				machinesChan <- &m
			}
		}(addr)
		addr = addr.Next()
	}

	go func() {
		wg.Wait()
		close(machinesChan)
	}()

	for m := range machinesChan {
		s.Machines = append(s.Machines, m)
	}

	s.SendMsg(table.ScanDoneMsg{})

}

func GetWhatsminerInfo(m *miner.Miner) error {
	info, err := m.API.WM.Read.MinerInfo()
	if err == nil {
		m.IP = info.Msg.IP
		m.Mac = info.Msg.Mac
	} else {
		return err
	}

	summary, err := m.API.WM.Read.Summary()
	if err == nil {
		if len(summary.SUMMARY) > 0 {
			m.Hashrate = summary.SUMMARY[0].HSRT
			m.Power = float64(summary.SUMMARY[0].Power)
			m.PowerLimit = float64(summary.SUMMARY[0].PowerLimit)
			m.Efficiency = summary.SUMMARY[0].PowerRate
			m.UpTime = time.Duration(summary.SUMMARY[0].Uptime)
		} else {
			m.Hashrate = 0
			m.Power = 0
			m.PowerLimit = 0
			m.Efficiency = 0
			m.UpTime = 0
		}
	} else {
		return err
	}

	pools, err := m.API.WM.Read.Pools()
	if err == nil {
		if len(pools.POOLS) > 0 {
			m.Pool = pools.POOLS[0].URL
		} else {
			m.Pool = "Blank"
		}
	} else {
		return err
	}

	return nil
}

func GetEpicInfo(m *miner.Miner) error {
	network, err := m.API.Epic.Get.Network()
	if err == nil {
		m.IP = network.Dhcp.Address
		m.Mac = network.Dhcp.MacAddress
	} else {
		return err
	}

	summary, err := m.API.Epic.Get.Summary()
	if err == nil {
		m.Power = summary.PowerSupplyStats.OutputVoltage
		m.PowerLimit = 0
		m.Efficiency = 0
		m.UpTime = time.Duration(summary.Session.Uptime)
	} else {
		return err
	}

	// hashrate, err := m.EAPI.Get.Hashrate()
	// if err == nil {
	// 	h := 0
	// 	for i := 0; i >= len(hashrate); i++ {
	// 		if len(hashrate[i].Total) > 0 {
	// 			h += int(hashrate[i].Total[0])
	// 		}
	// 	}
	// 	m.Miner.Hashrate = float64(h)
	// }

	return nil
}

func (s *Scanner) RefreshLoop(ctx context.Context) {
	t := time.NewTicker(s.RefreshTime)

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			s.RefreshMachineInfo()
		}
	}
}

func (s *Scanner) RefreshMachineInfo() {
	var wg sync.WaitGroup

	ch := make(chan *miner.Miner, len(s.Machines))

	for _, m := range s.Machines {
		switch m.Type {
		case 1:
			wg.Go(func() {
				err := GetWhatsminerInfo(m)
				if err != nil {
					s.SendLog(true, err.Error())
					return
				}
				ch <- m
			})
		case 2:
			wg.Go(func() {
				err := GetEpicInfo(m)
				if err != nil {
					s.SendLog(true, err.Error())
					return
				}
				ch <- m
			})
		}
	}

	wg.Wait()
	close(ch)

	s.Machines = []*miner.Miner{}
	for m := range ch {
		s.Machines = append(s.Machines, m)
	}

	msg := table.MinerUpdateMsg{}
	for _, m := range s.Machines {
		msg.Miners = append(msg.Miners, m)
	}

	s.SendMsg(msg)
}

func DetermineMinerType(m *miner.Miner, ip string) {
	err := isIPEpic(m, ip)
	if err == nil {
		m.Type = 2
		return
	}

	err = isIpWhatsminer(m, ip)
	if err == nil {
		m.Type = 1
		return
	}

	m.Type = 0
}

func isIPEpic(m *miner.Miner, ip string) error {
	m.API.Epic = epicapi.New(ip, "4028", "letmein")
	if err := m.API.Epic.Post.Authenticate(); err != nil {
		m.API.Epic = nil
		return err
	}

	return nil
}

func isIpWhatsminer(m *miner.Miner, ip string) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, 4028), 1*time.Second)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return fmt.Errorf("connection timed out for ip %s", ip)
		}
		return fmt.Errorf("error connecting: %w", err)
	}
	conn.Close()

	a, err := wmapi.NewWhatsminerAPI(ip, 4028, "admin")
	if err != nil {
		return fmt.Errorf("error creating middleware for ip %s: %w", ip, err)
	}

	m.API.WM = a
	return nil
}
