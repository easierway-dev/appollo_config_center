package cworker

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	jsoniter "github.com/json-iterator/go"
	"github.com/shima-park/agollo"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/capi"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/cconsul"
	"gitlab.mobvista.com/voyager/abtesting"
)

// Worker 工作者接口
type CWorker struct {
	AgolloClient agollo.Agollo
	WkInfo       WorkInfo
}

type WorkInfo struct {
	AppID     string
	Cluster   string
	Namespace []string
	Tag       string
}

type (
	BidForce struct {
		BidForceDevice map[string]*BidForceDeviceType `toml:"BidForceDeviceType"` // key="describe"

		TargetAdxDevice map[string]*DeviceKV //key=adx
	}
	DeviceKV struct {
		DeviceIds map[string]BidForceInfo //key=deviceId
	}
	BidForceDeviceType struct {
		DeviceId    []string `toml:"DeviceId"`
		DeviceIdMd5 []string `toml:"DeviceIdMd5"`
		Adx         []string `toml:"Adx"`
		BidForceInfo
	}
	BidForceInfo struct {
		TargetCampaign      int64   `toml:"TargetCampaign"`
		TargetTemplate      int32   `toml:"TargetTemplate"`
		TargetTemplates     []int32 `toml:"TargetTemplates"`
		TargetPrice         float64 `toml:"TargetPrice"`
		TargetRtToken       string  `toml:"TargetRtToken"`
		TargetRtTriggerItem string  `toml:"TargetRtTriggerItem"`
		User                string
	}
)

func (info *WorkInfo) Key() string {
	if info.Tag == "" {
		tag := ""
		sort.Strings(info.Namespace)
		for i, namespace := range info.Namespace {
			if i == 0 {
				tag = namespace
			} else {
				tag = fmt.Sprintf("%s_%s", tag, namespace)
			}
		}
		info.Tag = fmt.Sprintf("%s_%s_%s", info.AppID, info.Cluster, tag)
	}
	return info.Tag
}

// setup workder
func Setup(wInfo WorkInfo) (*CWorker, error) {
	var work *CWorker
	newAgo, err := agollo.New(
		ccommon.AgolloConfiger.ConfigServerURL,
		wInfo.AppID,
		agollo.Cluster(wInfo.Cluster),
		agollo.PreloadNamespaces(wInfo.Namespace...),
		agollo.AutoFetchOnCacheMiss(),
		agollo.FailTolerantOnBackupExists(),
	)
	if err != nil {
		return work, err
	}
	work = &CWorker{
		AgolloClient: newAgo,
		WkInfo:       wInfo,
	}
	return work, nil
}

func RemoveDuplicatesSlice(elements []string) []string {
	if len(elements) <= 1 {
		return elements
	}
	anyMap := make(map[string]struct{}, len(elements))
	ret := make([]string, 0, len(elements))
	for _, ele := range elements {
		if _, ok := anyMap[ele]; ok {
			continue
		}
		ret = append(ret, ele)
		anyMap[ele] = struct{}{}
	}
	return ret
}

