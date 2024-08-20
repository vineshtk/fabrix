package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// networkCmd represents the network command
var downCmd = &cobra.Command{
	Use:     "down",
	Aliases: []string{"dn"},
	Args:    cobra.ExactArgs(1),
	Short:   "Use this command to down a network and remove all the files associated",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {

		scriptPath := fmt.Sprintf("./fabrix/%v/Network/stopNetwork.sh", args[0])
        scriptDir := fmt.Sprintf("./fabrix/%v/Network/", args[0])
		
		err := os.Chmod(scriptPath, 0755)
		if err != nil {
			fmt.Printf("Error making script executable: %v\n", err)
			return
		}

		command := exec.Command("/bin/bash", "stopNetwork.sh")
        command.Dir = scriptDir
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		err = command.Run()
		if err != nil {
			fmt.Printf("Error executing script: %v\n", err)
			return
		}

		fmt.Println("Script executed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
	// networkCmd.PersistentFlags().String("foo", "", "A help for foo")
	// networkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
