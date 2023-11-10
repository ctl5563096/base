package library

import (
	"time"

	"go.uber.org/zap"
)

type HttpClientConfig struct {
	Name                      string //名称
	RequestTimeoutSecond      int    //请求超时时间
	DialTimeoutSecond         int    //连接超时
	DialKeepAliveSecond       int    //开启长连接
	MaxIdleConnections        int    //最大空闲连接数
	MaxIdleConnectionsPerHost int    //单Host最大空闲连接数
	IdleConnTimeoutSecond     int    // 空闲连接超时
	EnableSkyWalking          bool   // 是否开链路追踪
}

type HttpClientCacheConfig struct {
	Name                      string // 名称
	RequestTimeoutSecond      int    // 请求超时时间
	MaxIdleConnections        int    // 最大空闲连接数
	MaxIdleConnectionsPerHost int    // 单Host最大空闲连接数
	IdleConnTimeoutSecond     int    // 空闲连接超时
	EnableSkyWalking          bool   // 是否开链路追踪
	DialConf                  DialConfig
}

type DialConfig struct {
	Logger              *zap.Logger
	DialTimeoutSecond   int // 连接超时
	DialKeepAliveSecond int // 开启长连接
	DnsCacheNums        int
	DnsCacheTime        time.Duration
}
