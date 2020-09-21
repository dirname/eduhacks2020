package database

import (
	"eduhacks2020/Go/pkg/setting"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

// ORM 包含一个 *gorm.DB 对象
type ORM struct {
	DB *gorm.DB
}

// Close 关闭 ORM 的连接
func (o *ORM) Close() error {
	sqlDB, err := o.DB.DB()
	if err != nil {
		log.Error(err.Error())
	}
	return sqlDB.Close()
}

// Init 初始化 ORM 的连接
func (o *ORM) Init() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai", setting.Database.PgsqlHost, setting.Database.PgsqlUser, setting.Database.PgsqlPwd, setting.Database.PgsqlDB)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error(err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Error(err.Error())
	}
	// 设置连接池的配置
	sqlDB.SetMaxIdleConns(10)           // 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxOpenConns(100)          // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour) //设置连接可复用的最大时间
	o.DB = db
}
