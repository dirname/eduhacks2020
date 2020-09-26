package utils

import (
	"eduhacks2020/Go/crypto"
	"eduhacks2020/Go/pkg/setting"
	"errors"
	uuid "github.com/satori/go.uuid"
	"strings"
)

// GenUUID 生成uuid
func GenUUID() string {
	uuidFunc := uuid.NewV4()
	uuidStr := uuidFunc.String()
	uuidStr = strings.Replace(uuidStr, "-", "", -1)
	uuidByt := []rune(uuidStr)
	return string(uuidByt[8:24])
}

// GenClientID 对称加密IP和端口，当做clientId
func GenClientID() string {
	raw := []byte(setting.GlobalSetting.LocalHost + ":" + setting.CommonSetting.RPCPort)
	str, err := crypto.Encrypt(raw, []byte(setting.CommonSetting.CryptoKey))
	if err != nil {
		panic(err)
	}

	return str
}

// ParseRedisAddrValue 解析redis的地址格式
func ParseRedisAddrValue(redisValue string) (host string, port string, err error) {
	if redisValue == "" {
		err = errors.New("parsing address error")
		return
	}
	addr := strings.Split(redisValue, ":")
	if len(addr) != 2 {
		err = errors.New("parsing address error")
		return
	}
	host, port = addr[0], addr[1]

	return
}

// IsAddrLocal 判断地址是否为本机
func IsAddrLocal(host string, port string) bool {
	return host == setting.GlobalSetting.LocalHost && port == setting.CommonSetting.RPCPort
}

// IsCluster 是否集群
func IsCluster() bool {
	return setting.CommonSetting.Cluster
}

// GetAddrInfoAndIsLocal 获取client key地址信息
func GetAddrInfoAndIsLocal(clientID string) (addr string, host string, port string, isLocal bool, err error) {
	//解密ClientId
	addr, err = crypto.Decrypt(clientID, []byte(setting.CommonSetting.CryptoKey))
	if err != nil {
		return
	}

	host, port, err = ParseRedisAddrValue(addr)
	if err != nil {
		return
	}

	isLocal = IsAddrLocal(host, port)
	return
}

// GenGroupKey 生成群租密匙
func GenGroupKey(systemID, groupName string) string {
	return systemID + ":" + groupName
}
