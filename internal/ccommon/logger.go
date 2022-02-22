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
	dingKeys = AppConfiger.DingKeys
	dingusers = AppConfiger.DingUsers
	userMap = make(map[string]string)
	isAtallTmp := AppConfiger.IsAtAll
	//uniq appid config
	if AppConfiger.AppConfigMap != nil {
		if _, ok := AppConfiger.AppConfigMap[appid]; ok {
			dingKeys, dingusers, userMap, isAtallTmp = InitAppConfigMap(AppConfiger.AppConfigMap, appid,isAtallTmp)
			for k, v := range userMap {
				fmt.Println("userMapKey为:" + k + " " + "userMapValue为" + v + "\n")
			}
		}
	}
	for _, v := range dingKeys {
		fmt.Println("dingKey为:" + v + "\n")
	}
	for _, v := range dingusers {
		fmt.Println("dingusers为:" + v + "\n")
	}
	for k, v := range userMap {
		fmt.Println("userMapKey为:" + k + " " + "userMapValue为" + v + "\n")
	}
	fmt.Println("isAtAllTmp为:" + string(isAtallTmp))
	//apollo global_config
	if DyAgolloConfiger != nil {
		if dyAgoCfg, ok := DyAgolloConfiger[namespace]; ok {
			//default config
			if dyAgoCfg.AppConfig != nil {
				dingKeys, dingusers, userMap, isAtallTmp = InitDyAppConfigMap(dyAgoCfg.AppConfig, appid,isAtallTmp)
			}
			for _, v := range dingKeys {
				fmt.Println("dingKey为:" + v + "\n")
			}
			for _, v := range dingusers {
				fmt.Println("dingusers为:" + v + "\n")
			}
			for k, v := range userMap {
				fmt.Println("userMapKey为:" + k + " " + "userMapValue为" + v + "\n")
			}
			fmt.Println("isAtAllTmp为:" + string(isAtallTmp))
			//uniq appid config
			if dyAgoCfg.AppConfig.AppConfigMap != nil {
				if _, ok := dyAgoCfg.AppConfig.AppConfigMap[appid]; ok {
					dingKeys, dingusers, userMap, isAtallTmp = InitAppConfigMap(dyAgoCfg.AppConfig.AppConfigMap, appid,isAtallTmp)
				}
			}
		}
	}
	for _, v := range dingKeys {
		fmt.Println("dingKey为:" + v + "\n")
	}
	for _, v := range dingusers {
		fmt.Println("dingusers为:" + v + "\n")
	}
	for k, v := range userMap {
		fmt.Println("userMapKey为:" + k + " " + "userMapValue为" + v + "\n")
	}
	fmt.Println("isAtAllTmp为:" + string(isAtallTmp))
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
func InitAppConfigMap(appConfigMap map[string]ConfigInfo, appid string, isAtAllTmp int) (dingKeys []string, dingUsers []string, userMap map[string]string, isAtAll int) {

	if len(appConfigMap[appid].DingKeys) > 0 {
		dingKeys = appConfigMap[appid].DingKeys
	}
	if len(appConfigMap[appid].DingUsers) > 0 {
		dingUsers = appConfigMap[appid].DingUsers
	}
	for key, value := range appConfigMap[appid].DingUserMap {
		userMap[key] = value
	}
	if appConfigMap[appid].IsAtAll != 0 {
		isAtAllTmp = appConfigMap[appid].IsAtAll
	}
	return
}
func InitDyAppConfigMap(dyAppConfigMap *AppCfg, appid string,isAtAllTmp int) (dingKeys []string, dingUsers []string, userMap map[string]string, isAtAll int) {
	if len(dyAppConfigMap.AppConfigMap[appid].DingKeys) > 0 {
		dingKeys = dyAppConfigMap.AppConfigMap[appid].DingKeys
	}
	if len(dyAppConfigMap.AppConfigMap[appid].DingUsers) > 0 {
		dingUsers = dyAppConfigMap.AppConfigMap[appid].DingUsers
	}
	for key, value := range dyAppConfigMap.AppConfigMap[appid].DingUserMap {
		userMap[key] = value
	}
	if dyAppConfigMap.AppConfigMap[appid].IsAtAll != 0 {
		isAtAllTmp = dyAppConfigMap.AppConfigMap[appid].IsAtAll
	}
	return
}
