package library

import (
    "context"
    "go.uber.org/zap"
    "net"
    "talkcheap.xiaoeknow.com/xiaoetong/eframe/contract"
    "time"
)

type DnsResolverInterface interface {
    LookupHost(ctx context.Context, host string) (addrs []string, err error)
}

type DnsResolver struct {
    net.Resolver
    cacheTime time.Duration
    cacheMap *LocalCache
    logger *zap.Logger
}

func (r *DnsResolver) LookupHost(ctx context.Context, host string) (addrs []string, err error) {
    if val, ok := r.cacheMap.Get(host); ok {
        return val.([]string), nil
    }

    addrs, err = r.Resolver.LookupHost(ctx, host)
    if r.logger != nil {
        r.logger.Info("LookupHost: addrs",
            zap.String(contract.TraceId, getTraceId(ctx)),
            zap.Any("ips", addrs),
            zap.Error(err))
    }

    if err != nil {
        return
    }

    if len(addrs) > 0 {
        _ = r.cacheMap.Put(host, addrs, r.cacheTime)
    }

    return
}

