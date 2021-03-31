package ccommon

import (
	"gitlab.mobvista.com/voyager/zlog"
	"errors"
)

var CLogger *ConfigCenterLogger

type LogCfg struct {
	Runtime      *zlog.Ops `toml:"runtime_log"`
}

type ConfigCenterLogger struct {
	Runtime      zlog.Logger
}

func NewconfigCenterLogger(logCfg *LogCfg) (*ConfigCenterLogger, error) {
	var err error
	logger := &ConfigCenterLogger{}
	if logCfg == nil {
		return nil,errors.New("logCfg is nil")
	}
	if logger.Runtime, err = zlog.NewZLog(logCfg.Runtime); err != nil {
		return nil, err
	}
	return logger, nil
}
