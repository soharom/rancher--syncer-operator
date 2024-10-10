package internal

import (
	"net/http"
)

const (
	ClustersEndpoint = "/clusters"
	//RancherAPI          = "https://127.0.0.1:53954/v3"
	//Token               = "token-dc5bj:wcl58sh4259b955n9fghmbnfrlzcdqxmcw7vq5pb7g5x22cvdqwbvj"
	ResourceTypeCluster = "cluster"
)

type ClusterData struct {
	Id             string         `json:"id"`
	Name           string         `json:"Name"`
	ClusterActions ClusterActions `json:"actions"`
}

type ClusterActions struct {
	GenerateKubeconfigEndpoint string `json:"generateKubeconfig"`
}

type Clusters struct {
	ClusterDatas []ClusterData `json:"data"`
	ResourceType string        `json:"resourceTypes`
}

type GeneratedKubeconfig struct {
	Config string `json:"config"`
}

type Client struct {
	Token        string
	hc           http.Client
	SkipTlsCheck bool
}

type ClusterSecretData struct {
	Name       string
	Kubeconfig string
}
