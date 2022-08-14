package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries

var testDB *sql.DB

const (
	dbDriver = "postgres"
	dbSource = "postgres://postgres:postgres@localhost:4321/bank?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
