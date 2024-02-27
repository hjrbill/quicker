package service_hub

import etcd "go.etcd.io/etcd/client/v3"

const SERVICE_ROOT_PATH = "/quicker/index/"

type IHub interface {
	Close()                                                                                   // 关闭，释放占用资源
	Register(serviceName string, endpoint string, leaseID etcd.LeaseID) (etcd.LeaseID, error) // 服务注册
	UnRegister(serviceName string, endpoint string) error                                     // 服务注销
	GetEndpoints(serviceName string) []string                                                 // 获取服务节点列表（自行实现节点选择）
	GetEndpoint(serviceName string) string                                                    // 获取服务节点（可调用我方提供负载均衡或自行实现）
}
