package capi

import (
    "fmt"
    "encoding/json"
    "gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/chttp"
)

type ItemInfo struct {
    Key                         string     `toml:"key"`
    Value                       string     `toml:"value"`
    DataChangeCreatedBy         string     `toml:"dataChangeCreatedBy"`
    DataChangeLastModifiedBy    string     `toml:"dataChangeLastModifiedBy"`
    DataChangeCreatedTime       string     `toml:"dataChangeCreatedTime"`
    DataChangeLastModifiedTime  string     `toml:"dataChangeLastModifiedTime"`
}

type NamespaceInfo struct {
    AppId                       string     `toml:"appId"`
    ClusterName                 string     `toml:"clusterName"`
    NamespaceName               string     `toml:"namespaceName"`
    Comment                     string     `toml:"comment"`
    Format                      string     `toml:"format"`
    IsPublic                    bool     `toml:"isPublic"`
    DataChangeCreatedBy         string     `toml:"dataChangeCreatedBy"`
    DataChangeLastModifiedBy    string     `toml:"dataChangeLastModifiedBy"`
    DataChangeCreatedTime       string     `toml:"dataChangeCreatedTime"`
    DataChangeLastModifiedTime  string     `toml:"dataChangeLastModifiedTime"`
    Items                       []ItemInfo  `toml:"items"`
}

func GetNamespaceInfo(url,token string) (resp_body *NamespaceInfo, err error) {
    body, err := chttp.HttpGet(url,token)
    if err == nil {
        err = json.Unmarshal([]byte(body), &resp_body)
        //fmt.Println("\nxxdebugresp_body=",resp_body, "\nxxdebugerr=",err)
        if err != nil {
            fmt.Println("Unmarshal NamespaceInfo err=", err)
            return nil, err
        }
    }
    return
}