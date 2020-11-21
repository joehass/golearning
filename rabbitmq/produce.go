package rabbitmq

import (
	"encoding/json"
	"golearning/rabbitmq/mq/broker"
)

type ProducerOptions struct {
	ExchangeOpt *ExchangeOption
	BrokerURL   string
}

type Producer struct {
	broker broker.Broker
}

func NewProducer(opt *ProducerOptions) *Producer {
	broker := broker.NewAmqpBroker(&broker.AmqpBrokerOptions{
		Url:          opt.BrokerURL,
		Exchange:     opt.ExchangeOpt.Name,
		ExchangeType: opt.ExchangeOpt.Type,
	})

	return &Producer{broker: broker}
}

func (p *Producer) Publish(key string, data interface{}) error {
	var body []byte
	switch d := data.(type) {
	case string:
		body = []byte(d)
	default:
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		body = b
	}

	return p.broker.Publish(key, body)
}
