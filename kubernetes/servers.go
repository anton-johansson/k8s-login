package kubernetes

// Server indicates a specific configured Kubernetes context.
type Server struct {
	Name    string
	Address string
}

// GetServers gets a list of all configured contexts.
func GetServers(kubeconfig *KubeConfig) []Server {
	var servers []Server
	for key := range kubeconfig.config.Clusters {
		item := Server{
			Name:    key,
			Address: kubeconfig.config.Clusters[key].Server,
		}
		servers = append(servers, item)
	}
	return servers
}
