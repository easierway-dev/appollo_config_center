package cworker

import (
	"os"	

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

// AddWorker 添加 workder
func (cw *CWorker) AddWorker(worker CWorker) {
	cw.workers = append(cw.workers, worker)
}

func (cw *CWorker) Setup(appid, cluster string)(*CWorker,error){
	newAgo, err := agollo.New(
		ccommon.AgolloConfiger.ConfigServerURL,
		AppID,
		agollo.Cluster(cluster),
		agollo.PreloadNamespaces(ccommon.DyAgolloConfiger.AppConfig.Namespace),
		agollo.AutoFetchOnCacheMiss(),
		agollo.FailTolerantOnBackupExists(),
		agollo.WithLogger(agollo.NewLogger(agollo.LoggerWriter(os.Stdout))),
	)
	if err != nil {
		return nil, err
	}
	work := CWorker{
		AgolloClient:  newAgo,
		Cluster:        cluster,
		AppID:    appid,
	}
	return work, nil
}

func (cw *CWorker) Run(worker *CWorker){
	errorCh := worker.AgolloClient.Start()
	watchCh := worker.AgolloClient.Watch()
	go func(worker Worker) {
		for {
			select {
			case <-s.ctx.Done():
				ccommon.CLogger.Runtime.Infof(worker.Cluster, "watch quit...")
				return
			case err := <-errorCh:
				ccommon.CLogger.Runtime.Errorf("Error:", err)
			case update := <-watchCh:
				for path, value := range update.NewValue {
					v, _ := value.(string)
					err := cconsul.WriteOne(ccommon.DyAgolloConfiger.ClusterConfig.ClusterMap[worker.ClusterID].ConsulAddr, path, v)
					if err != nil {
						ccommon.CLogger.Runtime.Errorf("consul_addr[%s], err[%v]\n", ccommon.DyAgolloConfiger.ClusterConfig.ClusterMap[worker.ClusterID].ConsulAddr, err)
					}
				}
				ccommon.CLogger.Runtime.Infof("Apollo cluster(%s) namespace(%s) old_value:(%v) new_value:(%v) error:(%v)\n",
					worker.Cluster, update.Namespace, update.OldValue, update.NewValue, update.Error)
			}
		}
		//			s.wg.Done()
	}(worker)
}

func (cw *CWorker) Stop(worker *CWorker){
	worker.AgolloClient.Stop()
}
