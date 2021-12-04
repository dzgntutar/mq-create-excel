package main

import (
	"database/sql"
	"fmt"
	"mq-create-excel/models"
	"mq-create-excel/rabbit"
	"mq-create-excel/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

type app struct {
	fiber      *fiber.App
	db         *sql.DB
	connection *amqp.Connection
	channel    *amqp.Channel
}

var (
	dbError error
)

func main() {
	app := app{}

	app.initialize()

	app.routes()

	app.run()

}

func (app *app) initialize() {
	//Db
	var host, port, user, password, dbName = "localhost", "5432", "postgres", "Password1*", "my-db"
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	app.db, dbError = sql.Open("postgres", connectionString)

	if dbError != nil {
		panic(dbError)
	}

	//Rabbit
	app.connection = rabbit.CreateConnection()
	defer app.connection.Close()
	app.channel = rabbit.CreateChannel(app.connection)
	defer app.channel.Close()

	rabbit.CreateQueue("create-excel", app.channel)

	//Fiber
	app.fiber = fiber.New()

	app.fiber.Use(
		logger.New(),
	)
}

func (app *app) routes() {
	app.fiber.Get("/", func(c *fiber.Ctx) error {
		// message := amqp.Publishing{
		// 	ContentType: "text/plain",
		// 	Body:        []byte(c.Query("product")),
		// }

		// if err := myApp.Channel.Publish(
		// 	"",
		// 	"create-excel",
		// 	false,
		// 	false,
		// 	message,
		// ); err != nil {
		// 	fmt.Println(err)
		// 	return err
		// }

		var repository = repository.ProductRepository{
			DB: app.db,
		}

		products, err := repository.GetAll()
		if err != nil {
			return err
		}
		c.JSON(products)

		return nil
	})

	app.fiber.Post("/product", func(c *fiber.Ctx) error {
		p := models.Product{}

		if err := c.BodyParser(&p); err != nil {
			return err
		}

		var repository = repository.ProductRepository{
			DB: app.db,
		}

		repository.Insert(p)
		return nil
	})
}

func (app *app) run() {
	app.fiber.Listen(":3000")
}