func UpdateConsul(appid, namespace, cluster, key, value, mode string) {
	if value == "" {
		//ccommon.CLogger.Warn(ccommon.DefaultDingType,"value is nil !!! consul_addr[",consulAddr,"],key[",key,"]\n")
		fmt.Println("value is nil, will not update consul!!! cluster[", cluster, "],key[", key, "]\n")
		return
	}
	if ccommon.DyAgolloConfiger == nil {
		ccommon.CLogger.Warn(ccommon.DefaultDingType, "ccommon.DyAgolloConfiger = nil")
		return
	}
	dyAgoCfg, ok := ccommon.DyAgolloConfiger[namespace]
	if !ok {
		namespace = ccommon.DefaultNamespace
		if dyAgoCfg, ok = ccommon.DyAgolloConfiger[namespace]; !ok {
			ccommon.CLogger.Warn(ccommon.DefaultDingType, namespace, " not in ccommon.DyAgolloConfiger[", ccommon.DyAgolloConfiger, "]")
			return
		}
	}
	if dyAgoCfg.ClusterConfig == nil {
		ccommon.CLogger.Warn(ccommon.DefaultDingType, "consulAddr get failed ccommon.DyAgolloConfiger[", namespace, "]=", dyAgoCfg)
		return
	}
	if dyAgoCfg.ClusterConfig.ClusterMap == nil {
		ccommon.CLogger.Warn(ccommon.DefaultDingType, "consulAddr get failed ccommon.DyAgolloConfiger.ClusterConfig[", namespace, "]=", dyAgoCfg.ClusterConfig)
		return
	}
	if _, ok := dyAgoCfg.ClusterConfig.ClusterMap[cluster]; !ok {
		ccommon.CLogger.Warn(ccommon.DefaultDingType, "cluster:", cluster, "not in  ccommon.DyAgolloConfiger[", namespace, "].ClusterConfig")
		return
	}
	consulAddrList := dyAgoCfg.ClusterConfig.ClusterMap[cluster].ConsulAddr
	//err := cconsul.WriteOne(consulAddr, strings.Replace(key, ".", "/", -1), value)
	for _, consulAddr := range consulAddrList {
		err := cconsul.WriteOne(consulAddr, key, value, mode)
		if err != nil {
			ccommon.CLogger.Error(ccommon.DefaultDingType, "consul_addr[", consulAddr, "],key[", key, "], err[", err, "]\n")
		}
	}
	return
}

//func GetAppInfo(appid, namespace string) (enUpdate, enDelete int, accessToken string) {
//	//local config
//	if ccommon.AppConfiger.AppConfigMap == nil && ccommon.DyAgolloConfiger == nil {
//		return
//	}
//	if _, ok := ccommon.AppConfiger.AppConfigMap[appid]; !ok {
//		return
//	}
//	enUpdate = ccommon.AppConfiger.AppConfigMap[appid].EnUpdateConsul
//	enDelete = ccommon.AppConfiger.AppConfigMap[appid].EnDelConsul
//	accessToken = ccommon.AppConfiger.AppConfigMap[appid].AccessToken
//	dyAgoCfg, ok := ccommon.DyAgolloConfiger[namespace]
//	if !ok {
//		namespace = ccommon.DefaultNamespace
//		if dyAgoCfg, ok = ccommon.DyAgolloConfiger[namespace]; !ok {
//			return
//		}
//	}
//	if dyAgoCfg.AppConfig == nil {
//		return
//	}
//	if dyAgoCfg.AppConfig.AppConfigMap == nil {
//		return
//	}
//	enUpdate = dyAgoCfg.AppConfig.EnUpdateConsul
//	enDelete = dyAgoCfg.AppConfig.EnDelConsul
//
//	if _, ok := dyAgoCfg.AppConfig.AppConfigMap[appid]; !ok {
//		return
//	}
//
//	if dyAgoCfg.AppConfig.AppConfigMap[appid].EnUpdateConsul != 0 {
//		enUpdate = dyAgoCfg.AppConfig.AppConfigMap[appid].EnUpdateConsul
//	}
//	if dyAgoCfg.AppConfig.AppConfigMap[appid].EnDelConsul != 0 {
//		enDelete = dyAgoCfg.AppConfig.AppConfigMap[appid].EnDelConsul
//	}
//	if dyAgoCfg.AppConfig.AppConfigMap[appid].AccessToken != "" {
//		accessToken = dyAgoCfg.AppConfig.AppConfigMap[appid].AccessToken
//	}
//	return
//}

func GetModifyInfo(nsinfo *capi.NamespaceInfo, key string) (modifier string) {
	if nsinfo == nil {
		fmt.Println("NamespaceInfo is nil")
		return
	}
	for _, item := range nsinfo.Items {
		if key == item.Key {
			modifier = item.DataChangeLastModifiedBy
			break
		}
	}
	return
}

