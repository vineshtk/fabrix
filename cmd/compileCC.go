package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/vineshtk/fabrix/pkg/configs"
)

// networkCmd represents the network command
var compileCmd = &cobra.Command{
	Use:     "compile",
	Aliases: []string{"dcc"},
	Short:   "Use this command to deploy chaincode.",
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

		ccPath, err := cmd.Flags().GetString("path")
		if err != nil {
			fmt.Println("Error retrieving domain flag:", err)
			os.Exit(1)
		}

		ccLang, err := cmd.Flags().GetString("lang")
		if err != nil {
			fmt.Println("Error retrieving domain flag:", err)
			os.Exit(1)
		}
		configs.CompileChaincode(ccPath, ccLang)
		// fmt.Println("Selected Domain:", selectedDomain)
	},
}

func init() {

	rootCmd.AddCommand(compileCmd)
	// networkCmd.Flags().BoolVarP(&option, "option", "o", false, "Modify option")
	compileCmd.Flags().StringVarP(&ccPath, "path", "p", "", "specify the chaincode path")
	// compileCmd.Flags().StringVarP(&ccversion, "version", "v", "", "specify the version of chaincode")
	compileCmd.Flags().StringVarP(&ccLang, "lang", "l", "", "specify the chaincode language")

	// // Mark the flag as required
	compileCmd.MarkFlagRequired("path")
	// compileCmd.MarkFlagRequired("version")
	compileCmd.MarkFlagRequired("lang")

}
