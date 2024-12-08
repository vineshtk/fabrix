package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// networkCmd represents the network command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "Use this command to list all available domains.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		chooseDomain()
		domainOptions()
	},
}

var choosenDomain string

func init() {
	rootCmd.AddCommand(listCmd)
}

type model2 struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model2) Init() tea.Cmd {
	return nil
}

func (m model2) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model2) View() string {
	if m.choice != "" {
		choosenDomain = m.choice
		return quitTextStyle.Render(fmt.Sprintf("Ok, Choosen Network: %s", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Exit? ... See you later")
	}
	return "\n" + m.list.View()
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

			case "Go Back":
				rootCmd.SetArgs([]string{"list"})
				rootCmd.Execute()

			case "Exit":
				fmt.Println("Exiting...")
				os.Exit(0)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model3) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Ok let's continue", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Exit? ... See you later")
	}
	return "\n" + m.list.View()
}

func chooseDomain() {

	rootDir := "./fabrix"
	// Read the directory
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
	}
	var domains []string
	// List directories
	for _, entry := range entries {
		if entry.IsDir() {
			domains = append(domains, entry.Name())
		}
	}

	items := []list.Item{}

	// Append domain names as selectable items
	for _, domain := range domains {
		items = append(items, item(domain))
	}

	const defaultWidth = 20
	const listHeight = 10

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Please select the Domain:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	// l.Styles.HelpStyle = helpStyle

	m := model2{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func domainOptions() {
	items := []list.Item{
		item("Start network"),
		item("Info"),
		item("Deploy chaincode"),
		item("Down network"),
		item("Remove domain"),
		item("Go Back"),
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

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}





// func domainOptions(domainName string) {
// 	options := []string{"Start network", "Info", "Deploy chaincode", "Down network", "Remove domain", "Go Back", "Exit"}

// 	var selectedOption string
// 	prompt := &survey.Select{
// 		Message: "Choose an option:",
// 		Options: options,
// 	}

// 	if err := survey.AskOne(prompt, &selectedOption); err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}

// switch selectedOption {

// case "Start network":
// 	rootCmd.SetArgs([]string{"up", domainName})
// 	rootCmd.Execute()

// case "Info":
// 	rootCmd.SetArgs([]string{"list"})
// 	rootCmd.Execute()

// case "Deploy chaincode":
// 	rootCmd.SetArgs([]string{"deploy"})
// 	rootCmd.Execute()

// case "Down network":
// 	rootCmd.SetArgs([]string{"down", domainName})
// 	rootCmd.Execute()

// case "Remove domain":
// 	confirm := false
// 	prompt := &survey.Confirm{
// 		Message: fmt.Sprintf("Do you want to remove '%s'?", domainName),
// 	}
// 	survey.AskOne(prompt, &confirm)
// 	if confirm {
// 		rootCmd.SetArgs([]string{"remove", domainName})
// 		rootCmd.Execute()
// 	}
// 	rootCmd.SetArgs([]string{"list"})
// 	rootCmd.Execute()

// case "Go Back":
// 	rootCmd.SetArgs([]string{"list"})
// 	rootCmd.Execute()

// case "Exit":
// 	fmt.Println("Exiting...")
// 	os.Exit(0)
// }
// }
