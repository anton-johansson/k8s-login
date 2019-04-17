package main

import (
	"github.com/anton-johansson/k8s-login/cmd"
)

func main() {
	cmd.Execute()

	/*
		k8s_client "k8s.io/client-go/tools/clientcmd"
		rules := k8s_client.NewDefaultClientConfigLoadingRules()
		config, err := rules.Load()
		if err != nil {
			fmt.Println("Could not load kubeconfig")
			return
		}

		if len(config.Clusters) == 0 {
			fmt.Println("No servers found")
			return
		}

		for key := range config.Clusters {
			cluster := config.Clusters[key]
			fmt.Println(key + ": (" + cluster.Server + ")")
		}
	*/
}
