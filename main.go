package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var ChosenFarm FarmStruct

type model int

var Logging *log.Logger

const (
	TableView model = iota
	ListView
	FormView
	PoolView
	PowerView
)

var Models []tea.Model
var SelectedModel tea.Model

func PullSavedFarm() FarmStruct {
	// PUll FarmStruct From SAVEDFARM.json
	file, err := os.Open("SAVEDFARM.json")
	if err != nil {
		fmt.Println(err)
	}
	var farm FarmStruct

	byteValue, _ := io.ReadAll(file)

	json.Unmarshal(byteValue, &farm)

	if farm.Name == "" && farm.Start == "" && farm.End == "" {
		farm.Name = "Default"
		farm.Start = "192.168.0.0"
		farm.End = "192.168.0.255"
		return farm
	} else {
		return farm
	}
}

func SetSavedFarm(farm FarmStruct) {
	// PUSH FarmStruct To SAVEDFARM.json
	jsonBytes, _ := json.Marshal(farm)
	err := os.WriteFile("SAVEDFARM.json", jsonBytes, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	ChosenFarm = PullSavedFarm()
	Models = []tea.Model{NewTable(), NewList(), NewForm(), NewPoolForm(), NewPower()}
	SelectedModel = Models[TableView]
	p := tea.NewProgram(SelectedModel)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	}
}
