package main

import (
	"flag"
	"fmt"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccompare"
	"os"
	"os/signal"
	"syscall"

	"net/http"
	_ "net/http/pprof"

	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/ccommon"
	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/cserver"
)

// AppVersion 版本信息
var AppVersion = "unknown"

func handleKillSignal() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	<-sigchan
	ccommon.CLogger.Info(ccommon.InitDingType,"get shutdown signal.")
	os.Exit(0)
}

func main() {
	version := flag.Bool("v", false, "print current version")
	flag.Parse()
	if *version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}
	go func() {
		err := http.ListenAndServe(":8686", nil)
		if err != nil {
			panic(err)
		}
	}()

	var server *cserver.AgolloServer
	var err error
	server = cserver.NewAgolloServer()
	if err = cserver.Init(server); err != nil {
		panic(err)
	}
	server.Run()
	fmt.Println("agollo_server start success !!! will listen appolo update ...")
	ccommon.CLogger.Info(ccommon.InitDingType,"agollo_server start success !!! will listen appolo update ...")
	ccompare.ApolloCompareWithConsul()
	handleKillSignal()
	server.GracefulStop()
}
