package broker

import (
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/n0w4/gombackoff"
	"github.com/n0w4/gomj2k/model"
)

type kafkaBroker struct {
	bootstrapServer string
	producer        *kafka.Producer
}

func NewKafka(bootstrapServer string) *kafkaBroker {
	broker := kafkaBroker{
		bootstrapServer: bootstrapServer,
	}
	return &broker
}

func (k *kafkaBroker) WithProducer() *kafkaBroker {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":     k.bootstrapServer,
		"enable.idempotence":    true,
		"request.required.acks": -1,
	})
	if err != nil {
		log.Printf("%v", err)
	}

	k.producer = p

	go producerErrorListener(p)

	return k
}

func (k *kafkaBroker) Publish(messages ...model.StructuredMessage) {
	retryPolicy := gombackoff.NewRetryPolicy().ExponentialBackoff(5).Times(3)

	deliveryChan := make(chan kafka.Event, 1000)
	finishChan := make(chan bool)

	go listenDeliveryReport(deliveryChan, finishChan)

	headers := []kafka.Header{}

	for _, message := range messages {
		if message.Topic == "" {
			log.Printf("topic is empty, skipping message: %v", string(message.Payload))
			continue
		}
		for _, h := range message.Headers {
			headers = append(headers, kafka.Header{Key: h.Key, Value: []byte(h.Value)})
		}

		for {

			err := k.producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &message.Topic, Partition: kafka.PartitionAny},
				Value:          message.Payload,
				Headers:        headers,
				Key:            message.Key,
				Timestamp:      time.Now(),
			}, deliveryChan)

			if err != nil {
				if err.(kafka.Error).Code() == kafka.ErrQueueFull {
					if ok := retryPolicy.Wait(); ok {
						continue
					}
				}
				log.Printf("Failed to produce message: %v\n", err)
			}

			break
		}
	}

	for k.producer.Flush(5000) > 0 {
		log.Println("Flushing")
	}

	close(deliveryChan)
	<-finishChan
}

func producerErrorListener(p *kafka.Producer) {
	for e := range p.Events() {
		switch ev := e.(type) {
		case kafka.Error:
			log.Printf("Error: %v\n", ev)
		default:
			log.Printf("Ignored event: %s\n", ev)
		}
	}
}

func listenDeliveryReport(deliveryChan <-chan kafka.Event, finishChan chan<- bool) {
	for e := range deliveryChan {
		switch ev := e.(type) {
		case *kafka.Message:
			m := ev
			if m.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
			} else {
				log.Printf("Delivered message to topic %s [%d] at offset %v\n",
					*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
			}
		default:
			log.Printf("Ignored event: %s\n", ev)
		}
	}
	finishChan <- true
}
