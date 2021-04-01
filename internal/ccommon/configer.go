package ccommon

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"path/filepath"
)

var CConfiger *ConfigCenterInfo

const (
	ServerName    = "mvbjqa"
	SubServerName = "configCenter"
)

const (
	DirFlag = "configs"
)

const (
	CenterConfig = "cluster.toml"
	LogConfig  = "log.toml"
)

type BaseConf struct {
	LogCfg  *LogCfg
	CenterCfg *ConfigCenterInfo
}

type ConfigCenterInfo struct {
    ConfigServerUrl string `toml:"ip_port"`
    AppClusterMap map[string][]string `toml:"app_cluster_map"`
    ClusterMap map[string]*ClusterInfo `toml:"cluster_map"`
}

type ClusterInfo struct {
    ClusterDetail map[string]string
}

func ParseBaseConfig(configDir string) (*BaseConf, error) {
	cfg := &BaseConf{}
	centerCfg, err := parseConfigCenterConf(filepath.Join(configDir, CenterConfig))
	if err != nil {
		return nil, fmt.Errorf("parse logConfig error, err[%s]", err.Error())
	}

	cfg.CenterCfg = centerCfg

	logCfg, err := parseLogConfig(filepath.Join(configDir, LogConfig))
	if err != nil {
		return nil, fmt.Errorf("parse logConfig error, err[%s]", err.Error())
	}
	cfg.LogCfg = logCfg

	return cfg, nil
}

func parseLogConfig(fileName string) (*LogCfg, error) {
	logCfg := &LogCfg{}
	if err := parseTomlConfig(fileName, logCfg); err != nil {
		return logCfg, err
	}
	return logCfg, nil
}

func parseConfigCenterConf(fileName string) (*ConfigCenterInfo, error) {
	logCfg := &ConfigCenterInfo{}
	if err := parseTomlConfig(fileName, logCfg); err != nil {
		return logCfg, err
	}
	return logCfg, nil
}

func parseTomlConfig(fileName string, config interface{}) (err error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("readFile[%s], %s", fileName, err.Error())
	}

	if _, err = toml.Decode(string(data), config); err != nil {
		return fmt.Errorf("decodeFile[%s], %s", fileName, err.Error())
	}

	return nil
}
