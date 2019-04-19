package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCommand = &cobra.Command{
	Use:   "k8s-login",
	Short: "Handles logging in to a Kubernetes cluster",
}

// Execute executes the CLI
func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
