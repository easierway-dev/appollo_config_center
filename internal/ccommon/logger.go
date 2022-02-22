package ccommon

import (
	"errors"
	"fmt"

	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/cnotify"
	"gitlab.mobvista.com/voyager/zlog"
)

var CLogger *ccLogger

const (
	DefaultDingType     = ""
	InitDingType        = "init"
	DefaultPollDingType = "poll"
)

type LogCfg struct {
	Runtime *zlog.Ops `toml:"Runtime_log"`
}

type ccLogger struct {
	Runtime zlog.Logger
}

func NewconfigCenterLogger(logCfg *LogCfg) (*ccLogger, error) {
	var err error
	if logCfg == nil {
		return nil, errors.New("logCfg is nil")
	}
	logger := &ccLogger{}
	if logger.Runtime, err = zlog.NewZLog(logCfg.Runtime); err != nil {
		return nil, err
	}
	return logger, nil
}

func GetDingInfo(appid string, itype string) (dingKeys []string, dingusers []string, userMap map[string]string, isAtAll bool) {
	if appid == "" && itype == "info" {
		return
	}
	//local config
	namespace := DefaultNamespace
	//default config
	dingKeys = AppConfiger.DingKeys
	dingusers = AppConfiger.DingUsers
	userMap = AppConfiger.DingUserMap
	isAtAllTmp := AppConfiger.IsAtAll
	//uniq appid config

	dingKeys, dingusers, userMap, isAtAllTmp = InitAppConfigMap(AppConfiger.AppConfigMap, appid, isAtAllTmp)
	if DyAgolloConfiger == nil {
		return
	}
	//apollo global_config
	dyAgoCfg, ok := DyAgolloConfiger[namespace]
	if !ok {
		return
	}
	//default config
	dingKeys, dingusers, userMap, isAtAllTmp = InitDyAppConfigMap(dyAgoCfg.AppConfig, appid, isAtAllTmp)
	//uniq appid config
	dingKeys, dingusers, userMap, isAtAllTmp = InitAppConfigMap(dyAgoCfg.AppConfig.AppConfigMap, appid, isAtAllTmp)
	if isAtAllTmp == 1 {
		isAtAll = true
	}
	return
}

func (this *ccLogger) Info(args ...interface{}) {
	if this == nil || this.Runtime == nil {
		return
	}
	//dingkeys,dingusers,_,isatall := GetDingInfo(args[0].(string), "info")
	//cnotify.SendText(dingkeys,fmt.Sprintf("%s",args),dingusers,isatall)
	this.Runtime.Info(args)
}

func (this *ccLogger) Warn(args ...interface{}) {
	if this == nil || this.Runtime == nil {
		return
	}
	if _, ok := args[0].(string); ok {
		dingkeys, dingusers, _, isatall := GetDingInfo(args[0].(string), "warn")

		cnotify.SendText(dingkeys, fmt.Sprintf("%s", args), dingusers, isatall)
		this.Runtime.Warn(args)
	} else {
		dingKeys, dingUsers, userMap, isAtAll := GetDingInfo(args[1].(string), "warn")
		switch t := args[0].(type) {
		case []string:
			var keyStringValues []string
			for _, username := range t {
				if userphone, ok := userMap[username]; ok {
					keyStringValues = append(keyStringValues, userphone)
				}
			}
			dingUsers = append(dingUsers, keyStringValues...)
		default:
			fmt.Println("dingusers type error , need []string")
		}
		cnotify.SendText(dingKeys, fmt.Sprintf("%s", args[1:]), dingUsers, isAtAll)
		this.Runtime.Warn(args)
	}
}

func (this *ccLogger) Error(args ...interface{}) {
	if this == nil || this.Runtime == nil {
		return
	}
	dingkeys, dingusers, _, isatall := GetDingInfo(args[0].(string), "err")
	cnotify.SendText(dingkeys, fmt.Sprintf("%s", args), dingusers, isatall)
	this.Runtime.Error(args)
}

func (this *ccLogger) Infof(format string, args ...interface{}) {
	if this == nil || this.Runtime == nil {
		return
	}
	this.Runtime.Infof(format, args)
}

func (this *ccLogger) Warnf(format string, args ...interface{}) {
	if this == nil || this.Runtime == nil {
		return
	}
	this.Runtime.Warnf(format, args)
}

func (this *ccLogger) Errorf(format string, args ...interface{}) {
	if this == nil || this.Runtime == nil {
		return
	}
	this.Runtime.Errorf(format, args)
}
func InitAppConfigMap(appConfigMap map[string]ConfigInfo,appid string,isAtAllTmp int) (dingKeys []string, dingusers []string, userMap map[string]string, isAtAll int){
	if appConfigMap == nil {
		return
	}
	if _, ok := appConfigMap[appid]; !ok {
		return
	}
	if len(appConfigMap[appid].DingKeys) > 0 {
		dingKeys = AppConfiger.AppConfigMap[appid].DingKeys
	}
	if len(appConfigMap[appid].DingUsers) > 0 {
		dingusers = AppConfiger.AppConfigMap[appid].DingUsers
	}
	for key, value := range appConfigMap[appid].DingUserMap {
		if userMap == nil {
			userMap = map[string]string{}
		}
		userMap[key] = value
	}
	if appConfigMap[appid].IsAtAll != 0 {
		isAtAllTmp = AppConfiger.AppConfigMap[appid].IsAtAll
	}
	return
}
func InitDyAppConfigMap(dyAppConfigMap *AppCfg,appid string,isAtAllTmp int) (dingKeys []string, dingusers []string, userMap map[string]string, isAtAll int){
	if dyAppConfigMap == nil {
		return
	}
	if len(dyAppConfigMap.AppConfigMap[appid].DingKeys) > 0 {
		dingKeys = AppConfiger.AppConfigMap[appid].DingKeys
	}
	if len(dyAppConfigMap.AppConfigMap[appid].DingUsers) > 0 {
		dingusers = AppConfiger.AppConfigMap[appid].DingUsers
	}
	for key, value := range dyAppConfigMap.AppConfigMap[appid].DingUserMap {
		if userMap == nil {
			userMap = map[string]string{}
		}
		userMap[key] = value
	}
	if dyAppConfigMap.AppConfigMap[appid].IsAtAll != 0 {
		isAtAllTmp = dyAppConfigMap.AppConfigMap[appid].IsAtAll
	}
	return
}