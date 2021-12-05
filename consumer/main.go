package main

import (
	"log"
	"mq-create-excel/rabbit"
)

func main() {
	con := rabbit.CreateConnection()
	defer con.Close()

	channel := rabbit.CreateChannel(con)
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
			if string(message.Body) == "product-excel" {
				log.Printf("Gelen Mesaj: %s\n", message.Body)
				log.Printf("%s\n", "Excel işlemi başlıyor..")

			}
		}
	}()

	<-forever
}
