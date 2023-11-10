package library

import (
	es6 "github.com/olivere/elastic/v6"
)

type ElasticV7Config struct {
	Receiver            **ElasticV7
	ConnectionName      string                     // 连接名
	Addr                string                     // es地址 eg:http://127.0.0.1:9000
	Username            string                     // 用户名
	Password            string                     // 密码
	HealthCheckInterval int                        // 健康检查频率
	IsGzip              bool                       // 是否gizp压缩
	Retry               elastic.Retrier            // 重试方法
	InfoLogger          elastic.Logger             // 消息日志
	ErrLogger           elastic.Logger             // 错误日志
	Ext                 []elastic.ClientOptionFunc //拓展方法自定义
}

type ElasticV6Config struct {
	Receiver            **ElasticV6
	ConnectionName      string                 // 连接名
	Addr                string                 // es地址 eg:http://127.0.0.1:9000
	Username            string                 // 用户名
	Password            string                 // 密码
	HealthCheckInterval int                    // 健康检查频率
	IsGzip              bool                   // 是否gizp压缩
	Retry               es6.Retrier            // 重试方法
	InfoLogger          elastic.Logger         // 消息日志
	ErrLogger           elastic.Logger         // 错误日志
	Ext                 []es6.ClientOptionFunc //拓展方法自定义
}
