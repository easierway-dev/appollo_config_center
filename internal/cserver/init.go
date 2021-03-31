package cserver

import (
        "os"
        "log"
	"github.com/spf13/viper"
        "appollo_config_center/internal/ccommon"
)


func Init() {
        log.SetFlags(log.Lshortfile | log.LstdFlags)
        //init config
        cfg, err := ccommon.ParseBaseConfig(viper.GetString(ccommon.DirFlag))
        if err != nil {
                log.Printf("ParseConfig error: %s\n", err.Error())
                os.Exit(1)
        }
        //ccommon.CConfiger =  cfg.centerCfg
        // init log
        cl, err := ccommon.NewconfigCenterLogger(cfg.LogCfg)
        if err != nil {
                log.Println("Load Logger err: ", err)
                os.Exit(1)
        }
        ccommon.CLoger = cl
        cl.Runtime.Infof("Config=[%v],", cfg)

	// server
	server = NewAgolloServer()
	for AppId, cNameList := range cfg.centerCfg.AppClusterMap {
		for _, cName := range cNameList {
			cNameArr := SplitN(cName, "_", 1)
			if len(cNameArr) == 2 {
				cluster := SplitN(cName, "_", 1)[1]
				newAgo, err := agollo.New(
					cfg.ConfigServerUrl,
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
				work := &Worker{
					Agollo_client:  aClient,
					Cluster:        cluster,
					Consul_addr:    consulAddr,
				}
				server.AddWorker(work)
			} else {
				continue
			}
		}
	}
}
