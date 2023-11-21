package simplelist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dentych/rejseplan/rejseplan"
)

type Model struct {
	Items []rejseplan.StopLocation

	currentChoice int
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down", "j":
			m.currentChoice += 1
			if m.currentChoice >= len(m.Items) {
				m.currentChoice = len(m.Items) - 1
			}
		case "up", "k":
			m.currentChoice -= 1
			if m.currentChoice < 0 {
				m.currentChoice = 0
			}
		}
	}
	return m, nil
}

// View implements tea.Model.
func (m Model) View() string {
	var output string
	for i, item := range m.Items {
		if i == m.currentChoice {
			output += "> "
		} else {
			output += "  "
		}

		output += item.Name + "\n"
	}

	return output
}

func (m Model) Current() rejseplan.StopLocation {
	return m.Items[m.currentChoice]
}
