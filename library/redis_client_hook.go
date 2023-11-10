package library

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/SkyAPM/go2sky"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	"talkcheap.xiaoeknow.com/xiaoetong/eframe/contract"
)

// skyWalking接入hook
type SkyWalkingRedisHook struct {
	Peer string // redis地址
}

func (s SkyWalkingRedisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	span, err := s.CreateSpan(ctx, cmd.FullName(), cmd.Args())
	if err != nil {
		return ctx, nil
	}
	return context.WithValue(ctx, contract.Sw8Span, span), nil
}

func (s SkyWalkingRedisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	var (
		span go2sky.Span
		ok   bool
	)
	val := ctx.Value(contract.Sw8Span)
	if span, ok = val.(go2sky.Span); !ok {
		return nil
	}

	if cmd.Err() != nil {
		s.EndSpan(&span, cmd.Err())
	}
	s.EndSpan(&span, nil)
	return nil
}

func (s SkyWalkingRedisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	builder := strings.Builder{}
	for _, cmd := range cmds {
		builder.WriteString(fmt.Sprintf("%+v\n", cmd.Args()))
	}

	span, err := s.CreateSpan(ctx, "pipeline", builder.String())
	if err != nil {
		return ctx, nil
	}
	return context.WithValue(ctx, contract.Sw8Span, span), nil
}

func (s SkyWalkingRedisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	var (
		span go2sky.Span
		ok   bool
	)
	val := ctx.Value(contract.Sw8Span)
	if span, ok = val.(go2sky.Span); !ok {
		return nil
	}

	var errs []error
	for _, cmd := range cmds {
		if cmd.Err() != nil {
			errs = append(errs, cmd.Err())
		}
	}
	s.EndSpan(&span, errs...)
	return nil
}

func (s SkyWalkingRedisHook) CreateSpan(ctx context.Context, fullName, args interface{}) (go2sky.Span, error) {
	tracker := go2sky.GetGlobalTracer()
	if tracker == nil {
		return nil, nil
	}

	span, errT := tracker.CreateExitSpan(ctx, "redis", s.Peer, func(key, value string) error {
		return nil
	})
	if errT != nil {
		return nil, nil
	}

	span.SetComponent(contract.ComponentIDGoRedis)
	span.SetOperationName(fmt.Sprintf("%s->%v", getFuncName(5), fullName))
	span.SetSpanLayer(agentv3.SpanLayer_Cache)
	span.Tag(go2sky.TagDBType, "redis")
	span.Tag(go2sky.TagDBStatement, fmt.Sprintf("%s", args))
	return span, nil
}

func (s SkyWalkingRedisHook) EndSpan(span *go2sky.Span, errs ...error) {
	if span == nil {
		return
	}

	for _, err := range errs {
		if err != nil {
			(*span).Error(time.Now(), err.Error())
		}
	}
	(*span).End()
}

// 获取调用函数
func getFuncName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "runtime.Caller() failed"
	}

	return runtime.FuncForPC(pc).Name()
}
