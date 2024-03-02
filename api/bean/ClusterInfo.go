package bean

import "github.com/devtron-labs/devtron/pkg/serverConnection/bean"

type ClusterInfo struct {
	ClusterId              int                              `json:"clusterId"`
	ClusterName            string                           `json:"clusterName"`
	BearerToken            string                           `json:"bearerToken"`
	ServerUrl              string                           `json:"serverUrl"`
	InsecureSkipTLSVerify  bool                             `json:"insecureSkipTLSVerify"`
	KeyData                string                           `json:"keyData"`
	CertData               string                           `json:"certData"`
	CAData                 string                           `json:"CAData"`
	ServerConnectionConfig *bean.ServerConnectionConfigBean `json:"serverConnectionConfig"`
}
