package ccompare

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

var AgolloConfiger *AgolloCfg
var DyAgolloConfiger map[string]*DyAgolloCfg

var AppConfiger *AppCfg
var Configer *ConfigInfo
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

func ParseClusterConfig(data string) (*ClusterCfg, error) {
	cfg := &ClusterCfg{}
	if err := parseTomlStringConfig(data, cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func ParseAppConfig(data string) (*AppCfg, error) {
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
func HttpGet(url, token string) ([]byte, error) {

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0")
	req.Header.Set("Authorization", token)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	//'Accept-Encoding':'gzip, deflate, br'
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// 是否有 gzip
	gzipFlag := false
	for k, v := range resp.Header {
		if strings.ToLower(k) == "content-encoding" && strings.ToLower(v[0]) == "gzip" {
			gzipFlag = true
		}
	}
	defer func() { _ = resp.Body.Close() }()
	if resp != nil {
		if gzipFlag {
			// 创建 gzip.Reader
			gr, err := gzip.NewReader(resp.Body)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			defer gr.Close()
			respBody, err := ioutil.ReadAll(gr)
			return respBody, err
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		return respBody, err
	}
	return nil, errors.New("no response")
}

func HttpPostForm(url, token string, data map[string]interface{}) (resp_body string, err error) {
	client := &http.Client{}
	bytesData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", token)
	resp, _ := client.Do(req)
	if err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()
	if resp != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		resp_body = string(body)
	}
	return
}
func NewClient(addr string) (*api.Client, error) {
	conf := api.DefaultConfig()
	if addr != "" {
		conf.Address = addr
	}
	client, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}
	return client, nil
}