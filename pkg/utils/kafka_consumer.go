package utils

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"my_destributed_project/configs"
	"my_destributed_project/pkg/log"
)

// ConsumerHandler 处理 Kafka 消息
type ConsumerHandler struct {
	ready chan bool
}

// NewConsumerHandler 创建一个新的消费者处理器
func NewConsumerHandler() *ConsumerHandler {
	return &ConsumerHandler{
		ready: make(chan bool),
	}
}

func (h *ConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	close(h.ready)
	// 使用正确的格式记录Claims
	claims := make(map[string][]int32)
	for topic, partitions := range session.Claims() {
		claims[topic] = partitions
	}

	log.Logger.Info("消费者设置完成",
		zap.String("memberID", session.MemberID()),
		zap.Any("claims", claims),
	)
	return nil
}

func (h *ConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	h.ready = make(chan bool)
	log.Logger.Info("消费者清理完成", zap.String("memberID", session.MemberID()))
	return nil
}

func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log.Logger.Info("开始消费分区",
		zap.String("topic", claim.Topic()),
		zap.Int32("partition", claim.Partition()),
		zap.Int64("initialOffset", claim.InitialOffset()),
	)

	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Logger.Info("消息通道已关闭",
					zap.String("topic", claim.Topic()),
					zap.Int32("partition", claim.Partition()))
				return nil
			}

			// 使用字符串类型记录消息内容，避免zap直接序列化二进制数据
			keyStr := ""
			if message.Key != nil {
				keyStr = string(message.Key)
			}

			valueStr := ""
			if message.Value != nil {
				valueStr = string(message.Value)
			}

			log.Logger.Info("Kafka 消息消费成功",
				zap.String("topic", message.Topic),
				zap.Int32("partition", message.Partition),
				zap.Int64("offset", message.Offset),
				zap.String("key", keyStr),
				zap.String("message", valueStr),
				zap.Time("timestamp", message.Timestamp),
			)

			// 标记该消息已处理（确认已被消费）
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			log.Logger.Info("消费会话已结束",
				zap.String("topic", claim.Topic()),
				zap.Int32("partition", claim.Partition()))
			return nil
		}
	}
}

// StartConsumer 启动 Kafka 消费者
func StartConsumer(ctx context.Context, groupID string) {
	// 如果未指定消费者组ID，使用默认值
	if groupID == "" {
		groupID = "my_default_consumer_group"
	}

	config := sarama.NewConfig()
	config.Version = sarama.V2_8_1_0

	// 设置消费者组策略
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(), // 使用轮询策略，确保分配到分区
	}

	// 从最早的消息开始消费，确保不会漏掉任何消息
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// 设置重试策略
	config.Consumer.Retry.Backoff = 1 * time.Second

	// 启用自动提交
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	// 获取 Kafka 连接地址
	brokers := []string{fmt.Sprintf("%s:%s", configs.AppConfig.Server.KafkaHost, configs.AppConfig.Server.KafkaPort)}

	log.Logger.Info("连接Kafka",
		zap.Strings("brokers", brokers),
		zap.String("groupID", groupID),
	)

	// 创建 Kafka 客户端
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Logger.Fatal("创建 Kafka 客户端失败", zap.Error(err))
		return
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Logger.Error("关闭 Kafka 客户端失败", zap.Error(err))
		}
	}()

	// 获取所有可用的 Kafka 主题
	allTopics, err := client.Topics()
	if err != nil {
		log.Logger.Fatal("获取 Kafka topics 失败", zap.Error(err))
		return
	}

	// 过滤掉系统主题（以下划线开头的），只消费业务主题
	var topics []string
	for _, topic := range allTopics {
		if !strings.HasPrefix(topic, "_") {
			topics = append(topics, topic)
		}
	}

	if len(topics) == 0 {
		log.Logger.Fatal("没有可用的业务主题")
		return
	}

	log.Logger.Info("开始消费 Kafka topics", zap.Strings("topics", topics))

	// 创建 Kafka 消费者组
	consumerGroup, err := sarama.NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		log.Logger.Fatal("Kafka 消费者创建失败", zap.Error(err))
		return
	}
	defer func() {
		if err := consumerGroup.Close(); err != nil {
			log.Logger.Error("关闭 Kafka 消费者失败", zap.Error(err))
		}
	}()

	// 检查消费者组错误
	go func() {
		for err := range consumerGroup.Errors() {
			log.Logger.Error("消费者组错误", zap.Error(err))
		}
	}()

	handler := NewConsumerHandler()

	// 消费循环
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Logger.Info("Kafka 消费者收到退出信号，停止消费")
				return
			default:
				// 消费消息
				log.Logger.Info("开始一轮消费", zap.Strings("topics", topics))

				err := consumerGroup.Consume(ctx, topics, handler)
				if err != nil {
					if errors.Is(err, context.Canceled) {
						//if err == context.Canceled {
						log.Logger.Info("消费上下文已取消")
						return
					}

					log.Logger.Error("Kafka 消费失败，2秒后重试", zap.Error(err))

					// 等待2秒再重试，避免频繁重试
					select {
					case <-time.After(2 * time.Second):
						// 继续重试
					case <-ctx.Done():
						return
					}
				} else {
					// 这个分支通常在rebalance后执行
					log.Logger.Info("消费会话已结束，准备开始新的消费会话")
				}
			}
		}
	}()

	// 等待外部上下文结束
	<-ctx.Done()
	log.Logger.Info("Kafka 消费服务正在关闭")
}
