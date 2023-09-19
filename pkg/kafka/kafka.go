package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
)

type KafkaClient struct {
	consumer sarama.Consumer
	producer sarama.SyncProducer
}

func NewKafkaClient(brokerURLs []string) (*KafkaClient, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	consumer, err := sarama.NewConsumer(brokerURLs, config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kafka consumer: %v", err)
	}

	producer, err := sarama.NewSyncProducer(brokerURLs, config)
	if err != nil {
		return nil, fmt.Errorf("error creating Kafka producer: %v", err)
	}

	return &KafkaClient{
		consumer: consumer,
		producer: producer,
	}, nil
}

func (k *KafkaClient) SendMessageToTopic(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err := k.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("error sending message to topic %s: %v", topic, err)
	}

	return nil
}

func (k *KafkaClient) ConsumeFromTopic(topic string) (<-chan *sarama.ConsumerMessage, <-chan *sarama.ConsumerError) {
	partitionConsumer, err := k.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("error creating Kafka partition consumer for topic %s: %v", topic, err)
		return nil, nil
	}

	return partitionConsumer.Messages(), partitionConsumer.Errors()
}

func (k *KafkaClient) Close() {
	k.consumer.Close()
	k.producer.Close()
}
