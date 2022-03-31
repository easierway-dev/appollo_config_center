package ccompare

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

const Application = "application"

type Value interface {
	GetValue(path string) (interface{}, error)
	CompareValue() bool
}
type KeyInfo interface {
	GetCluster(key string) string
	GetNamespace(key string) string
}
type KValue struct {
	Cluster   string
	NameSpace map[string]*compareKey
}
type compareKey struct {
	notExistKey []string
	notEqualKey []string
}

var apolloInfo map[string][]*KValue

func (kvalue *KValue) GetCluster(key string) string {
	if key == "" {
		return ""
	}
	// 每个业务线
	for _, appIdProperty := range appIdsProperty {
		// 每个集群下的nameSpace
		for clusterName, namespace := range appIdProperty.NameSpace {
			for i := 0; i < len(namespace); i++ {
				for j := 0; j < len(namespace[i].Items); j++ {
					if namespace[i].Items[j].Key == key {
						return clusterName
					}
				}
			}
		}
	}
	return ""
}
func (kvalue *KValue) GetNamespace(key string) string {
	if key == "" {
		return ""
	}
	// 每个业务线
	for _, appIdProperty := range appIdsProperty {
		// 每个集群下的nameSpace
		for _, namespace := range appIdProperty.NameSpace {
			for i := 0; i < len(namespace); i++ {
				for j := 0; j < len(namespace[i].Items); j++ {
					if namespace[i].Items[j].Key == key {
						return namespace[i].NamespaceName
					}
				}
			}
		}
	}
	return ""
}

type ApolloValue struct {
	Items map[string][]*KValue
}
type ConsulValue struct {
}

func (apolloValue *ApolloValue) GetValue(path string) (interface{}, error) {
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
					fmt.Println("namespace is nil:", namespace[i].NamespaceName, "clusterName", clusterName, "AppId", appId)
					continue
				}
				kv := make(map[string]string)
				// 将单个namespace赋值到map中
				for j := 0; j < len(namespace[i].Items); j++ {
					kv[namespace[i].Items[j].Key] = namespace[i].Items[j].Value
				}
				if _, ok := kv["consul_key"]; ok {
					fmt.Println("key is consul_key")
					fmt.Println("namespace:", namespace[i].NamespaceName, "clusterName", clusterName, "AppId", appId)
					continue
				}
				comkey := &compareKey{}
				for k, v := range kv {
					consulValue1, err := consulValue.GetValue(client, k)
					if err != nil || consulValue1.(string) == "" {
						comkey.notExistKey = append(comkey.notExistKey, k)
						continue
					}
					if consulValue1 == v {
						continue
					}
					comkey.notEqualKey = append(comkey.notEqualKey, k)
				}
				kValue.NameSpace = make(map[string]*compareKey)
				kValue.NameSpace[namespace[i].NamespaceName] = comkey
			}
			kValue.Cluster = clusterName
			kValues = append(kValues, kValue)
		}
		apolloInfo[appId] = kValues
	}
	//fmt.Println("apolloInfo appid =", apolloInfo)
	//for appid, kval := range apolloInfo {
	//	fmt.Println("apolloInfo appid =", appid)
	//	//fmt.Println("apolloInfo kv =", kv)
	//	for _, val := range kval {
	//		fmt.Println("apolloInfo Cluster =", val.Cluster)
	//		fmt.Println("apolloInfo NameSpace =", val.NameSpace)
	//		fmt.Println("apolloInfo notEqualKey =", val.notEqualKey)
	//		fmt.Println("apolloInfo notExistKey =", val.notExistKey)
	//	}
	//}
	return apolloValue.Items, nil
}

//func (apolloValue *ApolloValue) CompareValue() bool {
//	fmt.Println("init consul client")
//	consulValue := &ConsulValue{}
//	client, _ := NewClient(ADDR)
//	value, _ := apolloValue.GetValue("")
//	m := value.(map[string][]*KValue)
//	fmt.Println("m =", m)
//	for appid, kv := range m {
//		fmt.Println("业务线:", appid)
//		// 每个业务线具体信息
//		for i := 0; i < len(kv); i++ {
//			// 每个集群对应的namespace
//			for _, namespace := range kv[i].NameSpace {
//				fmt.Println("集群:", kv[i].Cluster)
//				for i := 0; i < len(namespace); i++ {
//					fmt.Println("namespace:", namespace[i])
//					for k, v := range kv[i].KV {
//						//fmt.Println("k:", k)
//						if k == "" || k == "consul_key" {
//							continue
//						}
//						consulValue1, err := consulValue.GetValue(client, k)
//						if err != nil || consulValue1 == "" {
//							kv[i].notExistKey = append(kv[i].notExistKey, k)
//							continue
//						}
//						if consulValue1 == v {
//							continue
//						}
//						kv[i].notEqualKey = append(kv[i].notEqualKey, k)
//					}
//					fmt.Println("notEqualKey = ", kv[i].notEqualKey)
//					fmt.Println("notExistKey = ", kv[i].notExistKey)
//				}
//			}
//		}
//	}
//	return false
//}
func (consulValue *ConsulValue) GetValue(client *api.Client, path string) (interface{}, error) {
	value, err := GetConsulKV(client, path)
	if err != nil {
		return "", err
	}
	return value, nil
}
func (consulValue *ConsulValue) CompareValue() bool {
	return false
}

//func ApolloCompareWithConsul() error {
//	apolloKV, err := GetSingleNameSpaceInfo(Application)
//	if err != nil {
//		return err
//	}
//	consulKV, err := GetConsulKV()
//	if err != nil {
//		return err
//	}
//	for k, apolloValue := range apolloKV {
//		consulValue, ok := consulKV[k]
//		if !ok {
//			notExistKey = append(notExistKey, k)
//		}
//		if apolloValue != consulValue {
//			notEqualKey = append(notEqualKey, k)
//		}
//	}
//	fmt.Println("notEqualKey=", notEqualKey)
//	fmt.Println("notExistKey=", notExistKey)
//	return nil
//}
