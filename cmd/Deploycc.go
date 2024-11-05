package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vineshtk/fabrix/pkg/configs"
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
		fmt.Println(ccPath, ccVersion)
		configs.InstallChaincode(ccPath, ccLang)

	},
}

func init() {

	rootCmd.AddCommand(deployCmd)
	// networkCmd.Flags().BoolVarP(&option, "option", "o", false, "Modify option")
	deployCmd.Flags().StringVarP(&ccpath, "path", "p", "", "specify the chaincode path")
	deployCmd.Flags().StringVarP(&ccversion, "version", "v", "", "specify the version of chaincode")
	deployCmd.Flags().StringVarP(&cclang, "lang", "l", "", "specify the chaincode language")

	// Mark the flag as required
	deployCmd.MarkFlagRequired("path")
	deployCmd.MarkFlagRequired("version")
	deployCmd.MarkFlagRequired("lang")

}
