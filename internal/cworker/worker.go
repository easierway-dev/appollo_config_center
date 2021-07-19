package cworker

import (
	"fmt"
	"sort"
	"context"

        "github.com/shima-park/agollo"
        "gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
        "gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/cconsul"
	"gitlab.mobvista.com/voyager/abtesting"
	jsoniter "github.com/json-iterator/go"
)

const (
	ABTest = "abtesting"
	DefaultNamespace = "application"

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

func UpdateConsul(namespace, cluster, key, value string){
	if ccommon.DyAgolloConfiger != nil {
		if _,ok := ccommon.DyAgolloConfiger[namespace];!ok {
			namespace = DefaultNamespace
		}
		if _,ok := ccommon.DyAgolloConfiger[namespace];ok {
			if ccommon.DyAgolloConfiger[namespace].ClusterConfig != nil && ccommon.DyAgolloConfiger[namespace].ClusterConfig.ClusterMap != nil {
				if _,ok := ccommon.DyAgolloConfiger[namespace].ClusterConfig.ClusterMap[cluster];ok {
					consulAddr := ccommon.DyAgolloConfiger[namespace].ClusterConfig.ClusterMap[cluster].ConsulAddr
					if value == "" {
						ccommon.CLogger.Warnf("value is nil !!! consul_addr[",consulAddr,"],key[",key,"]\n")
						return
					}
					err := cconsul.WriteOne(consulAddr, key, value)
					if err != nil {
						ccommon.CLogger.Errorf("consul_addr[",consulAddr,"],key[",key,"], err[", err,"]\n")
					}
				} else {
					ccommon.CLogger.Warnf("cluster:",cluster,"not in  ccommon.DyAgolloConfiger[",namespace,"].ClusterConfig")
					return
				}
			} else {
				ccommon.CLogger.Warnf("consulAddr get failed ccommon.DyAgolloConfiger[",namespace,"=",ccommon.DyAgolloConfiger[namespace])
				return
			}
		} else {
			ccommon.CLogger.Warnf(namespace," not in ccommon.DyAgolloConfiger[",ccommon.DyAgolloConfiger,"]")
			return
		}
	} else {
		ccommon.CLogger.Warnf("ccommon.DyAgolloConfiger = nil")
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
				ccommon.CLogger.Infof(cw.WkInfo.Cluster, "watch quit...")
				return
			case err := <-errorCh:
				ccommon.CLogger.Warnf("Error:", err)
			case update := <-watchCh:
				skipped_keys := "iamstart"
				if update.Namespace == ABTest {
					abtest_valuelist := make([]*abtesting.AbInfo,0)
					path := ""
					for key, value := range update.NewValue {
						if key == "consul_key" {
							path = value.(string)
							continue
						}
						var abtest_value abtesting.AbInfo
						err := jsoniter.Unmarshal([]byte(value.(string)), &abtest_value)
						if err == nil {
							abtest_valuelist = append(abtest_valuelist, &abtest_value)
						} else {
							ccommon.CLogger.Errorf("jsoniter.Unmarshal(abtest_value failed, err:", err)
						}
					}
					if path != "" {
						v, err := jsoniter.Marshal(abtest_valuelist)
						if err != nil {
							ccommon.CLogger.Errorf("jsoniter.Marshal(abtest_valuelist) failed, err:", err)
						} else {
							UpdateConsul(update.Namespace, cw.WkInfo.Cluster, path, string(v))
						}
					}
				} else {
					for path, value := range update.NewValue {
						v, _ := value.(string)
						if ovalue, ok := update.OldValue[path]; ok {
							ov, _ := ovalue.(string)
							if ov == v {
								skipped_keys = fmt.Sprintf("%s,%s", skipped_keys, path)
								continue
							}
						}
						UpdateConsul(update.Namespace, cw.WkInfo.Cluster, path, v) 
					}
				}
				ccommon.CLogger.Infof("Apollo cluster(",cw.WkInfo.Cluster,") namespace(",update.Namespace,") old_value:("update.OldValue,") new_value:(",update.NewValue,") skipped_keys:[",skipped_keys,"] error:(",update.Error,")\n")
			}
		}
	}(cw)
}

//work stop
func (cw *CWorker) Stop(){
	cw.AgolloClient.Stop()
}
