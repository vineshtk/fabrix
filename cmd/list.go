package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
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
		rootDir := "./fabrix"
		// Read the directory
		entries, err := os.ReadDir(rootDir)
		if err != nil {
			fmt.Printf("Error reading directory: %v\n", err)
			return
		}
		var domains []string
		// List directories
		for _, entry := range entries {
			if entry.IsDir() {
				domains = append(domains, entry.Name())
			}
		}

		// Variable to store the user's selection
		var selectedDomain string

		// Prompt user to select a domain
		prompt := &survey.Select{
			Message: "Choose a domain:",
			Options: domains,
		}

		err = survey.AskOne(prompt, &selectedDomain)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		domainOptions(selectedDomain)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func domainOptions(domainName string) {
	options := []string{"Start network", "Info", "Deploy chaincode", "Down network", "Remove domain", "Go Back", "Exit"}

	var selectedOption string
	prompt := &survey.Select{
		Message: "Choose an option:",
		Options: options,
	}

	if err := survey.AskOne(prompt, &selectedOption); err != nil {
		fmt.Println("Error:", err)
		return
	}

	switch selectedOption {

	case "Start network":
		rootCmd.SetArgs([]string{"up", domainName})
		rootCmd.Execute()

	case "Info":
		rootCmd.SetArgs([]string{"list"})
		rootCmd.Execute()

	case "Deploy chaincode":
		rootCmd.SetArgs([]string{"deploy"})
		rootCmd.Execute()

	case "Down network":
		rootCmd.SetArgs([]string{"down", domainName})
		rootCmd.Execute()

	case "Remove domain":
		confirm := false
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Do you want to remove '%s'?", domainName),
		}
		survey.AskOne(prompt, &confirm)
		if confirm {
			rootCmd.SetArgs([]string{"remove", domainName})
			rootCmd.Execute()
		}
		rootCmd.SetArgs([]string{"list"})
		rootCmd.Execute()

	case "Go Back":
		rootCmd.SetArgs([]string{"list"})
		rootCmd.Execute()

	case "Exit":
		fmt.Println("Exiting...")
		os.Exit(0)
	}
}
