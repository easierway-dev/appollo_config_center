package ccompare

import (
	"fmt"
	"math/rand"
)

func Init() error {
	//init config
	cfg, err := ParseBaseConfig(DirFlag)
	if err != nil {
		fmt.Println("ParseConfig error: %s\n", err.Error())
		return err
	}
	AgolloConfiger = cfg.AgolloCfg
	AppConfiger = cfg.AppCfg
	ChklogRamdom = rand.Float64()
	ChklogRate = AppConfiger.ChklogRate
	// init log
	cl, err := NewconfigCenterLogger(cfg.LogCfg)
	if err != nil {
		fmt.Println("Load Logger err: ", err)
		return err
	}
	CLogger = cl
	CLogger.Info(DefaultDingType, "Config=", AgolloConfiger)
	DyAgolloConfiger = make(map[string]*DyAgolloCfg)
	//get global_config
	//return BuildGlobalAgollo(ccommon.AgolloConfiger, server)
	return nil
}
