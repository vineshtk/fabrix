package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github/vineshtk/fabrix/pkg/menu"
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
		menu.ShowMainMenu()
		var userInput string
		fmt.Scan(&userInput)
		if userInput == "N"{
			menu.GetInputsFromUser()
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
