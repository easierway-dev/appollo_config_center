package ccompare

import (
	"errors"
	"github.com/hashicorp/consul/api"
)

const ADDR = "47.252.4.203:8500"

func GetConsulKV(client *api.Client,path string) (string, error) {
	kvPair, err := GetData(client, path)
	if err != nil {
		return "",err
	}
	if kvPair == nil{
		return "",errors.New("consul key not value")
	}
	return string(kvPair.Value), nil
}
func NewClient(addr string) (*api.Client, error) {
	conf := api.DefaultConfig()
	if addr != "" {
		conf.Address = addr
	}
	client, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetData(client *api.Client, path string) (*api.KVPair, error) {
	kv := client.KV()
	KVPair, _, err := kv.Get(path, nil)
	if err != nil {
		return nil, err
	}
	return KVPair, nil
}
