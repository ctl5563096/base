package library

import (
	"context"
	"github.com/ctl5563096/base/contract"
	"net"
	"time"

	"go.uber.org/zap"
)

type Dialer struct {
	resolver DnsResolverInterface
	dialer   net.Dialer
	logger   *zap.Logger
}

func NewDialer(conf *DialConfig) *Dialer {
	return &Dialer{
		logger: conf.Logger,
		resolver: &DnsResolver{
			Resolver:  net.Resolver{},
			cacheTime: conf.DnsCacheTime,
			cacheMap:  NewLocalCache(conf.DnsCacheNums),
		},
		dialer: net.Dialer{
			Timeout:   time.Duration(conf.DialTimeoutSecond) * time.Second,
			KeepAlive: time.Duration(conf.DialKeepAliveSecond) * time.Second,
		},
	}
}

// DialContext connects to the address on the named network using
// the provided context.
func (d *Dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	ips, err := d.resolver.LookupHost(ctx, host)
	for _, ip := range ips {
		conn, err := d.dialer.DialContext(ctx, network, ip+":"+port)
		if err == nil {
			return conn, nil
		}

		if d.logger != nil {
			d.logger.Error("DialContext failed：",
				zap.String(contract.TraceId, getTraceId(ctx)),
				zap.String("network", network),
				zap.String("address", address), zap.Error(err))
		}
	}
	return d.dialer.DialContext(ctx, network, address)
}

// 获取全局追踪id
func getTraceId(ctx context.Context) string {
	values, ok := ctx.Value(contract.Ctx).(map[string]string)
	if ok { // 全局唯一标识
		return values[contract.TraceId]
	}
	return ""
}
