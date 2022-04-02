package ccompare

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type Value interface {
	CompareValue()
	Print(id ...interface{})
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

//var apolloInfo map[string][]*KValue

type ApolloValue struct {
	ApolloInfo map[string][]*KValue
}
type ConsulValue struct {
}

func (apolloValue *ApolloValue) CompareValue() {
	fmt.Println("start compare")
	consulValue := &ConsulValue{}
	var client []*api.Client
	apolloValue.ApolloInfo = make(map[string][]*KValue)
	// 这里地址写死了,可以动态获取apollo的值
	// 每个业务线
	for appId, appIdProperty := range Property.AppIdsProperty {
		// 暂时跳过dsp_abtest,bidforce
		if appId == "dsp_abtest" || appId == "bidforce" {
			continue
		}
		var kValues []*KValue
		// 每个集群下的nameSpace
		fmt.Println("appid ==", appId, "appIdProperty.NameSpace = ", appIdProperty.NameSpace)
		for clusterName, namespace := range appIdProperty.NameSpace {
			kValue := &KValue{}
			if consulAddr, ok := GlobalConfiger.ClusterMap[clusterName]; !ok {
				fmt.Println("content:", "cluster not correspond consulAddr AppId:", appId, "\tclusterName:", clusterName)
				continue
			} else {
				for i := 0; i < len(consulAddr.ConsulAddr); i++ {
					fmt.Println("consulAddr:", consulAddr.ConsulAddr[i])
					cli, _ := NewClient(consulAddr.ConsulAddr[i])
					if cli == nil {
						fmt.Println("consulAddr:", consulAddr.ConsulAddr[i]+" connect failed")
						continue
					}
					client = append(client, cli)
				}
			}
			for i := 0; i < len(namespace); i++ {
				//
				// namespace为空的时候，继续下一次循环
				if len(namespace[i].Items) == 0 {
					fmt.Println("namespace is nil AppId", appId, "\tclusterName", clusterName, "\tnamespace:", namespace[i].NamespaceName)
					continue
				}
				kv := make(map[string]*ItemInfo)
				// 将单个namespace赋值到map中
				for j := 0; j < len(namespace[i].Items); j++ {
					kv[namespace[i].Items[j].Key] = getItemInfo(namespace[i].Items[j])
				}
				if _, ok := kv["consul_key"]; ok {
					fmt.Println("content:", "key contain consul_key AppId", appId, "\tclusterName", clusterName, "\tnamespace:", namespace[i].NamespaceName)
					continue
				}
				comkey := &CompareKey{}
				comkey.NotExistKey = make(map[string]*ItemInfo)
				comkey.NotEqualKey = make(map[string]*ItemInfo)
				// 某个集群下consulAddr可能有多个
				for i := 0; i < len(client); i++ {
					for k, v := range kv {
						consulKValue, err := consulValue.GetValue(client[i], k)
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
				}
				kValue.NameSpace = make(map[string]*CompareKey)
				kValue.NameSpace[namespace[i].NamespaceName] = comkey
			}
			kValue.Cluster = clusterName
			kValues = append(kValues, kValue)
		}
		// 每个业务线对应的具体信息
		apolloValue.ApolloInfo[appId] = kValues
	}
}
func getItemInfo(item ItemInfo) *ItemInfo {
	return &ItemInfo{Value: item.Value, DataChangeLastModifiedBy: item.DataChangeLastModifiedBy}
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
func (apolloValue *ApolloValue) Print(appId ...interface{}) {
	fmt.Println("start.....")
	for _, value := range appId {
		if value == nil {
			printAll(apolloValue.ApolloInfo)
			break
		} else {
			printAppId(apolloValue.ApolloInfo, value)
		}
	}
}
func printAll(apolloKV map[string][]*KValue) {
	for _, value := range apolloKV {
		for _, val := range value {
			//fmt.Println("apolloInfo kv =", kv)
			fmt.Println("dsp apolloInfo Cluster =", val.Cluster)
			for namespace, keys := range val.NameSpace {
				fmt.Println("dsp apolloInfo NameSpace =", namespace)
				fmt.Println("dsp apolloInfo notExistKey =", keys.NotExistKey)
				fmt.Println("dsp apolloInfo NotEqualKey =", keys.NotEqualKey)
				//for k, v := range keys.NotExistKey {
				//	fmt.Println("dsp apolloInfo notExistKey =", k)
				//	fmt.Println("dsp apolloInfo DataChangeLastModifiedBy =", v.DataChangeLastModifiedBy)
				//}
				//for k, v := range keys.NotEqualKey {
				//	fmt.Println("dsp apolloInfo NotEqualKey =", k)
				//	fmt.Println("dsp apolloInfo DataChangeLastModifiedBy =", v.DataChangeLastModifiedBy)
				//}
			}
		}
	}
}
func printAppId(apolloKV map[string][]*KValue, appId ...interface{}) {
	fmt.Println("开始")
	for _, id := range appId {
		for _, val := range apolloKV[id.(string)] {
			//fmt.Println("apolloInfo kv =", kv)
			fmt.Println("dsp apolloInfo Cluster =", val.Cluster)
			for namespace, keys := range val.NameSpace {
				fmt.Println("dsp apolloInfo NameSpace =", namespace)
				fmt.Println("dsp apolloInfo notExistKey =", keys.NotExistKey)
				fmt.Println("dsp apolloInfo NotEqualKey =", keys.NotEqualKey)
			}
		}
	}
}

// 获取某个key隶属那个集群和namespace
func (k *Key) GetInfo(key string) map[string]string {
	if key == "" {
		return nil
	}
	info := make(map[string]string)
	// 每个业务线
	for _, appIdProperty := range Property.AppIdsProperty {
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
