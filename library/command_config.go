package library

import "context"

type CommandConfig struct {
	Signature   string                                            //脚本运行参数
	Description string                                            //脚本描述信息
	IsEndless   bool                                              //是否需要在执行完成后退出
	HandleFunc  func(ctx context.Context, cmd *ExecCommand) error //脚本逻辑函数
}
