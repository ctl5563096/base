package library

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ctl5563096/base/contract"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/ctl5563096/base/helpers/network"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

type HttpClient struct {
	*http.Client
}

func NewHttpClient(config *HttpClientConfig) (httpClient *HttpClient) {
	helpe
	httpClient = &HttpClient{
		Client: &http.Client{
			Timeout: time.Duration(config.RequestTimeoutSecond) * time.Second,
			Transport: &http.Transport{
				DisableKeepAlives:   config.DialKeepAliveSecond < 0,
				MaxIdleConns:        config.MaxIdleConnections,
				MaxIdleConnsPerHost: config.MaxIdleConnectionsPerHost,
				IdleConnTimeout:     time.Duration(config.IdleConnTimeoutSecond) * time.Second,
				Proxy:               http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   time.Duration(config.DialTimeoutSecond) * time.Second,
					KeepAlive: time.Duration(config.DialKeepAliveSecond) * time.Second,
				}).DialContext,
			},
		},
	}
	return
}

func NewHttpClientCache(config *HttpClientCacheConfig) (httpClient *HttpClient) {
	httpClient = &HttpClient{
		Client: &http.Client{
			Timeout: time.Duration(config.RequestTimeoutSecond) * time.Second,
			Transport: &http.Transport{
				DisableKeepAlives:   config.DialConf.DialKeepAliveSecond < 0,
				MaxIdleConns:        config.MaxIdleConnections,
				MaxIdleConnsPerHost: config.MaxIdleConnectionsPerHost,
				IdleConnTimeout:     time.Duration(config.IdleConnTimeoutSecond) * time.Second,
				Proxy:               http.ProxyFromEnvironment,
				DialContext:         NewDialer(&config.DialConf).DialContext,
			},
		},
	}
	return
}

func (c *HttpClient) Get(ctx context.Context, url string, header map[string]string, logger contract.RequestLoggerInterface) (response []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	var clientResp *http.Response
	var paramsString string
	var beginTime = time.Now()
	defer recordLog(req, &clientResp, &paramsString, &response, &err, beginTime, logger)
	if err != nil {
		return nil, err
	}

	addXeHeader(ctx, req)
	addCustomHeader(req, header)

	clientResp, err = c.Do(req)

	if err != nil {
		err = fmt.Errorf("response is nil: %w", err)
		return nil, err
	}

	if clientResp == nil {
		err = fmt.Errorf("response is nil: %w", err)
		return nil, err
	}
	defer clientResp.Body.Close()

	//2开头的状态码都是OK
	if clientResp.StatusCode < 200 || clientResp.StatusCode >= 300 {
		resBody, _ := getBytesFromHttpResponse(clientResp)
		err = &contract.HttpResponseError{
			Code:         clientResp.StatusCode,
			Msg:          fmt.Sprintf("response error, code %d", clientResp.StatusCode),
			ResponseBody: resBody,
		}
		return nil, err
	}

	response, err = getBytesFromHttpResponse(clientResp)

	return response, err
}

