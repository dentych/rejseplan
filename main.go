package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/dentych/rejseplan/rejseplan"
	"github.com/dentych/rejseplan/simplelist"
)

type CurrentModel interface {
	Update(tea.Msg) (CurrentModel, tea.Cmd)
	View() string
}

type Step int

const (
	StepOne Step = iota
	StepTwo
	StepThree
	StepError
)

type Model struct {
	currentModel CurrentModel
	step         Step
	stops        []rejseplan.Stop
	chosenStop   rejseplan.Stop
	departures   []rejseplan.Departure

	input textinput.Model
	list  simplelist.Model[rejseplan.Stop]
	table *table.Table

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
	case StepOne:
		if !m.input.Focused() {
			return m, m.input.Focus()
		}
		if keymsg, ok := msg.(tea.KeyMsg); ok && keymsg.String() == "enter" {
			stops := rejseplan.GetStops(m.input.Value())
			if len(stops) == 1 {
				m.step = StepThree
				m.chosenStop = stops[0]
				return m, nil
			} else if len(stops) > 1 {
				m.step = StepTwo
				m.stops = stops
				m.input.Reset()
				m.list = simplelist.Model[rejseplan.Stop]{Items: m.stops}
				return m, nil
			} else {
				m.error = "No stops found"
				m.step = StepError
				return m, nil
			}
		}
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	case StepTwo:
		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "enter" {
			m.chosenStop = m.list.Current()
			m.step = StepThree
			m.departures = rejseplan.GetDepartures(m.chosenStop)
			m.table = table.New().Border(lipgloss.NormalBorder()).
				BorderLeft(false).BorderRight(false).BorderTop(false).BorderBottom(false).
				Headers("LINE", "DIRECTION", "TIME").
				Rows(mapDepartures(m.departures)...)
			return m, tea.Quit
		}
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	// m.currentModel, cmd = m.currentModel.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.step == StepError {
		return m.error + "\n\nPress q to exit."
	}

	switch m.step {
	case StepOne:
		return "Find stop " + m.input.View()
	case StepTwo:
		s := lipgloss.NewStyle().Background(lipgloss.Color("#140c7d")).Foreground(lipgloss.Color("#fff")).MarginLeft(2).Padding(1).Render("Choose a stop")
		s2 := lipgloss.NewStyle().MarginTop(1).Render(m.list.View())
		return lipgloss.JoinVertical(lipgloss.Top, s, s2)
	case StepThree:
		return m.table.Render() + "\n"
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

func newStopList(stops []rejseplan.Stop) stopList {
	return stopList{stops: stops}
}

type stopList struct {
	stops []rejseplan.Stop
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
	rejseplan.Stop
}

func (i item) FilterValue() string {
	return i.Stop.Name
}

func mapDepartures(departures []rejseplan.Departure) [][]string {
	var output [][]string
	for _, departure := range departures {
		output = append(output, []string{departure.Line, departure.Direction, departure.PlannedTime})
	}
	return output
}

func mapStopsToItems(stops []rejseplan.Stop) []string {
	var output []string
	for i := range stops {
		output = append(output, stops[i].Name)
	}
	return output
}
