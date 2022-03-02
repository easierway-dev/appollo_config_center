package ccompare

import (
	"fmt"
	"github.com/shima-park/agollo"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/capi"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
)

type GlobalConfig struct {
	AppClusterMap map[string]ccommon.AppClusterInfo `toml:"app_cluster_map"`
	ClusterMap map[string]ccommon.ClusterInfo `toml:"cluster_map"`
}
var appID []string

func getAppID()  {
	url2 := fmt.Sprintf("http://%s/openapi/v1/apps", ccommon.AgolloConfiger.PortalURL)
	appInfo, _ := capi.GetAppInfo(url2, ccommon.Configer.AccessToken)
	fmt.Println("url2=", url2)
	fmt.Println("appInfo=", appInfo)
	for k,v :=range appInfo{
		if v.AppId == "dsp"{
			continue
		}
		appID[k] = v.AppId
	}
	fmt.Println("appID=",appID)

}
func getEnvClustersInfo(){
	url1 := fmt.Sprintf("http://%s/openapi/v1/apps/%s/envclusters", ccommon.AgolloConfiger.PortalURL, ccommon.AgolloConfiger.AppID)
	ecinfo, _ := capi.GetEnvClustersInfo(url1, ccommon.Configer.AccessToken)
	fmt.Println("url1=", url1)
	fmt.Println("ecinfo=", ecinfo)
}
func GetApolloGlobalConfig() (globalConfig *GlobalConfig){
	var agollo1 agollo.Agollo
	fmt.Println("Namespace=",ccommon.AgolloConfiger.Namespace)
	for _, ns := range ccommon.AgolloConfiger.Namespace {
		dycfg, err := ccommon.ParseDyConfig(agollo1.Get("cluster_map", agollo.WithNamespace(ns)), agollo1.Get("app_config_map", agollo.WithNamespace(ns)))
		if err != nil {
			ccommon.CLogger.Error(ccommon.DefaultDingType, "ParseDyConfig error: ", err.Error())
			panic(err)
		}
		cfg, err := ccommon.ParseAppClusterConfig(agollo1.Get("app_cluster_map", agollo.WithNamespace(ns)))
		if err != nil {
			ccommon.CLogger.Error(ccommon.DefaultDingType, "ParseAppClusterConfig error: ", err.Error())
			panic(err)
		}
		globalConfig = &GlobalConfig{AppClusterMap: cfg.AppClusterMap, ClusterMap: dycfg.ClusterConfig.ClusterMap}
	}
	fmt.Println("globalConfig=",globalConfig)
	return
}