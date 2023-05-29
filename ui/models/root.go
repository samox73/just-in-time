package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

type RootModel interface {
	Size() tea.WindowSizeMsg
	SetSize(size tea.WindowSizeMsg)

	View() string
	Height() int
}