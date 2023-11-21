package simplelist

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Item interface {
	Title() string
}

type Model[T Item] struct {
	Items []T

	currentChoice int
}

// Update implements tea.Model.
func (m Model[T]) Update(msg tea.Msg) (Model[T], tea.Cmd) {
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
func (m Model[T]) View() string {
	var output string
	for i, item := range m.Items {
		if i == m.currentChoice {
			output += "> "
		} else {
			output += "  "
		}

		output += item.Title() + "\n"
	}

	return output
}

func (m Model[T]) Current() T {
	return m.Items[m.currentChoice]
}
