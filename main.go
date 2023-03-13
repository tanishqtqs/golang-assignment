package main

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"z_test/db"
	"z_test/router"
)

func main() {

	db := db.DBConnect()
	driver, dbErr := postgres.WithInstance(db.DB, &postgres.Config{})
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migration", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	er := m.Up()
	if er == migrate.ErrNoChange {
		//
	}
	logrus.Print("MIGRATIONS UP...")
	r := router.Router()
	logrus.Print("SERVER UP...")
	err1 := http.ListenAndServe(":8080", r)
	if err1 != nil {
		log.Fatal("Error")
		return
	}
}
