package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/hjrbill/quicker/demo/route"
	"github.com/hjrbill/quicker/internal/kvdb"
	"github.com/hjrbill/quicker/pkg/util"
	"strconv"
)

var (
	mode         = flag.Int("mode", 1, "启动哪类服务。1-standalone web server, 2-grpc index server, 3-distributed web server")
	rebuildIndex = flag.Bool("index", false, "server 启动时是否需要重建索引")
	port         = flag.Int("port", 0, "server 的工作端口")
	dbPath       = flag.String("dbPath", "", "正排索引数据的存放路径")
	totalWorkers = flag.Int("totalWorkers", 0, "分布式环境中一共有几台 index worker")
	workerIndex  = flag.Int("workerIndex", 0, "本机是第几台 index worker(从 0 开始编号)")
)

var (
	dbType      = kvdb.BOLT                             //正排索引使用哪种 KV 数据库
	csvFile     = util.RootPath + "data/bili_video.csv" //原始的数据文件，由它来创建索引
	etcdServers = []string{"127.0.0.1:2379"}            //etcd 集群的地址
)

func StartGin() {
	engine := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	route.SetRoute(engine)
	engine.Run("127.0.0.1:" + strconv.Itoa(*port))
}

func main() {
	flag.Parse()

	switch *mode {
	case 1, 3:
		WebServerMain(*mode) //1：单机模式，索引功能嵌套在 web server 内部。3：分布式模式，web server 内持有一个哨兵，通过哨兵去访问各个 grpc index server
		StartGin()
	case 2:
		GrpcIndexerMain() //以 grpc server 的方式启动索引服务
	}
}

// go run ./demo/cmd -mode=1 -index=false -port=5678 -dbPath=temp/local_db/video_bolt
// go run ./demo/cmd -mode=2 -index=false -port=5600 -dbPath=temp/local_db/video_bolt -totalWorkers=2 -workerIndex=0
// go run ./demo/cmd -mode=2 -index=false -port=5601 -dbPath=temp/local_db/video_bolt -totalWorkers=2 -workerIndex=1
// go run ./demo/cmd -mode=3 -index=false -port=5678
