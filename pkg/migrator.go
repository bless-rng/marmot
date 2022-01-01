package marmot

import (
	"database/sql"
	"fmt"
	"github.com/bless-rng/marmot/internal/commands"
	"github.com/bless-rng/marmot/internal/dialect"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

var migrationTemplate = `-- +Up

-- +Down

-- +End`

type migrator struct {
	sqlDialect dialect.SqlDialect
	db         *sql.DB
	directory  string
	driver     string
}

func NewMigrator(dialect dialect.SqlDialect, db *sql.DB, directory string) (*migrator, error) {
	m := &migrator{
		sqlDialect: dialect,
		db:         db,
		directory:  directory,
	}
	err := m.prepare()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *migrator) prepare() error {
	_, err := m.db.Exec(m.sqlDialect.CreateVersionTableSql())
	if err != nil {
		return err
	}
	return nil
}

func (m *migrator) CreateMigration(description string) (migration string, error error) {
	description = strings.Join(strings.Fields(description), " ")
	if len(description) > 200 {
		panic("Max description len is 200")
	}
	if len(description) > 0 {
		description = "." + strings.Replace(description, " ", "_", -1)
	}

	migrationName := getMigrationPrefix() + description + ".sql"

	//check or create directory
	if _, err := os.Stat(m.directory); os.IsNotExist(err) {
		err := os.Mkdir(m.directory, os.ModePerm)
		if err != nil {
			return migrationName, err
		}
	}

	fileName := path.Join(m.directory, migrationName)
	err := os.WriteFile(fileName, []byte(migrationTemplate), 0666)
	if err != nil {
		return migrationName, err
	}
	return migrationName, nil
}

func getMigrationPrefix() string {
	return time.Now().Format("20060102-150405")
}

func (m *migrator) Up(migrationName string) (migration string, error error) {
	filePath := m.directory + "/" + migrationName + ".sql"

	migrationCommands, err := commands.GetCommandsByFile(filePath, commands.Up)
	if err != nil {
		return migrationName, err
	}

	transaction, err := m.db.Begin()
	if err != nil {
		return migrationName, err
	}

	_, err = transaction.Exec(m.sqlDialect.AddVersionSql(migrationName))
	if err != nil {
		return migrationName, err
	}

	for i, command := range migrationCommands {
		log.Println(fmt.Sprintf("Try execute command #%d: %s", i+1, command))
		_, err = m.db.Exec(command)
		if err != nil {
			_ = transaction.Rollback()
			return migrationName, err
		}
	}

	_, err = transaction.Exec(m.sqlDialect.VerifyVersion(migrationName))
	if err != nil {
		_ = transaction.Rollback()
		return migrationName, err
	}

	err = transaction.Commit()
	if err != nil {
		_ = transaction.Rollback()
		return migrationName, err
	}
	return migrationName, nil
}

func (m *migrator) Down(migration string) (version string, error error) {
	migrationCommands, err := commands.GetCommandsByFile(path.Join(m.directory, migration+".sql"), commands.Down)

	var v interface{}
	err = m.db.QueryRow(m.sqlDialect.FindMigration(migration)).Scan(&v)
	if err != nil {
		return migration, fmt.Errorf("migratoin %s does not exist it list of migrations", migration)
	}

	transaction, err := m.db.Begin()
	if err != nil {
		return migration, err
	}

	for i, command := range migrationCommands {
		log.Println(fmt.Sprintf("Try execute command #%d: %s", i+1, command))
		_, err = transaction.Exec(command)
		if err != nil {
			_ = transaction.Rollback()
			return migration, err
		}
	}

	_, err = transaction.Exec(m.sqlDialect.DeleteVersionSql(migration))
	if err != nil {
		_ = transaction.Rollback()
		return migration, err
	}

	err = transaction.Commit()
	if err != nil {
		_ = transaction.Rollback()
		return migration, err
	}
	return migration, nil
}
