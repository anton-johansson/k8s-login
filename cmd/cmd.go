package cmd

import (
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:   "k8s-login",
	Short: "Handles logging in to a Kubernetes cluster",
}
