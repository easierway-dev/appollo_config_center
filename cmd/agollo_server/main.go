package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gitlab.mobvista.com/mvbjqa/appollo_config_center/internal/cserver"
	"net/http"
	_ "net/http/pprof"
)

// AppVersion 版本信息
var AppVersion = "unknown"
var PprofPort *string

func main() {
	version := flag.Bool("v", false, "print current version")
	PprofPort = flag.String("p", "6666", "监控端口")
	flag.Parse()
	if *version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}
	go func() {
		http.ListenAndServe(fmt.Sprintf("localhost:%v", *PprofPort), nil)
	}()

	var server *cserver.AgolloServer
	var err error
	if server, err = cserver.Init(); err != nil {
		panic(err)
	}

	go server.Run()
	fmt.Println("agollo_server start success !!! will listen appolo update ...")
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	server.GracefulStop()
}