func (c *HttpClient) GetV2(ctx context.Context, url string, query map[string]string, header map[string]string, logger contract.RequestLoggerInterface) (response []byte, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	var clientResp *http.Response
	var paramsString string
	var beginTime = time.Now()
	defer recordLog(req, &clientResp, &paramsString, &response, &err, beginTime, logger)
	if err != nil {
		return nil, err
	}

	addXeHeader(ctx, req)
	addCustomHeader(req, header)

	// 拼接参数
	q := req.URL.Query()
	for key, val := range query {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()

	clientResp, err = c.Do(req)

	if err != nil {
		err = fmt.Errorf("response is nil: %w", err)
		return nil, err
	}

	if clientResp == nil {
		err = fmt.Errorf("response is nil: %w", err)
		return nil, err
	}
	defer clientResp.Body.Close()

	//2开头的状态码都是OK
	if clientResp.StatusCode < 200 || clientResp.StatusCode >= 300 {
		resBody, _ := getBytesFromHttpResponse(clientResp)
		err = &contract.HttpResponseError{
			Code:         clientResp.StatusCode,
			Msg:          fmt.Sprintf("response error, code %d", clientResp.StatusCode),
			ResponseBody: resBody,
		}
		return nil, err
	}

	response, err = getBytesFromHttpResponse(clientResp)

	return response, err
}

func (c *HttpClient) Post(ctx context.Context, url string, params []byte, header map[string]string, logger contract.RequestLoggerInterface) (response []byte, err error) {
	var clientResp *http.Response
	var paramsString string
	var beginTime = time.Now()
	var body io.Reader
	if params != nil {
		body = bytes.NewReader(params)
		paramsString = string(params)
	}
	req, e := http.NewRequestWithContext(ctx, "POST", url, body)
	defer recordLog(req, &clientResp, &paramsString, &response, &err, beginTime, logger)
	if e != nil {
		err = e
		return nil, e
	}

	addXeHeader(ctx, req)
	addCustomHeader(req, header)

	clientResp, err = c.Do(req)

	if err != nil {
		err = fmt.Errorf("response is nil: %w", err)
		return nil, err
	}

	if clientResp == nil {
		err = fmt.Errorf("response is nil: %w", err)
		return nil, err
	}
	defer clientResp.Body.Close()
	//2开头的状态码都是OK
	if clientResp.StatusCode < 200 || clientResp.StatusCode >= 300 {
		resBody, _ := getBytesFromHttpResponse(clientResp)
		err = &contract.HttpResponseError{
			Code:         clientResp.StatusCode,
			Msg:          fmt.Sprintf("response error, code %d", clientResp.StatusCode),
			ResponseBody: resBody,
		}
		return nil, err
	}

	response, err = getBytesFromHttpResponse(clientResp)

	return response, err
}
func (c *HttpClient) GetJsonWithHeader(ctx context.Context, url string, header map[string]string, response interface{}, logger contract.RequestLoggerInterface) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	var responseBytes []byte
	var clientResp *http.Response
	var params string
	var beginTime = time.Now()
	defer recordLog(req, &clientResp, &params, &responseBytes, &err, beginTime, logger)
	if err != nil {
		return err
	}

	addJsonHeader(req)
	addXeHeader(ctx, req)
	addCustomHeader(req, header)

	clientResp, err = c.Do(req)

	if err != nil {
		err = fmt.Errorf("response is nil: %w", err)
		return err
	}

	if clientResp == nil {
		err = fmt.Errorf("response is nil: %w", err)
		return err
	}
	defer clientResp.Body.Close()

	if clientResp.StatusCode != http.StatusOK {
		err = &contract.HttpResponseError{
			Code: clientResp.StatusCode,
			Msg:  fmt.Sprintf("response error, code %d", clientResp.StatusCode),
		}
		return err
	}

	responseBytes, err = getBytesFromHttpResponse(clientResp)
	if err != nil {
		return err
	}

	if response == nil {
		response = &(map[string]interface{}{})
	}

	if err = json.Unmarshal(responseBytes, response); err != nil {
		return err
	}

	return err
}

func (c *HttpClient) GetJson(ctx context.Context, url string, response interface{}, logger contract.RequestLoggerInterface) error {
	return c.GetJsonWithHeader(ctx, url, nil, response, logger)
}

func (c *HttpClient) PostJsonWithHeader(ctx context.Context, url string, params interface{}, header map[string]string, response interface{}, logger contract.RequestLoggerInterface) error {
	var responseBytes []byte
	var clientResp *http.Response
	var paramsString string
	var beginTime = time.Now()
	var body io.Reader
	var err error
	if params != nil {
		paramsBytes, e := json.Marshal(params)
		if e != nil {
			return e
		}

		if paramsBytes != nil {
			body = bytes.NewReader(paramsBytes)
			paramsString = string(paramsBytes)
		}

	}
	req, e := http.NewRequestWithContext(ctx, "POST", url, body)
	defer recordLog(req, &clientResp, &paramsString, &responseBytes, &err, beginTime, logger)
	if e != nil {
		err = e
		return err
	}
	addJsonHeader(req)
	addXeHeader(ctx, req)
	addCustomHeader(req, header)

	clientResp, err = c.Do(req)

	if err != nil {
		err = fmt.Errorf("response is nil: %w", err)
		return err
	}

	if clientResp == nil {
		err = fmt.Errorf("response is nil: %w", err)
		return err
	}
	defer clientResp.Body.Close()

	if clientResp.StatusCode != http.StatusOK {
		err = &contract.HttpResponseError{
			Code: clientResp.StatusCode,
			Msg:  fmt.Sprintf("response error, code %d", clientResp.StatusCode),
		}
		return err
	}

	responseBytes, err = getBytesFromHttpResponse(clientResp)
	if err != nil {
		return err
	}

	if response == nil {
		response = &(map[string]interface{}{})
	}

	if err = json.Unmarshal(responseBytes, response); err != nil {
		return err
	}

	return err
}

func (c *HttpClient) PostJson(ctx context.Context, url string, params interface{}, response interface{}, logger contract.RequestLoggerInterface) error {
	return c.PostJsonWithHeader(ctx, url, params, nil, response, logger)
}

