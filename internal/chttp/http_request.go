package chttp

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
)

func httpGet(url,token string) (resp_body string, err error) {
    client := &http.Client{}
    req,_ := http.NewRequest("GET",url,nil)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Authorization",token)
    resp,_ := client.Do(req)
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "{}", err
    }
    resp_body = string(body)
    return resp_body, nil
}

func httpPostForm(url string, token, data map[string]interface{}) {
    client := &http.Client{}
    bytesData, _ := json.Marshal(data)
    req, _ := http.NewRequest("POST",url,bytes.NewReader(bytesData))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Authorization",token)
    resp, _ := client.Do(req)
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "{}", err
    }
    resp_body = string(body)
    return string(resp_body), nil
}
