package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// networkCmd represents the network command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rn"},
	Args:    cobra.ExactArgs(1),
	Short:   "Use this command to remove all files for a network",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Define the folder to be removed
		folderToRemove := fmt.Sprintf("./fabrix/%v", args[0])

		// Check if the folder exists
		if _, err := os.Stat(folderToRemove); os.IsNotExist(err) {
			fmt.Printf("Folder does not exist: %s\n", folderToRemove)
			return
		}

		// Remove the folder and its contents
		err := os.RemoveAll(folderToRemove)
		if err != nil {
			fmt.Printf("Error removing folder: %v\n", err)
			return
		}

		fmt.Println("Network configurations removed successfully!")

		rootCmd.SetArgs([]string{"list"})
		rootCmd.Execute()
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	// networkCmd.PersistentFlags().String("foo", "", "A help for foo")
	// networkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