// 重新封装Do逻辑、添加skyWalking代码
func (c *HttpClient) Do(req *http.Request) (resp *http.Response, err error) {
	// 无tracer
	tracer := go2sky.GetGlobalTracer()
	if tracer == nil {
		resp, err = c.Client.Do(req)
		return
	}

	operateName := fmt.Sprintf("/%s%s", req.Method, req.URL.Path)
	// 创建tracer span失败
	span, err := tracer.CreateExitSpan(req.Context(), operateName, req.URL.Host, func(key, value string) error {
		req.Header.Set(key, value)
		return nil
	})
	if err != nil {
		fmt.Printf("Get CreateExitSpan err: %v", err)
		resp, err = c.Client.Do(req)
		return
	}
	defer span.End()

	// skywalking、执行Do参数
	/*	for k, v := range t.extraTags {
		span.Tag(go2sky.Tag(k), v)
	}*/
	span.SetComponent(contract.ComponentIDGOHttpClient)
	span.Tag(go2sky.TagHTTPMethod, req.Method)
	span.Tag(go2sky.TagURL, req.URL.String())
	span.SetSpanLayer(agentv3.SpanLayer_Http)
	resp, err = c.Client.Do(req)
	if err != nil {
		span.Error(time.Now(), err.Error())
		return
	}

	span.Tag(go2sky.TagStatusCode, strconv.Itoa(resp.StatusCode))
	if resp.StatusCode >= http.StatusBadRequest {
		span.Error(time.Now(), "Errors on handling client")
	}
	return
}

func (c *HttpClient) Close() error {
	c.CloseIdleConnections()
	return nil
}

func addCustomHeader(req *http.Request, header map[string]string) {
	for k, v := range header {
		req.Header.Add(k, v)
	}
}

func getBytesFromHttpResponse(response *http.Response) (b []byte, err error) {
	if response == nil {
		return nil, errors.New("http response is nil")
	}
	if response.ContentLength > 0 {
		b = make([]byte, response.ContentLength)
		readBytes := 0
		for {
			n, e := response.Body.Read(b[readBytes:])
			readBytes = readBytes + n
			if e != nil {
				break
			}
		}
	} else {
		b, err = io.ReadAll(response.Body)
	}
	return
}

func addJsonHeader(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
}

func addXeHeader(ctx context.Context, req *http.Request) {
	values, ok := ctx.Value(contract.Ctx).(map[string]string)
	if ok {
		//Xe灰度标识
		if header := values[contract.XeTagHeader]; header != "" {
			req.Header.Add(contract.XeTagHeader, header)
		}

		//SkyWalking标识
		if header := values[contract.Sw8Header]; header != "" {
			req.Header.Add(contract.Sw8Header, header)
		}

		if header := values[contract.Sw8CorrelationHeader]; header != "" {
			req.Header.Add(contract.Sw8CorrelationHeader, header)
		}

		// 全局唯一标识
		if header := values[contract.TraceId]; header != "" {
			req.Header.Add(contract.TraceId, header)
		}
	}

}

func recordLog(req *http.Request, resp **http.Response, params *string, response *[]byte, err *error, beginTime time.Time, logger contract.RequestLoggerInterface) {
	if logger != nil && req != nil {
		record := contract.XiaoeHttpRequestRecord{}
		record.Sw8 = req.Header.Get(contract.Sw8Header)
		record.Sw8Correlation = req.Header.Get(contract.Sw8CorrelationHeader)
		record.XeTag = req.Header.Get(contract.XeTagHeader)
		record.TraceId = req.Header.Get(contract.TraceId)
		record.TargetUrl = req.URL.String()
		record.Method = req.Method
		if params != nil {
			record.Params = *params
		}
		if err != nil && *err != nil {
			record.Msg = (*err).Error()
		}

		record.ClientIp = network.GetInternalIp()
		if resp != nil && *resp != nil {
			record.HttpStatus = (*resp).StatusCode
			record.ServerIp = (*resp).Request.RemoteAddr
			record.UserAgent = req.UserAgent()
		}
		if response != nil && *response != nil {
			record.Response = string(*response)
		}
		record.BeginTime = beginTime.Format("2006-01-02 15:04:05.000")
		end := time.Now()
		spend := end.UnixNano() - beginTime.UnixNano()
		record.EndTime = end.Format("2006-01-02 15:04:05.000")
		record.CostTime = int(spend / 1000000)
		logger.HttpRequestLog(&record)
	}
}
