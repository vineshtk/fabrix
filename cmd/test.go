package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	appStyle      = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

type stopMsg struct{}

type model struct {
	spinner  spinner.Model
	quitting bool
}

func newModel() model {
	s := spinner.New()
	s.Style = spinnerStyle
	return model{
		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case stopMsg:
		m.quitting = true
		return m, tea.Quit

	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	default:
		return m, nil
	}
}

func (m model) View() string {
	var s string

	if m.quitting {
		s += "Quitting...!"
	} else {
		s += m.spinner.View() + "Please wait! Clearing the network for you!!..."
	}
	s += "\n\n"

	if !m.quitting {
		s += helpStyle.Render("Press any key to exit")
	}

	if m.quitting {
		s += "\n"
	}

	return appStyle.Render(s)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Use this command to down a network and remove all the files associated",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(newModel())

		// Simulate activity
		// go func() {
		// 	for {
		// 		pause := time.Duration(rand.Int63n(899)+100) * time.Millisecond // nolint:gosec
		// 		time.Sleep(pause)

		// 		// Send the Bubble Tea program a message from outside the
		// 		// tea.Program. This will block until it is ready to receive
		// 		// messages.
		// 		p.Send(resultMsg{food: randomFood(), duration: pause})
		// 	}
		// }()

		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
