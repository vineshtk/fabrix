package cmd

// import (
// 	"fmt"
// 	"math"
// 	"strconv"
// 	"strings"
// 	"time"

// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// 	"github.com/fogleman/ease"
// 	"github.com/lucasb-eyer/go-colorful"
// 	"github.com/spf13/cobra"
// )

// // import (
// // 	"fmt"
// // 	"io"
// // 	"os"
// // 	"strings"

// // 	"github.com/charmbracelet/bubbles/list"
// // 	tea "github.com/charmbracelet/bubbletea"
// // 	"github.com/charmbracelet/lipgloss"
// // 	"github.com/spf13/cobra"
// // )

// // const listHeight = 14

// // var (
// // 	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
// // 	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
// // 	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
// // 	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
// // 	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
// // 	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
// // )

// // type item string

// // func (i item) FilterValue() string { return "" }

// // type itemDelegate struct{}

// // networkCmd represents the network command
// var testCmd = &cobra.Command{
// 	Use:     "test",
// 	Aliases: []string{"tst"},
// 	Short:   "Use this command to remove all files for a network",
// 	Long: `A longer description that spans multiple lines and likely contains examples
// and usage of using your command.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		initialModel := model{0, false, 10, 0, 0, false, false}
// 		p := tea.NewProgram(initialModel)
// 		if _, err := p.Run(); err != nil {
// 			fmt.Println("could not start program:", err)
// 		}
// 	},
// }

// func init() {
// 	rootCmd.AddCommand(testCmd)
// 	// networkCmd.PersistentFlags().String("foo", "", "A help for foo")
// 	// networkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
// }

// // func (d itemDelegate) Height() int                             { return 1 }
// // func (d itemDelegate) Spacing() int                            { return 0 }
// // func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
// // func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
// // 	i, ok := listItem.(item)
// // 	if !ok {
// // 		return
// // 	}

// // 	str := fmt.Sprintf("%d. %s", index+1, i)

// // 	fn := itemStyle.Render
// // 	if index == m.Index() {
// // 		fn = func(s ...string) string {
// // 			return selectedItemStyle.Render("> " + strings.Join(s, " "))
// // 		}
// // 	}

// // 	fmt.Fprint(w, fn(str))
// // }

// // type model struct {
// // 	list     list.Model
// // 	choice   string
// // 	quitting bool
// // }

// // func (m model) Init() tea.Cmd {
// // 	return nil
// // }

// // func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// // 	switch msg := msg.(type) {
// // 	case tea.WindowSizeMsg:
// // 		m.list.SetWidth(msg.Width)
// // 		return m, nil

// // 	case tea.KeyMsg:
// // 		switch keypress := msg.String(); keypress {
// // 		case "q", "ctrl+c":
// // 			m.quitting = true
// // 			return m, tea.Quit

// // 		case "enter":
// // 			i, ok := m.list.SelectedItem().(item)
// // 			if ok {
// // 				m.choice = string(i)
// // 			}
// // 			return m, tea.Quit
// // 		}
// // 	}

// // 	var cmd tea.Cmd
// // 	m.list, cmd = m.list.Update(msg)
// // 	return m, cmd
// // }

// // func (m model) View() string {
// // 	if m.choice != "" {
// // 		return quitTextStyle.Render(fmt.Sprintf("%s? Sounds good to me.", m.choice))
// // 	}
// // 	if m.quitting {
// // 		return quitTextStyle.Render("Not hungry? That’s cool.")
// // 	}
// // 	return "\n" + m.list.View()
// // }

// // func main() {
// // 	items := []list.Item{
// // 		item("Ramen"),
// // 		item("Tomato Soup"),
// // 		item("Hamburgers"),
// // 		item("Cheeseburgers"),
// // 		item("Currywurst"),
// // 		item("Okonomiyaki"),
// // 		item("Pasta"),
// // 		item("Fillet Mignon"),
// // 		item("Caviar"),
// // 		item("Just Wine"),
// // 	}

// // 	const defaultWidth = 20

// // 	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
// // 	l.Title = "What do you want for dinner?"
// // 	l.SetShowStatusBar(false)
// // 	l.SetFilteringEnabled(false)
// // 	l.Styles.Title = titleStyle
// // 	l.Styles.PaginationStyle = paginationStyle
// // 	l.Styles.HelpStyle = helpStyle

// // 	m := model{list: l}

// // 	if _, err := tea.NewProgram(m).Run(); err != nil {
// // 		fmt.Println("Error running program:", err)
// // 		os.Exit(1)
// // 	}
// // }

// // An example demonstrating an application with multiple views.
// //
// // Note that this example was produced before the Bubbles progress component
// // was available (github.com/charmbracelet/bubbles/progress) and thus, we're
// // implementing a progress bar from scratch here.

// const (
// 	progressBarWidth  = 71
// 	progressFullChar  = "█"
// 	progressEmptyChar = "░"
// 	dotChar           = " • "
// )

