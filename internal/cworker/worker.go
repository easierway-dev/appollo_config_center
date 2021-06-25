package cworker

import (
	"fmt"
	"sort"
	"context"

        "github.com/shima-park/agollo"
        "gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
        "gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/cconsul"
)

const (
	ABTest = "abtesting"

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

func UpdateConsul(namespace, cluster, key, value){
	if ccommon.DyAgolloConfiger != nil {
		if _,ok := ccommon.DyAgolloConfiger[namespace];ok {
			if ccommon.DyAgolloConfiger[namespace].ClusterConfig != nil && ccommon.DyAgolloConfiger[namespace].ClusterConfig.ClusterMap != nil {
				if _,ok := ccommon.DyAgolloConfiger[namespace].ClusterConfig.ClusterMap[cluster];ok {
					consulAddr := ccommon.DyAgolloConfiger[namespace].ClusterConfig.ClusterMap[cluster].ConsulAddr
					err := cconsul.WriteOne(consulAddr, key, value)
					if err != nil {
						ccommon.CLogger.Runtime.Errorf("consul_addr[%s], err[%v]\n", consulAddr, err)
					}
				} else {
					ccommon.CLogger.Runtime.Warnf("cluster:%s not in  ccommon.DyAgolloConfiger[%s].ClusterConfig", cluster,namespace)
					return
				}
			} else {
				ccommon.CLogger.Runtime.Warnf("consulAddr get failed ccommon.DyAgolloConfiger[%s]=%v",namespace,ccommon.DyAgolloConfiger[namespace])
				return
			}
		} else {
			ccommon.CLogger.Runtime.Warnf("%s not in ccommon.DyAgolloConfiger[%v]",namespace,ccommon.DyAgolloConfiger)
			return
		}
	} else {
		ccommon.CLogger.Runtime.Warnf("ccommon.DyAgolloConfiger = nil")
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
				ccommon.CLogger.Runtime.Infof(cw.WkInfo.Cluster, "watch quit...")
				return
			case err := <-errorCh:
				ccommon.CLogger.Runtime.Warnf("Error:", err)
			case update := <-watchCh:
				skipped_keys := "iamstart"
				abtest_valuelist = make([]interface)
				if update.Namespace == ABTest {
					path := ""
					for key, value := range update.NewValue {
						if key == "consul_key" {
							path = value
							continue
						}
						abtest_valuelist = append(abtest_valuelist, value)					 
					}
					if path != "" {
						UpdateConsul(update.Namespace, cw.WkInfo.Cluster, path, abtest_valuelist)
					}				
				} else {		
					for path, value := range update.NewValue {
						v, _ := value.(string)
						if ovalue, ok := update.OldValue[path]; ok {
							if ovalue.(string) == v {
								skipped_keys = fmt.Sprintf("%s,%s", skipped_keys, path)
								continue
							}
						}			
						UpdateConsul(update.Namespace, cw.WkInfo.Cluster, path, v) 
					}
				}
				ccommon.CLogger.Runtime.Infof("Apollo cluster(%s) namespace(%s) old_value:(%v) new_value:(%v) skipped_keys:[%s] error:(%v)\n",
					cw.WkInfo.Cluster, update.Namespace, update.OldValue, update.NewValue, skipped_keys, update.Error)
			}
		}
	}(cw)
}

//work stop
func (cw *CWorker) Stop(){
	cw.AgolloClient.Stop()
}
