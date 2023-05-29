package scenes

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"

	"github.com/samox73/just-in-time/ui/styles"
	"github.com/samox73/just-in-time/ui/components"
	"github.com/samox73/just-in-time/ui/models"
)

var _ tea.Model = (*mainMenu)(nil)

type mainMenu struct {
	root   models.RootModel
	list   list.Model
	error  *components.ErrorComponent
	banner string
}

const logoBanner = `
███████╗██╗████████╗
╚════██║██║╚══██╔══╝
     ██║██║   ██║
    ██╔╝██║   ██║
█████╔╝ ██║   ██║
╚════╝  ╚═╝   ╚═╝`

func NewMainMenu(root models.RootModel) tea.Model {
	trimmedBanner := strings.TrimSpace(logoBanner)
	var finalBanner strings.Builder

	for i, s := range strings.Split(trimmedBanner, "\n") {
		if i > 0 {
			finalBanner.WriteRune('\n')
		}

		foreground := styles.LogoForegroundStyles[i]
		background := styles.LogoBackgroundStyles[i]

		for _, c := range s {
			if c == '█' {
				finalBanner.WriteString(foreground.Render("█"))
			} else if c != ' ' {
				finalBanner.WriteString(background.Render(string(c)))
			} else {
				finalBanner.WriteRune(c)
			}
		}
	}

	model := mainMenu{
		root:   root,
		banner: finalBanner.String(),
	}

	items := []list.Item{
		components.ListItem[mainMenu]{
			ItemTitle: "Log",
			Activate: func(msg tea.Msg, currentModel mainMenu) (tea.Model, tea.Cmd) {
				newModel := NewTaskLog(root, currentModel)
				return newModel, newModel.Init()
			},
		},
		components.ListItem[mainMenu]{
			ItemTitle: "Exit",
			Activate: func(msg tea.Msg, currentModel mainMenu) (tea.Model, tea.Cmd) {
				return nil, tea.Quit
			},
		},
	}

	model.list = list.New(items, list.NewDefaultDelegate(), root.Size().Width, root.Size().Height-root.Height())
	model.list.Title = "Main Menu"
	model.list.Styles = styles.ListStyles

	return model
}

func (m mainMenu) Init() tea.Cmd {
	return nil
}

func (m mainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			i, ok := m.list.SelectedItem().(components.ListItem[mainMenu])
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
			switch string(msg.Runes) {
			case "q":
				return m, tea.Quit
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

func (m mainMenu) View() string {
	header := m.root.View()

	banner := lipgloss.NewStyle().Margin(2, 0, 0, 2).Render(m.banner)

	commit := viper.GetString("commit")
	if len(commit) > 8 {
		commit = commit[:8]
	}

	version := "\n"
	version += styles.LabelStyle.Render("Version: ")
	version += viper.GetString("version") + " - " + commit

	header = lipgloss.JoinVertical(lipgloss.Left, version, header)

	totalHeight := lipgloss.Height(header) + len(m.list.Items()) + lipgloss.Height(banner) + 5
	if totalHeight < m.root.Size().Height {
		header = lipgloss.JoinVertical(lipgloss.Left, banner, header)
	}

	if m.error != nil {
		err := m.error.View()
		m.list.SetSize(m.list.Width(), m.root.Size().Height-lipgloss.Height(header)-lipgloss.Height(err))
		return lipgloss.JoinVertical(lipgloss.Left, header, err, m.list.View())
	}

	m.list.SetSize(m.list.Width(), m.root.Size().Height-lipgloss.Height(header))
	return lipgloss.JoinVertical(lipgloss.Left, header, m.list.View())
}
