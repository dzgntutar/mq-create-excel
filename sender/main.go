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

type MyApp struct {
	app        *fiber.App
	db         *sql.DB
	connection *amqp.Connection
	channel    *amqp.Channel
}

var (
	dbError error
)

func main() {
	myApp := MyApp{}

	var host, port, user, password, dbName = "localhost", "5432", "postgres", "Password1*", "my-db"

	myApp.initialize(host, port, user, password, dbName)
	defer myApp.connection.Close()
	defer myApp.channel.Close()

	myApp.routes()

	myApp.app.Listen(":3000")

}

func (app *MyApp) initialize(host, port, user, password, dbName string) {
	//Db
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	app.db, dbError = sql.Open("postgres", connectionString)

	if dbError != nil {
		panic(dbError)
	}

	//Rebbit
	app.connection = rabbit.CreateConnection()
	app.channel = rabbit.CreateChannel(app.connection)

	rabbit.CreateQueue("create-excel", app.channel)

	//Fiber
	app.app = fiber.New()

	app.app.Use(
		logger.New(),
	)
}

func (app *MyApp) routes() {
	app.app.Get("/create-excel", func(c *fiber.Ctx) error {
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

		var product = models.Product{
			Name:  "Telefon",
			Price: 10000,
			Stock: 5,
		}

		var productRepository = repository.ProductRepository{
			DB: app.db,
		}

		productRepository.Insert(product)

		return nil
	})
}
