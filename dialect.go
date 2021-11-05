package marmot

import "fmt"

type dbDialect interface {
	createVersionTableSql() string
	addVersionSql(name string) string
	deleteVersionSql(name string) string
	allMigrationsSql() string
}

var versionTable = "marmot_migrations"

var (
	CreateVersionTableTmpl = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(id %s, name varchar(256), executed_at TIMESTAMP DEFAULT NOW());", versionTable, "%s")
	InsertVersionTmpl      = fmt.Sprintf("INSERT INTO %s(name) VALUES('%s');", versionTable, "%s")
	DeleteVersionTmpl      = fmt.Sprintf("DELETE FROM %s WHERE name = '%s';", versionTable, "%s")
	GetListOfMigrations    = fmt.Sprintf("SELECT name FROM %s order by executed_at DESC;", versionTable)
)

type dialect struct{}

func (db dialect) addVersionSql(name string) string {
	return fmt.Sprintf(InsertVersionTmpl, name)
}

func (db dialect) deleteVersionSql(name string) string {
	return fmt.Sprintf(DeleteVersionTmpl, name)
}

func (db dialect) allMigrationsSql() string {
	return GetListOfMigrations
}

type MysqlDialect struct {
	dialect
}

func (db MysqlDialect) createVersionTableSql() string {
	return fmt.Sprintf(CreateVersionTableTmpl, "INT AUTO_INCREMENT")
}

type PostgresqlDialect struct {
	dialect
}

func (db PostgresqlDialect) createVersionTableSql() string {
	return fmt.Sprintf(CreateVersionTableTmpl, "SERIAL")
}

type SQLiteDialect struct {
	dialect
}


func (db SQLiteDialect) createVersionTableSql() string {
	return fmt.Sprintf(CreateVersionTableTmpl, "AUTOINCREMENT")
}

type MsSqlDialect struct {
	dialect
}


func (db MsSqlDialect) createVersionTableSql() string {
	return fmt.Sprintf(CreateVersionTableTmpl, "AUTOINCREMENT")
}
