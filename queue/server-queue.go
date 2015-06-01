package queue

import "github.com/streadway/amqp"

//ListenQueue ...
func ListenQueue(url, queueName string) (*amqp.Connection, <-chan amqp.Delivery, error) {
	conn, err := amqp.Dial(url)

	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return conn, nil, err
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return conn, nil, err
	}

	err = ch.Qos(
		3,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return conn, nil, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return conn, nil, err
	}
	return conn, msgs, nil
}
