package ccompare

import (
	"errors"
	"fmt"
	"github.com/shima-park/agollo"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/capi"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
	"os"
	"sort"
)

const (
	DEV      = "DEV"
	DSPALIVG = "dsp_ali_vg"
	DSP      = "dsp"
)

type GlobalConfig struct {
	AppClusterMap map[string]ccommon.AppClusterInfo `toml:"app_cluster_map"`
	ClusterMap    map[string]ccommon.ClusterInfo    `toml:"cluster_map"`
	AccessToken   map[string]string                 `toml:"access_token"`
}
type ApolloProperty struct {
	AppID       string
	Env         string
	ClusterName string
}

var (
	NotContain = errors.New("not contain")
)
var appID []string
var envClustersInfo []*capi.EnvClustersInfo
var apolloProperty *ApolloProperty
var globalConfig *GlobalConfig
func getAppID() error{
	url2 := fmt.Sprintf("http://%s/openapi/v1/apps", ccommon.AgolloConfiger.PortalURL)
	token, err := getDspToken(globalConfig.AccessToken)
	if err != nil {
		return err
	}
	appInfo, _ := capi.GetAppInfo(url2,token)
	fmt.Println("url2=", url2)
	fmt.Println("appInfo=", appInfo)
	if len(appInfo) == 0 {
		fmt.Println("appInfo is nil ")
		return errors.New("appInfo is nil ")
	}
	for k, v := range appInfo {
		appID[k] = v.AppId
	}
	fmt.Println("appID=", appID)
	return nil
}
func getEnvClustersInfo() error{
	url1 := fmt.Sprintf("http://%s/openapi/v1/apps/%s/envclusters", ccommon.AgolloConfiger.PortalURL, ccommon.AgolloConfiger.AppID)
	token, err := getDspToken(globalConfig.AccessToken)
	if err != nil {
		return err
	}
	envClustersInfo, _ = capi.GetEnvClustersInfo(url1,token)
	if len(envClustersInfo) == 0 {
		fmt.Println("appInfo is nil ")
		return errors.New("appInfo is nil ")
	}
	fmt.Println("url1=", url1)
	fmt.Println("ecinfo=", envClustersInfo)
	return nil
}

// 获取全局配置
func GetApolloGlobalConfig() {
	// 生成一个agolloServer
	server, _ := NewAgolloServer(ccommon.AgolloConfiger)
	globalConfig = &GlobalConfig{}
	for _, ns := range ccommon.AgolloConfiger.Namespace {
		dyCfg, err := ccommon.ParseDyConfig(server.Get("cluster_map", agollo.WithNamespace(ns)), server.Get("app_config_map", agollo.WithNamespace(ns)))
		if err != nil {
			ccommon.CLogger.Error(ccommon.DefaultDingType, "ParseDyConfig error: ", err.Error())
			panic(err)
		}
		fmt.Println("dyCfg=",dyCfg)
		globalConfig.ClusterMap = dyCfg.ClusterConfig.ClusterMap
		for key, info := range dyCfg.AppConfig.AppConfigMap {
			globalConfig.AccessToken[key] = info.AccessToken
		}
		cfg, err := ccommon.ParseAppClusterConfig(server.Get("app_cluster_map", agollo.WithNamespace(ns)))
		if err != nil {
			ccommon.CLogger.Error(ccommon.DefaultDingType, "ParseAppClusterConfig error: ", err.Error())
			panic(err)
		}
		globalConfig.AppClusterMap = cfg.AppClusterMap
	}
	fmt.Println("globalConfig=", globalConfig)
	return
}
func NewAgolloServer(agolloCfg *ccommon.AgolloCfg) (newAgo agollo.Agollo, err error) {
	newAgo, err = agollo.New(
		agolloCfg.ConfigServerURL,
		agolloCfg.AppID,
		agollo.Cluster(agolloCfg.Cluster),
		agollo.PreloadNamespaces(agolloCfg.Namespace...),
		agollo.AutoFetchOnCacheMiss(),
		agollo.FailTolerantOnBackupExists(),
		agollo.WithLogger(agollo.NewLogger(agollo.LoggerWriter(os.Stdout))),
	)
	if err != nil {
		fmt.Println("Build_Global_Agollo err: %s\n", err.Error())
		return nil, err
	}
	return
}
func in(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}
func applyProperty() error {
	var isContainDEV bool
	isContainDsp := in(DSP, appID)
	if !isContainDsp {
		return errors.New("not contain Dsp")
	}
	for i := 0; i < len(envClustersInfo); i++ {
		if envClustersInfo[i].Env == DEV {
			isContainDEV = true
			isContainDspALiVg := in(DSPALIVG, envClustersInfo[i].Clusters)
			if !isContainDspALiVg {
				return errors.New("not contain DspALiVg")
			}
		}
		if !isContainDEV {
			return errors.New("not contain DEV")
		}
	}
	apolloProperty = &ApolloProperty{AppID: DSP, Env: DEV, ClusterName: DSPALIVG}
	return nil
}

// 通过集群名，appID，namespace查找对应的信息：获取集群下所有Namespace信息接口，在进行细分每一个namespace
func (apolloProperty *ApolloProperty) GetNameSpaceInfo() {
	url := fmt.Sprintf("http://%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces", ccommon.AgolloConfiger.PortalURL, apolloProperty.Env, apolloProperty.AppID, apolloProperty.ClusterName)
	nSAllInfo, _ := capi.GetNamespaceInfo(url, globalConfig.AccessToken[apolloProperty.AppID])
	fmt.Println("nSAllInfo=", nSAllInfo)
}
func NewNameSpaceInfo()  {
	// 获取全局配置
	GetApolloGlobalConfig()
	// 获取全局AppID
	getAppID()
	// 获取全局集群信息
	getEnvClustersInfo()
	// 验证DSP并赋值
	applyProperty()
	// 获取对应namespace的信息
	apolloProperty.GetNameSpaceInfo()
}
func getDspToken(m map[string]string) (token string ,err error){
	for _, token = range m{
		if _,ok:=m[DSP];ok{
			return
		}
	}
	return "",errors.New("not contain DspToken")
}