package miner

import (
	"time"

	"github.com/GridlessCompute/epicapi"
	"github.com/GridlessCompute/wmapi"
)

type Miner struct {
	API        API
	Selected   bool
	Type       int
	IP         string
	Mac        string
	Status     string
	Errors     string
	UpTime     time.Duration
	Hashrate   float64
	Efficiency float64
	Power      float64
	PowerLimit float64
	Pool       string
}

type API struct {
	Epic *epicapi.EpicApi
	WM   *wmapi.WhatsminerMiddleware
}

func NewMiner(ip string, t int, api any) *Miner {
	switch api := api.(type) {
	case epicapi.EpicApi:
		if t == 2 {
			mnr := Miner{
				API: API{
					Epic: &api,
				},
				Type: 2,
				IP:   ip,
			}
			return &mnr
		}
		return nil
	case wmapi.WhatsminerMiddleware:
		if t == 1 {
			mnr := Miner{
				API: API{
					WM: &api,
				},
				Type: 1,
				IP:   ip,
			}
			return &mnr
		}
		return nil
	default:
		return nil
	}
}

func NewTestMiner(t int, ip string) *Miner {
	mnr := Miner{
		Selected:   false,
		Type:       t,
		IP:         ip,
		Mac:        "AB:CD:EF:01:23:45",
		Status:     "test",
		Errors:     "test",
		UpTime:     time.Duration(5 * time.Hour),
		Hashrate:   400.0,
		Efficiency: 27.4,
		Power:      12,
		PowerLimit: 3500,
		Pool:       "why?",
	}

	return &mnr
}
