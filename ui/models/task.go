package models

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samox73/just-in-time/internal/models"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type Task struct {
	root       RootModel
	parent     tea.Model
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func NewTask(root RootModel, parent tea.Model) tea.Model {
	t := &Task{
		root:   root,
		inputs: make([]textinput.Model, 2),
		parent: parent,
	}
	for i := range t.inputs {
		ti := textinput.New()
		ti.CursorStyle = cursorStyle
		ti.CharLimit = 32
		switch i {
		case 0:
			ti.Placeholder = "Title"
			ti.Focus()
			ti.PromptStyle = focusedStyle
			ti.TextStyle = focusedStyle
		case 1:
			ti.Placeholder = "Description"
			ti.CharLimit = 64
		}
		t.inputs[i] = ti
	}
	return t
}

func (t Task) Init() tea.Cmd {
	return textinput.Blink
}

func (t Task) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log := log.Default()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return t, tea.Quit
		case tea.KeyEnter, tea.KeyTab, tea.KeyDown:
			// submit
			if msg.Type == tea.KeyEnter && t.focusIndex == len(t.inputs) {
				task := models.Task{
					Name:        t.inputs[0].Value(),
					Description: t.inputs[1].Value(),
				}
				db := t.root.GetDBConnector()
				statement := `insert into "tasks"("name", "description") values($1, $2)`
				_, err := db.Exec(statement, task.Name, task.Description)
				if err != nil {
					panic(err)
				}
				return t.parent, nil
			}
			t.focusIndex++
			if t.focusIndex > len(t.inputs) {
				log.Print("encountered overflow")
				t.focusIndex = 0
			}

			cmds := make([]tea.Cmd, len(t.inputs))
			for i := range t.inputs {
				if i == t.focusIndex {
					// set focused state
					cmds[i] = t.inputs[i].Focus()
					t.inputs[i].PromptStyle = focusedStyle
					t.inputs[i].TextStyle = focusedStyle
					continue
				}
				// remove focus
				t.inputs[i].Blur()
				t.inputs[i].PromptStyle = noStyle
				t.inputs[i].TextStyle = noStyle
			}
			return t, tea.Batch(cmds...)
		}
	}

	cmd := t.updateInputs(msg)
	return t, cmd
}

func (t *Task) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(t.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range t.inputs {
		t.inputs[i], cmds[i] = t.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (t Task) View() string {
	var b strings.Builder

	for i := range t.inputs {
		b.WriteString(t.inputs[i].View())
		if i < len(t.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if t.focusIndex == len(t.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(t.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}
