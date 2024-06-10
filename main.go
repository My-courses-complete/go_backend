package main

import (
	"context"
	"log"

	"github.com/My-courses-complete/go_backend.git/api"
	"github.com/My-courses-complete/go_backend.git/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := util.LoadConfig(".")
	if err!= nil {
        log.Fatalf("cannot load config: %v", err)
    }
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}

	server := api.NewServer(conn)

	err = server.Run(config.ServerAddress)
	if err!= nil {
		log.Fatalf("cannot run server: %v", err)
	}
}
