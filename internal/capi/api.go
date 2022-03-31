package capi

import (
	"encoding/json"
	"fmt"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/chttp"
)

type ItemInfo struct {
	Key                        string `toml:"key"`
	Value                      string `toml:"value"`
	DataChangeCreatedBy        string `toml:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `toml:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `toml:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `toml:"dataChangeLastModifiedTime"`
}

type NamespaceInfo struct {
	AppId                      string     `toml:"appId"`
	ClusterName                string     `toml:"clusterName"`
	NamespaceName              string     `toml:"namespaceName"`
	Comment                    string     `toml:"comment"`
	Format                     string     `toml:"format"`
	IsPublic                   bool       `toml:"isPublic"`
	DataChangeCreatedBy        string     `toml:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string     `toml:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string     `toml:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string     `toml:"dataChangeLastModifiedTime"`
	Items                      []ItemInfo `toml:"items"`
}
type EnvClustersInfo struct {
	Env      string   `toml:"env"`
	Clusters []string `toml:"clusters"`
}
type AppInfo struct {
	Name                       string `toml:"name"`
	AppId                      string `toml:"appId"`
	OrgId                      string `toml:"orgId"`
	OrgName                    string `toml:"orgName"`
	OwnerName                  string `toml:"ownerName"`
	OwnerEmail                 string `toml:"ownerEmail"`
	DataChangeCreatedBy        string `toml:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `toml:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `toml:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `toml:"dataChangeLastModifiedTime"`
}

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