package cmd

import (
	"fmt"
	"github.com/anton-johansson/k8s-login/version"
	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Prints the current version of this tool",
	Run: func(command *cobra.Command, args []string) {
		fmt.Println(version.Version())
	},
}

func init() {
	rootCommand.AddCommand(versionCommand)
}
