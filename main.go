package main

import (
	"context"
	"log"

	"github.com/My-courses-complete/go_backend.git/api"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:password@localhost:5432/go_course_bank?sslmode=disable"
	serverAddr = "0.0.0.0:8080"
)

func main() {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	server := api.NewServer(conn)

	err = server.Run(serverAddr)
	if err!= nil {
		log.Fatalf("cannot run server: %v", err)
	}
}
