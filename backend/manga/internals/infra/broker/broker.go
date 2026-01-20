package broker

import (
	"log"
	"qaanii/shared/utils"

	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	BROKER_CHANNEL    string = "@BROKER/CHANNEL"
	BROKER_CONNECTION string = "@BROKER/CONNECTION"
)

func Broker(app *fiber.App) (*amqp.Connection, *amqp.Channel) {
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

	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Locals(BROKER_CONNECTION, conn)
		ctx.Locals(BROKER_CHANNEL, channel)

		return ctx.Next()
	})

	return conn, channel
}
