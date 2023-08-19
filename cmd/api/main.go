package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"user/data"
	user "user/users"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"google.golang.org/grpc"
)

const webPort = "80"
const grpcPort = "5002"

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

	go app.GRPCListen()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}

func (app *Config) GRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))

	if err != nil {
		log.Fatal("Listener error", err)
	}

	s := grpc.NewServer()

	user.RegisterUserServiceServer(s, &UserServer{
		Models: app.Models,
	})

	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: ", err)
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
