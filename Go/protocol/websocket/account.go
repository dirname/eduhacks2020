package websocket

import (
	"eduhacks2020/Go/define"
	"eduhacks2020/Go/pkg/etcd"
	"eduhacks2020/Go/utils"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

type accountInfo struct {
	SystemID     string `json:"systemId"`
	RegisterTime int64  `json:"registerTime"`
}

// SystemMap 系统 Map
var SystemMap sync.Map

// Register 注册
func Register(systemID string) (err error) {
	//校验是否为空
	if len(systemID) == 0 {
		return errors.New("系统ID不能为空")
	}

	accountInfo := accountInfo{
		SystemID:     systemID,
		RegisterTime: time.Now().Unix(),
	}

	if utils.IsCluster() {
		//判断是否被注册
		resp, err := etcd.Get(define.EtcdPrefixAccountInfo + systemID)
		if err != nil {
			return err
		}

		if resp.Count > 0 {
			return errors.New("该系统ID已被注册")
		}

		jsonBytes, _ := json.Marshal(accountInfo)

		//注册
		err = etcd.Put(define.EtcdPrefixAccountInfo+systemID, string(jsonBytes))
		if err != nil {
			panic(err)
			return err
		}
	} else {
		if _, ok := SystemMap.Load(systemID); ok {
			return errors.New("该系统ID已被注册")
		}

		SystemMap.Store(systemID, accountInfo)
	}

	return nil
}
