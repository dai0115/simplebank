package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/simplebank/util"
)

// can be referenced from anywhere in the same package
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatalf("failed to load config file %#v", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("failed to connect to database %#v", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
