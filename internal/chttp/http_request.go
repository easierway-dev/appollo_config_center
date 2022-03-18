package chttp

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func HttpGet(url, token string) ([]byte, error) {

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0")
	req.Header.Set("Authorization", token)
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	//'Accept-Encoding':'gzip, deflate, br'
	resp, err := client.Do(req)
	if err != nil {
		return nil,err
	}
	// 是否有 gzip
	gzipFlag := false
	for k, v := range resp.Header {
		if strings.ToLower(k) == "content-encoding" && strings.ToLower(v[0]) == "gzip" {
			gzipFlag = true
		}
	}
	defer func() { _ = resp.Body.Close() }()
	if resp != nil {
		if gzipFlag {
			// 创建 gzip.Reader
			gr, err := gzip.NewReader(resp.Body)
			if err != nil {
				fmt.Println(err.Error())
				return nil,err
			}
			defer gr.Close()
			respBody, err := ioutil.ReadAll(gr)
			return respBody,err
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		return respBody,err
	}
	return nil,errors.New("no response")
}

func HttpPostForm(url, token string, data map[string]interface{}) (resp_body string, err error) {
	client := &http.Client{}
	bytesData, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", token)
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
