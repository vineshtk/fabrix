package cmd

import (
	"github.com/spf13/cobra"
)

// networkCmd represents the network command
var networkCmd = &cobra.Command{
	Use:     "Network",
	Aliases: []string{"N"},
	Short:   "Use this command to build a new network.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(networkCmd)
	// networkCmd.PersistentFlags().String("foo", "", "A help for foo")
	// networkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
