package library

import (
	"fmt"
	"time"

	"base/contract"
	"github.com/SkyAPM/go2sky"
	"gorm.io/gorm"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

// skyWalking接入hook
type SkyWalkingGormHook struct {
	Peer string // redis地址
}

func (s SkyWalkingGormHook) BeforeCallback(db *gorm.DB) {
	tracker := go2sky.GetGlobalTracer()
	if tracker == nil {
		return
	}

	span, errT := tracker.CreateExitSpan(db.Statement.Context, "gorm", s.Peer, func(key, value string) error {
		return nil
	})
	if errT != nil {
		return
	}
	span.SetComponent(contract.ComponentIDGoGorm)
	span.SetOperationName(fmt.Sprintf("%s->%v", getFuncName(4), db.Statement.Name()))
	span.SetSpanLayer(agentv3.SpanLayer_Database)
	span.Tag(go2sky.TagDBType, "gorm")
	db.Set(contract.Sw8Span, span)
	return
}

func (s SkyWalkingGormHook) AfterCallback(db *gorm.DB) {
	val, ok := db.Get(contract.Sw8Span)
	if ok == false {
		return
	}
	if span, ok := val.(go2sky.Span); ok {
		span.Tag(go2sky.TagDBStatement, fmt.Sprintf("%s", db.Statement.SQL.String()))
		if db.Error != nil {
			span.Error(time.Now(), db.Error.Error())
		}
		span.End()
	}
	return
}
