package index_service

import (
	"context"
	"fmt"
	pb "github.com/hjrbill/quicker/gen"
	"github.com/hjrbill/quicker/internal/kvdb"
	qlog "github.com/hjrbill/quicker/pkg/log"
	"github.com/hjrbill/quicker/pkg/net"
	"github.com/hjrbill/quicker/pkg/service_hub"
	"strconv"
	"time"
)

const (
	INDEX_SERVICE = "index_service" // 索引服务 worker 的服务名
)

var _ pb.IndexServiceServer = (*IndexServiceWorker)(nil)

// IndexServiceWorker 一个 grpc server，是本框架分布式部署的基本单位
type IndexServiceWorker struct {
	Indexer *Indexer
	// 服务注册有关的配置
	selfIP string
	hub    service_hub.IHub // 服务注册中心，可以选用基础的 ServiceHub 或其代理 HubProxy
}

// Init 初始化正排索引和倒排索引
func (work *IndexServiceWorker) Init(cap int, dbType kvdb.DBType, dbPath string) error {
	work.Indexer = new(Indexer)
	return work.Indexer.Init(cap, dbType, dbPath)
}

// Register
// @param servicePort 服务运行端口
// @param etcdServers etcd 集群地址，如果想使用单例模式，可以传 nil
// @param loadBalancer 负载均衡策略
func (work *IndexServiceWorker) Register(servicePort int, etcdServers []string, loadBalancer service_hub.LoadBalancer) error {
	if len(etcdServers) > 0 {
		if servicePort < 1024 {
			return fmt.Errorf("监听端口号 %d 为公认端口，不应使用，请使用大于 1024 的端口", servicePort)
		} else if servicePort > 49152 {
			return fmt.Errorf("监听端口号 %d 为动态端口，不建议使用，建议使用小于 49152 的端口", servicePort)
		}

		selfLocalIp, err := net.GetLocalIPWithHardware() // 获取本机 IP
		if err != nil {
			panic(err)
		}
		selfLocalIp = "127.0.0.1" // TODO 单机模拟分布式时，把 selfLocalIp 写死为 127.0.0.1
		work.selfIP = selfLocalIp + ":" + strconv.Itoa(servicePort)

		heartbeat := int64(3)
		work.hub = service_hub.GetServiceHub(etcdServers, heartbeat, loadBalancer)
		leaseID, err := work.hub.Register(INDEX_SERVICE, work.selfIP, 0)
		if err != nil {
			return err
		}
		go func() {
			for {
				leaseID, err = work.hub.Register(INDEX_SERVICE, work.selfIP, leaseID)
				if err != nil {
					qlog.Warnf("续约 index service 失败, err: %v", err)
				}
				time.Sleep(time.Duration(heartbeat)*time.Second - 100*time.Millisecond) // 提前心跳截止 100 毫秒进行注册，留出空余
			}
		}()
	}
	return nil
}

func (work *IndexServiceWorker) Close() error {
	if work.Indexer != nil {
		err := work.Indexer.Close()
		if err != nil {
			return err
		}
	}
	if work.hub != nil {
		err := work.hub.UnRegister(INDEX_SERVICE, work.selfIP)
		if err != nil {
			return err
		}
	}
	return nil
}

func (work *IndexServiceWorker) AddDocument(ctx context.Context, document *pb.Document) (*pb.AffectedCount, error) {
	cnt, err := work.Indexer.AddDoc(*document)
	return &pb.AffectedCount{Count: int32(cnt)}, err
}

func (work *IndexServiceWorker) DeleteDocument(ctx context.Context, id *pb.ID) (*pb.AffectedCount, error) {
	cnt, err := work.Indexer.DeleteDoc(id.ID)
	return &pb.AffectedCount{Count: int32(cnt)}, err
}

func (work *IndexServiceWorker) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	docs, err := work.Indexer.Search(req.TermQuery, req.OnFlag, req.OffFlag, req.OrFlag)
	return &pb.SearchResponse{Documents: docs}, err
}

func (work *IndexServiceWorker) Count(ctx context.Context, req *pb.CountRequest) (*pb.AffectedCount, error) {
	return &pb.AffectedCount{Count: work.Indexer.Count()}, nil
}
