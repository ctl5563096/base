package library

import "github.com/Shopify/sarama"

type KafkaProducerConfig struct {
	Receiver      **KafkaSyncProducer //实例接受对象
	Name          string              //连接名称(自定义)
	Version       string              //版本
	BrokerAddress []string            //消息代理服务器地址
	ConsoleDebug  bool                //是否进入命令行终端debug模式
	ExtraConfig   *sarama.Config      //额外的配置项
}
