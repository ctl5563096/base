package library

type RedisConfig struct {
    Receiver         **RedisClient
    ConnectionName   string // 连接名称自定义
    Addr             string // 地址
    Port             int    // 端口
    Password         string // 密码
    DB               int
    PoolSize         int    // 连接池大小
    EnableSkyWalking bool   // 开启skyWalking追踪
}
