package cworker

import (
	"fmt"
	"sort"
	"context"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/shima-park/agollo"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/cconsul"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/capi"	
	"gitlab.mobvista.com/voyager/abtesting"
	jsoniter "github.com/json-iterator/go"
)

// Worker 工作者接口
type CWorker struct {
        AgolloClient agollo.Agollo
        WkInfo      WorkInfo
}

type WorkInfo struct {
	AppID string
	Cluster string
	Namespace []string
	Tag string
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
		TargetCampaign  int64   `toml:"TargetCampaign"`
		TargetTemplate  int32   `toml:"TargetTemplate"`
		TargetTemplates []int32 `toml:"TargetTemplates"`
		TargetPrice     float64 `toml:"TargetPrice"`
		TargetRtToken   string  `toml:"TargetRtToken"`
		TargetRtTriggerItem string   `toml:"TargetRtTriggerItem"`
		User string
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
	  tag = fmt.Sprintf("%s_%s",tag, namespace)
        }
    }
    info.Tag = fmt.Sprintf("%s_%s_%s",info.AppID, info.Cluster, tag)
  }
  return info.Tag
}
// setup workder
func Setup(wInfo WorkInfo)(*CWorker,error){
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
		AgolloClient:  newAgo,
		WkInfo:      wInfo,
	}
	return work, nil
}

func UpdateConsul(appid, namespace, cluster, key, value , mode string){
	if ccommon.DyAgolloConfiger != nil {
		if _,ok := ccommon.DyAgolloConfiger[namespace];!ok {
			namespace = ccommon.DefaultNamespace
		}
		if dyAgoCfg,ok := ccommon.DyAgolloConfiger[namespace];ok {
			if dyAgoCfg.ClusterConfig != nil && dyAgoCfg.ClusterConfig.ClusterMap != nil {
				if _,ok := dyAgoCfg.ClusterConfig.ClusterMap[cluster];ok {
					if value == "" {
						//ccommon.CLogger.Warn(ccommon.DefaultDingType,"value is nil !!! consul_addr[",consulAddr,"],key[",key,"]\n")
						fmt.Println("value is nil, will not update consul!!! cluster[",cluster,"],key[",key,"]\n")
						return
					}
					consulAddrList := dyAgoCfg.ClusterConfig.ClusterMap[cluster].ConsulAddr
					//err := cconsul.WriteOne(consulAddr, strings.Replace(key, ".", "/", -1), value)
					for _,consulAddr := range consulAddrList {
						err := cconsul.WriteOne(consulAddr, key, value, mode)
						if err != nil {
							ccommon.CLogger.Error(ccommon.DefaultDingType,"consul_addr[",consulAddr,"],key[",key,"], err[", err,"]\n")
						} 
					}
				} else {
					ccommon.CLogger.Warn(ccommon.DefaultDingType,"cluster:",cluster,"not in  ccommon.DyAgolloConfiger[",namespace,"].ClusterConfig")
					return
				}
			} else {
				ccommon.CLogger.Warn(ccommon.DefaultDingType,"consulAddr get failed ccommon.DyAgolloConfiger[",namespace,"]=",dyAgoCfg)
				return
			}
		} else {
			ccommon.CLogger.Warn(ccommon.DefaultDingType,namespace," not in ccommon.DyAgolloConfiger[",ccommon.DyAgolloConfiger,"]")
			return
		}
	} else {
		ccommon.CLogger.Warn(ccommon.DefaultDingType,"ccommon.DyAgolloConfiger = nil")
	}
	return
}

func GetAppInfo(appid, namespace string) (enUpdate, enDelete int,accessToken string) {
	//local config
	if ccommon.AppConfiger.AppConfigMap != nil {
		if _,ok := ccommon.AppConfiger.AppConfigMap[appid];ok {
			enUpdate = ccommon.AppConfiger.AppConfigMap[appid].EnUpdateConsul
			enDelete = ccommon.AppConfiger.AppConfigMap[appid].EnDelConsul
			accessToken = ccommon.AppConfiger.AppConfigMap[appid].AccessToken
		} 		
	}
	//apollo global config
	if ccommon.DyAgolloConfiger != nil {
		if _,ok := ccommon.DyAgolloConfiger[namespace];!ok {
			namespace = ccommon.DefaultNamespace
		}
		if dyAgoCfg,ok := ccommon.DyAgolloConfiger[namespace];ok {
			if dyAgoCfg.AppConfig != nil {
				enUpdate = dyAgoCfg.AppConfig.EnUpdateConsul
				enDelete = dyAgoCfg.AppConfig.EnDelConsul
				if dyAgoCfg.AppConfig.AppConfigMap != nil {
					if _,ok := dyAgoCfg.AppConfig.AppConfigMap[appid];ok{
						if dyAgoCfg.AppConfig.AppConfigMap[appid].EnUpdateConsul != 0 {
							enUpdate = dyAgoCfg.AppConfig.AppConfigMap[appid].EnUpdateConsul
						}
						if dyAgoCfg.AppConfig.AppConfigMap[appid].EnDelConsul != 0 {
							enDelete = dyAgoCfg.AppConfig.AppConfigMap[appid].EnDelConsul
						}
						if dyAgoCfg.AppConfig.AppConfigMap[appid].AccessToken != "" {
							accessToken = dyAgoCfg.AppConfig.AppConfigMap[appid].AccessToken
						}
					}
				}
			}
		}
	}
	return
}

