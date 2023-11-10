package gin_plugin

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)
// InternalMiddleware 内网访问限制中间件 environment 指定当前的系统运行环境
// internalIps 指定内网允许访问的ip段 如:
// 100.64.0.0/10  127.0.0.1/32 172.16.0.0/12 192.168.0.0/16 9.0.0.0/8 11.0.0.0/8 30.0.0.0/8 等等
func InternalMiddleware(environment string, internalAllowIps []string) func(ginCtx *gin.Context) {
	return func(ginCtx *gin.Context) {
		if environment != "development" {
			clientIp := ginCtx.ClientIP()
			if clientIp == "localhost" || clientIp == "127.0.0.1" ||
				clientIp == "::1" || checkBatchIpAddress(clientIp, internalAllowIps) {
			} else {
				ginCtx.AbortWithStatusJSON(500, map[string]interface{}{
					"msg": "内网中间件限制：abort",
				})
			}
		}
	}
}


func checkBatchIpAddress(requestIp string, checkIps []string) bool {
	for _, checkIp := range checkIps {
		if hasIpAddress(requestIp, checkIp) {
			return true
		}
	}
	return false
}

func hasIpAddress(ip, cidr string) (has bool) {
	ips := strings.Split(ip, ".")
	if len(ips) != 4 {
		return false
	}

	cidrArr := strings.Split(cidr, "/")
	if len(cidrArr) != 2 {
		return false
	}

	cidrIps := strings.Split(cidrArr[0], ".")
	ipType, _ := strconv.Atoi(cidrArr[1])
	mask := uint32(0xFFFFFFFF) << (32 - ipType)
	if len(cidrIps) != 4 {
		return false
	}

	var ipAddr uint32
	var cidrIp uint32
	for i, v := range ips {
		ip, _ := strconv.Atoi(v)
		ipAddr += uint32(ip << (8 * (3 - i)))
	}

	for i, v := range cidrIps {
		ip, _ := strconv.Atoi(v)
		cidrIp += uint32(ip << (8 * (3 - i)))
	}

	return ipAddr&mask == cidrIp&mask
}