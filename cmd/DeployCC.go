package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var ccpath string
var ccversion string
var cclang string

// networkCmd represents the network command
var deployCmd = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"dcc"},
	Short:   "Use this command to deploy chaincode.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {

		file := ""
		prompt := &survey.Input{
			Message: "inform a file to save:",
			Suggest: func(toComplete string) []string {
				files, _ := filepath.Glob(toComplete + "*")
				return files
			},
		}

		survey.AskOne(prompt, &file)

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
		prompt1 := &survey.Select{
			Message: "Choose a domain:",
			Options: domains,
		}

		err = survey.AskOne(prompt1, &selectedDomain)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		name := false
		prompt2 := &survey.Confirm{
			Message: "Do you like pie?",
		}
		survey.AskOne(prompt2, &name)

		// Retrieve the value of the flag
		ccPath, err := cmd.Flags().GetString("path")
		if err != nil {
			fmt.Println("Error retrieving domain flag:", err)
			os.Exit(1)
		}

		ccVersion, err := cmd.Flags().GetString("version")
		if err != nil {
			fmt.Println("Error retrieving domain flag:", err)
			os.Exit(1)
		}

		ccLang, err := cmd.Flags().GetString("lang")
		if err != nil {
			fmt.Println("Error retrieving domain flag:", err)
			os.Exit(1)
		}
		// add the logic here
		fmt.Println(ccPath, ccVersion, ccLang, selectedDomain)
		// configs.InstallChaincode(selectedDomain, ccPath, ccLang)
	},
}

func init() {

	rootCmd.AddCommand(deployCmd)
	// networkCmd.Flags().BoolVarP(&option, "option", "o", false, "Modify option")
	// deployCmd.Flags().StringVarP(&ccpath, "path", "p", "", "specify the chaincode path")
	// deployCmd.Flags().StringVarP(&ccversion, "version", "v", "", "specify the version of chaincode")
	// deployCmd.Flags().StringVarP(&cclang, "lang", "l", "", "specify the chaincode language")

	// // Mark the flag as required
	// deployCmd.MarkFlagRequired("path")
	// deployCmd.MarkFlagRequired("version")
	// deployCmd.MarkFlagRequired("lang")

}
