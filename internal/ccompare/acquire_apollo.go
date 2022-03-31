package ccompare

import (
	"errors"
	"fmt"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/capi"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
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
	NameSpace   map[string][]*capi.NamespaceInfo
	AccessToken string
}
// 全部的业务线
var appID []string

// 集群信息
var envClustersInfoMap map[string][]*capi.EnvClustersInfo

// 各个业务线对应的token
var appIdAccessToken map[string]string

// 各个业务线对应的集群信息
var appIdsProperty map[string]*AppIdProperty

func New() *AppIdProperty {
	return &AppIdProperty{}
}

// 获取所有业务线
func getAllAppID() error {
	url2 := fmt.Sprintf("http://%s/openapi/v1/apps", ccommon.AgolloConfiger.PortalURL)
	_, err := getAccessToken(GlobalConfiger.AppConfigMap)
	fmt.Println(appIdAccessToken)
	//token, err := getDspToken(globalConfig.AccessToken)
	if err != nil {
		return err
	}
	// 只要获取某个业务线的token就可以，这里以dsp的token为例
	appInfo, _ := capi.GetAppInfo(url2, appIdAccessToken["dsp"])
	//fmt.Println("url2=", url2)
	//fmt.Println("appInfo=", appInfo)
	if len(appInfo) == 0 {
		fmt.Println("appInfo is nil ")
		return errors.New("appInfo is nil ")
	}
	for _, v := range appInfo {
		appID = append(appID, v.AppId)
	}
	//fmt.Println("appID=", appID)
	return nil
}

// 获取所有业务线的集群信息
func getEnvClustersInfo(appID string) (map[string][]*capi.EnvClustersInfo, error) {
	url1 := fmt.Sprintf("http://%s/openapi/v1/apps/%s/envclusters", ccommon.AgolloConfiger.PortalURL, appID)
	envClustersInfoMap = make(map[string][]*capi.EnvClustersInfo)
	for _, token := range appIdAccessToken {
		envClustersInfo, _ := capi.GetEnvClustersInfo(url1, token)
		envClustersInfoMap[appID] = envClustersInfo
	}
	if len(envClustersInfoMap) == 0 {
		fmt.Println("EnvClustersInfoMap is nil ")
		return nil, errors.New("EnvClustersInfoMap is nil ")
	}
	fmt.Println("ecinfo=", envClustersInfoMap)

	return envClustersInfoMap, nil
}

func (appIdProperty *AppIdProperty) applyProperty() error {
	//fmt.Println("appId length = ", len(appID))
	// index为下标索引，id为具体的业务线
	appIdsProperty = make(map[string]*AppIdProperty)
	// 各个业务线
	for _, id := range appID {
		envClustersMap, _ := getEnvClustersInfo(id)
		// 各个业务线下的集群信息
		for in, _ := range envClustersMap[id] {
			if envClustersMap[id][in].Env == DEV {
				appIdProperty = &AppIdProperty{AppId: id, Env: DEV, ClusterName: envClustersMap[id][in].Clusters, AccessToken: appIdAccessToken[id]}
				fmt.Println("appIdProperty AppId= ", appIdProperty.AppId)
				fmt.Println("appIdProperty Env= ", appIdProperty.Env)
				fmt.Println("appIdProperty ClusterName= ", appIdProperty.ClusterName)
				// 获取每个业务线的下的每个集群的所有namesapce
				for id, cluster := range envClustersMap[id][in].Clusters {
					nameSpaceInfo, err := appIdProperty.getNameSpaceInfo(id)
					// 如果当前集群没有namespace，直接跳过
					if err != nil || nameSpaceInfo == nil {
						continue
					}
					if len(appIdProperty.NameSpace) == 0 {
						appIdProperty.NameSpace = map[string][]*capi.NamespaceInfo{}
					}
					// 每个集群下对应的namespace
					appIdProperty.NameSpace[cluster] = nameSpaceInfo
				}
				fmt.Println("appIdProperty NameSpace= ", appIdProperty.NameSpace)
				appIdsProperty[id] = appIdProperty
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
func (apolloProperty *AppIdProperty) getNameSpaceInfo(id int) (respBody []*capi.NamespaceInfo, err error) {
	url := fmt.Sprintf("http://%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces", ccommon.AgolloConfiger.PortalURL, apolloProperty.Env, apolloProperty.AppId, apolloProperty.ClusterName[id])
	//fmt.Println("url=", url)
	//if apolloProperty.AccessToken == "" {
	//	return nil, errors.New("AccessToken is nil")
	//}
	// 暂时使用默认的accessToken,后面可以修改为apolloProperty.AccessToken
	nSAllInfo, _ := capi.GetAllNamespaceInfo(url, "280c6b92cd8ee4f1c5833b4bd22dfe44a4778ab5")
	if nSAllInfo == nil {
		return nil, errors.New("nSAllInfo is nil")
	}
	//fmt.Println("nSAllInfo=", nSAllInfo)
	return nSAllInfo, nil
}
func GetAppIdsProperty() (err error) {
	// 获取全局AppID
	err = getAllAppID()
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
func getAccessToken(m map[string]ccommon.ConfigInfo) (map[string]string, error) {
	if len(m) == 0 {
		return nil, errors.New("not contain DspToken")
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
			for k,v:= range keys.NotExistKey{
				fmt.Println("dsp apolloInfo notExistKey =", k)
				fmt.Println("dsp apolloInfo DataChangeLastModifiedBy =", v.DataChangeLastModifiedBy)
			}
			for k,v:= range keys.NotEqualKey{
				fmt.Println("dsp apolloInfo NotEqualKey =", k)
				fmt.Println("dsp apolloInfo DataChangeLastModifiedBy =", v.DataChangeLastModifiedBy)
			}
		}
	}
}
func Start(server *AgolloServer) {
	apollo := &ApolloValue{}
	// 获取全局配置
	GetApolloGlobalConfig(server)
	// 每个业务线的具体信息
	GetAppIdsProperty()
	// 对比
	apollo.CompareValue()
	readCon()
}
