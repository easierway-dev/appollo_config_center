package ccommon

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var AgolloConfiger *AgolloCfg
var DyAgolloConfiger map[string]*DyAgolloCfg

var DyDingKey string

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
	AgolloCfg *AgolloCfg
}


type DyAgolloCfg struct {
	ClusterConfig *ClusterCfg
	AppConfig *AppCfg
} 

type AgolloCfg struct {
	ConfigServerURL string                 `toml:"ipport"`
	AppID string                 		`toml:"appid"`
	Cluster string                 		`toml:"cluster"`
	Namespace []string                 	`toml:"namespace"`
	CyclePeriod int                		`toml:"cycleperiod"`
	DingKey string				`toml:"dingkey"`
}

type AppClusterCfg struct {
	Namespace	[]string `toml:"namespace"`
	AppClusterMap   map[string]AppClusterInfo    `toml:"app_cluster_map"`
}

type ClusterCfg struct {
	ClusterMap      map[string]ClusterInfo `toml:"cluster_map"`
}

type AppCfg struct {
	DingKey       string `toml:"ding_key"`
	AppConfigMap      map[string]ConfigInfo `toml:"app_config_map"`
}

type AppClusterInfo struct {
        Cluster []string `toml:"cluster"`
	Namespace       []string `toml:"namespace"`
}

type ClusterInfo struct {
	ConsulAddr string `toml:"consul_addr"`
}

type ConfigInfo struct {
	DingKey string `toml:"ding_key"`
}

func ParseBaseConfig(configDir string) (*BaseConf, error) {
	cfg := &BaseConf{}
	agolloCfg, err := ParseAgolloConfig(filepath.Join(configDir, AgolloConfig))
	if err != nil {
		return nil, fmt.Errorf("Parse logConfig error, err[%s]", err.Error())
	}

	cfg.AgolloCfg = agolloCfg

	logCfg, err := parseLogConfig(filepath.Join(configDir, LogConfig))
	if err != nil {
		return nil, fmt.Errorf("Parse logConfig error, err[%s]", err.Error())
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

func ParseAgolloConfig(fileName string) (*AgolloCfg, error) {
	cfg := &AgolloCfg{}
	if err := parseTomlConfig(fileName, cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func ParseAppClusterConfig(data string) (*AppClusterCfg, error) {
        cfg := &AppClusterCfg{}
        if err := parseTomlStringConfig(data, cfg); err != nil {
                return cfg, err
        }
        return cfg, nil
}

func ParseDyConfig(clusterConfig, appConfig string) (*DyAgolloCfg, error) {
        cfg := &DyAgolloCfg{}
	if clusterConfig != "" {
		clusterCfg, err := parseClusterConfig(clusterConfig)
		if err == nil {
				cfg.ClusterConfig = clusterCfg
		} else {
			return nil, fmt.Errorf("ParseClusterConfig error, err[%s]", err.Error())
		}
	}
	if appConfig != "" {
		appCfg, err := parseAppConfig(appConfig)
		if err == nil {
				cfg.AppConfig = appCfg
		} else {
			return nil, fmt.Errorf("ParseAppConfig error, err[%s]", err.Error())
		}
	}
        return cfg, nil
}

func parseClusterConfig(data string) (*ClusterCfg, error) {
        cfg := &ClusterCfg{}
        if err := parseTomlStringConfig(data, cfg); err != nil {
                return cfg, err
        }
        return cfg, nil
}

func parseAppConfig(data string) (*AppCfg, error) {
        cfg := &AppCfg{}
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
