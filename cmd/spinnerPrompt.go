package cmd

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	appStyle     = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

type stopMsg struct{}

type model4 struct {
	spinner  spinner.Model
	quitting bool
	message  string
}

func newModel(msg string) model4 {
	s := spinner.New()
	s.Style = spinnerStyle
	return model4{
		spinner: s,
        message: msg,
	}
}

func (m model4) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model4) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case stopMsg:
		m.quitting = true
		return m, tea.Quit

	case tea.KeyMsg:
        switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
        }
        return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	default:
		return m, nil
	}
}

func (m model4) View() string {
	var s string

	if m.quitting {
		s += "Quitting...!"
	} else {
		s += m.spinner.View() + m.spinner.View() + m.spinner.View() + m.message
	}
	s += "\n\n"

	// if !m.quitting {
	// 	s += helpStyle.Render("Press any key to exit")
	// }

	// if m.quitting {
	// 	s += "\n"
	// }

	return appStyle.Render(s)
}
