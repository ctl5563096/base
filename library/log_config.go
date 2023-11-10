package library

// 兼容旧逻辑
type LoggerConfig struct {
	Receiver       **Log  //实例化接收对象指针
	Level          string //日志级别
	Name           string //日志名称
	Path           string //日志路径
	MaxAgeDay      int    //最大保存天数
	MaxFileSize    int    //单日志文件最大大小M
	MaxBackups     int    //最大历史文件个数
	DebugMode      bool   //是否开启调试模式（会记录日志所在位置文件和行号）
	PrintInConsole bool   //在控制台显示
	ServiceName    string //系统服务名
	CompressFile   bool   //压缩文件 lumberjack
}

type BaseLoggerConfig struct {
	Receiver       **Log  //实例化接收对象指针
	Level          string //日志级别
	Name           string //日志名称
	Path           string //日志路径
	MaxAgeDay      int    //最大保存天数
	MaxFileSize    int    //单日志文件最大大小M
	MaxBackups     int    //最大历史文件个数
	DebugMode      bool   //是否开启调试模式（会记录日志所在位置文件和行号）
	PrintInConsole bool   //在控制台显示
	ServiceName    string //系统服务名
}

// 根据文件大小切割配置
type LumberLoggerConfig struct {
	BaseLoggerConfig
	CompressFile bool //压缩文件 lumberjack
}

// 根据时间切割日志
type RotateLoggerConfig struct {
	BaseLoggerConfig
	RotationTime  int  //切割日志时间 file-rotatelogs
	UseMaxBackups bool // fileRotate MaxAgeDay 和 MaxBackups，不能共存
}
