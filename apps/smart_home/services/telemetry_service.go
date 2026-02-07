package services

import (
	"fmt"

	"github.com/IBM/sarama"
)

// KafkaProducer - структура для асинхронного продюсера Kafka
type KafkaProducer struct {
	producer sarama.AsyncProducer
}

// NewKafkaProducer - конструктор для создания нового KafkaProducer
func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	// Конфигурация продюсера
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // Дожидаться подтверждения от всех брокеров
	config.Producer.Partitioner = sarama.NewRandomPartitioner // Случайная партиция
	config.Producer.Return.Successes = true                   // Возвращать успешные сообщения

	// Подключаемся к Kafka брокерам
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при создании продюсера: %v", err)
	}

	return &KafkaProducer{producer: producer}, nil
}

// SendMessage - метод отправки сообщения в Kafka (асинхронно)
func (k *KafkaProducer) SendMessage(topic, message string, headers []sarama.RecordHeader) {
	msg := &sarama.ProducerMessage{
		Topic:   topic,
		Value:   sarama.StringEncoder(message),
		Headers: headers,
	}

	// Отправляем сообщение в асинхронном режиме
	k.producer.Input() <- msg
}

// StartListener - метод для прослушивания канала успехов и ошибок
func (k *KafkaProducer) StartListener() {
	// Канал для успешных сообщений
	go func() {
		for {
			select {
			case successMsg := <-k.producer.Successes():
				fmt.Printf("Сообщение успешно отправлено в партицию %d с оффсетом %d\n", successMsg.Partition, successMsg.Offset)
			case errMsg := <-k.producer.Errors():
				fmt.Printf("Ошибка при отправке сообщения: %v\n", errMsg)
			}
		}
	}()
}

// Close - метод для закрытия продюсера
func (k *KafkaProducer) Close() error {
	return k.producer.Close()
}
