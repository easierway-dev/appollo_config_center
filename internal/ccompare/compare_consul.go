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
type CompareKey struct {
	NotExistKey map[string]*ItemInfo
	NotEqualKey map[string]*ItemInfo
}
// apollo具体值
type ApolloValue struct {
	ApolloInfo map[string][]*KValue
}

// consul具体值
type ConsulValue struct {
	ConsulInfo map[string]*ApolloValue
}

var Consul *ConsulValue

func (apolloValue *ApolloValue) CompareValue() {
	fmt.Println("start compare")
	consulValue := &ConsulValue{}
	Consul = &ConsulValue{}
	//client := make(map[string]*api.Client)
	apolloValue.ApolloInfo = make(map[string][]*KValue)
	consulValue.ConsulInfo = make(map[string]*ApolloValue)
	// 获取consul连接
	GetConsulClient()
	// 这里每个consul地址，然后进行每个业务线的遍历，最后每个地址下对应一个apollo配置
	for consulAddress, client := range ConsulClient {
		// 遍历每一个业务线
		for appId, appIdProperty := range Property.AppIdsProperty {
			// 暂时跳过dsp_abtest,bidforce,key中含有consul_key
			if appId == "dsp_abtest" || appId == "bidforce" {
				continue
			}
			var kValues []*KValue
			// 每个集群下的nameSpace
			fmt.Println("appid ==", appId, "appIdProperty.NameSpace = ", appIdProperty.NameSpace)
			// 遍历每一个业务线的集群信息
			for clusterName, namespace := range appIdProperty.NameSpace {
				for i := 0; i < len(namespace); i++ {
					kValue := &KValue{}
					// namespace为空的时候，继续下一次循环
					if len(namespace[i].Items) == 0 {
						fmt.Println("namespace is nil AppId", appId, "\tclusterName", clusterName, "\tnamespace:", namespace[i].NamespaceName)
						continue
					}
					kv := make(map[string]*ItemInfo)
					// 将单个namespace中的key,value赋值到map中，为了判断key中是否包含consul_key
					for j := 0; j < len(namespace[i].Items); j++ {
						kv[namespace[i].Items[j].Key] = getItemInfo(namespace[i].Items[j])
					}
					if _, ok := kv["consul_key"]; ok {
						fmt.Println("content:", "key contain consul_key AppId", appId, "\tclusterName", clusterName, "\tnamespace:", namespace[i].NamespaceName)
						continue
					}
					// 具体集群对应的consul地址
					consulAddr := GlobalConfiger.ClusterMap[clusterName]
					// 当前集群的consul地址检测
					for j := 0; j < len(consulAddr.ConsulAddr); j++ {
						if consulAddr.ConsulAddr[j] != consulAddress {
							continue
						}
						comkey := &CompareKey{}
						comkey.NotExistKey = make(map[string]*ItemInfo)
						comkey.NotEqualKey = make(map[string]*ItemInfo)
						for k, v := range kv {
							consulKValue, err := consulValue.GetValue(client, k)
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
						kValue.Cluster = clusterName
						fmt.Println("nameSpace =", namespace[i].NamespaceName)
						kValue.NameSpace[namespace[i].NamespaceName] = comkey
						kValues = append(kValues, kValue)
					}
				}
			}
			// 每个业务线对应的具体信息
			apolloValue.ApolloInfo[appId] = kValues
		}
		// 每个consul地址对应的的apollo信息
		consulValue.ConsulInfo[consulAddress] = apolloValue
		Consul = consulValue
	}
}

// 获取具体某个key的信息
func getItemInfo(item ItemInfo) *ItemInfo {
	return &ItemInfo{Value: item.Value, DataChangeLastModifiedBy: item.DataChangeLastModifiedBy}
}

// 获取consul的key的值
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

// 当apolloKV == nil，打印所有信息
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

// 打印具体某个业务线信息，例如dsp
func printAppId(apolloKV map[string][]*KValue, appId ...interface{}) {
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
