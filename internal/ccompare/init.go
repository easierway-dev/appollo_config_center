package ccompare

import (
	"fmt"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
	"math/rand"
)

func Init() error {
	//init config
	cfg, err := ccommon.ParseBaseConfig(ccommon.DirFlag)
	if err != nil {
		fmt.Println("ParseConfig error: %s\n", err.Error())
		return err
	}
	ccommon.AgolloConfiger = cfg.AgolloCfg
	ccommon.AppConfiger = cfg.AppCfg
	ccommon.ChklogRamdom = rand.Float64()
	ccommon.ChklogRate = ccommon.AppConfiger.ChklogRate
	// init log
	cl, err := ccommon.NewconfigCenterLogger(cfg.LogCfg)
	if err != nil {
		fmt.Println("Load Logger err: ", err)
		return err
	}
	ccommon.CLogger = cl
	ccommon.CLogger.Info(ccommon.DefaultDingType, "Config=", ccommon.AgolloConfiger)
	ccommon.DyAgolloConfiger = make(map[string]*ccommon.DyAgolloCfg)
	//get global_config
	//return BuildGlobalAgollo(ccommon.AgolloConfiger, server)
	return nil
}
