package cserver

import (
	"log"
	"os"
	"strings"

	"github.com/shima-park/agollo"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
)

func init_dyagolloconfiger(agolloCfg *ccommon.AgolloCfg, server *AgolloServer) {
	ccommon.DyAgolloConfiger = &ccommon.DyAgolloCfg{}
	newAgo, err := agollo.New(
		agolloCfg.ConfigServerURL,
		agolloCfg.AppID,
		agollo.Cluster(agolloCfg.Cluster),
		agollo.PreloadNamespaces(agolloCfg.Namespace),
		agollo.AutoFetchOnCacheMiss(),
		agollo.FailTolerantOnBackupExists(),
		agollo.WithLogger(agollo.NewLogger(agollo.LoggerWriter(os.Stdout))),
	)
	if err != nil {
		panic(err)
	}
	cfg, err := ccommon.ParseAppClusterConfig(newAgo.Get("app_cluster_map"))
	if err != nil {
			log.Printf("ParseAppClusterConfig error: %s\n", err.Error())
			panic(err)
	} else {
		ccommon.DyAgolloConfiger.AppClusterConfig = cfg
	}
	
	cfg1, err := ccommon.ParseClusterConfig(newAgo.Get("cluster_map"))
	if err != nil {
			log.Printf("ParseClusterConfig error: %s\n", err.Error())
			panic(err)
	} else {
		ccommon.DyAgolloConfiger.ClusterConfig = cfg1
	}
	
	cfg2, err := ccommon.ParseAppConfig(newAgo.Get("app_config_map"))
	if err != nil {
			log.Printf("ParseAppConfig error: %s\n", err.Error())
			panic(err)
	} else {
		ccommon.DyAgolloConfiger.AppConfig = cfg2
	}
	
	errorCh := newAgo.Start()
	watchCh := newAgo.Watch()

	go func(cluster string) {
		for {
			select {
			case <-server.ctx.Done():
			    ccommon.CLogger.Runtime.Errorf(cluster, "watch quit...")
			    return
			case err := <-errorCh:
				 ccommon.CLogger.Runtime.Errorf("Error:", err)
			case update := <-watchCh:
				if value, ok := update.NewValue["app_cluster_map"]; ok {
					cfg, err = ccommon.ParseAppClusterConfig(value.(string))
					if err != nil {
							log.Printf("ParseAppClusterConfig error: %s\n", err.Error())
							panic(err)
					} else {
						ccommon.DyAgolloConfiger.AppClusterConfig = cfg
					}
				}
				if value, ok := update.NewValue["cluster_map"]; ok {
					cfg1, err := ccommon.ParseClusterConfig(value.(string))
					if err != nil {
							log.Printf("ParseClusterConfig error: %s\n", err.Error())
					} else {
						ccommon.DyAgolloConfiger.ClusterConfig = cfg1
					}
				}
				if value, ok := update.NewValue["app_config_map"]; ok {
					cfg2, err = ccommon.ParseAppConfig(value.(string))
					if err != nil {
							log.Printf("ParseAppConfig error: %s\n", err.Error())
					} else {
						ccommon.DyAgolloConfiger.AppConfig = cfg2
					}
				}
				ccommon.CLogger.Runtime.Infof("Apollo cluster(%s) namespace(%s) old_value:(%v) new_value:(%v) error:(%v)\n",
					cluster, update.Namespace, update.OldValue, update.NewValue, update.Error)
			}
		}
	}(agolloCfg.Cluster)
}

func NewAgolloWorker(server *AgolloServer){
    for appID, cNameList := range ccommon.DyAgolloConfiger.AppClusterConfig.AppClusterMap {
            for _, cName := range cNameList {
                    cNameArr := strings.SplitN(cName, "_", 2)
                    if len(cNameArr) == 2 {
                            cluster := cNameArr[1]
                            newAgo, err := agollo.New(
                                    ccommon.AgolloConfiger.ConfigServerURL,
                                    appID,
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
				    AppID: appID,
				    Cluster: cluster,
                                    ClusterID: cName,
                            }
                            server.AddWorker(work)
                    } else {
                            ccommon.CLogger.Runtime.Errorf("invalue appClusterInfo AppClusterMap=", ccommon.DyAgolloConfiger.AppClusterConfig)
                            continue
                    }
            }
    }
}

func Init(server *AgolloServer)  error {

	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//init config
	cfg, err := ccommon.ParseBaseConfig(ccommon.DirFlag)
	if err != nil {
		log.Printf("ParseConfig error: %s\n", err.Error())
		return err
	}
	ccommon.AgolloConfiger =  cfg.AgolloCfg
	// init log
	cl, err := ccommon.NewconfigCenterLogger(cfg.LogCfg)
	if err != nil {
		log.Println("Load Logger err: ", err)
		return err
	}
	ccommon.CLogger = cl
	cl.Runtime.Infof("Config=[%v],", cfg)
	//get global_config
	init_dyagolloconfiger(cfg.AgolloCfg, server)
	// server
	NewAgolloWorker(server)
	return nil
}
