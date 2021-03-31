package cconsul

import (
	"github.com/hashicorp/consul/api"
	"os"
)

func writeOne(addr, path, value string) error {
	client, err := newClient(addr)
	if err != nil {
		println(err)
		os.Exit(-1)
	}
	return writeData(client, path, value)
}

func newClient(addr string) (*api.Client, error) {
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

func writeData(client *api.Client, path, value string) error {
	kv := client.KV()
	_, err := kv.Put(&api.KVPair{
		Key:   path,
		Value: []byte(value),
	}, nil)
	return err
}
