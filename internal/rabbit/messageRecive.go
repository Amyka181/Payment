package rabbit

import (
	"Payment/infrastructure/postgres"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

// TODO: везде, где функция может вернуть ошибку, нужно возвращать ошибку
func MessageReceive(db *postgres.DB) {
	conn, err := amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		//TODO: как тут, нет вохврата ошибка
		// а что произойдет если к реббиту подключиться не получится
		log.Println("Не удалось подключиться к RabbitMQ: %v", err)
	} else {
		log.Println("Успешно подключено к RabbitMQ")
	}
	//TODO: блок else лишний, так как мы можем просто написать
	//conn, err := amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
	//if err != nil {
	//	log.Println("Не удалось подключиться к RabbitMQ: %v", err)
	//	return
	//}
	//log.Println("Успешно подключено к RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		//TODO: как тут, нет вохврата ошибка
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
		//TODO: как тут, нет вохврата ошибка
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
		//TODO: как тут, нет вохврата ошибка
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
	//TODO: не работает таким образом, forever - это канал, в который не записываются сигналы
	// если хочешь грейсфул шотдаун, то нужно использовать канал os.Signal

}
