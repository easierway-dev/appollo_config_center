package ccompare

import (
	"github.com/shima-park/agollo"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
)

type GlobalConfig struct {
	AppConfigMap map[string]ccommon.ConfigInfo  `toml:"app_config_map"`
	ClusterMap   map[string]ccommon.ClusterInfo `toml:"cluster_map"`
}

// 全局配置
var GlobalConfiger *GlobalConfig

// 获取全局配置
func GetApolloGlobalConfig(server *AgolloServer) {
	GlobalConfiger = &GlobalConfig{}
	for _, ns := range ccommon.AgolloConfiger.Namespace {
		dyCfg, err := ccommon.ParseDyConfig(server.gAgollo.Get("cluster_map", agollo.WithNamespace(ns)), server.gAgollo.Get("app_config_map", agollo.WithNamespace(ns)))
		if err != nil {
			//ccommon.CLogger.Error(ccommon.DefaultDingType, "ParseDyConfig error: ", err.Error())
			panic(err)
		}
		GlobalConfiger.AppConfigMap = dyCfg.AppConfig.AppConfigMap
		GlobalConfiger.ClusterMap = dyCfg.ClusterConfig.ClusterMap
	}
	return
}
