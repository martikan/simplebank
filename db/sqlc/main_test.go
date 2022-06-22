package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/martikan/simplebank/util"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {

	config, err := util.ConfigUtils.LoadConfig("../..")
	if err != nil {
		log.Fatal("Cannot load configuration file:", err)
	}

	testDB, err = sql.Open(config.DbDriver, config.DbUrl)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
