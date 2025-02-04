package rabbit

import (
	"Payment/infrastructure/postgres"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// TODO: везде, где функция может вернуть ошибку, нужно возвращать ошибку
func MessageReceive(db *postgres.DB) error {
	conn, err := amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		//TODO: как тут, нет вохврата ошибка
		// а что произойдет если к реббиту подключиться не получится
		return fmt.Errorf("Не удалось подключиться к RabbitMQ: %v", err)
	}
	log.Println("Успешно подключено к RabbitMQ")

	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		//TODO: как тут, нет вохврата ошибка
		return fmt.Errorf("Не удалось открыть канал: %v", err)
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
		return fmt.Errorf("Не удалось создать очередь: %v", err)
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
		return fmt.Errorf("Не удалось начать потребление: %v", err)
	}

	var UserUp *postgres.UpdateUser
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("Завершаем программу...")

	return nil

	//TODO: не работает таким образом, forever - это канал, в который не записываются сигналы
	// если хочешь грейсфул шотдаун, то нужно использовать канал os.Signal

}
