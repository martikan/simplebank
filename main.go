package main

import (
	"database/sql"
	"log"

	"github.com/martikan/simplebank/api"
	db "github.com/martikan/simplebank/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:aaa@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8085"
)

func main() {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}

}
