package library

import (
    "fmt"

    "github.com/SkyAPM/go2sky"
    "github.com/SkyAPM/go2sky/reporter"
)

type SkyWalking struct {
    *go2sky.Tracer
}

func NewSkyWalkingTracker(conf *SkyWalkingConfig) (sky *SkyWalking, err error) {
    // 创建reporter
    gRPCReporter, err2 := reporter.NewGRPCReporter(conf.Addr, conf.ReportOpts...)
    if err2 != nil {
        err = fmt.Errorf("创建上传gRpcReporter失败:%w", err2)
        return
    }

    // 创建tracer
    var tracerOpts []go2sky.TracerOption
    tracerOpts = append(tracerOpts, go2sky.WithReporter(gRPCReporter), go2sky.WithSampler(conf.Sample))
    tracerOpts = append(tracerOpts, conf.TracerOpts...)
    tracer, err2 := go2sky.NewTracer(conf.ServiceName, tracerOpts...)
    if err2 != nil {
        err = fmt.Errorf("创建 tracer 失败:%w", err2)
        return
    }

    sky = &SkyWalking{
        tracer,
    }

    return
}

func (sky *SkyWalking) SwitchTrace(isTrace bool) {
    if isTrace{
        go2sky.SetGlobalTracer(sky.Tracer)
        return
    }
    go2sky.SetGlobalTracer(nil)
}
