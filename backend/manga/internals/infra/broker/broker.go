package broker

import (
	"log"
	"qaanii/shared/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	BROKER_CHANNEL    string = "@BROKER/CHANNEL"
	BROKER_CONNECTION string = "@BROKER/CONNECTION"
)

func Broker() (*amqp.Connection, *amqp.Channel) {
	envs := utils.Utils{}.Envs()
	broker_url := envs["broker_url"]

	if len(broker_url) == 0 {
		log.Panicln("Broker | URL not found")
		return nil, nil
	}

	conn, err := amqp.Dial(broker_url)
	if err != nil {
		log.Panicf("Broker | Connection with broker failed, error: %+v", err)
		return nil, nil
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Panicf("Broker | Channel connection failed, error: %+v", err)
		return nil, nil
	}

	return conn, channel
}
