package ccompare

import (
	"fmt"
	"github.com/hashicorp/consul/api"
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
	NotExistKey map[string]*ItemInfo
	NotEqualKey map[string]*ItemInfo
}
type ItemInfo struct {
	Key                        string
	Value                      string
	DataChangeCreatedBy        string
	DataChangeLastModifiedBy   string
	DataChangeCreatedTime      string
	DataChangeLastModifiedTime string
}

var apolloInfo map[string][]*KValue

type ApolloValue struct {
	Items map[string][]*KValue
}
type ConsulValue struct {
}

func (apolloValue *ApolloValue) CompareValue() {
	fmt.Println("init consul client")
	consulValue := &ConsulValue{}
	apolloInfo = make(map[string][]*KValue)
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
			// namespace为空的时候，继续下一次循环
			for i := 0; i < len(namespace); i++ {
				//
				if len(namespace[i].Items) == 0 {
					fmt.Println("namespace is nil:  ", "AppId", appId, "\tclusterName", clusterName, "\tnamespace:", namespace[i].NamespaceName)
					continue
				}
				kv := make(map[string]*ItemInfo)
				// 将单个namespace赋值到map中
				for j := 0; j < len(namespace[i].Items); j++ {
					kv[namespace[i].Items[j].Key] = &ItemInfo{Value: namespace[i].Items[j].Value}
				}
				if _, ok := kv["consul_key"]; ok {
					fmt.Println("content:", "key contain consul_key  ", "AppId", appId, "\tclusterName", clusterName, "\tnamespace:", namespace[i].NamespaceName)
				}
				comkey := &CompareKey{}
				for k, v := range kv {
					consulValue1, err := consulValue.GetValue(client, k)
					if err != nil || consulValue1.(string) == "" {
						comkey.NotExistKey[k] = v
						continue
					}
					if consulValue1 == v.Value {
						continue
					}
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

func (consulValue *ConsulValue) GetValue(client *api.Client, path string) (interface{}, error) {
	value, err := GetConsulKV(client, path)
	if err != nil {
		return "", err
	}
	return value, nil
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
