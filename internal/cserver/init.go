package cserver

import (
	"fmt"
	"os"

	"github.com/shima-park/agollo"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
)

func BuildGlobalAgollo(agolloCfg *ccommon.AgolloCfg, server *AgolloServer) error {
	newAgo, err := agollo.New(
		agolloCfg.ConfigServerURL,
		agolloCfg.AppID,
		agollo.Cluster(agolloCfg.Cluster),
		agollo.PreloadNamespaces(agolloCfg.Namespace...),
		agollo.AutoFetchOnCacheMiss(),
		agollo.FailTolerantOnBackupExists(),
		agollo.WithLogger(agollo.NewLogger(agollo.LoggerWriter(os.Stdout))),
	)
	if err != nil {
		fmt.Println("Build_global_agollo err: %s\n", err.Error())
		return err
	}
	server.BuildGAgollo(newAgo)
	return nil
}


func Init(server *AgolloServer)  error {
	//init config
	cfg, err := ccommon.ParseBaseConfig(ccommon.DirFlag)
	if err != nil {
		fmt.Println("ParseConfig error: %s\n", err.Error())
		return err
	}
	ccommon.AgolloConfiger =  cfg.AgolloCfg
	// init log
	cl, err := ccommon.NewconfigCenterLogger(cfg.LogCfg)
	if err != nil {
		fmt.Println("Load Logger err: ", err)
		return err
	}
	ccommon.CLogger = cl
	ccommon.CLogger.Info("Config=", ccommon.AgolloConfiger)
	ccommon.DyAgolloConfiger = make(map[string]*ccommon.DyAgolloCfg)
	//get global_config
	return BuildGlobalAgollo(ccommon.AgolloConfiger, server)
}
