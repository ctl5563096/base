package library


var SoftVersion = SoftVersionStrc{}

type SoftVersionStrc struct {
    Name    string  // 服务名称
    Version string  // git版本号
}