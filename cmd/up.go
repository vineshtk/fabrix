package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// networkCmd represents the network command
var upCmd = &cobra.Command{
	Use: "up",
	// Aliases: []string{"sn"},
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

		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		err = command.Run()
		if err != nil {
			fmt.Printf("Error executing script: %v\n", err)
			return
		}

		// p := tea.NewProgram(newModel("Please wait! setting up the network for you!"))

		// go func() {
		// 	err = command.Start()
		// 	if err != nil {
		// 		fmt.Printf("Error starting script: %v\n", err)
		// 		return
		// 	}

		// 	// Wait for the script to finish
		// 	err = command.Wait()
		// 	if err != nil {
		// 		fmt.Printf("Error executing script: %v\n", err)
		// 	}

		// 	// Stop the spinner
		// 	p.Send(stopMsg{})
		// }()

		// if _, err := p.Run(); err != nil {
		// 	fmt.Println("could not run program:", err)
		// 	os.Exit(1)
		// }

		fmt.Println("Network started successfully!")
		rootCmd.SetArgs([]string{"dp"})
		rootCmd.Execute()
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	// networkCmd.PersistentFlags().String("foo", "", "A help for foo")
	// networkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
