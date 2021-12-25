package main

import (
	"database/sql"
	"github.com/bless-rng/marmot"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3606)/database")
	if err != nil {
		panic(err)
	}

	migrator := marmot.NewMigrator(marmot.MysqlDialect{}, db, "migrations")
	command := os.Args[1]
	log.Println(os.Args)
	var name string
	if len(os.Args) == 3 {
		name = os.Args[2]
	}

	switch command {
	case "up":
		if len(name) > 0 {
			migrator.UpSingle(name)
			break
		}
		migrator.Up()
	case "down":
		if len(name) < 1 {
			panic(err)
		}
		migrator.DownSingle(name)
	case "new":
		migrator.CreateMigration(os.Args[2])
	default:
		break

	}
}
