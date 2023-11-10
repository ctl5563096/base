package runtime

import "os"

// GetExecDirectory 获取当前执行程序目录
func GetExecDirectory() string  {
	file, err := os.Getwd()
	if err == nil {
		return file + "/"
	}
	return ""
}



