package chttp

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
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
    Items                       *ItemInfo  `toml:"items"`
}

func HttpGet(url,token string) (resp_body *NamespaceInfo, err error) {
    client := &http.Client{}
    req,_ := http.NewRequest("GET",url,nil)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Authorization",token)
    resp,_ := client.Do(req)
    body, err := ioutil.ReadAll(resp.Body)
    if err == nil {
        err = json.Unmarshal(body, &resp_body{})
        if err != nil {
            return nil, err
        }
    }
    return
}

func HttpPostForm(url, token string, data map[string]interface{})(resp_body *NamespaceInfo, err error) {
    client := &http.Client{}
    bytesData, _ := json.Marshal(data)
    req, _ := http.NewRequest("POST",url,bytes.NewReader(bytesData))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Authorization",token)
    resp, _ := client.Do(req)
    body, err := ioutil.ReadAll(resp.Body)
    if err == nil {
        err = json.Unmarshal(body, &resp_body{})
        if err != nil {
            return nil, err
        }
    }
    return
}
