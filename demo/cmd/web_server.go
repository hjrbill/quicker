package main

import (
	"github.com/hjrbill/quicker/demo/handler"
	"github.com/hjrbill/quicker/demo/job"
	"github.com/hjrbill/quicker/index_service"
	"github.com/hjrbill/quicker/pkg/service_hub"
	"os"
	"os/signal"
	"syscall"
)

func WebServerInit(mode int) {
	switch mode {
	case 1:
		standaloneIndexer := new(index_service.Indexer)                        // 单机索引
		if err := standaloneIndexer.Init(50000, dbType, *dbPath); err != nil { //初始化索引
			panic(err)
		}
		if *rebuildIndex {
			job.BuildIndexFromFile(csvFile, standaloneIndexer, 0, 0) // 重建索引
		} else {
			standaloneIndexer.LoadFromForwardIndexFile() // 直接从正排索引文件里加载
		}
		handler.Indexer = standaloneIndexer
	case 3:
		handler.Indexer = index_service.NewSentinel(etcdServers, &service_hub.RoundRobin{}, 100)
	default:
		panic("invalid mode")
	}
}

func WebServerTeardown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	handler.Indexer.Close() // 接收到 kill 信号时关闭索引
	os.Exit(0)              // 然后自杀
}

func WebServerMain(mode int) {
	go WebServerTeardown()
	WebServerInit(mode)
}
