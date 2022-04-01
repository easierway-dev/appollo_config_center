package ccompare

import (
	"errors"
	"fmt"
)

type Config interface {
	GetConfigInfo() error
}

type GlobalConfig struct {
	AppConfigMap map[string]ConfigInfo  `toml:"app_config_map"`
	ClusterMap   map[string]ClusterInfo `toml:"cluster_map"`
}
type AppIdClustersInfo struct {
	// 全部的业务线
	AppID []string
	// 集群信息
	EnvClustersInfoMap map[string][]*EnvClustersInfo
}

// 全局配置
var GlobalConfiger *GlobalConfig
var AppIdClusters *AppIdClustersInfo

func (globalConfig *GlobalConfig) GetConfigInfo() error {
	GlobalConfiger = &GlobalConfig{}
	url := fmt.Sprintf("http://%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s", AgolloConfiger.PortalURL, "DEV", AgolloConfiger.AppID, AgolloConfiger.Cluster, AgolloConfiger.Namespace[0])
	fmt.Println("url=", url)
	globalInfo, _ := GetNamespaceInfo(url, "280c6b92cd8ee4f1c5833b4bd22dfe44a4778ab5")
	if globalInfo == nil {
		return errors.New("globalInfo is nil")
	}
	for _, item := range globalInfo.Items {
		if item.Key == "cluster_map" {
			clusterConfig, _ := ParseClusterConfig(item.Value)
			globalConfig.ClusterMap = clusterConfig.ClusterMap
		}
		if item.Key == "app_config_map" {
			appConfig, _ := ParseAppConfig(item.Value)
			globalConfig.AppConfigMap = appConfig.AppConfigMap
		}
	}
	GlobalConfiger = globalConfig
	return nil
}
func (appIdClustersInfo *AppIdClustersInfo) GetConfigInfo() error {
	AppIdClusters = &AppIdClustersInfo{}
	url1 := fmt.Sprintf("http://%s/openapi/v1/apps", AgolloConfiger.PortalURL)
	_, err := getAccessToken(GlobalConfiger.AppConfigMap)
	fmt.Println(appIdAccessToken)
	//token, err := getDspToken(globalConfig.AccessToken)
	if err != nil {
		return err
	}
	// 只要获取某个业务线的token就可以，这里以dsp的token为例
	appInfo, _ := GetAppInfo(url1, appIdAccessToken["dsp"])
	//fmt.Println("url2=", url2)
	//fmt.Println("appInfo=", appInfo)
	if len(appInfo) == 0 {
		fmt.Println("appInfo is nil ")
		return errors.New("appInfo is nil ")
	}
	appIdClustersInfo.EnvClustersInfoMap = make(map[string][]*EnvClustersInfo)
	for _, v := range appInfo {
		//appIdClustersInfo.AppID = append(appIdClustersInfo.AppID, v.AppId)
		url2 := fmt.Sprintf("http://%s/openapi/v1/apps/%s/envclusters", AgolloConfiger.PortalURL, v.AppId)
		for _, token := range appIdAccessToken {
			envClustersInfo, _ := GetEnvClustersInfo(url2, token)
			appIdClustersInfo.EnvClustersInfoMap[v.AppId] = envClustersInfo
		}
		if len(appIdClustersInfo.EnvClustersInfoMap) == 0 {
			fmt.Println("EnvClustersInfoMap is nil ")
			return errors.New("EnvClustersInfoMap is nil ")
		}
		fmt.Println("ecinfo=", appIdClustersInfo.EnvClustersInfoMap)
	}
	AppIdClusters = appIdClustersInfo
	return nil
}

//// 获取所有业务线
//func getAllAppID() error {
//	url2 := fmt.Sprintf("http://%s/openapi/v1/apps", AgolloConfiger.PortalURL)
//	_, err := getAccessToken(GlobalConfiger.AppConfigMap)
//	fmt.Println(appIdAccessToken)
//	//token, err := getDspToken(globalConfig.AccessToken)
//	if err != nil {
//		return err
//	}
//	// 只要获取某个业务线的token就可以，这里以dsp的token为例
//	appInfo, _ := GetAppInfo(url2, appIdAccessToken["dsp"])
//	//fmt.Println("url2=", url2)
//	//fmt.Println("appInfo=", appInfo)
//	if len(appInfo) == 0 {
//		fmt.Println("appInfo is nil ")
//		return errors.New("appInfo is nil ")
//	}
//
//	//fmt.Println("appID=", appID)
//	return nil
//}
//
//// 获取所有业务线的集群信息
//func getEnvClustersInfo(appID string) (map[string][]*EnvClustersInfo, error) {
//	url1 := fmt.Sprintf("http://%s/openapi/v1/apps/%s/envclusters", AgolloConfiger.PortalURL, appID)
//	envClustersInfoMap = make(map[string][]*EnvClustersInfo)
//	for _, token := range appIdAccessToken {
//		envClustersInfo, _ := GetEnvClustersInfo(url1, token)
//		envClustersInfoMap[appID] = envClustersInfo
//	}
//	if len(envClustersInfoMap) == 0 {
//		fmt.Println("EnvClustersInfoMap is nil ")
//		return nil, errors.New("EnvClustersInfoMap is nil ")
//	}
//	fmt.Println("ecinfo=", envClustersInfoMap)
//
//	return envClustersInfoMap, nil
//}
