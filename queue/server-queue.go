package queue

import (
	"log"

	"github.com/streadway/amqp"
)

/*Queue ...*/
type Queue struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	Out     <-chan amqp.Delivery
}

//ListenQueue ...
func ListenQueue(url, queueName string) (*Queue, error) {
	conn, err := amqp.Dial(url)

	if err != nil {
		return &Queue{nil, nil, nil}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return &Queue{conn, nil, nil}, err
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
		return &Queue{conn, ch, nil}, err
	}

	err = ch.Qos(
		3,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return &Queue{conn, ch, nil}, err
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
		return &Queue{conn, ch, nil}, err
	}
	return &Queue{conn, ch, msgs}, nil
}

/*PublishQueue ...*/
func PublishQueue(q *Queue) {
	go func() {
		for msg := range q.Out {
			go func(m *amqp.Delivery) {
				log.Printf("Reply to: %s", m.ReplyTo)
				log.Printf("Correlation id: %s", m.CorrelationId)
				log.Printf("Received a message: %s", m.Body)

				err := q.Channel.Publish(
					"",        // exchange
					m.ReplyTo, // routing key
					false,     // mandatory
					false,     // immediate
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: m.CorrelationId,
						Body:          m.Body,
					})

				if err != nil {
					log.Println("Error publishing to queue: "+m.ReplyTo, err)
				}
			}(&msg)
		}
	}()

}
