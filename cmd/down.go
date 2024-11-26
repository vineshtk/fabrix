package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

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
		command.Stdout = io.Discard
		command.Stderr = io.Discard

		// Get stdout and stderr pipes
		stdout, err := command.StdoutPipe()
		if err != nil {
			return
		}
		stderr, err := command.StderrPipe()
		if err != nil {
			return
		}

		err = command.Run()
		if err != nil {
			fmt.Printf("Error executing script: %v\n", err)
			return
		}

		// Create a scanner to read stdout
		stdoutScanner := bufio.NewScanner(stdout)
		go func() {
			for stdoutScanner.Scan() {
				fmt.Printf("[INFO]: %s\n", stdoutScanner.Text()) // Replace with progress bar updates
			}
		}()

		// Create a scanner to read stderr
		stderrScanner := bufio.NewScanner(stderr)
		go func() {
			for stderrScanner.Scan() {
				fmt.Printf("[ERROR]: %s\n", stderrScanner.Text())
			}
		}()

		// Wait for the command to finish
		if err := command.Wait(); err != nil {
			return
		}

		return

	},
}

func init() {
	rootCmd.AddCommand(downCmd)
	// networkCmd.PersistentFlags().String("foo", "", "A help for foo")
	// networkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type scriptModel struct {
	progress []string
	errMsg   string
}

func (m scriptModel) Init() tea.Cmd {
	return func() tea.Msg {

		return "Script executed successfully."
	}
}

func (m scriptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case string:
		m.progress = append(m.progress, msg)
	}
	return m, nil
}

func (m scriptModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Top, m.progress...)
}
