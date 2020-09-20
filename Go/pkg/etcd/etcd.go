package etcd

import (
	"context"
	"eduhacks2020/Go/pkg/setting"
	"github.com/coreos/etcd/clientv3"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var etcdKvClient *clientv3.Client
var mu sync.Mutex

// GetInstance 获取实例
func GetInstance() *clientv3.Client {
	if etcdKvClient == nil {
		client, err := clientv3.New(clientv3.Config{
			Endpoints:   setting.EtcdSetting.Endpoints,
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			log.Error(err)
			return nil
		}
		//创建时才加锁
		mu.Lock()
		defer mu.Unlock()
		etcdKvClient = client
		return etcdKvClient

	}
	return etcdKvClient
}

// Put 置参数
func Put(key, value string) error {
	_, err := GetInstance().Put(context.Background(), key, value)
	return err
}

// Get 获取参数
func Get(key string) (resp *clientv3.GetResponse, err error) {
	resp, err = GetInstance().Get(context.Background(), key)
	return resp, err
}
