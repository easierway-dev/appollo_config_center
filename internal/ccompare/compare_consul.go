package ccompare

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/capi"
)

type Value interface {
	CompareValue()
}

type KeyInfo interface {
	GetInfo(key string) map[string]string
}
type Key struct{}
// 对比consul_kv之后,记录每个集群名和namesapce下不存在和不相等的key
type KValue struct {
	Cluster   string
	NameSpace map[string]*CompareKey
}
//
type CompareKey struct {
	NotExistKey map[string]*capi.ItemInfo
	NotEqualKey map[string]*capi.ItemInfo
}

var apolloInfo map[string][]*KValue

type ApolloValue struct {
}
type ConsulValue struct {
}

func (apolloValue *ApolloValue) CompareValue() {
	fmt.Println("init consul client")
	consulValue := &ConsulValue{}
	apolloInfo = make(map[string][]*KValue)
	// 这里地址写死了,可以动态获取apollo的值
	client, _ := NewClient(ADDR)
	// 每个业务线
	for appId, appIdProperty := range appIdsProperty {
		// 暂时跳过dsp_abtest,bidforce
		if appId == "dsp_abtest" || appId == "bidforce" {
			continue
		}
		var kValues []*KValue
		// 每个集群下的nameSpace
		fmt.Println("appid ==", appId, "appIdProperty.NameSpace = ", appIdProperty.NameSpace)
		for clusterName, namespace := range appIdProperty.NameSpace {
			kValue := &KValue{}
			for i := 0; i < len(namespace); i++ {
				//
				// namespace为空的时候，继续下一次循环
				if len(namespace[i].Items) == 0 {
					fmt.Println("namespace is nil:  ", "AppId", appId, "\tclusterName", clusterName, "\tnamespace:", namespace[i].NamespaceName)
					continue
				}
				kv := make(map[string]*capi.ItemInfo)
				// 将单个namespace赋值到map中
				for j := 0; j < len(namespace[i].Items); j++ {
					kv[namespace[i].Items[j].Key] = getItemInfo(namespace[i].Items[j])
				}
				if _, ok := kv["consul_key"]; ok {
					fmt.Println("content:", "key contain consul_key  ", "AppId", appId, "\tclusterName", clusterName, "\tnamespace:", namespace[i].NamespaceName)
					continue
				}
				comkey := &CompareKey{}
				comkey.NotExistKey = make(map[string]*capi.ItemInfo)
				comkey.NotEqualKey = make(map[string]*capi.ItemInfo)
				for k, v := range kv {
					consulKValue, err := consulValue.GetValue(client, k)
					if err != nil || consulKValue == nil {
						// 对比之后不存在值
						comkey.NotExistKey[k] = v
						continue
					}
					if string(consulKValue.Value) == v.Value {
						continue
					}
					// 对比之后不相等值
					comkey.NotEqualKey[k] = v
				}
				kValue.NameSpace = make(map[string]*CompareKey)
				kValue.NameSpace[namespace[i].NamespaceName] = comkey
			}
			kValue.Cluster = clusterName
			kValues = append(kValues, kValue)
		}
		apolloInfo[appId] = kValues
	}
}
func getItemInfo(item capi.ItemInfo) *capi.ItemInfo {
	return &capi.ItemInfo{Value: item.Value, DataChangeLastModifiedBy: item.DataChangeLastModifiedBy}
}
func (consulValue *ConsulValue) GetValue(client *api.Client, path string) (*api.KVPair, error) {
	kv := client.KV()
	KVPair, _, err := kv.Get(path, nil)
	if err != nil {
		return nil, err
	}
	return KVPair, nil
}
func (consulValue *ConsulValue) CompareValue() {
}

// 获取某个key隶属那个集群和namespace
func (k *Key) GetInfo(key string) map[string]string {
	if key == "" {
		return nil
	}
	info := make(map[string]string)
	// 每个业务线
	for _, appIdProperty := range appIdsProperty {
		// 每个集群下的nameSpace
		for clusterName, namespace := range appIdProperty.NameSpace {
			for i := 0; i < len(namespace); i++ {
				for j := 0; j < len(namespace[i].Items); j++ {
					if namespace[i].Items[j].Key == key {
						info[clusterName] = namespace[i].NamespaceName
					}
				}
			}
		}
	}
	return info
}
