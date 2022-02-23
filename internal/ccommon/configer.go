package ccommon

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var AgolloConfiger *AgolloCfg
var DyAgolloConfiger map[string]*DyAgolloCfg

var AppConfiger *AppCfg
var ConfigerInfo *ConfigInfo
var ChklogRamdom float64
var ChklogRate float64

const (
	ServerName    = "mvbjqa"
	SubServerName = "configCenter"
)

const (
	ABTestAppid      = "abtest"
	BidForceAppid    = "bidforce"
	ABTest           = "abtesting"
	BidForceDsp      = "dsp"
	BidForceRtDsp    = "rtdsp"
	BidForcePioneer  = "pioneer"
	DefaultNamespace = "application"
)

const (
	DirFlag = "configs"
)

const (
	AgolloConfig = "agollo.toml"
	LogConfig    = "log.toml"
	AppConfig    = "app.toml"
)

type BaseConf struct {
	LogCfg    *LogCfg
	AgolloCfg *AgolloCfg
	AppCfg    *AppCfg
}

type DyAgolloCfg struct {
	ClusterConfig *ClusterCfg
	AppConfig     *AppCfg
}

type AgolloCfg struct {
	ConfigServerURL string   `toml:"ipport"`
	PortalURL       string   `toml:"portalport"`
	AppID           string   `toml:"appid"`
	Cluster         string   `toml:"cluster"`
	Namespace       []string `toml:"namespace"`
	CyclePeriod     int      `toml:"cycleperiod"`
}

type AppClusterCfg struct {
	Namespace     []string                  `toml:"namespace"`
	AppClusterMap map[string]AppClusterInfo `toml:"app_cluster_map"`
}

type ClusterCfg struct {
	ClusterMap map[string]ClusterInfo `toml:"cluster_map"`
}

type AppCfg struct {
	DingKeys       []string              `toml:"ding_keys"`
	DingUsers      []string              `toml:"ding_users"`
	DingUserMap    map[string]string     `toml:"ding_user_map"`
	IsAtAll        int                   `toml:"is_at_all"`
	EnUpdateConsul int                   `toml:"enable_update_consul"`
	EnDelConsul    int                   `toml:"enable_delete_consul"`
	ChklogRate     float64               `toml:"log_rate"`
	AppConfigMap   map[string]ConfigInfo `toml:"app_config_map"`
}

type AppClusterInfo struct {
	Cluster   []string `toml:"cluster"`
	Namespace []string `toml:"namespace"`
}

type ClusterInfo struct {
	ConsulAddr []string `toml:"consul_addr"`
}

type ConfigInfo struct {
	DingKeys       []string          `toml:"ding_keys"`            //ding token
	DingUsers      []string          `toml:"ding_users"`           //default ding @list
	DingUserMap    map[string]string `toml:"ding_user_map"`        //config real editor ding @list
	IsAtAll        int               `toml:"is_at_all"`            //1: atall 2:not atall
	EnUpdateConsul int               `toml:"enable_update_consul"` //1: enable update consul 2:not
	EnDelConsul    int               `toml:"enable_delete_consul"` //1: enable delete consul 2:not
	ChklogRate     float64           `toml:"log_rate"`
	AccessToken    string            `toml:"access_token"` //apollo api auth token
}

func ParseBaseConfig(configDir string) (*BaseConf, error) {
	cfg := &BaseConf{}
	agolloCfg, err := ParseAgolloConfig(filepath.Join(configDir, AgolloConfig))
	if err != nil {
		return nil, fmt.Errorf("Parse agoConfig error, err[%s]", err.Error())
	}

	cfg.AgolloCfg = agolloCfg

	logCfg, err := parseLogConfig(filepath.Join(configDir, LogConfig))
	if err != nil {
		return nil, fmt.Errorf("Parse logConfig error, err[%s]", err.Error())
	}
	cfg.LogCfg = logCfg

	appCfg, err := parseBaseAppConfig(filepath.Join(configDir, AppConfig))
	if err != nil {
		return nil, fmt.Errorf("Parse appConfig error, err[%s]", err.Error())
	}
	cfg.AppCfg = appCfg

	return cfg, nil
}

func parseLogConfig(fileName string) (*LogCfg, error) {
	cfg := &LogCfg{}
	if err := parseTomlConfig(fileName, cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func parseBaseAppConfig(fileName string) (*AppCfg, error) {
	cfg := &AppCfg{}
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
func InitAppCfgMap(appConfig *AppCfg, appid ,namespace string) (cfgInfo *ConfigInfo) {
	// 本地配置文件初始化
	ConfigerInfo.InitDyConfigerInfo(appConfig, appid, appConfig.AppConfigMap)
	dyAgoCfg, ok := DyAgolloConfiger[namespace]
	if !ok {
		namespace = DefaultNamespace
		if dyAgoCfg, ok = DyAgolloConfiger[namespace]; !ok {
			return
		}
	}
	// Apollo global_config初始化
	ConfigerInfo.InitDyConfigerInfo(dyAgoCfg.AppConfig,appid,dyAgoCfg.AppConfig.AppConfigMap)
	return ConfigerInfo
}
func(configInfo ConfigInfo) InitDyConfigerInfo(appcfg *AppCfg,appid string,cfg map[string]ConfigInfo){
	if appcfg == nil{
		return
	}
	configInfo.DingKeys = appcfg.DingKeys
	configInfo.DingUsers = appcfg.DingUsers
	configInfo.DingUserMap = appcfg.DingUserMap
	configInfo.IsAtAll = appcfg.IsAtAll
	configInfo.EnUpdateConsul = appcfg.EnUpdateConsul
	configInfo.EnDelConsul = appcfg.EnDelConsul
	configInfo.ChklogRate = appcfg.ChklogRate
	if cfg == nil {
		return
	}
	if _,ok := cfg[appid];!ok{
		return
	}
	if len(cfg[appid].DingKeys) > 0 {
		configInfo.DingKeys = cfg[appid].DingKeys
	}
	if len(cfg[appid].DingUsers) > 0 {
		configInfo.DingUsers = cfg[appid].DingUsers
	}
	for key, value := range cfg[appid].DingUserMap {
		if configInfo.DingUserMap == nil {
			configInfo.DingUserMap = make(map[string]string)
		}
		configInfo.DingUserMap[key] = value
	}
	if cfg[appid].IsAtAll != 0 {
		configInfo.IsAtAll = cfg[appid].IsAtAll
	}
	if cfg[appid].EnUpdateConsul != 0 {
		configInfo.EnUpdateConsul = cfg[appid].EnUpdateConsul
	}
	if cfg[appid].EnDelConsul != 0 {
		configInfo.EnDelConsul = cfg[appid].EnDelConsul
	}
	if cfg[appid].AccessToken != "" {
		configInfo.AccessToken = cfg[appid].AccessToken
	}
}