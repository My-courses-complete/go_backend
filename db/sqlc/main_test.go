package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:password@localhost:5432/go_course_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *pgx.Conn

func TestMain(m *testing.M) {
	var err error
	ctx := context.Background()
	testDB, err = pgx.Connect(ctx, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer testDB.Close(ctx)
	
	testQueries = New(testDB)

	os.Exit(m.Run())

}