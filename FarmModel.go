package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type FarmStruct struct {
	Name  string
	Start string
	End   string
}

type FarmList struct {
	Farms []FarmStruct
}

type FarmModel struct {
	Farms         FarmList
	FarmSelection list.Model
}

// Farm Functions ///
func (f FarmStruct) FilterValue() string {
	return f.Name
}

func (f FarmStruct) Title() string {
	return f.Name
}

func (f FarmStruct) Description() string {
	return f.Start + " -> " + f.End
}

func (m FarmModel) SetFarm() tea.Msg {
	var f SetFarmMsg
	return f
}

// List Functions //
func (m *FarmModel) GetFarmsFromJson() {
	file, err := os.Open("FARMS.json")
	if err != nil {
		fmt.Println(err)
	}
	var farms FarmList
	var farmsList FarmList

	byteValue, _ := io.ReadAll(file)

	json.Unmarshal(byteValue, &farms)

	farmsList.Farms = append(farmsList.Farms, farms.Farms...)

	m.Farms = farmsList
}

func (m FarmModel) LoadList() tea.Msg {
	return m
}

func (m *FarmModel) InitFarms(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	defaultList.SetShowHelp(true)
	m.FarmSelection = defaultList
	m.FarmSelection.Title = "Farms"

	farms := m.Farms

	index := 0
	for _, farm := range farms.Farms {
		m.FarmSelection.InsertItem(index, farm)
		index = index + 1
	}

}

func NewList() *FarmModel {
	farm := &FarmModel{}
	farm.GetFarmsFromJson()
	farm.InitFarms(75, 50)
	return farm
}

func (m *FarmModel) NewFarm(name, root, length string) {
	newFarm := FarmStruct{Name: name, Start: root, End: length}
	m.Farms.Farms = append(m.Farms.Farms, newFarm)
	m.FarmSelection.InsertItem(len(m.FarmSelection.Items()), newFarm)
}

func (m *FarmModel) SaveFarm() {
	farmsList := m.Farms

	jsonBytes, _ := json.Marshal(farmsList)
	err := os.WriteFile("FARMS.json", jsonBytes, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (m *FarmModel) DeleteFarm() {
	var farmList FarmList
	var newList FarmList

	farmList = m.Farms

	farm := m.FindFarm(m.FarmSelection.SelectedItem())

	for _, f := range farmList.Farms {
		if f.Name != farm.Name {
			newList.Farms = append(newList.Farms, f)
		}
	}

	m.FarmSelection.RemoveItem(m.FarmSelection.Cursor())

	jsonBytes, _ := json.Marshal(newList)
	err := os.WriteFile("FARMS.json", jsonBytes, 0644)
	if err != nil {
		fmt.Println(err)
	}

}

func (m *FarmModel) FindFarm(i list.Item) FarmStruct {
	///FUCK, MIND DONT WORK!
	fv := i.FilterValue()

	for _, farm := range m.Farms.Farms {
		if farm.Name == fv {
			return farm
		}
	}
	return FarmStruct{}

}

func (m *FarmModel) Init() tea.Cmd {
	return nil
}

func (m *FarmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		// m.GetFarmsFromJson()
		// m.InitFarms(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			Models[ListView] = m
			return Models[TableView].Update(tea.WindowSizeMsg{})
		case "q", "ctrl+c":
			return m, tea.Quit
		case "d":
			m.DeleteFarm()
			return m, nil
		case "enter":
			farm := m.FindFarm(m.FarmSelection.SelectedItem())
			if farm.Name != "" {
				ChosenFarm = farm
				SetSavedFarm(farm)
			}
			Models[ListView] = m
			return Models[TableView], m.SetFarm
		}

	case FarmStruct:
		m.NewFarm(msg.Name, msg.Start, msg.End)
		m.SaveFarm()
		// m.ReloadFarms()
		Models[ListView] = m
		return Models[ListView].Update(nil)
	}

	m.FarmSelection, cmd = m.FarmSelection.Update(msg)
	return m, cmd
}

func (m *FarmModel) View() string {
	return m.FarmSelection.View()
}
