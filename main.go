package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/simplebank/api"
	db "github.com/simplebank/db/sqlc"
	"github.com/simplebank/util"
)

func main() {
	config, err := util.LoadConfig("./")
	if err != nil {
		log.Fatalf("failed to load config file %s", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("failed to connect to database %s", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("cannot start server: %#v", err)
	}
}
