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

func GetDingKey(appid string) string {
        dingKey := ""
        namespace := DefaultNamespace
	if DyAgolloConfiger != nil {
		if dyAgoCfg,ok := DyAgolloConfiger[namespace];ok {
			if dyAgoCfg.AppConfig != nil {
				dingKey = dyAgoCfg.AppConfig.DingKey
			}
			if dyAgoCfg.AppConfig.AppConfigMap != nil {
				if _,ok := dyAgoCfg.AppConfig.AppConfigMap[appid];ok {
					dingKey = dyAgoCfg.AppConfig.AppConfigMap[appid].DingKey
				} 
			}
		}
        }
	if dingKey != "" {
		return dingKey
	} else {
		return DdingConfiger.DingKey
	}
}

func (this *ccLogger) Info(args ...interface{}) {
	if this == nil || this.Runtime == nil {
		return
	}
	dingkey := GetDingKey(args[0].(string))
	cnotify.SendText(dingkey,fmt.Sprintf("%s",args),DdingConfiger.DingUsers)
	this.Runtime.Info(args)
}

func (this *ccLogger) Warn(args ...interface{}) {
	if this == nil || this.Runtime == nil {
		return
	}
	dingkey := GetDingKey(args[0].(string))
	cnotify.SendText(dingkey,fmt.Sprintf("%s",args),DdingConfiger.DingUsers)
	this.Runtime.Warn(args)
}

func (this *ccLogger) Error(args ...interface{}) {
	if this == nil || this.Runtime == nil {
		return
	}
	dingkey := GetDingKey(args[0].(string))
	cnotify.SendText(dingkey,fmt.Sprintf("%s",args),DdingConfiger.DingUsers)
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
