package cserver

import (
	"log"
	"os"
	"strings"

	"github.com/shima-park/agollo"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
)

func Init() (*AgolloServer, error) {
	var server *AgolloServer

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//init config
	cfg, err := ccommon.ParseBaseConfig(ccommon.DirFlag)
	if err != nil {
		log.Printf("ParseConfig error: %s\n", err.Error())
		return nil, err
	}
	ccommon.AgolloConfiger =  cfg
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
	for AppID, cNameList := range cfg.CenterCfg.AppClusterMap {
		for _, cName := range cNameList {
			cNameArr := strings.SplitN(cName, "_", 2)
			consulAddr := cfg.CenterCfg.ClusterMap[cName].ConsulAddr
			if len(cNameArr) == 2 && consulAddr != "" {
				cluster := cNameArr[1]
				newAgo, err := agollo.New(
					cfg.CenterCfg.ConfigServerURL,
					AppID,
					agollo.Cluster(cluster),
					agollo.PreloadNamespaces("juno"),
					agollo.AutoFetchOnCacheMiss(),
					agollo.FailTolerantOnBackupExists(),
					agollo.WithLogger(agollo.NewLogger(agollo.LoggerWriter(os.Stdout))),
				)
				if err != nil {
					panic(err)
				}
				work := Worker{
					AgolloClient: newAgo,
					Cluster:      cluster,
					ConsulAddr:   consulAddr,
				}
				server.AddWorker(work)
			} else {
				ccommon.CLogger.Runtime.Errorf("invalue appClusterInfo AppClusterMap=", cfg.CenterCfg.AppClusterMap, "consulAddr=", consulAddr)
				continue
			}
		}
	}
	return server, nil
}
