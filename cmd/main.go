package main

import (
	"Payment/infrastructure/postgres"
	"Payment/internal/rabbit"
)

func main() {

	db := postgres.NewDB()

	rabbit.MessageReceive(db)
}
