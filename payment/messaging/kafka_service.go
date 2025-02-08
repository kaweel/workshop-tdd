package messaging

type KafkaProducer interface {
	Publish(RequestPublish) error
}

type kafkaProducer struct {
}
type RequestPublish struct {
	Topic   string
	Key     string
	Message any
}

func NewKafkaProducer() KafkaProducer {
	return &kafkaProducer{}
}

func (s *kafkaProducer) Publish(r RequestPublish) error {
	return nil
}
