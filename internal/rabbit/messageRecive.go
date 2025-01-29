package rabbit

import (
	"Payment/infrastructure/postgres"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

func MessageReceive(db *postgres.DB) {
	conn, err := amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Println("Не удалось подключиться к RabbitMQ: %v", err)
	} else {
		log.Println("Успешно подключено к RabbitMQ")
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Println("Не удалось открыть канал: %v", err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"my_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Не удалось создать очередь: %v", err)
	}

	messages, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Не удалось начать потребление: %v", err)
	}
	var UserUp *postgres.UpdateUser
	forever := make(chan bool)
	go func() {
		for message := range messages {
			log.Printf("Получено сообщение: %s", message.Body)
			err = json.Unmarshal(message.Body, &UserUp)
			if err != nil {
				return
			}
			err = db.ChangeBalance(UserUp)
			if err != nil {
				return
			}
		}
	}()

	log.Println("Ожидание сообщений. Нажмите CTRL+C для выхода.")
	<-forever

}
