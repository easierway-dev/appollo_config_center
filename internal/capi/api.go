package capi

import (
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
    IsPublic                    string     `toml:"isPublic"`
    DataChangeCreatedBy         string     `toml:"dataChangeCreatedBy"`
    DataChangeLastModifiedBy    string     `toml:"dataChangeLastModifiedBy"`
    DataChangeCreatedTime       string     `toml:"dataChangeCreatedTime"`
    DataChangeLastModifiedTime  string     `toml:"dataChangeLastModifiedTime"`
    Items                       []*ItemInfo  `toml:"items"`
}

func GetNamespaceInfo(url,token string) (resp_body *capi.NamespaceInfo, err error) {
    body, err := chttp.HttpGet(url,token)
    if err == nil {
        err = json.Unmarshal(body, &resp_body)
        if err != nil {
            return nil, err
        }
    }
    return
}