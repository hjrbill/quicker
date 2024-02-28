package test

import (
	qlog "github.com/hjrbill/quicker/pkg/log"
	"github.com/hjrbill/quicker/pkg/service_hub"
	"testing"
	"time"
)

func TestHubProxy(t *testing.T) {
	qps := 10 //qps 限制为 10
	proxy := service_hub.GetServiceHubProxy(etcdServers, 3, &service_hub.Random{}, qps)

	endpoint := "127.0.0.1:5000"
	_, err := proxy.Register(serviceName, endpoint, 0)
	if err != nil {
		qlog.Warnf("register error %v", err)
		t.Fail()
	}
	defer proxy.UnRegister(serviceName, endpoint)
	endpoints := proxy.GetEndpoints(serviceName)
	qlog.Debugf("endpoints %v\n", endpoints)

	endpoint = "127.0.0.2:5000"
	_, err = proxy.Register(serviceName, endpoint, 0)
	if err != nil {
		qlog.Warnf("register error %v", err)
		t.Fail()
	}
	defer proxy.UnRegister(serviceName, endpoint)
	endpoints = proxy.GetEndpoints(serviceName)
	qlog.Debugf("endpoints %v\n", endpoints)

	endpoint = "127.0.0.3:5000"
	_, err = proxy.Register(serviceName, endpoint, 0)
	if err != nil {
		qlog.Warnf("register error %v", err)
		t.Fail()
	}
	defer proxy.UnRegister(serviceName, endpoint)
	endpoints = proxy.GetEndpoints(serviceName)
	qlog.Debugf("endpoints %v\n", endpoints)

	time.Sleep(1 * time.Second)  //暂停 1 秒钟，等待令牌桶装满
	for i := 0; i < qps+5; i++ { //桶里面有 10 个令牌，从第 11 次开始就需等待
		endpoints = proxy.GetEndpoints(serviceName)
		qlog.Debugf("%d endpoints %v\n", i, endpoints)
	}

	time.Sleep(1 * time.Second) //暂停 1 秒钟，等待令牌桶装满
	go func() {
		for i := 0; i < qps+5; i++ { // 桶里面有 10 个令牌
			endpoints = proxy.GetEndpoints(serviceName)
			qlog.Debugf("并发一 %d endpoints %v\n", i, endpoints)
		}
	}()
	for i := 0; i < qps+5; i++ { // 桶里面有 10 个令牌，但与并发一一同消费，应该会出现被拒绝访问的情况
		endpoints = proxy.GetEndpoints(serviceName)
		qlog.Debugf("并发二 %d endpoints %v\n", i, endpoints)
	}
}
