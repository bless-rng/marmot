package dialect

import "fmt"

type postgreSQLDialect struct {
	tableName string
}

func NewPostgreSQLDialect() *postgreSQLDialect {
	return &postgreSQLDialect{
		tableName: "marmot_migrations",
	}
}

func (d *postgreSQLDialect) VerifyVersion(version string) string {
	return fmt.Sprintf("UPDATE %s SET migrated_at = NOW()", d.tableName)
}

func (d *postgreSQLDialect) CreateVersionTableSql() string {
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(id INT GENERATED BY DEFAULT AS IDENTITY, version VARCHAR(256) UNIQUE NOT NULL, migrated_at TIMESTAMP)", d.tableName)
}

func (d *postgreSQLDialect) AddVersionSql(version string) string {
	return fmt.Sprintf("INSERT INTO %s(version) VALUES ('%s')", d.tableName, version)
}

func (d *postgreSQLDialect) DeleteVersionSql(version string) string {
	return fmt.Sprintf("DELETE FROM %s where version = '%s'", d.tableName, version)
}

func (d *postgreSQLDialect) FindMigration(version string) string {
	return fmt.Sprintf("SELECT version FROM %s where version = '%s'", d.tableName, version)
}

func (d *postgreSQLDialect) AllMigrationsSql() string {
	return fmt.Sprintf("SELECT version FROM %s", d.tableName)
}