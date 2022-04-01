package ccompare

import (
	"errors"
	"fmt"
)

const (
	DEV      = "DEV"
	DSPALIVG = "dsp_ali_vg"
	DSP      = "dsp"
)

type AppIdProperty struct {
	AppId       string
	Env         string
	ClusterName []string
	NameSpace   map[string][]*NamespaceInfo
	AccessToken string
}

// 各个业务线对应的token
var appIdAccessToken map[string]string

// 各个业务线对应的集群信息
var appIdsProperty map[string]*AppIdProperty

func New() *AppIdProperty {
	return &AppIdProperty{}
}

func (appIdProperty *AppIdProperty) applyProperty() error {
	//fmt.Println("appId length = ", len(appID))
	// index为下标索引，id为具体的业务线
	appIdsProperty = make(map[string]*AppIdProperty)
	// 各个业务线
	envClustersMap := AppIdClusters.EnvClustersInfoMap
	// 各个业务线下的集群信息
	for appid, envClusters := range envClustersMap {
		for i := 0; i < len(envClusters); i++ {
			if envClusters[i].Env == DEV {
				appIdProperty = &AppIdProperty{AppId: appid, Env: DEV, ClusterName: envClusters[i].Clusters, AccessToken: appIdAccessToken[appid]}
				fmt.Println("appIdProperty AppId= ", appIdProperty.AppId)
				fmt.Println("appIdProperty Env= ", appIdProperty.Env)
				fmt.Println("appIdProperty ClusterName= ", appIdProperty.ClusterName)
				// 获取每个业务线的下的每个集群的所有namesapce
				for id, cluster := range envClusters[i].Clusters {
					nameSpaceInfo, err := appIdProperty.getNameSpaceInfo(id)
					// 如果当前集群没有namespace，直接跳过
					if err != nil || nameSpaceInfo == nil {
						continue
					}
					if len(appIdProperty.NameSpace) == 0 {
						appIdProperty.NameSpace = map[string][]*NamespaceInfo{}
					}
					// 每个集群下对应的namespace
					appIdProperty.NameSpace[cluster] = nameSpaceInfo
				}
				fmt.Println("appIdProperty NameSpace= ", appIdProperty.NameSpace)
				appIdsProperty[appid] = appIdProperty
			} else {
				return errors.New("not DEV environment")
			}
		}
	}
	return nil
}

// 通过全局的appIdsProperty
// 循环获取Clusters的Namespace
// 通过集群名，appID，namespace查找对应的信息：获取集群下所有Namespace信息接口，在进行细分每一个namespace
func (apolloProperty *AppIdProperty) getNameSpaceInfo(id int) (respBody []*NamespaceInfo, err error) {
	url := fmt.Sprintf("http://%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces", AgolloConfiger.PortalURL, apolloProperty.Env, apolloProperty.AppId, apolloProperty.ClusterName[id])
	//fmt.Println("url=", url)
	//if apolloProperty.AccessToken == "" {
	//	return nil, errors.New("AccessToken is nil")
	//}
	nSAllInfo, _ := GetAllNamespaceInfo(url, "280c6b92cd8ee4f1c5833b4bd22dfe44a4778ab5")
	if nSAllInfo == nil {
		return nil, errors.New("nSAllInfo is nil")
	}
	//fmt.Println("nSAllInfo=", nSAllInfo)
	return nSAllInfo, nil
}
func GetAppIdsProperty() (err error) {
	// 获取全局AppID
	appIdClusters := AppIdClustersInfo{}
	appIdClusters.GetConfigInfo()
	//err = getAllAppID()
	if err != nil {
		return err
	}
	appIdProperty := New()
	// 为appIdProperty赋值
	err = appIdProperty.applyProperty()
	if err != nil {
		return err
	}
	//SetNameSpaceInfo()
	//readCon()
	fmt.Println("appIdsProperty=", appIdsProperty)
	return nil
}

// 获取appId对应的accessToken
func getAccessToken(m map[string]ConfigInfo) (map[string]string, error) {
	if len(m) == 0 {
		return nil, errors.New("config is nil")
	}
	appIdAccessToken = make(map[string]string, 6)
	for key, info := range m {
		appIdAccessToken[key] = info.AccessToken
	}
	return appIdAccessToken, nil
}
func readCon() {
	//
	for _, val := range apolloInfo["dsp"] {
		//fmt.Println("apolloInfo kv =", kv)
		fmt.Println("dsp apolloInfo Cluster =", val.Cluster)
		for namespace, keys := range val.NameSpace {
			fmt.Println("dsp apolloInfo NameSpace =", namespace)
			for k, v := range keys.NotExistKey {
				fmt.Println("dsp apolloInfo notExistKey =", k)
				fmt.Println("dsp apolloInfo DataChangeLastModifiedBy =", v.DataChangeLastModifiedBy)
			}
			for k, v := range keys.NotEqualKey {
				fmt.Println("dsp apolloInfo NotEqualKey =", k)
				fmt.Println("dsp apolloInfo DataChangeLastModifiedBy =", v.DataChangeLastModifiedBy)
			}
		}
	}
}
func Start() {
	// 初始化配置文件
	if err := Init(); err != nil {
		panic(err)
	}
	apollo := &ApolloValue{}
	globalconfig := &GlobalConfig{}
	// 获取全局配置
	globalconfig.GetConfigInfo()
	// 每个业务线的具体信息
	GetAppIdsProperty()
	// 对比
	apollo.CompareValue()
	readCon()
}
