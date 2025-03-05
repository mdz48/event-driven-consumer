package core

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Database struct {
	Conn *sql.DB
}

func NewDatabase() *Database {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, name)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	conn.SetMaxOpenConns(10)

	if err := conn.Ping(); err != nil {
		log.Fatal(err)
	}

	return &Database{
		Conn: conn,
	}
}