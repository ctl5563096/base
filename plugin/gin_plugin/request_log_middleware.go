package gin_plugin

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"io/ioutil"
	"talkcheap.xiaoeknow.com/xiaoetong/eframe/contract"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RequestLogMiddleware 用于记录请求日志 传入符合 contract.LoggerInterface 约束的日志对象
// maxByte	requestBody大于maxByte，请求日志不打印
func RequestLogMiddleware(logger contract.XiaoeRequestLoggerInterface, maxByte float64) func(ginCtx *gin.Context) {
	return func(ginCtx *gin.Context) {
		record := contract.XiaoeHttpRequestRecord{}
		begin := time.Now()
		record.AppId = ginCtx.Query("app_id")
		record.TraceId = ginCtx.GetHeader(contract.TraceId)
		record.Sw8 = ginCtx.GetHeader(contract.Sw8Header)
		record.Sw8Correlation = ginCtx.GetHeader(contract.Sw8CorrelationHeader)
		record.XeTag = ginCtx.GetHeader(contract.XeTagHeader)
		record.TargetUrl = ginCtx.Request.RequestURI
		record.Method = ginCtx.Request.Method
		record.ClientIp = ginCtx.ClientIP()
		record.ServerIp = ginCtx.Request.Host
		record.UserAgent = ginCtx.GetHeader("User-Agent")
		record.BeginTime = begin.Format("2006-01-02 15:04:05.000")
		headerByte, _ := json.Marshal(ginCtx.Request.Header)
		record.Header = cast.ToString(headerByte)

		var readBytes int64 = 0
		if record.Method == "POST" && ginCtx.Request.ContentLength > 0 {
			dataAll := make([]byte, ginCtx.Request.ContentLength)
			for {
				n, e := ginCtx.Request.Body.Read(dataAll[readBytes:])
				readBytes = readBytes + int64(n)
				if e != nil {
					break
				}
			}
			record.Params = string(dataAll)
			ginCtx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(dataAll))
		}

		if record.Method == "GET" || record.Method == "OPTION" {
			record.Params = ginCtx.Request.URL.RawQuery
		}

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ginCtx.Writer}
		ginCtx.Writer = blw
		ginCtx.Next()
		end := time.Now()
		spend := end.UnixNano() - begin.UnixNano()
		record.EndTime = end.Format("2006-01-02 15:04:05.000")
		record.CostTime = int(spend / 1000000)
		if ginCtx.Writer.Written() {
			record.Response = blw.body.String()
		}
		record.HttpStatus = ginCtx.Writer.Status()

		//大于maxLogByte不打印日志
		if readBytes < int64(maxByte) {
			logger.HttpRequestLog(&record)
		}
	}
}
