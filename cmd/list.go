package cmd

import (
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
		rootCmd.SetArgs([]string{"lp"})
		rootCmd.Execute()
		rootCmd.SetArgs([]string{"dp"})
		rootCmd.Execute()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
