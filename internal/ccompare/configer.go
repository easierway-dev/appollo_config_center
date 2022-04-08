package ccompare

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"strconv"
)

type Config interface {
	GetConfigInfo() error
}

type GlobalConfig struct {
	AppConfigMap map[string]ConfigInfo  `toml:"app_config_map"`
	ClusterMap   map[string]ClusterInfo `toml:"cluster_map"`
	Timeout      int                    `toml:"timeout"`
}
type AppIdClustersInfo struct {
	// 全部的业务线集群信息
	EnvClustersInfoMap map[string][]*EnvClustersInfo
}

const (
	Default      = "default"
	ConsulAdd    = "47.252.4.203:8500"
	ClusterMap   = "cluster_map"
	AppConfigMap = "app_config_map"
	Timeout      = "timeout"
	LoopTime     = 10
)

// 全局配置
var GlobalConfiger *GlobalConfig
var AppIdClusters *AppIdClustersInfo

// 各个业务线对应的token
var AppIdAccessToken map[string]string
var ConsulClient map[string]*api.Client

// 获取global_config的配置
func (globalConfig *GlobalConfig) GetConfigInfo() error {
	GlobalConfiger = &GlobalConfig{}
	url := fmt.Sprintf("http://%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s", AgolloConfiger.PortalURL, "DEV", AgolloConfiger.AppID, AgolloConfiger.Cluster, AgolloConfiger.Namespace[0])
	fmt.Println("url=", url)
	// 默认的token
	globalInfo, _ := GetNamespaceInfo(url, "99648d96fe042df1b2280d6e7fe2278a08ff55b5")
	if globalInfo == nil {
		return errors.New("globalInfo is nil")
	}
	for _, item := range globalInfo.Items {
		switch item.Key {
		case ClusterMap:
			clusterConfig, _ := ParseClusterConfig(item.Value)
			globalConfig.ClusterMap = clusterConfig.ClusterMap
			break
		case AppConfigMap:
			appConfig, _ := ParseAppConfig(item.Value)
			globalConfig.AppConfigMap = appConfig.AppConfigMap
			break
		case Timeout:
			if item.Value == "" {
				globalConfig.Timeout = LoopTime
			}
			globalConfig.Timeout, _ = strconv.Atoi(item.Value)
			break
		default:
			break
		}
	}
	if len(globalConfig.ClusterMap) == 0 {
		info := ClusterInfo{ConsulAddr: []string{ConsulAdd}}
		m := map[string]ClusterInfo{}
		m[Default] = info
		globalConfig.ClusterMap = m
	}
	if globalConfig.Timeout == 0 {
		globalConfig.Timeout = LoopTime
	}
	GlobalConfiger = globalConfig
	fmt.Println("GlobalConfiger.ClusterMap = ", GlobalConfiger.ClusterMap)
	return nil
}
func (appIdClustersInfo *AppIdClustersInfo) GetConfigInfo() error {
	AppIdClusters = &AppIdClustersInfo{}
	url1 := fmt.Sprintf("http://%s/openapi/v1/apps", AgolloConfiger.PortalURL)
	// 动态获取业务线的token
	SetAppIDAccessToken()
	fmt.Println(AppIdAccessToken)
	//token, err := getDspToken(globalConfig.AccessToken)
	if len(AppIdAccessToken) == 0 {
		return errors.New("appID not correspond AccessToken")
	}
	// 只要获取某个业务线的token就可以，这里以dsp的token为例(存在隐患)
	// 每一个token对应一个业务线
	appInfo, _ := GetAppInfo(url1, AppIdAccessToken["dsp"])
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
		if token, ok := AppIdAccessToken[v.AppId]; ok {
			envClustersInfo, _ := GetEnvClustersInfo(url2, token)
			appIdClustersInfo.EnvClustersInfoMap[v.AppId] = envClustersInfo
		} else {
			fmt.Println(v.AppId + " not token ")
			continue
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

// 获取appId对应的accessToken
func SetAppIDAccessToken() {
	AppIdAccessToken = make(map[string]string, 6)
	// 配置文件里的token
	for key, config := range AppConfiger.AppConfigMap {
		AppIdAccessToken[key] = config.AccessToken
	}
	// 动态获取的token
	for key, config := range GlobalConfiger.AppConfigMap {
		// 有就替换为新的token
		if _, ok := AppIdAccessToken[key]; ok {
			AppIdAccessToken[key] = config.AccessToken
		} else {
			AppIdAccessToken[key] = config.AccessToken
		}
	}
}
func GetConsulClient() {
	ConsulClient = make(map[string]*api.Client)
	for _, consulAddr := range GlobalConfiger.ClusterMap {
		for i := 0; i < len(consulAddr.ConsulAddr); i++ {
			if _, ok := ConsulClient[consulAddr.ConsulAddr[i]]; ok {
				continue
			}
			fmt.Println("consulAddr:", consulAddr.ConsulAddr[i])
			cli, _ := NewClient(consulAddr.ConsulAddr[i])
			if cli == nil {
				fmt.Println("consulAddr:", consulAddr.ConsulAddr[i]+" connect failed")
				continue
			}
			// 每个集群对应一个client
			ConsulClient[consulAddr.ConsulAddr[i]] = cli
		}
	}
}

//
func (globalConfig *GlobalConfig) GetConfigInfo1() error {
	GlobalConfiger = &GlobalConfig{}
	url := fmt.Sprintf("http://%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s", AgolloConfiger.PortalURL, "DEV", AgolloConfiger.AppID, AgolloConfiger.Cluster, AgolloConfiger.Namespace[0])
	fmt.Println("url=", url)
	// 默认的token
	globalInfo, _ := GetNamespaceInfo(url, "280c6b92cd8ee4f1c5833b4bd22dfe44a4778ab5")
	if globalInfo == nil {
		return errors.New("globalInfo is nil")
	}
	for _, item := range globalInfo.Items {
		switch item.Key {
		case ClusterMap:
			clusterConfig, _ := ParseClusterConfig(item.Value)
			globalConfig.ClusterMap = clusterConfig.ClusterMap
			break
		case AppConfigMap:
			appConfig, _ := ParseAppConfig(item.Value)
			globalConfig.AppConfigMap = appConfig.AppConfigMap
			break
		case Timeout:
			if item.Value == "" {
				globalConfig.Timeout = LoopTime
			}
			globalConfig.Timeout, _ = strconv.Atoi(item.Value)
			break
		default:
			break
		}
	}
	if len(globalConfig.ClusterMap) == 0 {
		info := ClusterInfo{ConsulAddr: []string{ConsulAdd}}
		m := map[string]ClusterInfo{}
		m[Default] = info
		globalConfig.ClusterMap = m
	}
	if globalConfig.Timeout == 0 {
		globalConfig.Timeout = LoopTime
	}
	GlobalConfiger = globalConfig
	fmt.Println("GlobalConfiger.ClusterMap = ", GlobalConfiger.ClusterMap)
	return nil
}