func GetModifyInfo(nsinfo *capi.NamespaceInfo, key string) (modifier string) {
	if nsinfo == nil{
		fmt.Println("NamespaceInfo is nil")
		return
	}
	for _,item := range nsinfo.Items {
		if key == item.Key {
			modifier = item.DataChangeLastModifiedBy
			break
		}
	} 
	return
}

func MergeUpdate(appID, cluster string, updateNewValue, updateOldValue map[string]interface{}, nsinfo *capi.NamespaceInfo) (updatecontent, updateconsulvalue, path string, updated_keys, modifier_list []string, willUpdateConsul bool) {
	modifier := ""
	bidforce_value := ""
	abtest_value := ""
	willUpdateConsul = true
	i := 0
	for key, value := range updateNewValue {
		i = i + 1
		v, _ := value.(string)
		skip := false
		ovalue, ok := updateOldValue[key]
		if ok {
			ov, _ := ovalue.(string)
			if ov == v {
				skip = true
			}
		}
		if key == "consul_key" {
			path = value.(string)
			continue
		}
		if ! skip {
			modifier = GetModifyInfo(nsinfo, key)
			updatecontent = fmt.Sprintf("%s\nkey=%s\nold=%s\nnew=%s\nchangedby=%s\n", updatecontent, key, ovalue, value, modifier)
			updated_keys = append(updated_keys, fmt.Sprintf("update_key=%s__changedby=%s",key, modifier))
			if modifier != "" {
				modifier_list = append(modifier_list, modifier)
			}								
		}

		if strings.Contains(appID, ccommon.ABTestAppid) {
			var abtest_valuemap abtesting.AbInfo
			err := jsoniter.Unmarshal([]byte(value.(string)), &abtest_valuemap)
			if err == nil {
				if i < len(updateNewValue) {
					abtest_value = abtest_value + value.(string) + ",\n"
				} else {
					abtest_value = abtest_value + value.(string) + "\n"
				}
			} else {
				willUpdateConsul = false
				ccommon.CLogger.Error(appID,"#",cluster,"#",key,":", "\njsoniter.Unmarshal(abtest_value failed, err:", err)
			}
			updateconsulvalue = "["+strings.Trim(strings.Trim(abtest_value, "\n"),",")+"]"
		} else if strings.Contains(appID, ccommon.BidForceAppid) {
			var bidforce_valuemap = BidForce{}
			if _, err := toml.Decode(value.(string), &bidforce_valuemap);err == nil {
				bidforce_value = bidforce_value + strings.TrimSpace(value.(string)) + "\n"
			} else {
				ccommon.CLogger.Error(appID,"#",cluster,"#",key,":", "\ntoml.Decode(bidforce_value failed, err:", err)
				continue
			}
			updateconsulvalue = bidforce_value
		}
	}
	return
}