func MergeUpdate(appID, cluster string, updateNewValue, updateOldValue map[string]interface{}, nsinfo *capi.NamespaceInfo) (updatecontent, updateconsulvalue, path string, updated_keys, modifier_list []string, willUpdateConsul bool) {
	modifier := ""
	bidForceValue := ""
	abtestValue := ""
	willUpdateConsul = true
	i := 0
	for key, value := range updateNewValue {
		i++
		skip := false
		oValue, ok := updateOldValue[key]
		if ok {
			if oValue.(string) == value.(string) {
				skip = true
			}
		}
		if key == "consul_key" {
			path = value.(string)
			continue
		}
		if !skip {
			modifier = GetModifyInfo(nsinfo, key)
			updatecontent = fmt.Sprintf("%s\nkey=%s\nold=%s\nnew=%s\nchangedby=%s\n", updatecontent, key, oValue, value, modifier)
			updated_keys = append(updated_keys, fmt.Sprintf("update_key=%s__changedby=%s", key, modifier))
			if modifier != "" {
				modifier_list = append(modifier_list, modifier)
			}
		}

		if strings.Contains(appID, ccommon.ABTestAppid) {
			var abTestValueMap abtesting.AbInfo
			err := jsoniter.Unmarshal([]byte(value.(string)), &abTestValueMap)
			if err != nil {
				willUpdateConsul = false
				ccommon.CLogger.Error(appID, "#", cluster, "#", key, ":", "\njsoniter.Unmarshal(abtest_value failed, err:", err)
				return
			}
			if i < len(updateNewValue) {
				abtestValue += value.(string) + ",\n"
			} else {
				abtestValue += value.(string) + "\n"
			}
			updateconsulvalue = "[" + strings.Trim(strings.Trim(abtestValue, "\n"), ",") + "]"
		} else if strings.Contains(appID, ccommon.BidForceAppid) {
			var bidForceValueMap = BidForce{}
			_, err := toml.Decode(value.(string), &bidForceValueMap)
			if err != nil {
				ccommon.CLogger.Error(appID, "#", cluster, "#", key, ":", "\ntoml.Decode(bidforce_value failed, err:", err)
				continue
			}
			bidForceValue += strings.TrimSpace(value.(string)) + "\n"
			updateconsulvalue = bidForceValue
		}
	}
	return
}

