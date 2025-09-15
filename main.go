package main

import (
	"WMTUI/internal/config"
	"WMTUI/internal/minerScanner"
	"WMTUI/internal/ui"
	"WMTUI/internal/ui/logging"
	"WMTUI/internal/ui/minertable"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// var ChosenFarm models.FarmStruct
//
// type model int
//
// var Logging *log.Logger
//
// const (
// 	TableView model = iota
// 	ListView
// 	FormView
// 	PoolView
// 	PowerView
// )
//
// var Models []tea.Model
// var SelectedModel tea.Model
//
// func PullSavedFarm() models.FarmStruct {
// 	// PUll FarmStruct From SAVEDFARM.json
// 	file, err := os.Open("SAVEDFARM.json")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	var farm models.FarmStruct
//
// 	byteValue, _ := io.ReadAll(file)
//
// 	json.Unmarshal(byteValue, &farm)
//
// 	if farm.Name == "" && farm.Start == "" && farm.End == "" {
// 		farm.Name = "Default"
// 		farm.Start = "192.168.0.0"
// 		farm.End = "192.168.0.255"
// 		return farm
// 	} else {
// 		return farm
// 	}
// }
//
// func SetSavedFarm(farm models.FarmStruct) {
// 	// PUSH FarmStruct To SAVEDFARM.json
// 	jsonBytes, _ := json.Marshal(farm)
// 	err := os.WriteFile("SAVEDFARM.json", jsonBytes, 0644)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }
//
// func main() {
// 	ChosenFarm = PullSavedFarm()
// 	Models = []tea.Model{models.NewTable(), models.NewList(), models.NewForm(), models.NewPoolForm(), models.NewPower()}
// 	SelectedModel = Models[TableView]
// 	p := tea.NewProgram(SelectedModel)
//
// 	if _, err := p.Run(); err != nil {
// 		fmt.Println("Error running program: ", err)
// 		os.Exit(1)
// 	}
// }

func main() {

	s := config.Site{
		Name:    "Test",
		IPRange: "10.20.0.0/24",
	}

	ms := minerScanner.Scanner{
		Conf: s,
	}

	mt := minertable.NewMinerTableModel()
	lm := logging.NewLogging()
	mm := ui.NewModel(&ms, lm, mt)
	p := tea.NewProgram(mm, tea.WithAltScreen())
	ms.SetProgram(p)

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("WORK IN PROGRESS")
}
