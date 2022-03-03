package ccompare

import (
	"errors"
	"fmt"
	"github.com/hashicorp/consul/api"
	"time"
)

const ADDR = "47.252.4.203:8500"

func GetConsulKV() (map[string]string, error) {
	var consulClient *api.Client
	fmt.Println("init consul client")
	conf := api.DefaultConfig()
	conf.Address = ADDR
	client, err := api.NewClient(conf)
	if err != nil {
		fmt.Println("consul client init failed")
		panic(err)
	}
	consulClient = client
	pair, _, _ := consulClient.KV().List("/", &api.QueryOptions{WaitTime: time.Minute})
	consulKV := make(map[string]string)
	for _, kv := range pair {
		consulKV[kv.Key] = string(kv.Value)
	}
	if len(consulKV) == 0 {
		return nil, errors.New("consul is nil")
	}
	return consulKV, nil
}
