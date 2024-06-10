package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/My-courses-complete/go_backend.git/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	if err!= nil {
        log.Fatalf("cannot load config: %v", err)
    }
	ctx := context.Background()
	testDB, err = pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer testDB.Close()
	
	testQueries = New()

	os.Exit(m.Run())

}