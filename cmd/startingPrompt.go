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
)

// networkCmd represents the network command
var startingPromptCmd = &cobra.Command{
	Use:     "startingPrompt",
	Aliases: []string{"sp"},
	Short:   "Use this command to remove all files for a network",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		startPrompt()
	},
}

func init() {
	rootCmd.AddCommand(startingPromptCmd)
}

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

func (m model1) View() string {
	if m.choice != "" {
		return quitTextStyle.Render(fmt.Sprintf("%s? Ok let's continue", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Exit? ... See you later")
	}
	return "\n" + m.list.View()
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
			// Save the choice and quit the Bubble Tea program
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

func startPrompt() {
	items := []list.Item{
		item("Create new domain"),
		item("Choose existing domain"),
		item("Exit"),
	}

	const defaultWidth = 20
	const listHeight = 10

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Choose an option:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle

	m := model1{list: l}

	// Run the Bubble Tea program
	prog, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if m, ok := prog.(model1); ok {
		switch m.choice {
		case "Create new domain":
			rootCmd.SetArgs([]string{"create"})
			rootCmd.Execute()

		case "Choose existing domain":
			rootCmd.SetArgs([]string{"list"})
			rootCmd.Execute()

		case "Exit":
			fmt.Println("Exiting...")
			os.Exit(0)
		}
	}
}
