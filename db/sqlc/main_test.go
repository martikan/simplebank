package db

import (
	"database/sql"
	"fmt"
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

	url := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.DbUsr, config.DbPass, config.DbUrl, config.DbDb)
	fmt.Println(url)

	testDB, err = sql.Open(config.DbDriver, url)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
