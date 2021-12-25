package marmot

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var migrationTemplate = `-- +Up

-- +Down

-- +End`

type migrator struct {
	dialect   dbDialect
	db        *sql.DB
	directory string
	driver    string
}

func NewMigrator(dialect dbDialect, db *sql.DB, directory string) *migrator {
	m := &migrator{
		dialect:   dialect,
		db:        db,
		directory: directory,
	}
	m.prepare()
	return m
}

type state string

const (
	Up   state = "UP"
	Down state = "DOWN"
)

func (m *migrator) getExecutedMigrations() map[string]bool {
	executedMigrations, err := m.db.Query(m.dialect.allMigrationsSql())
	if err != nil {
		log.Fatalf("Error when try get executed migrations: %s", err)
	}
	executed := make(map[string]bool)
	for executedMigrations.Next() {
		var name string
		err := executedMigrations.Scan(&name)
		if err != nil {
			log.Fatalf("Error when try scan migrations table: %s", err)
		}
		executed[name] = true
	}
	return executed
}

func (m *migrator) prepare() {
	_, err := m.db.Exec(m.dialect.createVersionTableSql())
	if err != nil {
		log.Fatalf("Can't create version table: %s", err)
	}

}

func (m *migrator) Up() {
	//Collect not executed migrations
	executedMigration := m.getExecutedMigrations()
	var files []string
	_ = filepath.Walk(m.directory, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			name := strings.Split(path, "/")[1]
			if executedMigration[name] != true {
				files = append(files, path)
			}
		}
		sort.Strings(files)
		return nil
	})
	for _, file := range files {
		m.UpSingle(file)
	}
}

func (m *migrator) UpSingle(name string) {
	filePath := m.directory + "/" + name
	migrationCommands := getCommandsByFile(filePath)[Up]

	transaction, err := m.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for i := range migrationCommands {
		_, err = transaction.Exec(migrationCommands[i])
		if err != nil {
			_ = transaction.Rollback()
			log.Fatalf("Up migration error: %s", err)
		}
	}

	_, err = transaction.Exec(m.dialect.addVersionSql(name))
	if err != nil {
		_ = transaction.Rollback()
		log.Fatalf("Add version error: %s", err)
	}

	err = transaction.Commit()
	if err != nil {
		_ = transaction.Rollback()
		log.Fatalf("Can't commit up migration changes: %s", err)
	}
	fmt.Println("Migration ", name, " successful added.")
}

func (m *migrator) DownSingle(fileName string) {
	migrationCommands := getCommandsByFile(path.Join(m.directory, fileName))[Down]
	transaction, err := m.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	for i := range migrationCommands {
		_, err = transaction.Exec(migrationCommands[i])
		if err != nil {
			_ = transaction.Rollback()
			log.Fatalf("Down migration error: %s", err)
		}
	}

	_, err = transaction.Exec(m.dialect.deleteVersionSql(fileName))
	if err != nil {
		_ = transaction.Rollback()
		log.Fatalf("Devele version error: %s", err)
	}

	err = transaction.Commit()
	if err != nil {
		_ = transaction.Rollback()
		log.Fatalf("Can't commit migration down changes: %s", err)
	}
	fmt.Println("Migration ", fileName, " successful removed.")
}

func (m *migrator) CreateMigration(description string) {
	//prepare description
	extension := ".sql"
	if len(description) < 1 {
		description = extension
	} else {
		description = fmt.Sprintf(".%s%s", description, extension)
	}

	//check or create directory
	if _, err := os.Stat(m.directory); os.IsNotExist(err) {
		err := os.Mkdir(m.directory, os.ModePerm)
		if err != nil {
			log.Fatalf("Can't create directory %s: %s", m.directory, err)
		}
	}

	//create migration
	dst := time.Now().Format("2006-01-02T15-04-05")
	fileName := path.Join(m.directory, dst+description)
	err := os.WriteFile(fileName, []byte(migrationTemplate), 0666)
	if err != nil {
		log.Fatalf("Error when try create new migration: %s", err)
	}
	fmt.Println("Success create migration", dst+description)
}
