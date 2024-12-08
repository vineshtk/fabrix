package cmd

import (
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
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

		// stdout, err := command.StdoutPipe()
		// if err != nil {
		// 	fmt.Printf("Error creating stdout pipe: %v\n", err)
		// 	return
		// }

		// stderr, err := command.StderrPipe()
		// if err != nil {
		// 	fmt.Printf("Error creating stderr pipe: %v\n", err)
		// 	return
		// }

		p := tea.NewProgram(newModel())

		// go func() {
		// 	scanner := bufio.NewScanner(stdout)
		// 	for scanner.Scan() {
		// 		fmt.Printf("[INFO]: %s\n", scanner.Text())
		// 	}
		// }()

		// go func() {
		// 	scanner := bufio.NewScanner(stderr)
		// 	for scanner.Scan() {
		// 		fmt.Printf("[ERROR]: %s\n", scanner.Text())
		// 	}
		// }()

		go func() {
			err = command.Start()
			if err != nil {
				fmt.Printf("Error starting script: %v\n", err)
				return
			}

			// Wait for the script to finish
			err = command.Wait()
			if err != nil {
				fmt.Printf("Error executing script: %v\n", err)
			}

			// Stop the spinner
			p.Send(stopMsg{})
		}()

		if _, err := p.Run(); err != nil {
			fmt.Println("could not run program:", err)
			os.Exit(1)
		}

		fmt.Println("Network started successfully!")
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
