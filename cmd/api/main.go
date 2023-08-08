package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"user/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB *sql.DB

	Models data.Models
}

func main() {
	log.Println("Starting API server on port", webPort)

	conn := connectToDb()

	if conn == nil {
		log.Fatal("Could not connect to database")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	fmt.Println("dsn user")
	fmt.Println(dsn)
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDb() *sql.DB {
	dsn := os.Getenv("DSN")

	fmt.Println("dsn users")
	fmt.Println(dsn)
	for {
		connection, err := openDB(dsn)

		if err != nil {
			log.Println("Postgres is unavailable - sleeping")
			counts++

		} else {
			log.Println("Postgres is available - continuing")
			return connection
		}

		if counts > 10 {
			log.Fatal("Could not connect to database", err)
			return nil
		}

		log.Println("Sleeping for two seconds...")

		time.Sleep(2 * time.Second)

		continue
	}
}
