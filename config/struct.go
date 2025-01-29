package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
)

type Config struct {
	BDVersion
	BDInfo
}
type BDVersion struct {
	Version string
}

type BDInfo struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func LoadEnv() (*Config, error) {

	//err := godotenv.Load()
	//if err != nil {
	//	return nil, fmt.Errorf("Ошибка загрузки .env файла: %v", err)
	//}

	config := &Config{
		BDVersion{Version: os.Getenv("VERSION")},
		BDInfo{Host: os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME")},
	}

	return config, nil
}

func ConnectDB(cfg *Config) (*pgx.Conn, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	return conn, nil

}