// // General stuff for styling the view
// var (
// 	keywordStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
// 	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
// 	ticksStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
// 	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
// 	progressEmpty = subtleStyle.Render(progressEmptyChar)
// 	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
// 	mainStyle     = lipgloss.NewStyle().MarginLeft(2)

// 	// Gradient colors we'll use for the progress bar
// 	ramp = makeRampStyles("#B14FFF", "#00FFA3", progressBarWidth)
// )

// type (
// 	tickMsg  struct{}
// 	frameMsg struct{}
// )

// func tick() tea.Cmd {
// 	return tea.Tick(time.Second, func(time.Time) tea.Msg {
// 		return tickMsg{}
// 	})
// }

// func frame() tea.Cmd {
// 	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
// 		return frameMsg{}
// 	})
// }

// type model struct {
// 	Choice   int
// 	Chosen   bool
// 	Ticks    int
// 	Frames   int
// 	Progress float64
// 	Loaded   bool
// 	Quitting bool
// }

// func (m model) Init() tea.Cmd {
// 	return tick()
// }

// // Main update function.
// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	// Make sure these keys always quit
// 	if msg, ok := msg.(tea.KeyMsg); ok {
// 		k := msg.String()
// 		if k == "q" || k == "esc" || k == "ctrl+c" {
// 			m.Quitting = true
// 			return m, tea.Quit
// 		}
// 	}

// 	// Hand off the message and model to the appropriate update function for the
// 	// appropriate view based on the current state.
// 	if !m.Chosen {
// 		return updateChoices(msg, m)
// 	}
// 	return updateChosen(msg, m)
// }

// // The main view, which just calls the appropriate sub-view
// func (m model) View() string {
// 	var s string
// 	if m.Quitting {
// 		return "\n  See you later!\n\n"
// 	}
// 	if !m.Chosen {
// 		s = choicesView(m)
// 	} else {
// 		s = chosenView(m)
// 	}
// 	return mainStyle.Render("\n" + s + "\n\n")
// }



// func checkbox(label string, checked bool) string {
// 	if checked {
// 		return checkboxStyle.Render("[x] " + label)
// 	}
// 	return fmt.Sprintf("[ ] %s", label)
// }

// func progressbar(percent float64) string {
// 	w := float64(progressBarWidth)

// 	fullSize := int(math.Round(w * percent))
// 	var fullCells string
// 	for i := 0; i < fullSize; i++ {
// 		fullCells += ramp[i].Render(progressFullChar)
// 	}

// 	emptySize := int(w) - fullSize
// 	emptyCells := strings.Repeat(progressEmpty, emptySize)

// 	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
// }

// // Utils

// // Generate a blend of colors.
// func makeRampStyles(colorA, colorB string, steps float64) (s []lipgloss.Style) {
// 	cA, _ := colorful.Hex(colorA)
// 	cB, _ := colorful.Hex(colorB)

// 	for i := 0.0; i < steps; i++ {
// 		c := cA.BlendLuv(cB, i/steps)
// 		s = append(s, lipgloss.NewStyle().Foreground(lipgloss.Color(colorToHex(c))))
// 	}
// 	return
// }

// // Convert a colorful.Color to a hexadecimal format.
// func colorToHex(c colorful.Color) string {
// 	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
// }

// // Helper function for converting colors to hex. Assumes a value between 0 and
// // 1.
// func colorFloatToHex(f float64) (s string) {
// 	s = strconv.FormatInt(int64(f*255), 16)
// 	if len(s) == 1 {
// 		s = "0" + s
// 	}
// 	return
// }


// package main

// // A simple program demonstrating the spinner component from the Bubbles
// // component library.

// import (
// 	"fmt"
// 	"os"

// 	"github.com/charmbracelet/bubbles/spinner"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// )

// type errMsg error

// type model struct {
// 	spinner  spinner.Model
// 	quitting bool
// 	err      error
// }

// func initialModel() model {
// 	s := spinner.New()
// 	s.Spinner = spinner.Dot
// 	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
// 	return model{spinner: s}
// }

// func (m model) Init() tea.Cmd {
// 	return m.spinner.Tick
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "q", "esc", "ctrl+c":
// 			m.quitting = true
// 			return m, tea.Quit
// 		default:
// 			return m, nil
// 		}

// 	case errMsg:
// 		m.err = msg
// 		return m, nil

// 	default:
// 		var cmd tea.Cmd
// 		m.spinner, cmd = m.spinner.Update(msg)
// 		return m, cmd
// 	}
// }

// func (m model) View() string {
// 	if m.err != nil {
// 		return m.err.Error()
// 	}
// 	str := fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n", m.spinner.View())
// 	if m.quitting {
// 		return str + "\n"
// 	}
// 	return str
// }

// func main() {
// 	p := tea.NewProgram(initialModel())
// 	if _, err := p.Run(); err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// }