package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/samuelowad/bank/pkg/util"
	"log"
	"os"
	"testing"
)

var testQueries *Queries

var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../../../")
	if err != nil {
		log.Fatal("can't load config'", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
