package index_service

import (
	"context"
	"fmt"
	"github.com/hjrbill/quicker/internal/kvdb"
	"github.com/hjrbill/quicker/pb"
	qlog "github.com/hjrbill/quicker/pkg/log"
	"github.com/hjrbill/quicker/pkg/service_hub"
	"github.com/hjrbill/quicker/pkg/util"
	"strconv"
	"time"
)

const (
	INDEX_SERVICE = "index_service"
)

var _ pb.IndexServiceServer = (*IndexServiceWork)(nil)

// IndexServiceWork 一个 grpc server，是本框架分布式部署的基本单位
type IndexServiceWork struct {
	Indexer *Indexer
	// 服务注册有关的配置
	selfIP string
	hub    service_hub.IHub // 服务注册中心，可以选用基础的 ServiceHub 或其代理 HubProxy
}

// Init 初始化正排索引和倒排索引
func (work *IndexServiceWork) Init(cap int, dbType kvdb.DBType, dbPath string) error {
	work.Indexer = new(Indexer)
	return work.Indexer.Init(cap, dbType, dbPath)
}

// Register
// @param servicePort 服务运行端口
// @param etcdServers etcd 集群地址，如果想使用单例模式，可以传 nil
// @param loadBalancer 负载均衡策略
func (work *IndexServiceWork) Register(servicePort int, etcdServers []string, loadBalancer service_hub.LoadBalancer) error {
	if len(etcdServers) > 0 {
		if servicePort < 1024 {
			return fmt.Errorf("监听端口号 %d 为公认端口，不应使用，请使用大于 1024 的端口", servicePort)
		} else if servicePort > 49152 {
			return fmt.Errorf("监听端口号 %d 为动态端口，不建议使用，建议使用小于 49152 的端口", servicePort)
		}

		selfLocalIp, err := util.GetLocalIPWithHardware() // 获取本机 IP
		if err != nil {
			panic(err)
		}
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

func (work *IndexServiceWork) Close() error {
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

func (work *IndexServiceWork) AddDocument(ctx context.Context, document *pb.Document) (*pb.AffectedCount, error) {
	cnt, err := work.Indexer.AddDoc(*document)
	return &pb.AffectedCount{Count: int32(cnt)}, err
}

func (work *IndexServiceWork) DeleteDocument(ctx context.Context, id *pb.ID) (*pb.AffectedCount, error) {
	cnt, err := work.Indexer.DeleteDoc(id.ID)
	return &pb.AffectedCount{Count: int32(cnt)}, err
}

func (work *IndexServiceWork) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	docs, err := work.Indexer.Search(req.TermQuery, req.OnFlag, req.OffFlag, req.OrFlag)
	return &pb.SearchResponse{Documents: docs}, err
}

func (work *IndexServiceWork) Count(ctx context.Context, req *pb.CountRequest) (*pb.AffectedCount, error) {
	return &pb.AffectedCount{Count: work.Indexer.Count()}, nil
}
