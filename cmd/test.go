package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vineshtk/fabrix/pkg/configs"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Use this command to down a network and remove all the files associated",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		configs.InstallChaincode("auto.com", "../../../pkg/configs/defaults/Chaincode", "golang", "kbaauto_1.0", "sample-chaincode", "1.0.2", "2")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
