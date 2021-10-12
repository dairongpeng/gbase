package gbase

import (
	"fmt"
	"github.com/Shopify/sarama"
)

// NewKafkaConfig 获取sarama默认配置
func NewKafkaConfig() *sarama.Config {
	// sarama默认给的配置
	c := sarama.NewConfig()
	// 默认消息不丢
	c.Producer.RequiredAcks = sarama.WaitForAll
	// 消费者默认从last消费位消费
	c.Consumer.Offsets.Initial = sarama.OffsetOldest
	return c
}

// NewKafkaClient 获取Kafka的Client
func NewKafkaClient(name string, config *sarama.Config) (sarama.Client, error) {
	addrs := Cfg().GetStringSlice(name + ".addrs")
	client, err := sarama.NewClient(addrs, config)
	if err != nil {
		return nil, fmt.Errorf("new sarama client: %w", err)
	}
	return client, nil
}

// NewSyncProducer 创建同步Producer
func NewSyncProducer(client sarama.Client) (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, fmt.Errorf("new sync producer from client: %w", err)
	}
	// return otelsarama.WrapSyncProducer(client.Config(), producer), nil
	return producer, nil
}

// NewAsyncProducer 返回异步Producer
func NewAsyncProducer(client sarama.Client) (sarama.AsyncProducer, error) {
	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		return nil, fmt.Errorf("new async producer from client: %w", err)
	}
	// return otelsarama.WrapAsyncProducer(client.Config(), producer), nil
	return producer, nil
}
