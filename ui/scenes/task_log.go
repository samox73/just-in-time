package scenes

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/samox73/just-in-time/ui/components"
	"github.com/samox73/just-in-time/ui/models"
)

var _ tea.Model = (*taskLog)(nil)

type taskLog struct {
	root   models.RootModel
	list   list.Model
	parent tea.Model
}

func NewTaskLog(root models.RootModel, parent tea.Model) tea.Model {
	l := list.New(tasksToList(root), list.NewDefaultDelegate(), root.Size().Width, root.Size().Height-root.Height())
	return &taskLog{
		root:   root,
		list:   l,
		parent: parent,
	}
}

func (m taskLog) Init() tea.Cmd {
	return nil
}

func (m taskLog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// List enables its own keybindings when they were previously disabled
	m.list.DisableQuitKeybindings()

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m taskLog) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.root.View(), m.list.View())
}

func tasksToList(root models.RootModel) []list.Item {
	items := []list.Item{
		components.ListItem[taskLog]{
			ItemTitle:       "Upgrade EBS Volumes",
			ItemDescription: "Affected customers: asdf, qwer, lkjh",
		},
		components.ListItem[taskLog]{
			ItemTitle:       "Work on Pboperator",
			ItemDescription: "Fix bugs",
		},
	}

	return items
}
