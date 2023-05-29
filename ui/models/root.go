package models

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
)

type RootModel interface {
	Size() tea.WindowSizeMsg
	SetSize(size tea.WindowSizeMsg)
	GetDBConnector() *sql.DB

	View() string
	Height() int
}
