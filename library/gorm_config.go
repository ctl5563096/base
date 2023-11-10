package library

import (
	"gorm.io/gorm"
)

type GormConfig struct {
	Receiver              **GormDB
	ConnectionName        string            //连接名称
	DBName                string            //db名称
	Host                  string            //地址
	Port                  string            //端口
	UserName              string            //用户名
	Password              string            //密码
	Charset               string            //字符集
	ParseTime             string            //解析时间
	Loc                   string            //时区
	MaxLifeTime           int               //空闲连接最大保持时长(秒)
	MaxOpenConn           int               //最大打开连接数
	MaxIdleConn           int               //最大空闲连接数
	EnableTransaction     bool              //是否开启事务
	ReadOnlySlavesConfigs []GormSlaveConfig //只读实例配置
	GormDetailConfig      *gorm.Config      // 配置详情
	EnableSkyWalking      bool              // 是否开启链路追终
}

type GormSlaveConfig struct {
	Host     string
	Port     string
	UserName string
	Password string
}
