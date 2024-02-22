package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "fabrix",
	Version: version,
	Short:   "fabrix is a tool to create a fabric network",
	Long: `fabrix is a tool,
	that helps chaincode developers to setup a fabric network easily,
	for deploying and testing the chaincode.`,
	Run: func(cmd *cobra.Command, args []string) { 
		fmt.Println("Welcome to Fabrix - The helper tool for chaincode developers to create fabric network, it takes away all the heavy lifting for you!!!")
		fmt.Println("You will guided during all Hyperledger Fabric deployment. Let's start...")
		fmt.Println("Plase choose from the menu")
		fmt.Println("MENU")
		fmt.Println("N - New network")
		fmt.Println("S - Select an existing network")
		fmt.Println("D - Docker status")
		fmt.Println("C - Clean all Docker resources")
		fmt.Println("Q - Quit")
		
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
