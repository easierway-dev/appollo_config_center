package ccompare

import (
	"encoding/json"
	"fmt"
)

func GetNamespaceInfo(url, token string) (respBody *NamespaceInfo, err error) {
	body, err := HttpGet(url, token)
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
	body, err := HttpGet(url, token)
	fmt.Println("body = ", string(body))
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
func GetEnvClustersInfo(url, token string) (respBody []*EnvClustersInfo, err error) {
	body, err := HttpGet(url, token)
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
	body, err := HttpGet(url, token)
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
