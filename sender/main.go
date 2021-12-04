package main

import (
	"fmt"
	"mq-create-excel/rabbit"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/streadway/amqp"
)

func main() {

	con := rabbit.CreateConnection()
	defer con.Close()

	channel := rabbit.CreateChannel(con)
	defer channel.Close()

	rabbit.CreateQueue("create-excel", channel)

	app := fiber.New()

	app.Use(
		logger.New(),
	)

	app.Get("/create-excel", func(c *fiber.Ctx) error {
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(c.Query("product")),
		}

		if err := channel.Publish(
			"",
			"create-excel",
			false,
			false,
			message,
		); err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})

	app.Listen(":3000")

}
