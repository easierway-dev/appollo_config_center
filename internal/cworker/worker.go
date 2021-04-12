package cworker

import (
	"fmt"
	"os"	
	"context"

        "github.com/shima-park/agollo"
        "gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
        "gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/cconsul"
)

// Worker 工作者接口
type CWorker struct {
        AgolloClient agollo.Agollo
        WkInfo      WorkInfo
}

type WorkInfo struct {
	AppID string
	Cluster string
	Namespace string
	Tag string
}

func (info *WorkInfo) Key() string {
  if info.Tag == "" {
    info.Tag = fmt.Sprintf("%s_%s_%s",info.AppID, info.Cluster, info.Namespace)
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
		agollo.PreloadNamespaces(wInfo.Namespace),
		agollo.AutoFetchOnCacheMiss(),
		agollo.FailTolerantOnBackupExists(),
		agollo.WithLogger(agollo.NewLogger(agollo.LoggerWriter(os.Stdout))),
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
				ccommon.CLogger.Runtime.Errorf("Error:", err)
			case update := <-watchCh:
				for path, value := range update.NewValue {
					v, _ := value.(string)
					err := cconsul.WriteOne(ccommon.DyAgolloConfiger.ClusterConfig.ClusterMap[cw.WkInfo.Cluster].ConsulAddr, path, v)
					if err != nil {
						ccommon.CLogger.Runtime.Errorf("consul_addr[%s], err[%v]\n", ccommon.DyAgolloConfiger.ClusterConfig.ClusterMap[cw.WkInfo.Cluster].ConsulAddr, err)
					}
				}
				ccommon.CLogger.Runtime.Infof("Apollo cluster(%s) namespace(%s) old_value:(%v) new_value:(%v) error:(%v)\n",
					cw.WkInfo.Cluster, update.Namespace, update.OldValue, update.NewValue, update.Error)
			}
		}
	}(cw)
}

//work stop
func (cw *CWorker) Stop(){
	cw.AgolloClient.Stop()
}
