package kubernetes

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

type clusterItem struct {
	Name    string  `yaml:"name"`
	Cluster cluster `yaml:"cluster"`
}

type cluster struct {
	ClusterAuthorityData string `yaml:"certificate-authority-data"`
	Server               string `yaml:"server"`
}

type userItem struct {
	Name string `yaml:"name"`
	User user   `yaml:"user"`
}

type user struct {
	ClientCertificteData string `yaml:"client-certificate-data"`
	Token                string `yaml:"token"`
	ClientKeyData        string `yaml:"client-key-data"`
}

type contextItem struct {
	Name    string  `yaml:"name"`
	Context context `yaml:"context"`
}

type context struct {
	Cluster   string `yaml:"cluster"`
	User      string `yaml:"user"`
	Namespace string `yaml:"namespace,omitempty"`
}

// KubeConfig stores information that is found in a kube config file
type KubeConfig struct {
	APIVersion     string                 `yaml:"apiVersion"`
	Clusters       []clusterItem          `yaml:"clusters"`
	Users          []userItem             `yaml:"users"`
	Contexts       []contextItem          `yaml:"contexts,omitempty"`
	CurrentContext string                 `yaml:"current-context,omitempty"`
	Kind           string                 `yaml:"kind,omitempty"`
	Preferences    map[string]interface{} `yaml:"preferences,omitempty"`
}

// ParseKubeConfig returns a KubeConfig struct built from input string kube config
func ParseKubeConfig(config string) (*KubeConfig, error) {
	if config == "" {
		return nil, fmt.Errorf("Cannot parse empty config")
	}

	var kubeConfig KubeConfig
	err := yaml.Unmarshal([]byte(config), &kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal YAML config with error %+v", err)
	}
	if len(kubeConfig.Clusters) <= 0 || len(kubeConfig.Users) <= 0 {
		return nil, fmt.Errorf("Config %+v contains no valid clusters or users", kubeConfig)
	}
	u := kubeConfig.Users[0].User
	if u.Token == "" && (u.ClientCertificteData == "" || u.ClientKeyData == "") {
		return nil, fmt.Errorf("Config requires either token or certificate auth for user %+v", u)
	}
	c := kubeConfig.Clusters[0].Cluster
	if c.Server == "" {
		return nil, fmt.Errorf("Config has invalid or non existent server for cluster %+v", c)
	}

	return &kubeConfig, nil
}
