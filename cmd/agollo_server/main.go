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
	go ccompare.Start()
	ticker := time.NewTicker(time.Minute * 10)
	for {
		select {
		case <-ticker.C:
			fmt.Println("10分钟到....")
			ccompare.Start()
		default:
		}
	}
	handleKillSignal()
	//server.GracefulStop()
}
