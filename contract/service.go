package contract

import "context"

type BaseService interface {
	Register() error
	LazyLoad() bool
	GetService() interface{}
	MustGetService() interface{}
	ServiceType() string
}

type RunnableEService interface {
	Run(ctx context.Context) error
}

type StoppableEService interface {
	Stop() error
}
