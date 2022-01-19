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
			enUpdate := false
			if dyAgoCfg.AppConfig != nil {
				enUpdate = dyAgoCfg.AppConfig.EnUpdateConsul
				if dyAgoCfg.AppConfig.AppConfigMap != nil {
					if _,ok := dyAgoCfg.AppConfig.AppConfigMap[appid];ok{
						enUpdate = dyAgoCfg.AppConfig.AppConfigMap[appid].EnUpdateConsul
					}
				}
			}
			if !enUpdate {
				ccommon.CLogger.Warn(appid, "is not permit to update consul")
				return
			}
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
				skipped_keys := ""
				if strings.Contains(cw.WkInfo.AppID, ccommon.ABTestAppid) {
					path := ""
					abtestvalue := ""
					i := 0
					abUpdate := true
					for key, value := range update.NewValue {
						i = i + 1
						v, _ := value.(string)
						if ovalue, ok := update.OldValue[key]; ok {
							ov, _ := ovalue.(string)
							if ov == v {
								skipped_keys = fmt.Sprintf("%s,%s,", skipped_keys, key)
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
							abUpdate = false
							ccommon.CLogger.Error(cw.WkInfo.AppID,"_",cw.WkInfo.Cluster,"_",key,":", "\njsoniter.Unmarshal(abtest_value failed, err:", err)
						}
					}
					if path != "" && abUpdate {
						UpdateConsul(cw.WkInfo.AppID, update.Namespace, cw.WkInfo.Cluster, path, "["+strings.Trim(strings.Trim(abtestvalue, "\n"),",")+"]")
					}
				} else if strings.Contains(cw.WkInfo.AppID, ccommon.BidForceAppid) {
					var bidforce_valuemap = BidForce{}
					path := ""
					bidforce_value := ""
					for key, value := range update.NewValue {
						v, _ := value.(string)
						if ovalue, ok := update.OldValue[key]; ok {
							ov, _ := ovalue.(string)
							if ov == v {
								skipped_keys = fmt.Sprintf("%s,%s,", skipped_keys, key)
							}
						}
						if key == "consul_key" {
							path = value.(string)
							continue
						}
						if _, err := toml.Decode(value.(string), &bidforce_valuemap);err == nil {
							bidforce_value = bidforce_value + strings.TrimSpace(value.(string)) + "\n"
						} else {
							ccommon.CLogger.Error(cw.WkInfo.AppID,"_",cw.WkInfo.Cluster,"_",key,":", "\ntoml.Decode(bidforce_value failed, err:", err)
							continue
						}
					}
					if path != "" {
						UpdateConsul(cw.WkInfo.AppID, update.Namespace, cw.WkInfo.Cluster, path, bidforce_value)
					}
				} else {
					for path, value := range update.NewValue {
						v, _ := value.(string)
						if ovalue, ok := update.OldValue[path]; ok {
							ov, _ := ovalue.(string)
							if ov == v {
								skipped_keys = fmt.Sprintf("%s,%s,", skipped_keys, path)
								continue
							}
						}
						UpdateConsul(cw.WkInfo.AppID, update.Namespace, cw.WkInfo.Cluster, path, v) 
					}
				}
				updatecontent := ""
				updatekey := ""		
				if len(update.NewValue) == 0 {
					updatecontent = fmt.Sprintf("clear_config or create_namesplace:%s[%s]",cw.WkInfo.Cluster,update.Namespace)
				}
				for k, v := range update.NewValue {
					if ! strings.Contains(skipped_keys, fmt.Sprintf(",%s,", k)) {
						if _,ok := update.OldValue[k]; ok{
							updatecontent = fmt.Sprintf("%s\nkey=%s\nold=%s\nnew=%s", updatecontent, k, update.OldValue[k], v)
						} else {
							updatecontent = fmt.Sprintf("%s\nkey=%s\nold=%s\nnew=%s", updatecontent, k, "", v)
						}
						updatekey = fmt.Sprintf("key=%s\n", k)
					}
				}
				if updatecontent == "" {
					updatecontent = fmt.Sprintf("delelte key=%s",skipped_keys)
				}
				if find := strings.Contains(cw.WkInfo.AppID, ccommon.ABTestAppid); ! find && updatekey != "" {
					updatecontent = updatekey
				}
				ccommon.CLogger.Info(ccommon.DefaultDingType,"Apollo cluster(",cw.WkInfo.Cluster,") namespace(",update.Namespace,") \nold_value:(", update.OldValue,") \nnew_value:(",update.NewValue,") \nskipped_keys:[",skipped_keys,"] error:(",update.Error,")\n")
				ccommon.CLogger.Warn(cw.WkInfo.AppID,"#",cw.WkInfo.Cluster,"#",update.Namespace,": \nupdatecontent:\n",updatecontent)
			}
		}
	}(cw)
}

//work stop
func (cw *CWorker) Stop(){
	cw.AgolloClient.Stop()
}
