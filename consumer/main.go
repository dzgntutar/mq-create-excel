package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {
	rabbitServerURL := "amqp://guest:guest@localhost:5672"

	connection, err := amqp.Dial(rabbitServerURL)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	messages, err := channel.Consume(
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
			log.Printf("Gelen Mesaj: %s\n", message.Body)
		}
	}()

	<-forever
}
