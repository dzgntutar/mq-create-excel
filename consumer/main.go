package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"github.com/xuri/excelize/v2"
	"log"
	"mq-create-excel/rabbit"
	"mq-create-excel/repository"
	"strconv"
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

				// create excel , save in a folder

				f := excelize.NewFile()
				f.NewSheet("Products")

				for i, product := range products {

					err = f.SetCellValue("Products", "A"+strconv.Itoa(i+1), product.Id) //name
					/*if err != nil {
						log.Println(err)
						log.Println("A" + string(i+1))
					}*/
					f.SetCellValue("Products", "B"+strconv.Itoa(i+1), product.Name)  //name
					f.SetCellValue("Products", "C"+strconv.Itoa(i+1), product.Price) //price
					f.SetCellValue("Products", "D"+strconv.Itoa(i+1), product.Stock) //stock
				}

				if err := f.SaveAs("products.xlsx"); err != nil {
					log.Fatal(err)
				}
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
