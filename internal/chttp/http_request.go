package chttp

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
)

func HttpGet(url,token string) (resp_body interface{}, err error) {
    client := &http.Client{}
    req,_ := http.NewRequest("GET",url,nil)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Authorization",token)
    resp,_ := client.Do(req)
    resp_body, err := ioutil.ReadAll(resp.Body)
    return 
}

func HttpPostForm(url, token string, data map[string]interface{})(resp_body interface{}, err error) {
    client := &http.Client{}
    bytesData, _ := json.Marshal(data)
    req, _ := http.NewRequest("POST",url,bytes.NewReader(bytesData))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Authorization",token)
    resp, _ := client.Do(req)
    resp_body, err := ioutil.ReadAll(resp.Body)
    return 
}
