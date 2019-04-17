package cmd

import (
	"fmt"
	"github.com/anton-johansson/k8s-login/kubernetes"
	"github.com/spf13/cobra"
)

var verbose bool
var kubeconfigFileName string

var serversCommand = &cobra.Command{
	Use:   "servers",
	Short: "Prints the available servers in your kubeconfig",
	Run: func(command *cobra.Command, args []string) {
		kubeconfig, error := kubernetes.GetKubeConfig(kubeconfigFileName)
		if error != nil {
			fmt.Println(error)
			return
		}

		servers := kubernetes.GetServers(kubeconfig)
		for _, server := range servers {
			if verbose {
				fmt.Println(server.Name + " (" + server.Address + ")")
			} else {
				fmt.Print(server.Name + " ")
			}
		}
		if !verbose {
			fmt.Println()
		}
	},
}

func init() {
	serversCommand.Flags().BoolVarP(&verbose, "verbose", "v", false, "Whether or not to use verbose output")
	serversCommand.Flags().StringVar(&kubeconfigFileName, "kubeconfig", kubernetes.GetDefaultKubeConfigFileName(), "The path to the kubeconfig")
	rootCommand.AddCommand(serversCommand)
}
