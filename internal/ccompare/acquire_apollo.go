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

// 各个业务线对应的集群信息

type Properties struct {
	AppIdsProperty map[string]*AppIdProperty
}

var Property *Properties

func New() *AppIdProperty {
	return &AppIdProperty{}
}

func (appIdProperty *AppIdProperty) applyProperty() error {
	//fmt.Println("appId length = ", len(appID))
	// index为下标索引，id为具体的业务线

	Property.AppIdsProperty = make(map[string]*AppIdProperty)
	// 各个业务线
	envClustersMap := AppIdClusters.EnvClustersInfoMap
	// 各个业务线下的集群信息
	for appid, envClusters := range envClustersMap {
		for i := 0; i < len(envClusters); i++ {
			if envClusters[i].Env == DEV {
				appIdProperty = &AppIdProperty{AppId: appid, Env: DEV, ClusterName: envClusters[i].Clusters, AccessToken: AppIdAccessToken[appid]}
				fmt.Println("appIdProperty AppId= ", appIdProperty.AppId)
				fmt.Println("appIdProperty Env= ", appIdProperty.Env)
				fmt.Println("appIdProperty ClusterName= ", appIdProperty.ClusterName)
				fmt.Println("appIdProperty AccessToken= ", appIdProperty.AccessToken)
				// 获取每个业务线的下的每个集群的所有namesapce
				for id, cluster := range envClusters[i].Clusters {
					nameSpaceInfo, err := appIdProperty.getNameSpaceInfo(id)
					// 如果当前集群没有namespace，直接跳过
					if err != nil || nameSpaceInfo == nil {
						fmt.Println("APPID:", appid, "\tclusterName:", cluster, "\terr:", err)
						continue
					}
					if len(appIdProperty.NameSpace) == 0 {
						appIdProperty.NameSpace = map[string][]*NamespaceInfo{}
					}
					// 每个集群下对应的namespace
					appIdProperty.NameSpace[cluster] = nameSpaceInfo
				}
				Property.AppIdsProperty[appid] = appIdProperty
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
	fmt.Println("url=", url)
	if apolloProperty.AccessToken == "" {
		return nil, errors.New("AccessToken is nil")
	}
	nSAllInfo, _ := GetAllNamespaceInfo(url, apolloProperty.AccessToken)
	if nSAllInfo == nil {
		return nil, errors.New("nSAllInfo is nil")
	}
	//fmt.Println("nSAllInfo=", nSAllInfo)
	return nSAllInfo, nil
}
func (property *Properties) getAppIdsProperty() (err error) {
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
	fmt.Println("appIdsProperty=", property.AppIdsProperty)
	return nil
}

func Start() {
	apollo := &ApolloValue{}
	Property = &Properties{}
	// 每个业务线的具体信息
	Property.getAppIdsProperty()
	// 对比
	apollo.CompareValue()
	//apollo.Print(nil)
	Consul.Print(nil)
}
