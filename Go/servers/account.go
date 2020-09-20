package servers

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
	SystemId     string `json:"systemId"`
	RegisterTime int64  `json:"registerTime"`
}

var SystemMap sync.Map

func Register(systemId string) (err error) {
	//校验是否为空
	if len(systemId) == 0 {
		return errors.New("系统ID不能为空")
	}

	accountInfo := accountInfo{
		SystemId:     systemId,
		RegisterTime: time.Now().Unix(),
	}

	if utils.IsCluster() {
		//判断是否被注册
		resp, err := etcd.Get(define.EtcdPrefixAccountInfo + systemId)
		if err != nil {
			return err
		}

		if resp.Count > 0 {
			return errors.New("该系统ID已被注册")
		}

		jsonBytes, _ := json.Marshal(accountInfo)

		//注册
		err = etcd.Put(define.EtcdPrefixAccountInfo+systemId, string(jsonBytes))
		if err != nil {
			panic(err)
			return err
		}
	} else {
		if _, ok := SystemMap.Load(systemId); ok {
			return errors.New("该系统ID已被注册")
		}

		SystemMap.Store(systemId, accountInfo)
	}

	return nil
}
