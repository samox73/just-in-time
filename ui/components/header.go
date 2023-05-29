package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/samox73/just-in-time/ui/models"
	"github.com/samox73/just-in-time/ui/styles"
)

var _ tea.Model = (*headerComponent)(nil)

type headerComponent struct {
	root       models.RootModel
	labelStyle lipgloss.Style
}

func NewHeaderComponent(root models.RootModel) tea.Model {
	return headerComponent{
		root:       root,
		labelStyle: styles.LabelStyle,
	}
}

func (h headerComponent) Init() tea.Cmd {
	return nil
}

func (h headerComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return h, nil
}

func (h headerComponent) View() string {
	out := h.labelStyle.Render("User: ")
	out += "None"

	return lipgloss.NewStyle().Margin(1, 0).Render(out)
}
