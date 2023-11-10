package gin_plugin

import (
    "fmt"
    "strconv"
    "time"

    "github.com/SkyAPM/go2sky"
    "github.com/gin-gonic/gin"
    agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
    "talkcheap.xiaoeknow.com/xiaoetong/eframe/contract"
    "talkcheap.xiaoeknow.com/xiaoetong/eframe/library"
)

// skyWalking 上报数据中间件
func SkyWalkingReport(engine *gin.Engine, sky *library.SkyWalking) gin.HandlerFunc {
    if engine == nil || sky == nil {
        return func(c *gin.Context) {
            c.Next()
        }
    }

    // 返回中间件处理
    return func(c *gin.Context) {
        // 过滤掉健康检查
        if c.Request.URL.Path == "/health" {
            c.Next()
            return
        }

        gTracer := go2sky.GetGlobalTracer()
        if gTracer == nil {
            c.Next()
            return
        }

        // 创建span
        span, ctx, err := gTracer.CreateEntrySpan(c.Request.Context(), getOperationName(c), func(key string) (string, error) {
            return c.Request.Header.Get(key), nil
        })
        if err != nil {
            c.Next()
            return
        }

        span.SetComponent(contract.ComponentIDGINHttpServer)
        span.Tag(go2sky.TagHTTPMethod, c.Request.Method)
        span.Tag(go2sky.TagURL, c.Request.Host+c.Request.URL.Path)
        span.SetSpanLayer(agentv3.SpanLayer_Http)

        c.Request = c.Request.WithContext(ctx)

        c.Next()

        if len(c.Errors) > 0 {
            span.Error(time.Now(), c.Errors.String())
        }
        span.Tag(go2sky.TagStatusCode, strconv.Itoa(c.Writer.Status()))
        span.End()
    }
}

func getOperationName(c *gin.Context) string {
    return fmt.Sprintf("/%s%s", c.Request.Method, c.FullPath())
}
