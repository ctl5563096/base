package library

type MysqlConfig struct {
	Receiver       **DB
	ConnectionName string
	DBName         string
	Host           string
	Port           string
	UserName       string
	Password       string
	Charset        string //字符集
	ParseTime      string //解析时间
	Loc            string //时区
	MaxLifeTime    int
	MaxOpenConn    int
	MaxIdleConn    int
}
