package ccommon

import (
	"errors"

	"gitlab.mobvista.com/voyager/zlog"
)

var CLogger *ccLogger
var CLogCfg *LogCfg

type LogCfg struct {
	Runtime *zlog.Ops `toml:"runtime_log"`
}

type ccLogger struct {
	Runtime zlog.Logger
}

func  (this *ccLogger) NewconfigCenterLogger(logCfg *LogCfg) error {
	var err error
	if logCfg == nil {
		return errors.New("logCfg is nil")
	}
	if this.Runtime, err = zlog.NewZLog(logCfg.Runtime); err != nil {
		return err
	}
	return nil
}

func (this *ccLoger) Infof(format string, args ...interface{}) {
	if this == nil || this.runtime == nil {
		return
	}
	this.runtime.Infof(format, args)
}

func (this *ccLoger) Warnf(format string, args ...interface{}) {
	if this == nil || this.runtime == nil {
		return
	}
	this.runtime.Warnf(format, args)
}

func (this *ccLoger) Errorf(format string, args ...interface{}) {
	if this == nil || this.runtime == nil {
		return
	}
	this.runtime.Errorf(format, args)
}

