package config

import (
	"database/sql"
	"fmt"
	"log"
)

func ConnectDB() *sql.DB {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Welcome to the Todo APP")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)

	db, err := sql.Open(cfg.Driver, psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Database")
	return db
}
