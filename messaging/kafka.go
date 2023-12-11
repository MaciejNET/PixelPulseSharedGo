package messaging

import (
	"github.com/IBM/sarama"
	"log"
	"sync"
)

type MessageHandler func(msg *sarama.ConsumerMessage)

type MessageTypeHandlers map[string]MessageHandler

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{producer: producer}, nil
}

func (kp *KafkaProducer) SendMessage(topic string, key string, value string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(value),
	}

	_, _, err := kp.producer.SendMessage(msg)
	return err
}

func (kp *KafkaProducer) Close() {
	kp.producer.Close()
}

type KafkaConsumer struct {
	consumer sarama.Consumer
	handlers MessageTypeHandlers
}

func NewKafkaConsumer(brokers []string, groupID string, topics []string, handlers MessageTypeHandlers) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{consumer: consumer, handlers: handlers}, nil
}

func (kc *KafkaConsumer) Consume(topics []string) {
	partitions, err := kc.consumer.Partitions(topics[0])
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(len(partitions))

	for _, partition := range partitions {
		var pc sarama.PartitionConsumer
		pc, err = kc.consumer.ConsumePartition(topics[0], partition, sarama.OffsetNewest)
		if err != nil {
			log.Fatal(err)
		}

		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for {
				select {
				case msg := <-pc.Messages():
					kc.handleMessage(msg)
				case err = <-pc.Errors():
					log.Println("Error: ", err)
				}
			}
		}(pc)
	}

	wg.Wait()
}

func (kc *KafkaConsumer) Close() {
	kc.consumer.Close()
}

func (kc *KafkaConsumer) handleMessage(msg *sarama.ConsumerMessage) {
	messageType := string(msg.Key)

	if handler, ok := kc.handlers[messageType]; ok {
		handler(msg)
	} else {
		log.Printf("No handler found for message type: %s\n", messageType)
	}
}
