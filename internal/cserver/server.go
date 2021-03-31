package cserver

import (
	"context"
	"sync"
	"github.com/shima-park/agollo"
)

// NewAgolloServer 创建一个新的 AgolloServer
func NewAgolloServer() *AgolloServer {
	s := &AgolloServer{}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	return s
}

// Worker 工作者接口
type Worker interface {
	Agollo_client agollo.Agollo
	Cluster string
	Consul_addr string
}

// AgolloServer server 服务
type AgolloServer struct {
	workers []Worker

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// AddWorker 添加 workder
func (s *AgolloServer) AddWorker(worker Worker) {
	s.workers = append(s.workers, worker)
}

// Run 运行 server
func (s *AgolloServer) Run() {
	for _, worker := range s.workers {
		errorCh := worker.Agollo_client.Start()
		watchCh := worker.Agollo_client.Watch()
		go func(worker Worker) {
			for {
				select {
				case <-ctx.Done():
				    ccommon.CLoger.Runtime.Infof(worker.Cluster, "watch quit...")
				    return
				case err := <-errorCh:
					ccommon.CLoger.Runtime.Errorf("Error:", err)
				case update := <-watchCh:
					for path, value := range update.NewValue {
						v, _ := value.(string)
						err := common.writeOne(worker.Consul_addr, path, v)
						if err != nil {
							ccommon.CLoger.Runtime.Errorf("consul_addr[%s], err[%v]\n", worker.Consul_addr, err)
						}
					}
					ccommon.CLoger.Runtime.Infof("Apollo cluster(%s) namespace(%s) old_value:(%v) new_value:(%v) error:(%v)\n",
						worker.Cluster, update.Namespace, update.OldValue, update.NewValue, update.Error)
				}
			}
			s.wg.Done()
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

