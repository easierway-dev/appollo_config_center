package cconsul

import (
	"os"

	"github.com/hashicorp/consul/api"
)

func WriteOne(addr, path, value , mode string) error {
	client, pair, err := NewClient(addr, path, value)
	if err != nil {
		println(err)
		os.Exit(-1)
	}
	switch mode {
	case "get":
	    return GetData(client, pair)
	case "del":
	    return DeleteData(client, pair)
	default:
	    return WriteData(client, pair)
	}	
}

func NewClient(addr, path, value string) (*api.Client, *api.KVPair, error) {
	conf := api.DefaultConfig()
	if addr != "" {
		conf.Address = addr
	}
	client, err := api.NewClient(conf)
	if err != nil {
		return nil, nil, err
	}
	//初始化一个kv
	pair := &api.KVPair{
		Key:   path,
		Value: []byte(value),
	}
	return client, pair, nil
}

func WriteData(client *api.Client, pair *api.KVPair) error {
	kv := client.KV()
	_, err := kv.Put(pair, nil)
	return err
}

func GetData(client *api.Client, pair *api.KVPair) error {
	kv := client.KV()
	_,_, err := kv.Get(pair.Key, nil)
	return err
}

func DeleteData(client *api.Client, pair *api.KVPair) error {
	kv := client.KV()
	_, err := kv.Delete(pair.Key, nil)
	return err
}