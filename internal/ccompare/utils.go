package ccompare

import (
	"errors"
	"fmt"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/capi"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
)

type GlobalConfig struct {
	AppConfigMap map[string]ccommon.ConfigInfo  `toml:"app_config_map"`
	ClusterMap   map[string]ccommon.ClusterInfo `toml:"cluster_map"`
}

// 全局配置
var GlobalConfiger *GlobalConfig

// 获取全局配置
func GetApolloGlobalConfig() error {
	GlobalConfiger = &GlobalConfig{}

	url := fmt.Sprintf("http://%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s", ccommon.AgolloConfiger.PortalURL, "DEV", ccommon.AgolloConfiger.AppID, ccommon.AgolloConfiger.Cluster, ccommon.AgolloConfiger.Namespace[0])
	fmt.Println("url=", url)
	globalInfo, _ := capi.GetNamespaceInfo(url, "280c6b92cd8ee4f1c5833b4bd22dfe44a4778ab5")
	if globalInfo == nil {
		return errors.New("globalInfo is nil")
	}
	for _, item := range globalInfo.Items {
		if item.Key == "cluster_map" {
			clusterConfig, _ := ccommon.ParseClusterConfig(item.Value)
			GlobalConfiger.ClusterMap = clusterConfig.ClusterMap
		}
		if item.Key == "app_config_map" {
			appConfig, _ := ccommon.ParseAppConfig(item.Value)
			GlobalConfiger.AppConfigMap = appConfig.AppConfigMap
		}
	}
	return nil
}
