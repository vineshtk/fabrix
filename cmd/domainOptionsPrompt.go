package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// networkCmd represents the network command
var domainOptionsPromptCmd = &cobra.Command{
	Use:     "domainOptionsPrompt",
	Aliases: []string{"dp"},
	Short:   "Use this command to remove all files for a network",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		domainOptions()
	},
}

func init() {
	rootCmd.AddCommand(domainOptionsPromptCmd)
}

type model3 struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model3) Init() tea.Cmd {
	return nil
}

func (m model3) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model3) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("You have choosen : %s", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Exit? ... See you later")
	}
	return "\n" + m.list.View()
}

func domainOptions() {
	items := []list.Item{
		item("Start network"),
		item("Info"),
		item("Deploy chaincode"),
		item("Down network"),
		item("Remove domain"),
		item("Home"),
		item("Exit"),
	}

	const defaultWidth = 20
	const domainsListHeight = 12

	l := list.New(items, itemDelegate{}, defaultWidth, domainsListHeight)
	l.Title = "Choose an option:"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	// l.Styles.HelpStyle = helpStyle

	m := model3{list: l}

	prog, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if m, ok := prog.(model3); ok {
		switch m.choice {
		case "Start network":
			rootCmd.SetArgs([]string{"up", choosenDomain})
			rootCmd.Execute()

		case "Info":
			rootCmd.SetArgs([]string{"list"})
			rootCmd.Execute()

		case "Deploy chaincode":
			rootCmd.SetArgs([]string{"deploy"})
			rootCmd.Execute()

		case "Down network":
			rootCmd.SetArgs([]string{"down", choosenDomain})
			rootCmd.Execute()

		case "Remove domain":
			rootCmd.SetArgs([]string{"remove", choosenDomain})
			rootCmd.Execute()
			
		case "Home":
			rootCmd.SetArgs([]string{"sp"})
			rootCmd.Execute()

		case "Exit":
			fmt.Println("Exiting...")
			os.Exit(0)
		}
	}
}
