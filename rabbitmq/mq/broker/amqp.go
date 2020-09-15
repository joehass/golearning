package broker

import (
	"errors"
	"github.com/streadway/amqp"
	"log"
	"sync"
	"time"
)

type AmqpBrokerOptions struct {
	Url          string
	Exchange     string
	ExchangeType string
}

type AmqpBroker struct {
	m        sync.Mutex
	conn     *amqp.Connection
	notifies map[string]chan *amqp.Error
	options  *AmqpBrokerOptions
}

func NewAmqpBroker(option *AmqpBrokerOptions) *AmqpBroker {
	conn, err := amqp.Dial(option.Url)
	if err != nil {
		panic(err)
	}

	ab := &AmqpBroker{
		conn:     conn,
		options:  option,
		notifies: make(map[string]chan *amqp.Error),
	}

	go ab.keepAlive()

	return ab
}

func (a *AmqpBroker) keepAlive() {
	if a.conn != nil {
		cc := a.conn.NotifyClose(make(chan *amqp.Error))
		log.Printf("amqp conn close: %v", <-cc)
	}

	var err error
	if a.conn, err = amqp.Dial(a.options.url); err != nil {
		log.Printf("amqp redial faild: %v", err)
		time.AfterFunc(5*time.Second, a.keepAlive)
		return
	}

	log.Println("amqp redial success...")
	for _, n := range a.notifies {
		n <- amqp.ErrClosed
	}
	a.keepAlive()
}

func (a *AmqpBroker) Consume(queue *Queue) error {
	if a.conn == nil {
		return errors.New("conn is nil")
	}
	channel, err := a.conn.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	if err := channel.ExchangeDeclare(
		a.options.Exchange,
		a.options.ExchangeType,
		true, false, falase, false, nil); err != nil {
		return err
	}

	if _, err := channel.QueueDeclare(queue.Name, true, false, false, false, nil); err != nil {
		return err
	}

	if err := channel.QueueBind(queue.Name, queue.RouteKey, a.options.Exchange, false, nil); err != nil {
		return err
	}

	delivery, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	notify := make(chan *amqp.Error)
	defer close(notify)

	a.m.Lock()
	a.notifies[queue.Name] = notify
	a.m.Unlock()

	for {
		select {
		case err := <-notify:
			return err
		case d := <-delivery:
			switch status := queue.Handle(d.Body); status {
			case Retry:
				if err := a.retry(queue, d); err != nil {
					d.Nack(false, true)
				} else {
					d.Ack(false)
				}
			default:
				d.Ack(false)
			}
		}
	}
}

func (a *AmqpBroker) retry(queue *Queue, d amqp.Delivery) error {
	channel, err := a.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	retryCount, _ := d.Headers["x-retry-count"].(int32)

	if int(retryCount) >= len(queue.RetryQueue) {
		return nil
	}

	delay := queue.RetryQueue[retryCount]

	delayDuration := time.Duration(delay) * time.Millisecond
	delayQ := fmt.Sprintf("delay.%s.%s.%s", delayDuration.String(), a.options.Exchange, queue.Name)

	if _, err := channel.QueueDeclare(delayQ,
		true, false, false, false, amqp.Table{
			"x-dead-letter-exchange":    a.options.Exchange,
			"x-dead-letter-routing-key": queue.RouteKey,
			"x-message-ttl":             delay,
			"x-expires":                 delay * 2,
		}); err != nil {
		return err
	}

	return channel.Publish("", delayQ, false, false, amqp.Publishing{
		Headers:      amqp.Table{"x-retry-count": retryCount + 1},
		Body:         d.Body,
		DeliveryMode: amqp.Persistent,
	})
}

func (a *AmqpBroker) Publish(key string, body []byte) error {
	channel, err := a.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	if err := channel.ExchangeDeclare(
		a.options.Exchange,
		a.options.ExchangeType,
		true, false, false, false, nil,
	); err != nil {
		return err
	}

	return channel.Publish(a.options.Exchange, key, false, false, amqp.Publishing{
		Headers:      amqp.Table{},
		ContentType:  "",
		Body:         body,
		DeliveryMode: amqp.Persistent,
	})
}
