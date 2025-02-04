package main

import (
	"Payment/config"
	"Payment/infrastructure/postgres"
	"Payment/internal/rabbit"
	"log"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println("Ошибка загрузки конфигурации: %v", err)
	}

	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Println(err)
	}

	err = rabbit.MessageReceive(db)
	if err != nil {
		log.Println(err)
	}

}
