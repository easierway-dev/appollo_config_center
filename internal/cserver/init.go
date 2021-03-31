package cserver

import (
	"strings"
        "os"
        "log"
	"github.com/spf13/viper"
        "gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
	"github.com/shima-park/agollo"
)


func Init()(*AgolloServer, error) { 
	var server *AgolloServer

        log.SetFlags(log.Lshortfile | log.LstdFlags)
        //init config
        cfg, err := ccommon.ParseBaseConfig(viper.GetString(ccommon.DirFlag))
        if err != nil {
                log.Printf("ParseConfig error: %s\n", err.Error())
                return nil, err
        }
        //ccommon.CConfiger =  cfg.CenterCfg
        // init log
        cl, err := ccommon.NewconfigCenterLogger(cfg.LogCfg)
        if err != nil {
                log.Println("Load Logger err: ", err)
                return nil, err
        }
        ccommon.CLogger = cl
        cl.Runtime.Infof("Config=[%v],", cfg)

	// server
	server = NewAgolloServer()
	for AppId, cNameList := range cfg.CenterCfg.AppClusterMap {
		for _, cName := range cNameList {
			cNameArr := strings.SplitN(cName, "_", 1)
			consulAddr := cfg.CenterCfg.ClusterMap[cName].ClusterDetail["consul_addr"]
			if len(cNameArr) == 2 {
				cluster := cNameArr[1]
				newAgo, err := agollo.New(
					cfg.CenterCfg.ConfigServerUrl,
					AppId,
					agollo.Cluster(cluster),
					agollo.PreloadNamespaces("application"),
					agollo.AutoFetchOnCacheMiss(),
					agollo.FailTolerantOnBackupExists(),
					agollo.WithLogger(agollo.NewLogger(agollo.LoggerWriter(os.Stdout))),
				)
				if err != nil {
					panic(err)
				}
				work := Worker{
					AgolloClient:  newAgo,
					Cluster:        cluster,
					ConsulAddr:    consulAddr,
				}
				server.AddWorker(work)
			} else {
				continue
			}
		}
	}
	return server, nil
}
