package main

import (
	"database/sql"
	"github.com/bless-rng/marmot/pkg"
	marmotDialect "github.com/bless-rng/marmot/pkg/dialect"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"os"
	"strings"
)

func main() {
	db, err := sql.Open("pgx", "postgres://user:password@localhost:9902/database")
	if err != nil {
		panic(err)
	}

	migrator, err := marmot.NewMigrator(marmotDialect.NewPostgreSQLDialect(), db, "migrations")
	if err != nil {
		panic(err)
	}

	command := os.Args[1]
	if command == "UA" {

		migrations, err := os.ReadDir("migrations")
		if err != nil {
			panic(err)
		}

		for _, file := range migrations {
			version := strings.TrimSuffix(file.Name(), ".sql")
			up, err := migrator.Up(version)
			if err != nil {
				panic(err)
			}
			log.Println("Success up migration: ", up)
		}
		return
	}
	name := os.Args[2]
	switch command {
	case "Up", "up", "UP":
		up, err := migrator.Up(name)
		if err != nil {
			panic(err)
		}
		log.Println("Migration Up:", up)
	case "Down", "DOWN", "down":
		down, err := migrator.Down(name)
		if err != nil {
			panic(err)
		}
		log.Println("Migration down:", down)
	case "Create ", "CREATE", "create":
		migration, err := migrator.CreateMigration(name)
		if err != nil {
			panic(err)
		}
		log.Println("Migrations created:", migration)
		log.Println(`For split command in SQL use 
-- +Then
delimiter between commands`)
	}
}