//work run
func (cw *CWorker) Run(ctx context.Context){
	errorCh := cw.AgolloClient.Start()
	watchCh := cw.AgolloClient.Watch()
	go func(cw *CWorker) {
		for {
			select {
			case <-ctx.Done():
				ccommon.CLogger.Info(ccommon.DefaultDingType,cw.WkInfo.Cluster, "watch quit...")
				return
			case err := <-errorCh:
				if ccommon.AppConfiger.AppConfigMap != nil {
					if _,ok := ccommon.AppConfiger.AppConfigMap[ccommon.DefaultPollDingType];ok {
						ccommon.ChklogRate = ccommon.AppConfiger.AppConfigMap[ccommon.DefaultPollDingType].ChklogRate
					}
				}
				if ccommon.ChklogRamdom < ccommon.ChklogRate {
					ccommon.CLogger.Info(ccommon.DefaultPollDingType,"Error:", err)
				}
			case update := <-watchCh:
				consulMode := "write"
				enConsul, enDelete, token := GetAppInfo(cw.WkInfo.AppID, update.Namespace)
				if enConsul != 1 {
					ccommon.CLogger.Warn(cw.WkInfo.AppID, "is not permit to update consul")
					ccommon.CLogger.Info(ccommon.DefaultDingType,"Apollo cluster(",cw.WkInfo.Cluster,") namespace(",update.Namespace,") \nold_value:(", update.OldValue,") \nnew_value:(",update.NewValue,") \n error:(",update.Error,")\n")
				} else {
					deleted_keys := []string{}
					updatecontent := ""
					updated_keys := []string{}
					modifier := ""
					willUpdateConsul := true
					modifier_list := []string{}
					url := fmt.Sprintf("http://%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s", ccommon.AgolloConfiger.PortalURL, "DEV", cw.WkInfo.AppID, cw.WkInfo.Cluster, update.Namespace)
					ns_info,_ := capi.GetNamespaceInfo(url, token)
					if strings.Contains(cw.WkInfo.AppID, ccommon.ABTestAppid) || strings.Contains(cw.WkInfo.AppID, ccommon.BidForceAppid) {
						updateconsulvalue := ""
						path := ""
						updatecontent, updateconsulvalue, path, updated_keys, modifier_list, willUpdateConsul = MergeUpdate(cw.WkInfo.AppID, cw.WkInfo.Cluster, update.NewValue, update.OldValue, ns_info)
						if path != "" {
							UpdateConsul(cw.WkInfo.AppID, update.Namespace, cw.WkInfo.Cluster, path, updateconsulvalue, consulMode)
						}
						//delete keys
						for k, _ := range update.OldValue {
							if _,ok := update.NewValue[k]; ! ok {
								deleted_keys = append(deleted_keys, k)
							}
						}
					} else {
						//新增、更新
						for path, value := range update.NewValue {
							v, _ := value.(string)
							ovalue, ok := update.OldValue[path]
							if ok {
								ov, _ := ovalue.(string)
								//未发生变化的key，跳过不更新
								if ov == v {
									continue
								}
							}
							modifier = GetModifyInfo(ns_info, path)
							//updatecontent = fmt.Sprintf("%s\nkey=%s\nold=%s\nnew=%s\nmodifier=%s\n", updatecontent, path, ovalue, value, modifier)
							updated_keys = append(updated_keys, fmt.Sprintf("update_key=%s__changedby=%s",path, modifier))
							if modifier != "" {
								modifier_list = append(modifier_list, modifier)
							}								
							UpdateConsul(cw.WkInfo.AppID, update.Namespace, cw.WkInfo.Cluster, path, v, consulMode) 
						}
						//删除
						if enDelete == 1 {
							for path, value := range update.OldValue {
								if _,ok := update.NewValue[path]; ! ok {
									deleted_keys = append(deleted_keys, path)
									v, _ := value.(string)
									consulMode = "del"
									UpdateConsul(cw.WkInfo.AppID, update.Namespace, cw.WkInfo.Cluster, path, v, consulMode)
								}
							}
						}
					}
					//只有abtest显示更新内容的详情，其他都只提示变更的key
					if find := strings.Contains(cw.WkInfo.AppID, ccommon.ABTestAppid); ! find && len(updated_keys) > 0 {
						updatecontent = strings.Join(updated_keys, "\n")
					}
					//记录删除的key
					if len(deleted_keys) > 0 {					
						updatecontent = fmt.Sprintf("%s\n\ndelelte_key=%s",updatecontent, strings.Join(deleted_keys, "#"))
					}
					ccommon.CLogger.Info(ccommon.DefaultDingType,"Apollo cluster(",cw.WkInfo.Cluster,") namespace(",update.Namespace,") \nold_value:(", update.OldValue,") \nnew_value:(",update.NewValue,") \n error:(",update.Error,")\n")
					if willUpdateConsul {
						if updatecontent == "" {
							updatecontent = fmt.Sprintf("nothing to update !!!\nisSupportDelete=",enDelete, " (1: support)")
						}
						if len(modifier_list) > 0 {
							ccommon.CLogger.Warn(modifier_list, cw.WkInfo.AppID,"#",cw.WkInfo.Cluster,"#",update.Namespace,": \nupdatecontent:\n",updatecontent)
						} else {
							ccommon.CLogger.Warn(cw.WkInfo.AppID,"#",cw.WkInfo.Cluster,"#",update.Namespace,": \nupdatecontent:\n",updatecontent)
						}
					}	else {
						ccommon.CLogger.Warn(cw.WkInfo.AppID,"#",cw.WkInfo.Cluster,"#",update.Namespace,": !!! invalid config will not update consul !!!")
					}				
				}
			}
		}
	}(cw)
}

//work stop
func (cw *CWorker) Stop(){
	cw.AgolloClient.Stop()
}
