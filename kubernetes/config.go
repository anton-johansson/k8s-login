package kubernetes

import (
	"os"
	client "k8s.io/client-go/tools/clientcmd"
	api "k8s.io/client-go/tools/clientcmd/api"
)

// KubeConfig holds a kubeconfig and its file name
type KubeConfig struct {
	// FileName is the file name of the kubeconfig file
	FileName string
	config *api.Config
}

// GetDefaultKubeConfigFileName gets the default file name of the kubeconfig
func GetDefaultKubeConfigFileName() string {
	return client.RecommendedHomeFile
}

// GetKubeConfig loads a kubeconfig from the given file or a default file
func GetKubeConfig(argument string) (*KubeConfig, error) {
	var fileName string
	var config *api.Config
	var error error

	rules := client.NewDefaultClientConfigLoadingRules()
	if argument != "" {
		fileName = argument
		if _, error = os.Stat(argument); os.IsNotExist(error) {
			error = nil
			config = api.NewConfig()
		} else if error == nil {
			rules.ExplicitPath = argument
			config, error = rules.Load()
		}
	} else {
		fileName = client.RecommendedHomeFile
		config, error = rules.Load()
	}

	return &KubeConfig{
		FileName: fileName,
		config: config,
	}, error
}

// UpdateKubeConfig updates the local kubeconfig (~/.kube/config) with the given user.
func UpdateKubeConfig(kubeconfig *KubeConfig, user UserData, updateContext bool, serverName string) {
	authInfo := api.NewAuthInfo()
	if current, ok := kubeconfig.config.AuthInfos[user.Name]; ok {
		authInfo = current
	}

	authInfo.AuthProvider = &api.AuthProviderConfig{
		Name: "oidc",
		Config: map[string]string{
			"client-id": user.ClientID,
			"client-secret": user.ClientSecret,
			"id-token": user.IDToken,
			"refresh-token": user.RefreshToken,
			"idp-issuer-url": user.IssuerURL,
		},
	}
	kubeconfig.config.AuthInfos[user.Name] = authInfo

	if updateContext {
		context := api.NewContext()
		if current, ok := kubeconfig.config.Contexts[serverName]; ok {
			context = current
		}
		context.Cluster = serverName
		context.AuthInfo = user.Name
		kubeconfig.config.Contexts[serverName] = context
		kubeconfig.config.CurrentContext = serverName
	}

	client.WriteToFile(*kubeconfig.config, kubeconfig.FileName)
}
