package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/samuelowad/bank/pkg/util"

	"github.com/samuelowad/bank/api"
	db "github.com/samuelowad/bank/pkg/db/sqlc"

	_ "github.com/golang/mock/mockgen/model"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load config'", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server")
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
