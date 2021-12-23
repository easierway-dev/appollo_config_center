package chttp

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "net/http"
)

func httpGet(url,token string){
	var resp_json interface{}
    client := &http.Client{}
    req,_ := http.NewRequest("GET",url,nil)
    req.Header.Add("Authorization",token)
    resp,_ := client.Do(req)
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return resp_json, err
    }
    err := json.Unmarshal(body, &resp_json)
    if err != nil {
    	return resp_json, err
    }
    return resp_json, nil
}

func httpPostForm(url string, token, data map[string]interface{}) {
	var resp_json interface{}
    client := &http.Client{}
    bytesData, _ := json.Marshal(data)
    req, _ := http.NewRequest("POST",url,bytes.NewReader(bytesData))
    req.Header.Add("Authorization",token)
    resp, _ := client.Do(req)
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return resp_json, err
    }
    err := json.Unmarshal(body, &resp_json)
    if err != nil {
    	return resp_json, err
    }
    return resp_json, nil
}
