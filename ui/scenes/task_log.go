package scenes

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/samox73/just-in-time/internal/models"
	"github.com/samox73/just-in-time/ui/components"
	uimodels "github.com/samox73/just-in-time/ui/models"
)

var _ tea.Model = (*taskLog)(nil)

type taskLog struct {
	root   uimodels.RootModel
	list   list.Model
	parent tea.Model
}

func NewTaskLog(root uimodels.RootModel, parent tea.Model) tea.Model {
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

func (m *taskLog) AddTask(task models.Task) {
	m.list.InsertItem(len(m.list.Items()), components.ListItem[uimodels.Task]{
		ItemTitle:       task.Name,
		ItemDescription: task.Description,
	})
}

func (m taskLog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			i, ok := m.list.SelectedItem().(components.ListItem[taskLog])
			if ok {
				if i.Activate != nil {
					newModel, cmd := i.Activate(msg, m)
					if newModel != nil || cmd != nil {
						if newModel == nil {
							newModel = m
						}
						return newModel, cmd
					}
					return m, nil
				}
			}
			return m, nil
		case tea.KeyRunes:
			switch keystring := msg.String(); keystring {
			case "q":
				return m.parent, nil
			case "n":
				newTaskModel := uimodels.NewTask(m.root, m)
				return newTaskModel, newTaskModel.Init()
			default:
				var cmd tea.Cmd
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			}
		default:
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}
	case tea.WindowSizeMsg:
		top, right, bottom, left := lipgloss.NewStyle().Margin(2, 2).GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom)
		m.root.SetSize(msg)
	}

	return m, nil
}

func (m taskLog) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.root.View(), m.list.View())
}

func tasksToList(root uimodels.RootModel) []list.Item {
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
