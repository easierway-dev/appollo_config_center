package ccompare

import (
	"context"
	"fmt"
	"github.com/shima-park/agollo"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
	"math/rand"
	"os"
	"sync"
)

// AgolloServer server 服务
type AgolloServer struct {
	regworkers     sync.Map //map[string]WorkerInfo
	runningworkers sync.Map ///map[string]Worker

	gAgollo agollo.Agollo
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

// NewAgolloServer 创建一个新的 AgolloServer
func NewAgolloServer() *AgolloServer {
	s := &AgolloServer{}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	return s
}
//func NewAgolloServer(agolloCfg *ccommon.AgolloCfg) (newAgo agollo.Agollo, err error) {
//	newAgo, err = agollo.New(
//		agolloCfg.ConfigServerURL,
//		agolloCfg.AppID,
//		agollo.Cluster(agolloCfg.Cluster),
//		agollo.PreloadNamespaces(agolloCfg.Namespace...),
//		agollo.AutoFetchOnCacheMiss(),
//		agollo.FailTolerantOnBackupExists(),
//		agollo.WithLogger(agollo.NewLogger(agollo.LoggerWriter(os.Stdout))),
//	)
//	if err != nil {
//		fmt.Println("Build_Global_Agollo err: %s\n", err.Error())
//		return nil, err
//	}
//	return
//}

func BuildGlobalAgollo(agolloCfg *ccommon.AgolloCfg, server *AgolloServer) error {
	newAgo, err := agollo.New(
		agolloCfg.ConfigServerURL,
		agolloCfg.AppID,
		agollo.Cluster(agolloCfg.Cluster),
		agollo.PreloadNamespaces(agolloCfg.Namespace...),
		agollo.AutoFetchOnCacheMiss(),
		agollo.FailTolerantOnBackupExists(),
		agollo.WithLogger(agollo.NewLogger(agollo.LoggerWriter(os.Stdout))),
	)
	if err != nil {
		fmt.Println("Build_global_agollo err: %s\n", err.Error())
		return err
	}
	server.BuildGAgollo(newAgo)
	return nil
}

func (s *AgolloServer) BuildGAgollo(agollo agollo.Agollo) {
	s.gAgollo = agollo
}
func Init(server *AgolloServer)  error {
	//init config
	cfg, err := ccommon.ParseBaseConfig(ccommon.DirFlag)
	if err != nil {
		fmt.Println("ParseConfig error: %s\n", err.Error())
		return err
	}
	ccommon.AgolloConfiger =  cfg.AgolloCfg
	ccommon.AppConfiger = cfg.AppCfg
	ccommon.ChklogRamdom = rand.Float64()
	ccommon.ChklogRate = ccommon.AppConfiger.ChklogRate
	// init log
	cl, err := ccommon.NewconfigCenterLogger(cfg.LogCfg)
	if err != nil {
		fmt.Println("Load Logger err: ", err)
		return err
	}
	ccommon.CLogger = cl
	ccommon.CLogger.Info(ccommon.DefaultDingType,"Config=", ccommon.AgolloConfiger)
	ccommon.DyAgolloConfiger = make(map[string]*ccommon.DyAgolloCfg)
	//get global_config
	return BuildGlobalAgollo(ccommon.AgolloConfiger, server)
}