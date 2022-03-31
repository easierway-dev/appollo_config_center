package main

import (
	"fmt"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccompare"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"
)

// AppVersion 版本信息
var AppVersion = "unknown"

func handleKillSignal() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	<-sigchan
	//ccommon.CLogger.Info(ccommon.InitDingType,"get shutdown signal.")
	os.Exit(0)
}

func main() {
	var server *ccompare.AgolloServer
	var err error
	server = ccompare.NewAgolloServer()
	if err = ccompare.Init(server); err != nil {
		panic(err)
	}
	go ccompare.Start(server)
	ticker := time.NewTicker(time.Minute * 10)
	for {
		select {
		case <-ticker.C:
			fmt.Println("10分钟到....")
			ccompare.Start(server)
		default:
		}
	}
	handleKillSignal()
	//server.GracefulStop()
}
