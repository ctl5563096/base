package gin_plugin

import (
	"context"
	"github.com/ctl5563096/base/contract"
	"github.com/ctl5563096/base/helpers/str"
	"github.com/gin-gonic/gin"
)

func XeSpecificContextSet(ginCtx *gin.Context) {
	if ginCtx.GetHeader(contract.TraceId) == "" {
		ginCtx.Request.Header.Set(contract.TraceId, str.RandStringBytesMaskImprSrcUnsafe(16))
	}
	ctx := ginCtx.Request.Context()
	xeSpecific := make(map[string]string, 4)
	xeSpecific[contract.TraceId] = ginCtx.GetHeader(contract.TraceId)
	xeSpecific[contract.XeTagHeader] = ginCtx.GetHeader(contract.XeTagHeader)
	xeSpecific[contract.Sw8Header] = ginCtx.GetHeader(contract.Sw8Header)
	xeSpecific[contract.Sw8CorrelationHeader] = ginCtx.GetHeader(contract.Sw8CorrelationHeader)
	xeCtx := context.WithValue(ctx, contract.XeCtx, xeSpecific)
	ginCtx.Request = ginCtx.Request.WithContext(xeCtx)
}
