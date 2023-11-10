package library

import (
	"context"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"sync"
)

type KafkaGroupConsumer struct {
	consumerGroup       sarama.ConsumerGroup
	messageHandle       MessageHandleFun       // 自动确认消息，当手动方法不存在时才会使用
	messageHandleByHand MessageHandleFunByHand // 手动确认消息
	consumeErrHandle    ConsumeErrHandleFunc
	setupHandle         SetupHandleFun
	cleanupHandle       CleanupHandleFun
	topics              []string
	lock                sync.Mutex
}

func NewKafkaGroupConsumer(config *KafkaGroupConsumerConfig) (kafkaGroupConsumer *KafkaGroupConsumer, err error) {
	if config.ExtraConfig == nil {
		config.ExtraConfig = sarama.NewConfig()
	}

	if config.Name != "" {
		config.ExtraConfig.ClientID = config.Name
	}

	if config.Version == "" {
		config.ExtraConfig.Version = sarama.V2_6_0_0
	}

	if config.Topics == nil || len(config.Topics) == 0 {
		err = errors.New("请指定需要消费的topic 以,号分隔")
		return nil, err
	}

	config.ExtraConfig.Version, err = sarama.ParseKafkaVersion(config.Version)
	if err != nil {
		err = fmt.Errorf("[%s] version string is err: %w", config.Name, err)
		return
	}

	if config.ConsoleDebug == true {
		sarama.Logger = log.New(os.Stdout, "["+config.Name+"]", log.LstdFlags)
	}

	if config.InitialOffset == sarama.OffsetOldest || config.InitialOffset == sarama.OffsetNewest {
		config.ExtraConfig.Consumer.Offsets.Initial = config.InitialOffset
	}

	consumerGroup, err := sarama.NewConsumerGroup(config.BrokerAddress, config.GroupName, config.ExtraConfig)
	if err != nil {
		return nil, err
	}

	return &KafkaGroupConsumer{
		consumerGroup: consumerGroup,
		topics:        config.Topics,
	}, nil
}

func NewGroupConsumerHandler(handlerFun MessageHandleFun, handlerByHand MessageHandleFunByHand, setupHandle SetupHandleFun, cleanupHandle CleanupHandleFun) *GroupConsumerHandler {
	return &GroupConsumerHandler{
		handleMessage:       handlerFun,
		handleMessageByHand: handlerByHand,
		handleSetup:         setupHandle,
		handleCleanup:       cleanupHandle,
	}
}

func (c *KafkaGroupConsumer) Close() error {
	return c.consumerGroup.Close()
}

func (c *KafkaGroupConsumer) SetSetupHandleFunc(f SetupHandleFun) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.setupHandle = f
}

func (c *KafkaGroupConsumer) SetCleanupHandleFunc(f CleanupHandleFun) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cleanupHandle = f
}

func (c *KafkaGroupConsumer) SetMessageHandleFunc(f MessageHandleFun) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.messageHandle = f
}

func (c *KafkaGroupConsumer) SetMessageHandleByHandFunc(f MessageHandleFunByHand) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.messageHandleByHand = f
}

func (c *KafkaGroupConsumer) SetConsumeErrHandleFunc(f ConsumeErrHandleFunc) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.consumeErrHandle = f
}

func (c *KafkaGroupConsumer) StartConsume(ctx context.Context) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.consumeErrHandle == nil {
		err = errors.New("请指定 ConsumeErrHandleFunc 消费异常处理逻辑")
		return
	}

	if c.messageHandle == nil {
		err = errors.New("请指定 MessageHandleFun 消息消费逻辑")
		return
	}

	handler := NewGroupConsumerHandler(c.messageHandle, c.messageHandleByHand, c.setupHandle, c.cleanupHandle)
	go func() {
		for {
			if err := c.consumerGroup.Consume(ctx, c.topics, handler); err != nil {
				c.consumeErrHandle(err)
				return
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()
	return nil
}

type ConsumeErrHandleFunc func(error)

type MessageHandleFun func(message *sarama.ConsumerMessage)

type MessageHandleFunByHand func(session *sarama.ConsumerGroupSession, message *sarama.ConsumerMessage)

type SetupHandleFun func(session *sarama.ConsumerGroupSession) error

type CleanupHandleFun func(session *sarama.ConsumerGroupSession) error

type GroupConsumerHandler struct {
	handleMessage       MessageHandleFun
	handleMessageByHand MessageHandleFunByHand
	handleSetup         SetupHandleFun
	handleCleanup       CleanupHandleFun
}

func (h *GroupConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	if h.handleSetup != nil {
		return h.handleSetup(&session)
	}
	return nil
}

func (h *GroupConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	if h.handleCleanup != nil {
		return h.handleCleanup(&session)
	}
	return nil
}

func (h *GroupConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			if h.handleMessageByHand != nil {
				h.handleMessageByHand(&session, msg)
				break
			}

			h.handleMessage(msg)
			session.MarkMessage(msg, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
