package rabbit

import (
	"fmt"

	"github.com/streadway/amqp"
)

// type Rabbit struct {
// 	connection *amqp.Connection
// 	channel    *amqp.Channel
// }

func CreateConnection() *amqp.Connection {
	rabbitServerURL := "amqp://guest:guest@localhost:5672"
	connection, err := amqp.Dial(rabbitServerURL)
	if err != nil {
		fmt.Println("Rabbit mq sunucusuna baglanılamadı..")
		panic(err)
	}

	return connection
}

func CreateChannel(connection *amqp.Connection) *amqp.Channel {
	channel, err := connection.Channel()
	if err != nil {
		fmt.Println("Kanal olusturulurken hata meydana geldi")
		panic(err)
	}

	return channel
}

func CreateQueue(name string, channel *amqp.Channel) {
	_, err := channel.QueueDeclare(name, true, false, false, false, nil)

	if err != nil {
		panic(err)
	}
}
