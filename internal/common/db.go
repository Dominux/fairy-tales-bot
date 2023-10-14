package common

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

func NewDB(dbUser string, dbPswd string, dbName string) *sqlx.DB {
	dsl := fmt.Sprintf("postgres://%s:%s@ft-db/%s?sslmode=disable", dbUser, dbPswd, dbName)

	for {
		db, err := sqlx.Connect("pgx", dsl)

		if err == nil {
			log.Print("Established db connection")
			return db
		}

		log.Print("Failed to establish db connection, trying again")
		time.Sleep(3 * time.Second)
	}
}
