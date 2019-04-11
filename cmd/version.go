package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Prints the current version of this tool",
	Run: func(command *cobra.Command, args []string) {
		fmt.Println("1.1.1")
	},
}

func init() {
	RootCommand.AddCommand(versionCommand)
}
