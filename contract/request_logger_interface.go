package contract

type RequestLoggerInterface interface {
	HttpRequestLog(record *HttpRequestRecord)
}
