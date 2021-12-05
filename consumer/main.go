package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"log"
	"mq-create-excel/rabbit"
	"mq-create-excel/repository"
)

type app struct {
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

	var (
		repository = repository.ProductRepository{
			DB: app.db,
		}
	)

	app.connection = rabbit.CreateConnection()
	defer app.connection.Close()

	app.channel = rabbit.CreateChannel(app.connection)
	defer app.channel.Close()

	messages, err := app.channel.Consume(
		"create-excel",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("Rabbit mq baglantisi basarili..")
	log.Println("Mesajlar okunmaya başlıyor..")

	forever := make(chan bool)

	go func() {
		for message := range messages {
			if string(message.Body) == "product-excel" {
				log.Printf("%s\n", "Excel işlemi başlıyor..")

				products, err := repository.GetAll()

				if err != nil {
					log.Println(err)
				}

				log.Println(len(products))
				// create excel , save in a folder after that insert in a db

			}
		}
	}()

	<-forever
}

func (app *app) initialize() {
	//Db
	var host, port, user, password, dbName = "localhost", "5432", "postgres", "Password1*", "my-db"
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	app.db, dbError = sql.Open("postgres", connectionString)

	if dbError != nil {
		panic(dbError)
	}
}
