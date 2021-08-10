package ccommon

import (
	"fmt"
	"errors"
	
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/cnotify"
	"gitlab.mobvista.com/voyager/zlog"
)

var CLogger *ccLogger

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

func (this *ccLogger) Info(args ...interface{}) {
	if this == nil || this.Runtime == nil {
		return
	}
	if cnotify.DyDingKey != "" {
		cnotify.SendText(cnotify.DyDingKey,fmt.Sprintf("%s",args),DdingConfiger.DingUsers)
	} else if DdingConfiger.DingKey != "" {
		cnotify.SendText(DdingConfiger.DingKey,fmt.Sprintf("%s",args),DdingConfiger.DingUsers)
	}
	this.Runtime.Info(args)
}

func (this *ccLogger) Warn(args ...interface{}) {
	if this == nil || this.Runtime == nil {
		return
	}
	this.Runtime.Warn(args)
}

func (this *ccLogger) Error(args ...interface{}) {
	if this == nil || this.Runtime == nil {
		return
	}
	if cnotify.DyDingKey != "" {
		cnotify.SendText(cnotify.DyDingKey,fmt.Sprintf("%s",args),DdingConfiger.DingUsers)
	} else if DdingConfiger.DingKey != "" {
		cnotify.SendText(DdingConfiger.DingKey,fmt.Sprintf("%s",args),DdingConfiger.DingUsers)
	}
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
