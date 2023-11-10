package library

import (
    "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConf struct {
    Receiver       **MongoClient
    ConnectionName string // 连接名
    Host           string // ip
    Port           string // 端口号
    Username       string // 用户名
    Password       string // 密码
    Timeout        int    // 超时时间
    option         *options.ClientOptions
}
