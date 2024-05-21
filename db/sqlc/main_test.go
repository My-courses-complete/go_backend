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

func TestMain(m *testing.M) {
	ctx := context.Background()
	db, err := pgx.Connect(ctx, dbSource)
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer db.Close(ctx)

	testQueries = New(db)

	os.Exit(m.Run())

}