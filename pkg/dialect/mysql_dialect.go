package dialect

import (
	"fmt"
)

type mySQLDialect struct {
	tableName string
}

func NewMySQLDialect() *mySQLDialect {
	return &mySQLDialect{
		tableName: "marmot_migrations",
	}
}

func (d *mySQLDialect) CreateVersionTableSql() string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(id INT NOT NULL AUTO_INCREMET, version VARCHAR(256) UNIQUE NOT NULL, migrated_at TIMESTAMP)", d.tableName)
}

func (d *mySQLDialect) VerifyVersion(version string) string {
	return fmt.Sprintf("UPDATE %s SET migrated_at = NOW()", d.tableName)
}

func (d *mySQLDialect) AddVersionSql(version string) string {
	return fmt.Sprintf("INSERT INTO %s(version) VALUES ('%s')", d.tableName, version)
}

func (d *mySQLDialect) FindMigration(version string) string {
	return fmt.Sprintf("SELECT version FROM %s where version = '%s'", d.tableName, version)
}

func (d *mySQLDialect) DeleteVersionSql(version string) string {
	return fmt.Sprintf("DELETE FROM %s where version = '%s'", d.tableName, version)
}

func (d *mySQLDialect) AllMigrationsSql() string {
	return fmt.Sprintf("SELECT version FROM %s", d.tableName)
}
