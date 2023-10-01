package migration

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
)

func RunMigrations(dbURL string) {
	fmt.Println("run migration")
	if err := getMigration(dbURL).Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

func getMigration(dbURL string) *migrate.Migrate {
	dir, _ := os.Getwd()
	if os.Getenv("env") == "local" {
		dir = strings.SplitAfter(dir, "rinha-backend")[0]
	}
	fmt.Println(dbURL)
	m, err := migrate.New(
		fmt.Sprintf("file://%s/migrations", dir),
		dbURL,
	)
	if err != nil {
		log.Fatal(err)
	}
	return m
}
