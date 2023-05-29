package ui

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/samox73/just-in-time/ui/components"
	"github.com/samox73/just-in-time/ui/scenes"
)

type rootModel struct {
	headerComponent tea.Model
	currentSize     tea.WindowSizeMsg
	db              *sql.DB
}

func newModel() *rootModel {
	connStr := "postgres://postgres:test@localhost/jit?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	m := &rootModel{
		currentSize: tea.WindowSizeMsg{
			Width:  20,
			Height: 14,
		},
		db: db,
	}

	m.headerComponent = components.NewHeaderComponent(m)

	return m
}

func (m *rootModel) GetDBConnector() *sql.DB {
	return m.db
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
