package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	p := tea.NewProgram(initModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type model struct {
	list    []string
	cursor  int
	input   textinput.Model
	isInput bool
}

func initModel() model {
	ti := textinput.New()
	ti.Placeholder = "To-do..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		list:    []string{"hello world", "salut ça va"},
		cursor:  0,
		input:   ti,
		isInput: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Keybinds in new todo modal
		if m.isInput {
			switch msg.String() {

			case "enter":
				if m.isInput {
					m = m.createTodo()
					return m, cmd
				}

			case "esc":
				if m.isInput {
					m.isInput = false
					m.input.SetValue("")
					m.input.Blur()
				}
			}
		} else { // Keybinds in main view
			switch msg.String() {

			case "q":
				return m, tea.Quit

			case "up", "k":
				if !m.isInput && m.cursor > 0 {
					m.cursor--
				}

			case "down", "j":
				if !m.isInput && m.cursor < len(m.list)-1 {
					m.cursor++
				}

			case "n":
				if !m.isInput {
					m.isInput = true
					m.input.Focus()
					return m, cmd
				}
			}
		}

		// Global keybinds
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	if m.isInput {
		m.input, cmd = m.input.Update(msg)
	}

	return m, cmd
}

func (m model) createTodo() model {
	todoText := m.input.Value()
	if todoText != "" {
		m.list = append(m.list, todoText)
	}
	m.input.SetValue("")
	m.isInput = false
	m.input.Blur()
	return m
}

func (m model) View() string {
	var styleTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#d787ff")).
		Margin(1, 4, 1, 4).
		Padding(0, 1, 0, 1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#d787ff"))
	var styleHint = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#808080"))

	s := styleTitle.Render("TODO")

	if m.isInput {
		s += "\n"
		s += m.input.View()
		s += styleHint.Render("\n\n[Press Enter to add, Esc to cancel]")
	} else {
		s += "\n"
		for i, choice := range m.list {

			// Is the cursor pointing at this choice?
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = ">" // cursor!
			}

			// Render the row
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		s += styleHint.Render("\n[n to add a todo - q to quit.]\n")
	}
	// Send the UI for rendering
	return s
}
