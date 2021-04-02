package ccommon

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var AgolloConfiger *AgolloConfig

const (
	ServerName    = "mvbjqa"
	SubServerName = "configCenter"
)

const (
	DirFlag = "configs"
)

const (
	AgolloConfig = "agollo.toml"
	LogConfig    = "log.toml"
)

type BaseConf struct {
	LogCfg    *LogCfg
	AgolloCfg *AgolloConfig
}

type DyAgolloConf struct {
	AppClusterConfig *AppClusterConfig
	ClusterConfig *ClusterConfig
	AppConfig *AppConfig
} 

type AgolloConfig struct {
	ConfigServerURL string                 `toml:"ipport"`
	AppID string                 		`toml:"appid"`
	Cluster string                 		`toml:"cluster"`
	Namespace []string                 	`toml:"namespace"`
}

type AppClusterConfig struct {
	AppClusterMap   map[string][]string    `toml:"app_cluster_map"`
}

type ClusterConfig struct {
	ClusterMap      map[string]ClusterInfo `toml:"cluster_map"`
}

type AppConfig struct {
	AppConfigMap      map[string]ConfigInfo `toml:"app_config_map"`
}

type ClusterInfo struct {
	ConsulAddr string `toml:"consul_addr"`
}

type ConfigInfo struct {
	DingKey string `toml:"ding_key"`
}

func ParseBaseConfig(configDir string) (*BaseConf, error) {
	cfg := &BaseConf{}
	agolloCfg, err := parseAgolloConfig(filepath.Join(configDir, AgolloConfig))
	if err != nil {
		return nil, fmt.Errorf("parse logConfig error, err[%s]", err.Error())
	}

	cfg.AgolloCfg = agolloCfg

	logCfg, err := parseLogConfig(filepath.Join(configDir, LogConfig))
	if err != nil {
		return nil, fmt.Errorf("parse logConfig error, err[%s]", err.Error())
	}
	cfg.LogCfg = logCfg

	return cfg, nil
}



func parseLogConfig(fileName string) (*LogCfg, error) {
	cfg := &LogCfg{}
	if err := parseTomlConfig(fileName, cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func parseAgolloConfig(fileName string) (*AgolloConfig, error) {
	cfg := &AgolloConfig{}
	if err := parseTomlConfig(fileName, cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func parseAppClusterConfig(data string) (*AppClusterConfig, error) {
        cfg := &AppClusterConfig{}
        if err := parseTomlStringConfig(data, cfg); err != nil {
                return cfg, err
        }
        return cfg, nil
}

func parseClusterConfig(data string) (*ClusterConfig, error) {
        cfg := &ClusterConfig{}
        if err := parseTomlStringConfig(data, cfg); err != nil {
                return cfg, err
        }
        return cfg, nil
}

func parseAppConfig(data string) (*AppConfig, error) {
        cfg := &AppConfig{}
        if err := parseTomlStringConfig(data, cfg); err != nil {
                return cfg, err
        }
        return cfg, nil
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

func parseTomlStringConfig(tomlData string, config interface{}) (err error) {

        if _, err = toml.Decode(string(tomlData), config); err != nil {
                return fmt.Errorf("decode %s, %s", tomlData, err.Error())
        }

        return nil
}
