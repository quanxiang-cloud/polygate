package kafka

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/quanxiang-cloud/polygate/pkg/config"
)

// Record Record
type Record struct {
	RequestID       string
	UserID          string
	UserName        string
	OperationTime   int64
	OperationUA     string
	OperationType   string
	OperationModule string
	IP              string
	Detail          string
}

// Client Client
type Client struct {
	topic    string
	producer sarama.SyncProducer
}

// New New
func New(producer sarama.SyncProducer, conf config.Handler) *Client {
	return &Client{
		topic:    conf.Topic,
		producer: producer,
	}
}

// Send send
func (c *Client) Send(record *Record) error {

	recordByte, err := json.Marshal(record)
	if err != nil {
		return err
	}

	_, _, err = c.producer.SendMessage(&sarama.ProducerMessage{
		Topic: c.topic,
		Value: sarama.ByteEncoder(recordByte),
	})

	return err
}

func pre(conf config.Kafka) *sarama.Config {
	config := sarama.NewConfig()

	// TLS
	config.Net.TLS.Enable = conf.Sarama.Net.TLS.Enable
	config.Net.TLS.Config = conf.Sarama.Net.TLS.Config

	return config
}

// NewSyncProducer new sync producer
func NewSyncProducer(conf config.Kafka) (sarama.SyncProducer, error) {
	config := pre(conf)
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(conf.Broker, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

// NewAsyncProducer new async producer
func NewAsyncProducer(conf config.Kafka) (sarama.AsyncProducer, error) {
	config := pre(conf)
	config.Producer.Return.Successes = conf.Sarama.Producer.Return.Successes

	producer, err := sarama.NewAsyncProducer(conf.Broker, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}
