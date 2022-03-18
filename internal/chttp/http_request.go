package chttp

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
)

func HttpGet(url,token string) (respBody []byte, err error) {
    client := &http.Client{}
    req,_ := http.NewRequest("GET",url,nil)
    //req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0")
    req.Header.Set("Authorization",token)
    req.Header.Set("Accept-Encoding", "gzip, deflate, br")
    //'Accept-Encoding':'gzip, deflate, br'
    resp,err := client.Do(req)
    if err != nil {
        return 
    }
    defer func() { _ = resp.Body.Close() }()
    if resp != nil {
        //result := make(map[string]interface{})
        body, _ := ioutil.ReadAll(resp.Body)
        //_ = json.Unmarshal(body, &result)
        respBody = body
    }
    return
}

func HttpPostForm(url, token string, data map[string]interface{})(resp_body string, err error) {
    client := &http.Client{}
    bytesData, _ := json.Marshal(data)
    req, _ := http.NewRequest("POST",url,bytes.NewReader(bytesData))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Authorization",token)
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