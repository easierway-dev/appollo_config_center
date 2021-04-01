package cconsul

import (
	"os"

	"github.com/hashicorp/consul/api"
)

func WriteOne(addr, path, value string) error {
	client, err := NewClient(addr)
	if err != nil {
		println(err)
		os.Exit(-1)
	}
	return WriteData(client, path, value)
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

func WriteData(client *api.Client, path, value string) error {
	kv := client.KV()
	_, err := kv.Put(&api.KVPair{
		Key:   path,
		Value: []byte(value),
	}, nil)
	return err
}
