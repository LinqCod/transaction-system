package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
)

func InitDB() (*sql.DB, error) {
	username := viper.GetString("POSTGRES_USER")
	password := viper.GetString("POSTGRES_PASSWORD")
	port := viper.GetString("POSTGRES_PORT")
	dbname := viper.GetString("POSTGRES_DB")

	url := fmt.Sprintf("postgres://%v:%v@db:%v/%v?sslmode=disable", username, password, port, dbname)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("validation of db parameters failed due to error: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to open db connection due to err: %v", err)
	}

	log.Println("postgres db connected successfully!")
	return db, nil
}
