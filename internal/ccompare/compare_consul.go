package ccompare

import (
	"fmt"
	"github.com/shima-park/agollo"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/capi"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
	"os"
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
	server, _ := NewAgolloServer(ccommon.AgolloConfiger)
	globalConfig = &GlobalConfig{}
	for _, ns := range ccommon.AgolloConfiger.Namespace {
		dycfg, err := ccommon.ParseDyConfig(server.Get("cluster_map", agollo.WithNamespace(ns)), server.Get("app_config_map", agollo.WithNamespace(ns)))
		if err != nil {
			ccommon.CLogger.Error(ccommon.DefaultDingType, "ParseDyConfig error: ", err.Error())
			panic(err)
		}
		globalConfig.ClusterMap = dycfg.ClusterConfig.ClusterMap
		cfg, err := ccommon.ParseAppClusterConfig(server.Get("app_cluster_map", agollo.WithNamespace(ns)))
		if err != nil {
			ccommon.CLogger.Error(ccommon.DefaultDingType, "ParseAppClusterConfig error: ", err.Error())
			panic(err)
		}
		globalConfig.AppClusterMap = cfg.AppClusterMap
	}
	fmt.Println("globalConfig=",globalConfig)
	return
}
func NewAgolloServer(agolloCfg *ccommon.AgolloCfg) (newAgo agollo.Agollo,err error){
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
		return nil,err
	}
	return
}