package database

import (
	"flag"
	"github.com/globalsign/mgo"
	"github.com/go-ini/ini"
	"log"
	"net"
	"sync"
)

type databaseConf struct {
	PgsqlHost string
	PgsqlUser string
	PgsqlPwd  string
	PgsqlDB   string
	MongoHost string
	MongoUser string
	MongoPwd  string
	MongoDB   string
	RedisHost string
	RedisPwd  string
}

// SettingDatabase 数据库的配置
var SettingDatabase = &databaseConf{}

type adminConf struct {
	Username string
	Password string
}

// AdminConf 后台管理员的配置
var AdminConf = &adminConf{}

type commonConf struct {
	Port      string
	RPCPort   string
	Cluster   bool
	CryptoKey string
}

// CommonSetting 通用的配置
var CommonSetting = &commonConf{}

type etcdConf struct {
	Endpoints []string
}

//EtcdSetting Etcd的集群配置
var EtcdSetting = &etcdConf{}

type global struct {
	LocalHost      string //本机内网IP
	ServerList     map[string]string
	ServerListLock sync.RWMutex
}

// GlobalSetting 一些全局性的配置
var GlobalSetting = &global{}

// DialInfo Mongo 的连接信息
var DialInfo = &mgo.DialInfo{}

var cfg *ini.File

// ReadConfigure 读取配置文件
func ReadConfigure() {
	configFile := flag.String("c", "conf/app.ini", "-c conf/app.ini")

	var err error

	cfg, err = ini.Load(*configFile)
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("common", CommonSetting)
	mapTo("etcd", EtcdSetting)
	mapTo("database", SettingDatabase)
	mapTo("admin", AdminConf)

	DialInfo.Addrs = []string{SettingDatabase.MongoHost}
	DialInfo.Source = SettingDatabase.MongoDB
	DialInfo.Username = SettingDatabase.MongoUser
	DialInfo.Password = SettingDatabase.MongoPwd

	GlobalSetting = &global{
		LocalHost:  getIntranetIp(),
		ServerList: make(map[string]string),
	}
}


// DefaultSetting 获取自动设置
func DefaultSetting() {
	CommonSetting = &commonConf{
		Port:      "6000",
		RPCPort:   "7000",
		Cluster:   false,
		CryptoKey: "785744acc225bf22",
	}

	SettingDatabase = &databaseConf{
		PgsqlHost: "go.htdocs.net",
		PgsqlUser: "educator",
		PgsqlPwd:  "EduHacks2020.*",
		PgsqlDB:   "education",
		MongoHost: "go.htdocs.net",
		MongoUser: "educator",
		MongoPwd:  "EduHacks2020.*",
		MongoDB:   "education",
		RedisHost: "go.htdocs.net:6379",
		RedisPwd:  "EduHacks2020.*",
	}

	AdminConf = &adminConf{
		Username: "dirname",
		Password: "admin",
	}

	GlobalSetting = &global{
		LocalHost:  getIntranetIp(),
		ServerList: make(map[string]string),
	}
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}

// getIntranetIp 获取本机内网IP
func getIntranetIp() string {
	adders, _ := net.InterfaceAddrs()

	for _, addr := range adders {
		// 检查ip地址判断是否回环地址
		if inet, ok := addr.(*net.IPNet); ok && !inet.IP.IsLoopback() {
			if inet.IP.To4() != nil {
				return inet.IP.String()
			}
		}
	}

	return ""
}
