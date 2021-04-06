package cserver

import (
	"context"
	"sync"

	"github.com/shima-park/agollo"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/cconsul"
)

// NewAgolloServer 创建一个新的 AgolloServer
func NewAgolloServer() *AgolloServer {
	s := &AgolloServer{}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	return s
}

// Worker 工作者接口
type Worker struct {
	AgolloClient agollo.Agollo
	Cluster      string
	AppID	string
}



// AgolloServer server 服务
type AgolloServer struct {
	regworkers sync.Map //map[string]Worker
	runningworkers sync.Map ///map[string]Worker

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}


// 根据globla_config.app_cluster_map注册worker
func (s *AgolloServer) update(worker Worker) {
	s.regworkers.Store(fmt.Sprintf("%s_%s", worker.AppID, worker.Cluster),worker)
}

func (s *AgolloServer) watch() {

        s.regworkers.Store(fmt.Sprintf("%s_%s", worker.AppID, worker.Cluster),worker)
}

// Run 运行 server
func (s *AgolloServer) Run() {
	for _, worker := range s.regworkers {
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
		s.wg.Add(1)
	}
	s.wg.Wait()
}

// GracefulStop 优雅退出
func (s *AgolloServer) GracefulStop() {
	s.cancel()
	s.wg.Wait()
}
