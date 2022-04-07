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
	ConsulAddr  string
	NotExistKey map[string]*ItemInfo
	NotEqualKey map[string]*ItemInfo
}

type ApolloValue struct {
	ApolloInfo map[string][]*KValue
}
type ConsulValue struct {
	ConsulInfo map[string]*ApolloValue
}

func (apolloValue *ApolloValue) CompareValue() {
	fmt.Println("start compare")
	consulValue := &ConsulValue{}
	client := make(map[string]*api.Client)
	apolloValue.ApolloInfo = make(map[string][]*KValue)
	consulValue.ConsulInfo = make(map[string]*ApolloValue)
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
			for i := 0; i < len(namespace); i++ {
				// namespace为空的时候，继续下一次循环
				if len(namespace[i].Items) == 0 {
					fmt.Println("namespace is nil AppId", appId, "\tclusterName", clusterName, "\tnamespace:", namespace[i].NamespaceName)
					continue
				}
				if consulAddr, ok := GlobalConfiger.ClusterMap[clusterName]; !ok {
					fmt.Println("content:", "cluster not correspond consulAddr AppId:", appId, "\tclusterName:", clusterName)
					continue
				} else {
					for i := 0; i < len(consulAddr.ConsulAddr); i++ {
						if _, ok := client[consulAddr.ConsulAddr[i]]; ok {
							continue
						}
						fmt.Println("consulAddr:", consulAddr.ConsulAddr[i])
						cli, _ := NewClient(consulAddr.ConsulAddr[i])
						if cli == nil {
							fmt.Println("consulAddr:", consulAddr.ConsulAddr[i]+" connect failed")
							continue
						}
						// 每个集群对应一个client
						client[consulAddr.ConsulAddr[i]] = cli
					}
				}
				kv := make(map[string]*ItemInfo)
				// 将单个namespace中的key,value赋值到map中
				for j := 0; j < len(namespace[i].Items); j++ {
					kv[namespace[i].Items[j].Key] = getItemInfo(namespace[i].Items[j])
				}
				if _, ok := kv["consul_key"]; ok {
					fmt.Println("content:", "key contain consul_key AppId", appId, "\tclusterName", clusterName, "\tnamespace:", namespace[i].NamespaceName)
					continue
				}
				// 某个集群下consulAddr可能有多个
				for addr, cli := range client {
					comkey := &CompareKey{}
					comkey.NotExistKey = make(map[string]*ItemInfo)
					comkey.NotEqualKey = make(map[string]*ItemInfo)
					for k, v := range kv {
						consulKValue, err := consulValue.GetValue(cli, k)
						if consulKValue == nil || err != nil {
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
					//comkey.ConsulAddr = addr
					kValue.NameSpace = make(map[string]*CompareKey)
					kValue.NameSpace[namespace[i].NamespaceName] = comkey
					kValue.Cluster = clusterName
					kValues = append(kValues, kValue)
					// 每个业务线对应的具体信息
					apolloValue.ApolloInfo[appId] = kValues
					consulValue.ConsulInfo[addr] = apolloValue
				}
			}
		}
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
	fmt.Println("start apolloValue.....")
	for _, value := range appId {
		if value == nil {
			printAll(apolloValue.ApolloInfo)
			break
		} else {
			printAppId(apolloValue.ApolloInfo, value)
		}
	}
}
func (consulValue *ConsulValue) Print(appId ...interface{}) {
	fmt.Println("start consulValue.....")
	for _, value := range appId {
		if value == nil {
			fmt.Println("consulValue.ConsulInfo = ", consulValue.ConsulInfo)
		}
	}
}
func printAll(apolloKV map[string][]*KValue) {
	for appid, value := range apolloKV {
		for _, val := range value {
			for namespace, keys := range val.NameSpace {
				for i := 0; i < len(GlobalConfiger.ClusterMap[val.Cluster].ConsulAddr); i++ {
					fmt.Println("appid =", appid, "\tapolloInfo Cluster =", val.Cluster, "\tapolloInfo NameSpace =", namespace,
						"\tconsulAddr: ", GlobalConfiger.ClusterMap[val.Cluster].ConsulAddr[i], "  apolloInfo notExistKey =", keys.NotExistKey)
					fmt.Println("appid =", appid, "\tapolloInfo Cluster =", val.Cluster, "\tapolloInfo NameSpace =", namespace,
						"\tconsulAddr: ", GlobalConfiger.ClusterMap[val.Cluster].ConsulAddr[i], "  apolloInfo NotEqualKey =", keys.NotEqualKey)
				}
			}
		}
	}
}
func printAppId(apolloKV map[string][]*KValue, appId ...interface{}) {
	fmt.Println("开始")
	for _, id := range appId {
		for _, val := range apolloKV[id.(string)] {
			if val == nil {
				return
			}
			//fmt.Println("apolloInfo kv =", kv)
			fmt.Println("dsp apolloInfo Cluster =", val.Cluster)
			for namespace, keys := range val.NameSpace {
				fmt.Println("dsp apolloInfo NameSpace =", namespace)
				for i := 0; i < len(GlobalConfiger.ClusterMap[val.Cluster].ConsulAddr); i++ {
					fmt.Println("consulAddr: ", GlobalConfiger.ClusterMap[val.Cluster].ConsulAddr[i], " dsp apolloInfo notExistKey =", keys.NotExistKey)
					fmt.Println("consulAddr: ", GlobalConfiger.ClusterMap[val.Cluster].ConsulAddr[i], " dsp apolloInfo NotEqualKey =", keys.NotEqualKey)
				}
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
