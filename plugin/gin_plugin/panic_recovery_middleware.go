package gin_plugin

import (
    "net"
    "net/http"
    "net/http/httputil"
    "os"
    "runtime/debug"
    "strings"

    "go.uber.org/zap"
    "talkcheap.xiaoeknow.com/xiaoetong/eframe/contract"
    "talkcheap.xiaoeknow.com/xiaoetong/eframe/library"

    "github.com/gin-gonic/gin"
)

// PanicRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func PanicRecovery(logger *library.Log, stack bool) gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // Check for a broken connection, as it is not really a
                // condition that warrants a panic stack trace.
                var brokenPipe bool
                if ne, ok := err.(*net.OpError); ok {
                    if se, ok := ne.Err.(*os.SyscallError); ok {
                        if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
                            brokenPipe = true
                        }
                    }
                }

                httpRequest, _ := httputil.DumpRequest(c.Request, false)
                if brokenPipe {
                    logger.Error(c.Request.URL.Path,
                        zap.Any("error", err),
                        zap.String("request", string(httpRequest)),
                    )
                    // If the connection is dead, we can't write a status to it.
                    c.Error(err.(error)) // nolint: errcheck
                    c.Abort()
                    return
                }

                if stack {
                    logger.Error("[Recovery from panic]",
                        zap.String("trace_id", c.GetHeader(contract.TraceId)),
                        zap.Any("error", err),
                        zap.String("request", string(httpRequest)),
                        zap.String("stack", string(debug.Stack())),
                    )
                } else {
                    logger.Error("[Recovery from panic]",
                        zap.String("trace_id", c.GetHeader(contract.TraceId)),
                        zap.Any("error", err),
                        zap.String("request", string(httpRequest)),
                    )
                }

                // 返回body中没有写数据、默认回写500
                if c.Writer.Size() == -1 || c.Writer.Size() == 0 {
                    c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]interface{}{
                        "code": http.StatusInternalServerError,
                        "msg":  err,
                        "data": nil,
                    })
                }
            }
        }()
        c.Next()
    }
}
