package etcd

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	log "github.com/sirupsen/logrus"
	"time"
)

// ServiceReg 注册租约服务
type ServiceReg struct {
	client        *clientv3.Client
	lease         clientv3.Lease //租约
	leaseResp     *clientv3.LeaseGrantResponse
	cancelFun     func()
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
}

// NewServiceReg 注册新的服务
func NewServiceReg(addr []string, timeNum int64) (*ServiceReg, error) {
	var (
		err    error
		client *clientv3.Client
	)

	if client, err = clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	}); err != nil {
		return nil, err
	}

	ser := &ServiceReg{
		client: client,
	}

	if err := ser.setLease(timeNum); err != nil {
		return nil, err
	}
	go ser.ListenLeaseRespChan()
	return ser, nil
}

//设置租约
func (r *ServiceReg) setLease(timeNum int64) error {
	lease := clientv3.NewLease(r.client)

	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	leaseResp, err := lease.Grant(ctx, timeNum)
	if err != nil {
		cancel()
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.TODO())
	leaseRespChan, err := lease.KeepAlive(ctx, leaseResp.ID)

	if err != nil {
		return err
	}

	r.lease = lease
	r.leaseResp = leaseResp
	r.cancelFun = cancelFunc
	r.keepAliveChan = leaseRespChan
	return nil
}

// ListenLeaseRespChan 监听续租情况
func (r *ServiceReg) ListenLeaseRespChan() {
	for {
		select {
		case leaseKeepResp := <-r.keepAliveChan:
			if leaseKeepResp == nil {
				log.Error("已经关闭续租功能")
				return
			}
			//log.Info("续租成功")
		}
	}
}

// PutService 注册租约
func (r *ServiceReg) PutService(key, val string) error {
	kv := clientv3.NewKV(r.client)
	_, err := kv.Put(context.TODO(), key, val, clientv3.WithLease(r.leaseResp.ID))
	return err
}

//RevokeLease 撤销租约
func (r *ServiceReg) RevokeLease() error {
	r.cancelFun()
	time.Sleep(2 * time.Second)
	_, err := r.lease.Revoke(context.TODO(), r.leaseResp.ID)
	return err
}
