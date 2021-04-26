package cserver

import (
	"time"
	"context"
	"sync"

	"github.com/shima-park/agollo"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/cworker"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
)

// NewAgolloServer 创建一个新的 AgolloServer
func NewAgolloServer() *AgolloServer {
	s := &AgolloServer{}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	return s
}

// AgolloServer server 服务
type AgolloServer struct {
	regworkers sync.Map //map[string]WorkerInfo
	runningworkers sync.Map ///map[string]Worker

	gAgollo agollo.Agollo
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func (s *AgolloServer) UpdateOne(cfg *ccommon.AppClusterCfg){
	//clear regworkers
	s.regworkers.Range(func(k, v interface{}) bool {
		if k.(string) != "" {
			s.regworkers.Delete(k)
		}
		return true
	})
	namespace := cfg.Namespace
	for appid, appclusterinfo := range cfg.AppClusterMap {
		if len(appclusterinfo.Namespace) != 0 {
			namespace = appclusterinfo.Namespace
		}
		for _, cluster := range appclusterinfo.Cluster {
			wInfo := cworker.WorkInfo{
				AppID : appid,
				Cluster : cluster,
				Namespace : namespace,
			}
			key := wInfo.Key()
			//store regworkers
		    	s.regworkers.Store(key,wInfo)
		}
	}
}

func (s *AgolloServer) BuildGAgollo (agollo agollo.Agollo){
	s.gAgollo = agollo
}

// 根据globla_config.app_cluster_map注册worker
func (s *AgolloServer) Update() {
	for _, ns := range ccommon.AgolloConfiger.Namespace {
		dycfg, err := ccommon.ParseDyConfig(s.gAgollo.Get("cluster_map", agollo.WithNamespace(ns)),s.gAgollo.Get("app_config_map", agollo.WithNamespace(ns)))
		if err != nil {
				ccommon.CLogger.Runtime.Errorf("ParseDyConfig error: %s\n", err.Error())
				panic(err)
		}
		ccommon.DyAgolloConfiger[ns] = dycfg

		cfg, err := ccommon.ParseAppClusterConfig(s.gAgollo.Get("app_cluster_map", agollo.WithNamespace(ns)))
		if err != nil {
				ccommon.CLogger.Runtime.Errorf("ParseAppClusterConfig error: %s\n", err.Error())
				panic(err)
		}	
		s.UpdateOne(cfg)
	}
	
	errorCh := s.gAgollo.Start()
	watchCh := s.gAgollo.Watch()

	go func(cluster string) {
		for {
			select {
			case <-s.ctx.Done():
			    ccommon.CLogger.Runtime.Errorf(cluster, "watch quit...")
			    return
			case err := <-errorCh:
				 ccommon.CLogger.Runtime.Warnf("Error:", err)
			case update := <-watchCh:
				clusterCfg := ""
				appCfg := ""
				if value, ok := update.NewValue["cluster_map"]; ok {
					clusterCfg = value.(string)
				}
				if value, ok := update.NewValue["app_config_map"]; ok {
					appCfg = value.(string)
				}
				dycfg, err := ccommon.ParseDyConfig(clusterCfg, appCfg)
				if err != nil {
						ccommon.CLogger.Runtime.Errorf("update ParseDyConfig error: %s\n", err.Error())
				} else {
					ccommon.DyAgolloConfiger[update.Namespace] = dycfg
				}
				if value, ok := update.NewValue["app_cluster_map"]; ok {
					cfg, err := ccommon.ParseAppClusterConfig(value.(string))
					if err != nil {
							ccommon.CLogger.Runtime.Errorf("update ParseAppClusterConfig error: %s\n", err.Error())
							panic(err)
					} else {
						s.UpdateOne(cfg)
					}
				}
				ccommon.CLogger.Runtime.Infof("Global Apollo cluster(%s) namespace(%s) old_value:(%v) new_value:(%v) error:(%v)\n",
					cluster, update.Namespace, update.OldValue, update.NewValue, update.Error)
			}
		}
	}(ccommon.AgolloConfiger.Cluster)
}

func (s *AgolloServer) Watch() {
	t := time.NewTicker(time.Duration(ccommon.AgolloConfiger.CyclePeriod)*time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			//ccommon.CLogger.Runtime.Infof("I am alive and watch change")
			//start
			s.regworkers.Range(func(k, v interface{}) bool {
				if _,ok := s.runningworkers.Load(k); !ok {
					worker,err := cworker.Setup(v.(cworker.WorkInfo))
					if err == nil {
						worker.Run(s.ctx)
						ccommon.CLogger.Runtime.Infof("will setup worker: %s", k.(string))
						s.wg.Add(1)
						s.runningworkers.Store(k,worker)
					} else {
						ccommon.CLogger.Runtime.Errorf("creeative worker failed !!! workerInfo=",v)
					}
				}
				return true	
			})
			//stop
			s.runningworkers.Range(func(k, v interface{}) bool {
				if _,ok := s.regworkers.Load(k); !ok {
					v.(*cworker.CWorker).Stop()
					ccommon.CLogger.Runtime.Infof("will stop woker: %s",k.(string), "wait 3s to envalid  !!!")
					s.runningworkers.Delete(k)
				}
				return true
			})			
		}
	}
}

func (s *AgolloServer) Run() {
	go s.Watch()
	s.Update()
	s.wg.Wait()
}

// GracefulStop 优雅退出
func (s *AgolloServer) GracefulStop() {
	s.cancel()
	s.wg.Wait()
}