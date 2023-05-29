package components

import tea "github.com/charmbracelet/bubbletea"

type ListItem[T tea.Model] struct {
	Activate  func(msg tea.Msg, currentModel T) (tea.Model, tea.Cmd)
	ItemTitle string
	ItemDescription string
}

func (n ListItem[any]) Title() string {
	return n.ItemTitle
}

func (n ListItem[any]) FilterValue() string {
	return n.ItemTitle
}

func (n ListItem[any]) Description() string {
	return n.ItemDescription
}
