package main

import (
	"fmt"
	"github.com/hjrbill/quicker/demo/job"
	"github.com/hjrbill/quicker/index_service"
	"github.com/hjrbill/quicker/pb"
	"github.com/hjrbill/quicker/pkg/service_hub"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var service *index_service.IndexServiceWorker //IndexWorker，是一个 grpc server

func GrpcIndexerInit() {
	// 监听本地端口
	lis, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(*port))
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	service = new(index_service.IndexServiceWorker)
	//初始化索引
	service.Init(50000, dbType, *dbPath+"_part"+strconv.Itoa(*workerIndex))
	if *rebuildIndex {
		log.Printf("totalWorkers=%d, workerIndex=%d", *totalWorkers, *workerIndex)
		job.BuildIndexFromFile(csvFile, service.Indexer, *totalWorkers, *workerIndex) //重建索引
	} else {
		service.Indexer.LoadFromForwardIndexFile() //直接从正排索引文件里加载
	}
	// 注册服务的具体实现
	pb.RegisterIndexServiceServer(server, service)
	// 启动服务
	fmt.Printf("start grpc server on port %d\n", *port)
	//向注册中心注册自己，并周期性续命
	service.Register(*port, etcdServers, &service_hub.RoundRobin{})
	err = server.Serve(lis) //Serve 会一直阻塞，所以放到一个协程里异步执行
	if err != nil {
		service.Close()
		fmt.Printf("start grpc server on port %d failed: %s\n", *port, err)
	}
}

func GrpcIndexerTeardown() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	service.Close() //接收到 kill 信号时关闭索引
	os.Exit(0)      //然后自杀
}

func GrpcIndexerMain() {
	go GrpcIndexerTeardown()
	GrpcIndexerInit()
}
