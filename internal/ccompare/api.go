package ccompare

import (
	"encoding/json"
	"fmt"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/chttp"
)

func GetNamespaceInfo(url, token string) (respBody *NamespaceInfo, err error) {
	body, err := chttp.HttpGet(url, token)
	if err != nil {
		fmt.Println("get body err:", err)
		return nil, err
	}
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		fmt.Println("Unmarshal NamespaceInfo err=", err)
		return nil, err
	}
	return
}
func GetAllNamespaceInfo(url, token string) (respBody []*NamespaceInfo, err error) {
	body, err := chttp.HttpGet(url, token)
	if err != nil {
		fmt.Println("get body err:", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &respBody)
	if len(respBody) == 0 {
		fmt.Println("no namespace under the cluster")
		return
	}
	if err != nil {
		fmt.Println("Unmarshal NamespaceInfo err=", err)
		return nil, err
	}
	return
}
func GetEnvClustersInfo(url, token string) (respBody []*EnvClustersInfo, err error) {
	body, err := chttp.HttpGet(url, token)
	if err != nil {
		fmt.Println("get body err:", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &respBody)
	if err != nil {
		fmt.Println("Unmarshal EnvClustersInfo err=", err)
		return nil, err
	}
	return
}
func GetAppInfo(url, token string) (respBody []*AppInfo, err error) {
	body, err := chttp.HttpGet(url, token)
	if err != nil {
		fmt.Println("get body err:", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(body), &respBody)
	if err != nil {
		fmt.Println("Unmarshal AppInfo err=", err)
		return nil, err
	}
	return
}
