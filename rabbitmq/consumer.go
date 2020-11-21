package rabbitmq

import (
	"errors"
	"golearning/rabbitmq/mq/broker"
	"golearning/rabbitmq/mq/help"
	"log"
	"time"
)

var RetryError = errors.New("job retry")

type ExchangeOption struct {
	Name string
	Type string
}

type ConsumerOptions struct {
	ExchangeOpt *ExchangeOption
	BrokerURL   string
}

type Consumer struct {
	broker broker.Broker
}

func NewConsumer(opt *ConsumerOptions) *Consumer {
	broker := broker.NewAmqpBroker(&broker.AmqpBrokerOptions{
		Url:          opt.BrokerURL,
		Exchange:     opt.ExchangeOpt.Name,
		ExchangeType: opt.ExchangeOpt.Type,
	})

	consumer := &Consumer{broker: broker}

	return consumer
}

type Job func([]byte) error

type params struct {
	retryQueue []int64
}

type Param func(*params)

func Retry(startegy help.RetryStrategy, retry help.Retry) Param {
	return func(p *params) {
		if startegy == help.CUSTOMQUEUE {
			for _, delay := range retry.Queue {
				d, err := time.ParseDuration(delay)
				if err != nil {
					panic(err)
				}
				p.retryQueue = append(p.retryQueue, int64(d/time.Millisecond))
			}
			return
		}
		d, err := time.ParseDuration(retry.Delay)
		if err != nil {
			panic(err)
		}
		p.retryQueue = help.GetRetryQueue(int64(d/time.Millisecond), retry.Max, startegy)
	}
}

func evaParam(param []Param) *params {
	ps := &params{}
	for _, p := range param {
		p(ps)
	}
	return ps
}

func (c *Consumer) LaunchJob(key, queue string, job Job, param ...Param) {
	ps := evaParam(param)

	q := &broker.Queue{
		Name:       queue,
		RouteKey:   key,
		RetryQueue: ps.retryQueue,
		Handle: func(body []byte) broker.Status {
			var err error
			switch err = job(body); err {
			case RetryError:
				return broker.Retry
			default:
				return broker.Success
			}
		},
	}

	for {
		log.Printf("job %s start consume...", queue)
		if err := c.broker.Consume(q); err != nil {
			log.Printf("job %s consume error: %v ,retrying consume after 30s", queue, err)
			time.Sleep(30 * time.Second)
		}
	}
}
