package library

import (
    "github.com/SkyAPM/go2sky"
    "github.com/SkyAPM/go2sky/reporter"
)

type SkyWalkingConfig struct {
    Receiver    **SkyWalking
    ServiceName string  // 当前服务名
    Addr        string  // 上报地址
    Sample      float64 // 采样率s

    ReportOpts []reporter.GRPCReporterOption
    TracerOpts []go2sky.TracerOption
}
