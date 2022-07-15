package main

import (
	"database/sql"
	"fmt"
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

	url := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.DbUsr, config.DbPass, config.DbUrl, config.DbDb)

	conn, err := sql.Open(config.DbDriver, url)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	store := db.NewStore(conn)

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot create server", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}

}
