package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// networkCmd represents the network command
var runCmd = &cobra.Command{
	Use: "up",
	// Aliases: []string{"up"},
	Args:  cobra.ExactArgs(1),
	Short: "Use this command to start a network and install chaincode",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {

		scriptPath := fmt.Sprintf("./fabrix/%v/Network/startNetwork.sh", args[0])
		scriptDir := fmt.Sprintf("./fabrix/%v/Network/", args[0])

		err := os.Chmod(scriptPath, 0755)
		if err != nil {
			fmt.Printf("Error making script executable: %v\n", err)
			return
		}

		command := exec.Command("/bin/bash", "startNetwork.sh")
		command.Dir = scriptDir
		command.Stdout = io.Discard
		command.Stderr = io.Discard
		// command.Stdout = os.Stdout
		// command.Stderr = os.Stderr

		err = command.Run()
		if err != nil {
			fmt.Printf("Error executing script: %v\n", err)
			return
		}

		fmt.Println("Network started successfully!")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	// networkCmd.PersistentFlags().String("foo", "", "A help for foo")
	// networkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
