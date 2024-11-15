package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
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
	},
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
	startPrompt()
}

func startPrompt() {
	options := []string{"Create new domain", "Choose existing domain", "Exit"}

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
}
