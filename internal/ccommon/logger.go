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

func GetDingInfo(appid string, itype string) (dingKeys []string, dingusers []string, userMap map[string]string, isAtall bool) {
	if appid == "" && itype == "info" {
		return
	}
	//local config
	namespace := DefaultNamespace
	//default config
	isAtallTmp := 0
	//uniq appid config
	dingKeys, dingusers, userMap, isAtallTmp = InitAppCfgMap(AppConfiger, appid)
	if DyAgolloConfiger != nil {
		if dyAgoCfg, ok := DyAgolloConfiger[namespace]; ok {
			dingKeys, dingusers, userMap, isAtallTmp = InitAppCfgMap(dyAgoCfg.AppConfig, appid)
		}
	}
	if isAtallTmp == 1 {
		isAtall = true
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
	if _, ok := interface{}(args[0]).(string); ok {
		dingkeys, dingusers, _, isatall := GetDingInfo(args[0].(string), "warn")

		cnotify.SendText(dingkeys, fmt.Sprintf("%s", args), dingusers, isatall)
		this.Runtime.Warn(args)
	} else {
		dingkeys, dingusers, usermap, isatall := GetDingInfo(args[1].(string), "warn")
		switch t := args[0].(type) {
		case []string:
			keyStringValues := []string{}
			for _, username := range t {
				if userphone, ok := usermap[username]; ok {
					keyStringValues = append(keyStringValues, userphone)
				}
			}
			dingusers = append(dingusers, keyStringValues...)
		default:
			fmt.Println("dingusers type error , need []string")
		}
		cnotify.SendText(dingkeys, fmt.Sprintf("%s", args[1:]), dingusers, isatall)
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
func InitAppCfgMap(appConfig *AppCfg, appid string) (dingKeys []string, dingUsers []string, userMap map[string]string, isAtAllTmp int) {
	if appConfig == nil{
		return
	}
	dingKeys = appConfig.DingKeys
	dingUsers = appConfig.DingUsers
	userMap = appConfig.DingUserMap
	isAtAllTmp = appConfig.IsAtAll
	if appConfig.AppConfigMap == nil {
		return
	}
	if _,ok := appConfig.AppConfigMap[appid];!ok{
		return
	}
	if len(appConfig.AppConfigMap[appid].DingKeys) > 0 {
		dingKeys = appConfig.AppConfigMap[appid].DingKeys
	}
	if len(appConfig.AppConfigMap[appid].DingUsers) > 0 {
		dingUsers = appConfig.AppConfigMap[appid].DingUsers
	}
	for key, value := range appConfig.AppConfigMap[appid].DingUserMap {
		if userMap == nil {
			userMap = make(map[string]string)
		}
		userMap[key] = value
	}
	if appConfig.AppConfigMap[appid].IsAtAll != 0 {
		isAtAllTmp = appConfig.AppConfigMap[appid].IsAtAll
	}
	return
}
