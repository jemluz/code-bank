package kafka

import ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

type KafkaProducer struct {
	Producer *ckafka.Producer
}

func NewKafkaProducer() KafkaProducer {
	return KafkaProducer{}
}

func (k *KafkaProducer) SetupProducer(bootstrapServer string) {
	configMap := &ckafka.configMap{
		"bootstrap.servers":bootstrapServer,
	}
	k.Producer, _ = ckafka.NewKafkaProducer(configMap)
}

func (k *KafkaProducer) Publish(msg string, topic string) error {
	message := &ckafka.Messsage{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Value: []byte(msg)
	}

	err := k.Producer.Produce(message, deliveryChan: nil)

	if err != nil {
		return err
	}

	return nil
}