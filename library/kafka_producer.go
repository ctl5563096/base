package library

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
)

type KafkaSyncProducer struct {
	sarama.SyncProducer
}

func NewKafkaSyncProducer(config *KafkaProducerConfig) (producer *KafkaSyncProducer, err error) {
	if config.ExtraConfig == nil {
		config.ExtraConfig = sarama.NewConfig()
	}
	if config.Version == "" {
		config.ExtraConfig.Version = sarama.V2_6_0_0
	}

	if config.ConsoleDebug == true {
		//TODO 并发安全问题
		sarama.Logger = log.New(os.Stdout, "["+config.Name+"]", log.LstdFlags)
	}

	config.ExtraConfig.Version, err = sarama.ParseKafkaVersion(config.Version)
	if err != nil {
		err = fmt.Errorf("[%s] version string is err: %w", config.Name, err)
		return
	}
	config.ExtraConfig.Producer.Return.Successes = true

	if config.BrokerAddress == nil || len(config.BrokerAddress) == 0 {
		err = fmt.Errorf("[%s] brokerAddress is empty", config.Name)
		return
	}

	syncProducer, e := sarama.NewSyncProducer(config.BrokerAddress, config.ExtraConfig)
	if e != nil {
		err = fmt.Errorf("[%s] new sync producer is error: %w", config.Name, e)
		return
	}

	producer = &KafkaSyncProducer{
		syncProducer,
	}
	return
}

func (producer *KafkaSyncProducer) Close() (err error) {
	err = producer.SyncProducer.Close()
	return
}
