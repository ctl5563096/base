package library

import "github.com/Shopify/sarama"

type KafkaGroupConsumerConfig struct {
	Receiver      **KafkaGroupConsumer //实例接受对象
	Name          string               //名称（自定义）
	GroupName     string               //消费组名称
	Version       string               //版本
	Topics        []string             //消费主题
	BrokerAddress []string             //消息代理服务器地址
	ConsoleDebug  bool                 //开启终端debug模式
	InitialOffset int64                //默认消费策略
	ExtraConfig   *sarama.Config       //额外配置项
}
