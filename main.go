package main

import (
	"database/sql"
	"log"

	"github.com/martikan/simplebank/api"
	db "github.com/martikan/simplebank/db/sqlc"
	"github.com/martikan/simplebank/util"

	_ "github.com/lib/pq"
)

func main() {

	config, err := util.ConfigUtils.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load configuration file:", err)
	}

	conn, err := sql.Open(config.DbDriver, config.DbUrl)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}

}
