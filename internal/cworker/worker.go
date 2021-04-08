package cworker

import (
	"os"	
	"context"

        "github.com/shima-park/agollo"
        "gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
        "gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/cconsul"
)

// Worker 工作者接口
type CWorker struct {
        AgolloClient agollo.Agollo
        Cluster      string
        AppID   string
}

type AgolloInfo struct {
	AppID string
	Cluster string
	Namespace string
}

// setup workder
func Setup(appid, cluster, namespace string)(*CWorker,error){
	var work *CWorker
	newAgo, err := agollo.New(
		ccommon.AgolloConfiger.ConfigServerURL,
		appid,
		agollo.Cluster(cluster),
		agollo.PreloadNamespaces(namespace),
		agollo.AutoFetchOnCacheMiss(),
		agollo.FailTolerantOnBackupExists(),
		agollo.WithLogger(agollo.NewLogger(agollo.LoggerWriter(os.Stdout))),
	)
	if err != nil {
		return work, err
	}
	work = &CWorker{
		AgolloClient:  newAgo,
		Cluster:        cluster,
		AppID:    appid,
	}
	return work, nil
}

//work run
func (cw *CWorker) Run(worker CWorker, ctx context.Context){
	errorCh := worker.AgolloClient.Start()
	watchCh := worker.AgolloClient.Watch()
	go func(worker CWorker) {
		for {
			select {
			case <-ctx.Done():
				ccommon.CLogger.Runtime.Infof(worker.Cluster, "watch quit...")
				return
			case err := <-errorCh:
				ccommon.CLogger.Runtime.Errorf("Error:", err)
			case update := <-watchCh:
				for path, value := range update.NewValue {
					v, _ := value.(string)
					err := cconsul.WriteOne(ccommon.DyAgolloConfiger.ClusterConfig.ClusterMap[worker.Cluster].ConsulAddr, path, v)
					if err != nil {
						ccommon.CLogger.Runtime.Errorf("consul_addr[%s], err[%v]\n", ccommon.DyAgolloConfiger.ClusterConfig.ClusterMap[worker.Cluster].ConsulAddr, err)
					}
				}
				ccommon.CLogger.Runtime.Infof("Apollo cluster(%s) namespace(%s) old_value:(%v) new_value:(%v) error:(%v)\n",
					worker.Cluster, update.Namespace, update.OldValue, update.NewValue, update.Error)
			}
		}
		//			s.wg.Done()
	}(worker)
}

//work stop
func (cw *CWorker) Stop(worker CWorker){
	worker.AgolloClient.Stop()
}
