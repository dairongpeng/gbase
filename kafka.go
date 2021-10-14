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

//// Consumer 消费者
//type Consumer struct {
//	Client         sarama.Client
//	ConsumerGroup  sarama.ConsumerGroup
//	ConsumeFunc    ConsumeFunc
//	MaxElapsedTime time.Duration
//
//	Ready chan struct{}
//}
//
//// ConsumeFunc 如果消费失败，返回error。返回error的消息会反复重试
//type ConsumeFunc func(context.Context, *sarama.ConsumerMessage) error
//
//// NewConsumer 指定你用来消费消息的函数, 如果返回的 err != nil, 会导致消息重试
//// 对于每组 (topic, partition) 会起一个协程来处理消息, 所以绝大多数时候你不需要额外的并发控制.
//func NewConsumer(addrs []string, groupID string, config *sarama.Config, msgConsumeFunc ConsumeFunc) (*Consumer, error) {
//	return newConsumer(addrs, groupID, config, msgConsumeFunc, 0)
//}
//
//// NewConsumerWithMaxElapsed 指定你用来消费消息的函数, 如果返回的 err != nil, 会导致消息重试
//// 对于每组 (topic, partition) 会起一个协程来处理消息, 所以绝大多数时候你不需要额外的并发控制.
//func NewConsumerWithMaxElapsed(addrs []string, groupID string, config *sarama.Config, msgConsumeFunc ConsumeFunc, maxElapsedSecond int) (*Consumer, error) {
//	return newConsumer(addrs, groupID, config, msgConsumeFunc, maxElapsedSecond)
//}
//
//func newConsumer(addrs []string, groupID string, config *sarama.Config, msgConsumeFunc ConsumeFunc, maxElapsedSecond int) (*Consumer, error) {
//	if config == nil {
//		config = sarama.NewConfig()
//	}
//	atLeastVersion(config)
//
//	client, err := sarama.NewClient(addrs, config)
//	if err != nil {
//		return nil, err
//	}
//
//	consumerGroup, err := sarama.NewConsumerGroup(addrs, groupID, config)
//	if err != nil {
//		return nil, err
//	}
//
//	consumer := &Consumer{
//		Client:         client,
//		ConsumerGroup:  consumerGroup,
//		ConsumeFunc:    msgConsumeFunc,
//		MaxElapsedTime: time.Duration(maxElapsedSecond) * time.Second,
//
//		Ready: make(chan struct{}),
//	}
//
//	return consumer, nil
//}
//
//func atLeastVersion(config *sarama.Config) {
//	if !config.Version.IsAtLeast(kafkaVersion) {
//		Debug(context.Background()).
//			Str("old_version", config.Version.String()).
//			Str("new_version", kafkaVersion.String()).
//			Msg("Change kafka version")
//		config.Version = kafkaVersion
//	}
//}
//
//// 默认使用kafka的版本
//var kafkaVersion = sarama.V2_7_0_0
//
//// StartConsumerWithContext 启动消费 Kafka 消息的应用服务
//func StartConsumerWithContext(ctx context.Context, c *Consumer, topics []string) {
//	startConsumer(ctx, c, topics)
//}
//
//func startConsumer(ctx context.Context, c *Consumer, topics []string) {
//	// TODO: panic recovery, k8s 的健康检测和自动重启?
//	initSaramaLog()
//
//	ctx, notifyStop := context.WithCancel(ctx)
//
//	wg := &sync.WaitGroup{}
//
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//
//		b := backoff.NewExponentialBackOff()
//		b.Multiplier = 1.1
//		// NOTE: 说不定哪天会成为坑
//		b.MaxElapsedTime = time.Hour * 24 * 30 * 12 * 10
//		err := backoff.RetryNotify(func() error {
//			// `Consume` should be called inside an infinite loop, when a
//			// server-side rebalance happens, the consumer session will need to be
//			// recreated to get the new claims
//			err := c.ConsumerGroup.Consume(ctx, topics, c)
//			if err != nil {
//				log.Error().
//					Err(err).
//					Strs("topics", topics).
//					Msg("Consume message with error")
//				return err
//			}
//			// check if context was cancelled, signaling that the consumer should stop
//			if ctx.Err() != nil {
//				return nil
//			}
//			c.Ready = make(chan struct{})
//			return errors.New("UnknownCase: Kafka consume with error and context was not cancelled")
//		}, b, func(err error, d time.Duration) {
//			log.Error().Err(err).Dur("backoff", d).Msg("Retry to start consumer")
//		})
//		if err != nil {
//			log.Error().Err(err).Msg("Start consumer with retry failed")
//		}
//	}()
//
//	<-c.Ready
//	Info(ctx).Msgf("Kafka consumer started with %v", topics)
//
//	waitCloseOrShutdown(ctx, func(ctx context.Context) error {
//		// Notify consumer to stop
//		notifyStop()
//
//		// Wait until consumer stopped
//		wg.Wait()
//
//		if err := c.ConsumerGroup.Close(); err != nil {
//			log.Error().Err(err).Msg("Close Kafka ConsumerGroup With error")
//			return err
//		}
//		return nil
//	})
//}
//
//// CloseFunc
//type CloseFunc func(context.Context) error
//
//// waitCloseOrShutdown 等待收到系统信号停止应用, 或者等待应用主动停止
//func waitCloseOrShutdown(ctx context.Context, closeFunc CloseFunc) {
//	// TODO: 超时时间可配置?
//	isCancel := shouldShutdown(ctx)
//	for _, flushFunc := range getFlushFunc(ctx) {
//		if err := shutdownServe(ctx, flushFunc, 3*time.Second); err != nil {
//			Warn(ctx).Err(err).Msg("Flush with error")
//		}
//	}
//	if isCancel {
//		if err := shutdownServe(ctx, closeFunc, 5*time.Second); err != nil {
//			Warn(ctx).Err(err).Msg("Close with error")
//		}
//	}
//}
//
//func shouldShutdown(ctx context.Context) (shutdown bool) {
//	// Wait for interrupt signal to gracefully shutdown the server with
//	// a timeout of 5 seconds.
//	quit := make(chan os.Signal, 1)
//
//	// kill (no param) default send syscall.SIGTERM
//	// kill -2 is syscall.SIGINT
//	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
//	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
//	select {
//	case <-ctx.Done():
//		log.Info().Msg("Consumer cancel")
//		return false
//	case <-quit:
//		return true
//	}
//}
//
//var flushFuncs = []CloseFunc{}
//
//func getFlushFunc(ctx context.Context) []CloseFunc {
//	return flushFuncs
//}
//
//func shutdownServe(ctx context.Context, closeFunc CloseFunc, timeout time.Duration) error {
//	ctx, cancel := context.WithTimeout(ctx, timeout)
//	defer cancel()
//
//	return closeFunc(ctx)
//}
//
//var saramaLoggerOnce = sync.Once{}
//
//type SaramaLogger struct{}
//
//func (SaramaLogger) Print(v ...interface{}) {
//	e := DebugWithoutCtx()
//	if e != nil {
//		e.Msg(fmt.Sprint(v...))
//	}
//}
//func (SaramaLogger) Printf(format string, v ...interface{}) {
//	DebugWithoutCtx().Msgf(format, v...)
//}
//func (SaramaLogger) Println(v ...interface{}) {
//	e := DebugWithoutCtx()
//	if e != nil {
//		e.Msg(fmt.Sprintln(v...))
//	}
//}
//
//// initSaramaLog 用 zerolog.Debug 替换调 sarama 使用的 log.
//func initSaramaLog() {
//	saramaLoggerOnce.Do(func() {
//		sarama.Logger = SaramaLogger{}
//	})
//}
