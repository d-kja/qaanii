package broker

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Broker(broker_url string) (*amqp.Connection, *amqp.Channel) {
	if len(broker_url) == 0 {
		log.Panicln("[SHARED/BROKER] - URL not found")
		return nil, nil
	}

	conn, err := amqp.Dial(broker_url)
	if err != nil {
		log.Panicf("[SHARED/BROKER] - Connection with broker failed, error: %+v", err)
		return nil, nil
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Panicf("[SHARED/BROKER] - Channel connection failed, error: %+v", err)
		return nil, nil
	}

	return conn, channel
}
