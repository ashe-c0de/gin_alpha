package utils

import (
	"fmt"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"my_destributed_project/configs"
	"my_destributed_project/pkg/log"
	"sync"
)

var (
	once          sync.Once
	KafkaProducer sarama.SyncProducer
)

// GetKafkaProducer 懒加载 Kafka 生产者
func GetKafkaProducer() (sarama.SyncProducer, error) {
	var err error
	once.Do(func() {
		config := sarama.NewConfig()
		config.Producer.Return.Successes = true

		url := fmt.Sprintf("%s:%s", configs.AppConfig.Server.KafkaHost, configs.AppConfig.Server.KafkaPort)
		KafkaProducer, err = sarama.NewSyncProducer([]string{url}, config)
	})

	return KafkaProducer, err
}

// ProduceMessage 发送 Kafka 消息
func ProduceMessage(topic, message string) error {
	producer, err := GetKafkaProducer()
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Logger.Info("Kafka 消息发送失败", zap.Any("err", err))
		return err
	}

	log.Logger.Info("Kafka 消息发送成功",
		zap.Int32("Partition", partition),
		zap.Int64("Offset", offset))
	return nil
}
