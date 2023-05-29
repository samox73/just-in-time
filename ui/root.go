package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/errors"
	"github.com/samox73/just-in-time/ui/components"
	"github.com/samox73/just-in-time/ui/scenes"
)

type rootModel struct {
	headerComponent tea.Model
	currentSize     tea.WindowSizeMsg
}

func newModel() *rootModel {
	m := &rootModel{
		currentSize: tea.WindowSizeMsg{
			Width:  20,
			Height: 14,
		},
	}

	m.headerComponent = components.NewHeaderComponent(m)

	return m
}

func (m *rootModel) Size() tea.WindowSizeMsg {
	return m.currentSize
}

func (m *rootModel) SetSize(size tea.WindowSizeMsg) {
	m.currentSize = size
}

func (m *rootModel) View() string {
	return m.headerComponent.View()
}

func (m *rootModel) Height() int {
	return lipgloss.Height(m.View()) + 1
}

func RunTea() error {
	if err := tea.NewProgram(scenes.NewMainMenu(newModel()), tea.WithAltScreen(), tea.WithMouseCellMotion()).Start(); err != nil {
		return errors.Wrap(err, "internal tea error")
	}
	return nil
}
