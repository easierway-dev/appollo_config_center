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
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/chttp"	
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

func UpdateConsul(appid, namespace, cluster, key, value string){
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
						err := cconsul.WriteOne(consulAddr, key, value)
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

func GetAppInfo(appid, namespace string) (enUpdate bool, accessToken string) {
	if ccommon.DyAgolloConfiger != nil {
		if _,ok := ccommon.DyAgolloConfiger[namespace];!ok {
			namespace = ccommon.DefaultNamespace
		}
		if dyAgoCfg,ok := ccommon.DyAgolloConfiger[namespace];ok {
			if dyAgoCfg.AppConfig != nil {
				enUpdate = dyAgoCfg.AppConfig.EnUpdateConsul
				if dyAgoCfg.AppConfig.AppConfigMap != nil {
					if _,ok := dyAgoCfg.AppConfig.AppConfigMap[appid];ok{
						enUpdate = dyAgoCfg.AppConfig.AppConfigMap[appid].EnUpdateConsul
						accessToken = dyAgoCfg.AppConfig.AppConfigMap[appid].AccessToken
					}
				}
			}
		}
	}
	return
}

func GetModifyInfo(nsinfo interface{}, key string) string {
	modifier := ""
	if itemsList,find := nsinfo["items"]; find {
		for _,item := range itemsList {
			if k,find := item["key"]; find && key == k {
				modifier = item["dataChangeLastModifiedBy"]
				break
			}
		} 
	}
	return modifier
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
				enConsul,token := EnableConsulUpdate(cw.WkInfo.AppID, update.Namespace)
				if ! enConsul {
					ccommon.CLogger.Warn(cw.WkInfo.AppID, "is not permit to update consul")
					ccommon.CLogger.Info(ccommon.DefaultDingType,"Apollo cluster(",cw.WkInfo.Cluster,") namespace(",update.Namespace,") \nold_value:(", update.OldValue,") \nnew_value:(",update.NewValue,") \n error:(",update.Error,")\n")
				} else {
					deleted_keys := ""
					updatecontent := ""
					updatekey := ""
					url := fmt.Sprintf("http://%s/openapi/v1/envs/%s/apps/%s/clusters/%s/namespaces/%s", ccommon.AgolloConfiger.PortalURL, "DEV", cw.WkInfo.AppID, update.Namespace)
					ns_info := chttp.httpGet(url, token)
					modifier_list := []string{}
					if strings.Contains(cw.WkInfo.AppID, ccommon.ABTestAppid) {
						path := ""
						abtestvalue := ""
						for key, value := range update.NewValue {
							v, _ := value.(string)
							skip := false
							ovalue, ok := update.OldValue[key]
							if ok {
								ov, _ := ovalue.(string)
								if ov == v {
									skip = true
								}
							}
							if ! skip {
								modifier : = GetModifyInfo(ns_info, key)
								updatecontent = fmt.Sprintf("%s\nkey=%s\nold=%s\nnew=%s\nmodifier=%s\n", updatecontent, key, ovalue, value, modifier)
								updatekey = fmt.Sprintf("key=%s,%s\nmodifier=%s\n", updatekey, key, modifier)
							  if modifier != "" {
									modifier_list = append(modifier_list, modifier)
								}								
							}
							if key == "consul_key" {
								path = value.(string)
								continue
							}

							var abtest_value abtesting.AbInfo
							err := jsoniter.Unmarshal([]byte(value.(string)), &abtest_value)
							if err == nil {
								if i < len(update.NewValue) {
									abtestvalue = abtestvalue + value.(string) + ",\n"
								} else {
									abtestvalue = abtestvalue + value.(string) + "\n"
								}
							} else {
								ccommon.CLogger.Error(cw.WkInfo.AppID,"jsoniter.Unmarshal(abtest_value failed, err:", err)
							}
						}
						if path != "" {
							UpdateConsul(cw.WkInfo.AppID, update.Namespace, cw.WkInfo.Cluster, path, "["+abtestvalue+"]")
						}
					} else if strings.Contains(cw.WkInfo.AppID, ccommon.BidForceAppid) {
						var bidforce_valuemap = BidForce{}
						path := ""
						bidforce_value := ""
						for key, value := range update.NewValue {
							v, _ := value.(string)
							skip := false
							ovalue, ok := update.OldValue[key]
							if ok {
								ov, _ := ovalue.(string)
								if ov == v {
									skip =true
								}						
							}
							if ! skip {
								modifier : = GetModifyInfo(ns_info, key)
								updatecontent = fmt.Sprintf("%s\nkey=%s\nold=%s\nnew=%s\nmodifier=%s\n", updatecontent, key, ovalue, value, modifier)
								updatekey = fmt.Sprintf("key=%s,%s\nmodifier=%s\n", updatekey, key, modifier)
							  if modifier != "" {
									modifier_list = append(modifier_list, modifier)
								}								
							}
							if key == "consul_key" {
								path = value.(string)
								continue
							}
							if _, err := toml.Decode(value.(string), &bidforce_valuemap);err == nil {
								bidforce_value = bidforce_value + strings.TrimSpace(value.(string)) + "\n"
							} else {
								ccommon.CLogger.Error(cw.WkInfo.AppID,"toml.Decode(bidforce_value failed, err:", err)
								continue
							}
						}
						if path != "" {
							UpdateConsul(cw.WkInfo.AppID, update.Namespace, cw.WkInfo.Cluster, path, bidforce_value)
						}
					} else {
						for path, value := range update.NewValue {
							v, _ := value.(string)
							ovalue, ok := update.OldValue[path]
							if ok {
								ov, _ := ovalue.(string)
								if ov == v {
									continue
								}
							}
							modifier : = GetModifyInfo(ns_info, key)
							updatecontent = fmt.Sprintf("%s\nkey=%s\nold=%s\nnew=%s\nmodifier=%s\n", updatecontent, key, ovalue, value, modifier)
							updatekey = fmt.Sprintf("%s\nkey=%s  modifier=%s\n", updatekey, key, modifier)
							if modifier != "" {
								modifier_list = append(modifier_list, modifier)
							}								
							UpdateConsul(cw.WkInfo.AppID, update.Namespace, cw.WkInfo.Cluster, path, v) 
						}
					}
					//delete keys
					for k, _ := range update.OldValue {
						if _,ok := update.NewValue[k]; ! ok {
							deleted_keys = fmt.Sprintf("%s,%s,", deleted_keys, k)
						}
					}
					//record updatekey except abtest
					if find := strings.Contains(cw.WkInfo.AppID, ccommon.ABTestAppid); ! find && updatekey != "" {
						updatecontent = updatekey
					}
					if deleted_keys != "" {					
						updatecontent = fmt.Sprintf("%s \ndelelte_key=%s",updatecontent, deleted_keys)
					}
					if len(modifier_list) > 0 {
						ccommon.CLogger.Warn(modifier_list, cw.WkInfo.AppID,"#",cw.WkInfo.Cluster,"#",update.Namespace,": \nupdatecontent:\n",updatecontent)
					} else {
						ccommon.CLogger.Warn(cw.WkInfo.AppID,"#",cw.WkInfo.Cluster,"#",update.Namespace,": \nupdatecontent:\n",updatecontent)
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