//work run
func (cw *CWorker) Run(ctx context.Context) {
	errorCh := cw.AgolloClient.Start()
	watchCh := cw.AgolloClient.Watch()
	go func(cw *CWorker) {
		for {
			select {
			case <-ctx.Done():
				ccommon.CLogger.Info(ccommon.DefaultDingType, cw.WkInfo.Cluster, "watch quit...")
				return
			case err := <-errorCh:
				if ccommon.AppConfiger.AppConfigMap != nil {
					if _, ok := ccommon.AppConfiger.AppConfigMap[ccommon.DefaultPollDingType]; ok {
						ccommon.ChklogRate = ccommon.AppConfiger.AppConfigMap[ccommon.DefaultPollDingType].ChklogRate
					}
				}
				if ccommon.ChklogRamdom < ccommon.ChklogRate {
					ccommon.CLogger.Info(ccommon.DefaultPollDingType, "Error:", err)
				}
			case update := <-watchCh:
				consulMode := "write"
				ccommon.Configer = ccommon.InitAppCfgMap(ccommon.AppConfiger, cw.WkInfo.AppID, update.Namespace)
				fmt.Println("ccommon.Configer=",ccommon.Configer)
				enConsul := ccommon.Configer.EnUpdateConsul
				enDelete := ccommon.Configer.EnDelConsul
				token := ccommon.Configer.AccessToken
				//enConsul, enDelete, token := GetAppInfo(cw.WkInfo.AppID, update.Namespace)
				if enConsul != 1 {
					ccommon.CLogger.Warn(cw.WkInfo.AppID, "is not permit to update consul")
					ccommon.CLogger.Info(ccommon.DefaultDingType, "Apollo cluster(", cw.WkInfo.Cluster, ") namespace(", update.Namespace, ") \nold_value:(", update.OldValue, ") \nnew_value:(", update.NewValue, ") \n error:(", update.Error, ")\n")
				} else {
					var deletedKeys []string
					updateContent := ""
					var updatedKeys []string
					modifier := ""
					willUpdateConsul := true
					var modifierList []string
					url := fmt.Sprintf("http://%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s", ccommon.AgolloConfiger.PortalURL, "DEV", cw.WkInfo.AppID, cw.WkInfo.Cluster, update.Namespace)
					nsInfo, _ := capi.GetNamespaceInfo(url, token)
					if strings.Contains(cw.WkInfo.AppID, ccommon.ABTestAppid) || strings.Contains(cw.WkInfo.AppID, ccommon.BidForceAppid) {
						updateConsulValue := ""
						path := ""
						updateContent, updateConsulValue, path, updatedKeys, modifierList, willUpdateConsul = MergeUpdate(cw.WkInfo.AppID, cw.WkInfo.Cluster, update.NewValue, update.OldValue, nsInfo)
						if path != "" && willUpdateConsul {
							UpdateConsul(cw.WkInfo.AppID, update.Namespace, cw.WkInfo.Cluster, path, updateConsulValue, consulMode)
						}
						//delete keys
						for k, _ := range update.OldValue {
							if _, ok := update.NewValue[k]; !ok {
								deletedKeys = append(deletedKeys, k)
							}
						}
					} else {
						//新增、更新
						for path, value := range update.NewValue {
							if oValue, ok := update.OldValue[path]; ok {
								//未发生变化的key，跳过不更新
								if oValue.(string) == value.(string) {
									continue
								}
							}
							modifier = GetModifyInfo(nsInfo, path)
							//updatecontent = fmt.Sprintf("%s\nkey=%s\nold=%s\dunnew=%s\nmodifier=%s\n", updatecontent, path, ovalue, value, modifier)
							updatedKeys = append(updatedKeys, fmt.Sprintf("update_key=%s__changedby=%s", path, modifier))
							if modifier != "" {
								modifierList = append(modifierList, modifier)
							}
							UpdateConsul(cw.WkInfo.AppID, update.Namespace, cw.WkInfo.Cluster, path, value.(string), consulMode)
						}
						//删除
						if enDelete == 1 {
							for path, value := range update.OldValue {
								if _, ok := update.NewValue[path]; !ok {
									deletedKeys = append(deletedKeys, path)
									consulMode = "del"
									UpdateConsul(cw.WkInfo.AppID, update.Namespace, cw.WkInfo.Cluster, path, value.(string), consulMode)
								}
							}
						}
					}
					//只有abtest显示更新内容的详情，其他都只提示变更的key
					if find := strings.Contains(cw.WkInfo.AppID, ccommon.ABTestAppid); !find && len(updatedKeys) > 0 {
						updateContent = strings.Join(updatedKeys, "\n")
					}
					//记录删除的key
					if len(deletedKeys) > 0 {
						updateContent = fmt.Sprintf("%s\n\ndelelte_key=%s", updateContent, strings.Join(deletedKeys, "#"))
					}
					ccommon.CLogger.Info(ccommon.DefaultDingType, "Apollo cluster(", cw.WkInfo.Cluster, ") namespace(", update.Namespace, ") \nold_value:(", update.OldValue, ") \nnew_value:(", update.NewValue, ") \n error:(", update.Error, ")\n")
					if willUpdateConsul {
						if updateContent == "" {
							updateContent = fmt.Sprintf("nothing to update !!!\nisSupportDelete=", string(enDelete), " (1: support)")
						}
						if len(modifierList) > 0 {
							ccommon.CLogger.Warn(RemoveDuplicatesSlice(modifierList), cw.WkInfo.AppID, "#", cw.WkInfo.Cluster, "#", update.Namespace, ": \nupdatecontent:\n", updateContent)
						} else {
							ccommon.CLogger.Warn(cw.WkInfo.AppID, "#", cw.WkInfo.Cluster, "#", update.Namespace, ": \nupdatecontent:\n", updateContent)
						}
					} else {
						ccommon.CLogger.Warn(cw.WkInfo.AppID, "#", cw.WkInfo.Cluster, "#", update.Namespace, ": !!! invalid config will not update consul !!!")
					}
				}
			}
		}
	}(cw)
}

//work stop
func (cw *CWorker) Stop() {
	cw.AgolloClient.Stop()
}
