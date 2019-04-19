package cmd

import (
	"fmt"
	"github.com/anton-johansson/k8s-login/version"
	"github.com/spf13/cobra"
)

var short bool

func init() {
	var command = &cobra.Command{
		Use:   "version",
		Short: "Prints the current version of this tool",
		Run: func(command *cobra.Command, args []string) {
			info := version.Info()
			if short {
				fmt.Println(info.Version)
			} else {
				fmt.Println(info.Version + " (go version: " + info.GoVersion + ", commit: " + info.Commit + ")")
			}
		},
	}

	command.Flags().BoolVarP(&short, "short", "s", false, "Whether or not to output the actual version only")
	rootCommand.AddCommand(command)
}
