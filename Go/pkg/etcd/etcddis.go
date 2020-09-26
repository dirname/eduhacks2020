package etcd

import (
	"context"
	"eduhacks2020/Go/pkg/setting"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	log "github.com/sirupsen/logrus"
	"time"
)

// ClientDis etcd 客户端
type ClientDis struct {
	client *clientv3.Client
}

// NewClientDis 新的客户端
func NewClientDis(addr []string) (*ClientDis, error) {
	conf := clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(conf)
	if err == nil {
		return &ClientDis{
			client: client,
		}, nil
	}
	return nil, err
}

// GetService 获取服务
func (d *ClientDis) GetService(prefix string) ([]string, error) {
	resp, err := d.client.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	adders := d.extractAdders(resp)

	go d.watcher(prefix)
	return adders, nil
}

func (d *ClientDis) watcher(prefix string) {
	rch := d.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for resp := range rch {
		for _, ev := range resp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				d.SetServiceList(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE:
				d.DelServiceList(string(ev.Kv.Key))
			}
		}
	}
}

func (d *ClientDis) extractAdders(resp *clientv3.GetResponse) []string {
	adders := make([]string, 0)
	if resp == nil || resp.Kvs == nil {
		return adders
	}
	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			d.SetServiceList(string(resp.Kvs[i].Key), string(resp.Kvs[i].Value))
			adders = append(adders, string(v))
		}
	}
	return adders
}

// SetServiceList 设置服务列表
func (d *ClientDis) SetServiceList(key, val string) {
	setting.GlobalSetting.ServerListLock.Lock()
	defer setting.GlobalSetting.ServerListLock.Unlock()
	setting.GlobalSetting.ServerList[key] = val
	log.Info("Discovery Service: ", key, " address: ", val)
}

// DelServiceList 下线服务
func (d *ClientDis) DelServiceList(key string) {
	setting.GlobalSetting.ServerListLock.Lock()
	defer setting.GlobalSetting.ServerListLock.Unlock()
	delete(setting.GlobalSetting.ServerList, key)
	log.Println("Service offline: ", key)
}
