package library

type EnvConfig struct {
	Receivers **Env
	Name string //文件名称
	FileType string //文件类型
	FilePath string
	FileName string
	EnableReReading bool //是否支持配置内容热加载
}