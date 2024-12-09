package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// networkCmd represents the network command
var downCmd = &cobra.Command{
	Use:     "down",
	Aliases: []string{"dn"},
	Args:    cobra.ExactArgs(1),
	Short:   "Use this command to down a network and remove all the files associated",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {

		scriptPath := fmt.Sprintf("./fabrix/%v/Network/stopNetwork.sh", args[0])
		scriptDir := fmt.Sprintf("./fabrix/%v/Network/", args[0])

		err := os.Chmod(scriptPath, 0755)
		if err != nil {
			fmt.Printf("Error making script executable: %v\n", err)
			return
		}

		command := exec.Command("/bin/bash", "stopNetwork.sh")
		command.Dir = scriptDir

		p := tea.NewProgram(newModel())

		go func() {
			err = command.Start()
			if err != nil {
				fmt.Printf("Error starting script: %v\n", err)
				return
			}

			// Wait for the script to finish
			err = command.Wait()
			if err != nil {
				fmt.Printf("Error executing script: %v\n", err)
			}

			// Stop the spinner
			p.Send(stopMsg{})
		}()

		if _, err := p.Run(); err != nil {
			fmt.Println("could not run program:", err)
			os.Exit(1)
		}

		fmt.Println("Network started successfully!")
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}

var (
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	appStyle     = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

type stopMsg struct{}

type model4 struct {
	spinner  spinner.Model
	quitting bool
}

func newModel() model4 {
	s := spinner.New()
	s.Style = spinnerStyle
	return model4{
		spinner: s,
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

func (m model4) View() string {
	var s string

	if m.quitting {
		s += "Quitting...!"
	} else {
		s += m.spinner.View() + m.spinner.View() + m.spinner.View() + " Please wait! Clearing the network for you!!..."
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
