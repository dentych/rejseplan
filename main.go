package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dentych/rejseplan/rejseplan"
)

type CurrentModel interface {
	Update(tea.Msg) (CurrentModel, tea.Cmd)
	View() string
}

type Step int

const (
	StepFirst Step = iota
	StepSecond
	StepThird
	StepError
)

type Model struct {
	currentModel CurrentModel
	step         Step
	stops        []rejseplan.StopLocation
	chosenStop   rejseplan.StopLocation
	departures   []rejseplan.Departure

	input textinput.Model
	list  list.Model
	table table.Model

	error string
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	if m.step == StepError {
		return m, nil
	}

	var cmd tea.Cmd
	switch m.step {
	case StepFirst:
		if !m.input.Focused() {
			return m, m.input.Focus()
		}
		if keymsg, ok := msg.(tea.KeyMsg); ok && keymsg.String() == "enter" {
			stops := rejseplan.GetStops(m.input.Value())
			if len(stops) == 1 {
				m.step = StepThird
				m.chosenStop = stops[0]
				return m, nil
			} else if len(stops) > 1 {
				m.step = StepSecond
				m.stops = stops
				m.input.Reset()
				defaultDelegate := list.NewDefaultDelegate()
				defaultDelegate.ShowDescription = false
				m.list = list.New(mapStopsToItems(stops), defaultDelegate, 70, 10)
				return m, nil
			} else {
				m.error = "No stops found"
				m.step = StepError
				return m, nil
			}
		}
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	case StepSecond:

	}

	// m.currentModel, cmd = m.currentModel.Update(msg)
	return m, cmd
}

func mapStopsToItems(stops []rejseplan.StopLocation) []list.Item {
	var output []list.Item
	for i := range stops {
		output = append(output, item{stops[i]})
	}
	return output
}

func (m Model) View() string {
	if m.step == StepError {
		return m.error + "\n\nPress q to exit."
	}

	switch m.step {
	case StepFirst:
		return "Find stop " + m.input.View()
	case StepSecond:
		var output string
		for _, stop := range m.stops {
			output += stop.Name + "\n"
		}
		return output
	}
	// return m.currentModel.View()
	return ""
}

func main() {
	model := Model{input: textinput.New()}
	t := tea.NewProgram(model)
	if _, err := t.Run(); err != nil {
		panic(err)
	}
}

func NewStationSearch() stationSearch {
	return stationSearch{
		input: textinput.New(),
	}
}

type stationSearch struct {
	input textinput.Model
}

func (m stationSearch) Update(msg tea.Msg) (CurrentModel, tea.Cmd) {
	if !m.input.Focused() {
		m.input.Focus()
	}
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "enter" {
		stops := rejseplan.GetStops(m.input.Value())
		if len(stops) > 1 {
			return m, nil
		}
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m stationSearch) View() string {
	return "Find stop " + m.input.View()
}

func newStopList(stops []rejseplan.StopLocation) stopList {
	return stopList{stops: stops}
}

type stopList struct {
	stops []rejseplan.StopLocation
}

type listDepartures struct {
}

func (m listDepartures) Update(msg tea.Msg) (CurrentModel, tea.Cmd) {
	return m, nil
}

func (m listDepartures) View() string {
	return ""
}

type item struct {
	rejseplan.StopLocation
}

func (i item) FilterValue() string {
	return i.StopLocation.Name
}
