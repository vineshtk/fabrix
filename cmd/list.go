package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// networkCmd represents the network command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "Use this command to list all the networks available.",
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

		// List directories
		for i, entry := range entries {
			if entry.IsDir() {
				fmt.Println(i, ".", entry.Name())
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
