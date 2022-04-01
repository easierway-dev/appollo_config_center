package ccompare

import (
	"github.com/hashicorp/consul/api"
)

const ADDR = "47.252.4.203:8500"

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
