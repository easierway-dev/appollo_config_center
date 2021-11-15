package ccommon

import (
	"fmt"
	"errors"
	
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/cnotify"
	"gitlab.mobvista.com/voyager/zlog"
)

var CLogger *ccLogger

const (
	DefaultDingType = ""
	DefaultDingType = ""
	InitDingType = "init"
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

func GetDingInfo(appid string, itype string) (dingKeys []string,dingusers []string) {
	if appid == "" && itype == "info"{
		return dingKeys, dingusers
	}

        namespace := DefaultNamespace
	dingKeys = AppConfiger.DingKeys
	dingusers = AppConfiger.DingUsers
	if DyAgolloConfiger != nil {
		if dyAgoCfg,ok := DyAgolloConfiger[namespace];ok {
			if dyAgoCfg.AppConfig != nil {
				if len(dyAgoCfg.AppConfig.DingKeys) > 0 {
					dingKeys = dyAgoCfg.AppConfig.DingKeys
				}
				if len(dyAgoCfg.AppConfig.DingUsers) > 0 {
					dingusers = dyAgoCfg.AppConfig.DingUsers
				}
			}
			if dyAgoCfg.AppConfig.AppConfigMap != nil {
				if _,ok := dyAgoCfg.AppConfig.AppConfigMap[appid];ok {
					if len(dyAgoCfg.AppConfig.AppConfigMap[appid].DingKeys) > 0 {
						dingKeys = dyAgoCfg.AppConfig.AppConfigMap[appid].DingKeys
					}
					if len(dyAgoCfg.AppConfig.AppConfigMap[appid].DingUsers) > 0 {
						dingusers = dyAgoCfg.AppConfig.AppConfigMap[appid].DingUsers
					}
				} 
			}
		}
        }
	return dingKeys, dingusers
}

func (this *ccLogger) Info(args ...interface{}) {
	if this == nil || this.Runtime == nil {
		return
	}
	dingkeys,dingusers := GetDingInfo(args[0].(string), "info")
	cnotify.SendText(dingkeys,fmt.Sprintf("%s",args),dingusers)
	this.Runtime.Info(args)
}

func (this *ccLogger) Warn(args ...interface{}) {
	if this == nil || this.Runtime == nil {
		return
	}
	dingkeys,dingusers := GetDingInfo(args[0].(string), "warn")
	cnotify.SendText(dingkeys,fmt.Sprintf("%s",args),dingusers)
	this.Runtime.Warn(args)
}

func (this *ccLogger) Error(args ...interface{}) {
	if this == nil || this.Runtime == nil {
		return
	}
	dingkeys,dingusers := GetDingInfo(args[0].(string), "err")
	cnotify.SendText(dingkeys,fmt.Sprintf("%s",args),dingusers)
	this.Runtime.Error(args)
}


func (this *ccLogger) Infof(format string, args ...interface{}) {
        if this == nil || this.Runtime == nil {
                return
        }
        this.Runtime.Infof(format,args)
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
