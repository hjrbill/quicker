package test

import (
	qlog "quicker/pkg/log"
	"quicker/pkg/service_hub"
	"testing"
)

var (
	serviceName = "test_service"
	etcdServers = []string{"127.0.0.1:2379"}
)

func TestServiceHub(t *testing.T) {
	hub := service_hub.GetServiceHub(etcdServers, 3, &service_hub.RoundRobin{})
	endpoint := "127.0.0.1:5000"
	_, err := hub.Register(serviceName, endpoint, 0)
	if err != nil {
		qlog.Debugf("register error %v", err)
		t.Fail()
	}
	defer hub.UnRegister(serviceName, endpoint)
	endpoints := hub.GetEndpoints(serviceName)
	qlog.Debugf("endpoints %v\n", endpoints)

	endpoint = "127.0.0.2:5000"
	_, err = hub.Register(serviceName, endpoint, 2024) // 故意给与一个虚假 leaseID
	if err != nil {
		qlog.Debugf("register error %v", err)
		t.Fail()
	}
	defer hub.UnRegister(serviceName, endpoint)
	endpoints = hub.GetEndpoints(serviceName)
	qlog.Debugf("endpoints %v\n", endpoints)

	endpoint = "127.0.0.3:5000"
	_, err = hub.Register(serviceName, endpoint, 0)
	if err != nil {
		qlog.Debugf("register error %v", err)
		t.Fail()
	}
	defer hub.UnRegister(serviceName, endpoint)
	endpoints = hub.GetEndpoints(serviceName)
	qlog.Debugf("endpoints %v\n", endpoints)
}
