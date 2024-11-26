package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/vineshtk/fabrix/pkg/prompts"
)

var version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "fabrix",
	Version: version,
	Short:   "fabrix is a tool to create a fabric network",
	Long: `fabrix is a tool,
	that helps chaincode developers to setup a fabric network easily,
	for deploying and testing the chaincode.`,

	Run: func(cmd *cobra.Command, args []string) {
		prompts.ShowMainMenu()
		startPrompt()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model1 struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model1) Init() tea.Cmd {
	return nil
}

func (m model1) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

			case "Create new domain":
				rootCmd.SetArgs([]string{"create"})
				rootCmd.Execute()
		
			case "Choose existing domain":
				rootCmd.SetArgs([]string{"list"})
				rootCmd.Execute()
		
			case "Go back":
				rootCmd.SetArgs([]string{"fabrix"})
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

func (m model1) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Ok let's continue", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Exit? ... See you later")
	}
	return "\n" + m.list.View()
}

func startPrompt() {
	items := []list.Item{
		item("Create new domain"),
		item("Choose existing domain"),
		item("Go back"),
		item("Exit"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Choose an option:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model1{list: l}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// func startPrompt() {
// 	options := []string{"Create new domain", "Choose existing domain", "Exit"}

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

	// case "Create new domain":
	// 	rootCmd.SetArgs([]string{"create"})
	// 	rootCmd.Execute()

	// case "Choose existing domain":
	// 	rootCmd.SetArgs([]string{"list"})
	// 	rootCmd.Execute()

	// case "Go back":
	// 	rootCmd.SetArgs([]string{"fabrix"})
	// 	rootCmd.Execute()

	// case "Exit":
	// 	fmt.Println("Exiting...")
	// 	os.Exit(0)
	// }
// }
