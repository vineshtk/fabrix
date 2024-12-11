package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// networkCmd represents the network command
var deployPromptCmd = &cobra.Command{
	Use:     "deployPrompt",
	Aliases: []string{"dccp"},
	Short:   "Use this command to deploy chaincode.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test all codes here")
		fp := filepicker.New()
		// fp.AllowedTypes = []string{".json", ".go"}
		fp.CurrentDirectory, _ = os.Getwd()
		fp.DirAllowed = true
		fp.FileAllowed = false
		fp.ShowSize = true
		fp.ShowPermissions = false
		m := model5{
			filepicker: fp,
		}
		tm, _ := tea.NewProgram(&m).Run()
		mm := tm.(model5)
		fmt.Println("\n  You selected: " + m.filepicker.Styles.Selected.Render(mm.selectedFolder) + "\n")
	},
}

func init() {
	rootCmd.AddCommand(deployPromptCmd)
}

type model5 struct {
	filepicker     filepicker.Model
	selectedFolder string
	quitting       bool
	err            error
}

type clearErrorMsg struct{}

// func clearErrorAfter(t time.Duration) tea.Cmd {
// 	return tea.Tick(t, func(_ time.Time) tea.Msg {
// 		return clearErrorMsg{}
// 	})
// }

func (m model5) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m model5) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the relative path of the selected file.
		currentDir, err := os.Getwd()
		fmt.Println("now at", currentDir)
		if err != nil {
			m.err = errors.New("unable to determine the current directory")
			return m, cmd
		}
		relativePath, err := filepath.Rel(currentDir, path)
		if err != nil {
			m.err = errors.New("unable to determine the relative path")
			return m, cmd
		}
 
		m.selectedFolder = relativePath
	}

	// // Did the user select a disabled file?
	// // This is only necessary to display an error to the user.
	// if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
	// 	// Let's clear the selectedFile and display an error.
	// 	m.err = errors.New(path + " is not valid.")
	// 	m.selectedFolder = ""
	// 	return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	// }

	return m, cmd
}

func (m model5) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selectedFolder == "" {
		s.WriteString("Pick the chaincode folder:")
	} else {
		s.WriteString("Selected folder: " + m.filepicker.Styles.Selected.Render(m.selectedFolder))
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return s.String()
}
