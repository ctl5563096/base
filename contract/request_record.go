package contract

type HttpRequestRecord struct {
	TraceId        string `json:"trace_id"`        //请求唯一id
	Sw8            string `json:"sw8"`             //SkyWalking链路标示 对应header sw8
	Sw8Correlation string `json:"sw8_correlation"` //SkyWalking链路标示 对应header sw8-correlation
	Uid            string `json:"uid"`             //用户id
	HttpStatus     int    `json:"http_status"`     //http状态码
	TargetUrl      string `json:"target_url"`      //请求地址
	Method         string `json:"method"`          //请求方式
	Msg            string `json:"msg"`             //日志说明
	ClientIp       string `json:"client_ip"`       //客户端ip
	ServerIp       string `json:"server_ip"`       //服务端ip
	UserAgent      string `json:"user_agent"`      //请求代理
	BeginTime      string `json:"begin_time"`      //开始时间 格式:2006-01-02 15:04:05.000
	EndTime        string `json:"end_time"`        //结束时间 格式:2006-01-02 15:04:05.000
	CostTime       int    `json:"cost_time"`       //请求花费时间 毫秒
	LogAt          string `json:"log_at"`          //日志记录时间 zap日志自带
	ServiceName    string `json:"service_name"`    //服务号名称 zap日志自带
	Params         string `json:"params"`          //请求参数
	Response       string `json:"response"`        //响应内容
	Header         string `json:"header"`          //请求头
}
