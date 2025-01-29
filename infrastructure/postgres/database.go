package postgres

import (
	"Payment/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
)

type DB struct {
	Conn *pgx.Conn
}

type User struct {
	ID      int
	Balance int
}

type UpdateUser struct {
	ID     int
	Change int
}

func NewDB() *DB {

	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}
	log.Println(cfg)
	conn, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	fmt.Println("Успешное подключение к базе данных!")
	return &DB{Conn: conn}
}

func (db *DB) ChangeBalance(UserUp *UpdateUser) error {
	tx, err := db.Conn.Begin(context.Background())
	if err != nil {
		log.Fatal("Unable to start transaction:", err)
	}

	_, err = db.ShowBalanceTx(tx, UserUp.ID)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	reqSql := "UPDATE public.users SET balance=balance+$1 WHERE id=$2"
	_, err = tx.Exec(context.Background(), reqSql, UserUp.Change, UserUp.ID)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	tx.Commit(context.Background())

	return nil

}

func (db *DB) ShowBalanceTx(tx pgx.Tx, id int) (int, error) {

	var user User
	reqSql := "SELECT id, balance FROM public.users WHERE id=$1"
	err := tx.QueryRow(context.Background(), reqSql, id).Scan(&user.ID, &user.Balance)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return 0, err
		}
		return 0, err
	}
	return user.Balance, nil
}
