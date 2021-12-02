package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/streadway/amqp"
)

func main() {
	//rabbitServerURL := os.Getenv("RABBIT_SERVER_URL")
	rabbitServerURL := "amqp://guest:guest@localhost:5672"
	connection, err := amqp.Dial(rabbitServerURL)

	if err != nil {
		fmt.Println("Rabbit mq sunucusuna baglanılamadı..")
		panic(err)
	}

	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		fmt.Println("Kanal olusturulurken hata meydana geldi")
		panic(err)

	}

	defer channel.Close()

	_, err = channel.QueueDeclare("create-excel", true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

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
			return err
		}

		return nil
	})

	app.Listen(":3000")

}
