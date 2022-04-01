package ccompare

type EnvClustersInfo struct {
	Env      string   `toml:"env"`
	Clusters []string `toml:"clusters"`
}
type AppInfo struct {
	Name                       string `toml:"name"`
	AppId                      string `toml:"appId"`
	OrgId                      string `toml:"orgId"`
	OrgName                    string `toml:"orgName"`
	OwnerName                  string `toml:"ownerName"`
	OwnerEmail                 string `toml:"ownerEmail"`
	DataChangeCreatedBy        string `toml:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `toml:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `toml:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `toml:"dataChangeLastModifiedTime"`
}
type NamespaceInfo struct {
	AppId                      string     `toml:"appId"`
	ClusterName                string     `toml:"clusterName"`
	NamespaceName              string     `toml:"namespaceName"`
	Comment                    string     `toml:"comment"`
	Format                     string     `toml:"format"`
	IsPublic                   bool       `toml:"isPublic"`
	DataChangeCreatedBy        string     `toml:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string     `toml:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string     `toml:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string     `toml:"dataChangeLastModifiedTime"`
	Items                      []ItemInfo `toml:"items"`
}
type ItemInfo struct {
	Key                        string `toml:"key"`
	Value                      string `toml:"value"`
	DataChangeCreatedBy        string `toml:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `toml:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `toml:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `toml:"dataChangeLastModifiedTime"`
}

type AppCfg struct {
	DingKeys       []string              `toml:"ding_keys"`
	DingUsers      []string              `toml:"ding_users"`
	DingUserMap    map[string]string     `toml:"ding_user_map"`
	IsAtAll        int                   `toml:"is_at_all"`
	EnUpdateConsul int                   `toml:"enable_update_consul"`
	EnDelConsul    int                   `toml:"enable_delete_consul"`
	ChklogRate     float64               `toml:"log_rate"`
	AppConfigMap   map[string]ConfigInfo `toml:"app_config_map"`
}

type AppClusterInfo struct {
	Cluster   []string `toml:"cluster"`
	Namespace []string `toml:"namespace"`
}

type ClusterInfo struct {
	ConsulAddr []string `toml:"consul_addr"`
}

type ConfigInfo struct {
	DingKeys       []string          `toml:"ding_keys"`            //ding token
	DingUsers      []string          `toml:"ding_users"`           //default ding @list
	DingUserMap    map[string]string `toml:"ding_user_map"`        //config real editor ding @list
	IsAtAll        int               `toml:"is_at_all"`            //1: atall 2:not atall
	EnUpdateConsul int               `toml:"enable_update_consul"` //1: enable update consul 2:not
	EnDelConsul    int               `toml:"enable_delete_consul"` //1: enable delete consul 2:not
	ChklogRate     float64           `toml:"log_rate"`
	AccessToken    string            `toml:"access_token"` //apollo api auth token
}

type AppClusterCfg struct {
	Namespace     []string                  `toml:"namespace"`
	AppClusterMap map[string]AppClusterInfo `toml:"app_cluster_map"`
}

type ClusterCfg struct {
	ClusterMap map[string]ClusterInfo `toml:"cluster_map"`
}
