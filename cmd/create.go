package cmd

import (
	"fmt"
	"os"

	"github.com/vineshtk/fabrix/pkg/menu"

	"github.com/spf13/cobra"
)

// networkCmd represents the network command
var networkCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Use this command to build a new network.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve the value of the flag
		channelName, err := cmd.Flags().GetString("channel")
		if err != nil {
			fmt.Println("Error retrieving domain flag:", err)
			os.Exit(1)
		}

		version, err := cmd.Flags().GetString("version")
		if err != nil {
			fmt.Println("Error retrieving domain flag:", err)
			os.Exit(1)
		}
		
		menu.ShowMainMenu()
		menu.GetInputsFromUser(channelName, version)
	},
}

func init() {

	rootCmd.AddCommand(networkCmd)
	// networkCmd.Flags().BoolVarP(&option, "option", "o", false, "Modify option")
	networkCmd.Flags().String("channel", "", "for custom channel(eg.mychannel)")
	networkCmd.Flags().String("version", "2.5.4", "for specific fabric version(eg.2.5.4)")
}
