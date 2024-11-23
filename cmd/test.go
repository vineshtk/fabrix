package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type model struct {
	altscreen  bool
	quitting   bool
	suspending bool
}

var (
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("204")).Background(lipgloss.Color("235"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

// networkCmd represents the network command
var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"tst"},
	Short:   "Use this command to remove all files for a network",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := tea.NewProgram(model{}).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	// networkCmd.PersistentFlags().String("foo", "", "A help for foo")
	// networkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.ResumeMsg:
		m.suspending = false
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "ctrl+z":
			m.suspending = true
			return m, tea.Suspend
		case " ":
			var cmd tea.Cmd
			if m.altscreen {
				cmd = tea.ExitAltScreen
			} else {
				cmd = tea.EnterAltScreen
			}
			m.altscreen = !m.altscreen
			return m, cmd
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.suspending {
		return ""
	}

	if m.quitting {
		return "Bye!\n"
	}

	const (
		altscreenMode = " altscreen mode "
		inlineMode    = " inline mode "
	)

	var mode string
	if m.altscreen {
		mode = altscreenMode
	} else {
		mode = inlineMode
	}

	return fmt.Sprintf("\n\n  You're in %s\n\n\n", keywordStyle.Render(mode)) +
		helpStyle.Render("  space: switch modes • ctrl-z: suspend • q: exit\n")
}

